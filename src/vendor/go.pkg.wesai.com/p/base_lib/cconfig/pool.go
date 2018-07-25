package cconfig

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 池配置。
type PoolConfig struct {
	ConfigCenterType ConfigCenterType // 配置中心的类型。
	MaxActive        uint64           // 池持有的连接/客户端的最大数量。
	MaxIdle          uint64           // 池能持有的空闲的连接/客户端的最大数量。
	IdleTimeout      time.Duration    // 空闲连接/客户端的超时时间。超时的空闲连接/客户端会被清除。
	ClientConfig     ClientConfig     // 客户端配置。
}

// 检查自身。
func (pconfig *PoolConfig) Check() error {
	switch pconfig.ConfigCenterType {
	case CC_TYPE_ETCD, CC_TYPE_CONSUL:
	case CC_TYPE_ZK:
		return errors.New(
			fmt.Sprintf("Oops! The config center for type '%s' has not been implemented!",
				pconfig.ConfigCenterType))
	default:
		return errors.New(
			fmt.Sprintf("Unsupported config center type '%s'!",
				pconfig.ConfigCenterType))
	}
	if pconfig.MaxActive == 0 {
		return errors.New("The max active is zero!")
	}
	if pconfig.MaxIdle > pconfig.MaxActive {
		return errors.New("Overlarge max idle!")
	}
	return pconfig.ClientConfig.Check()
}

// 代表用于新建客户端的函数的类型。
type newFunc func(id string) (Client, error)

// 代表用于生成客户端生成器的函数的类型。
type genNewFunc func(pconfig PoolConfig) newFunc

// 客户端池的接口。
type Pool interface {
	Get(blocked bool) (Client, error)
	Destroy() error
	Destroyed() bool
	Active() uint64
	MaxActive() uint64
	MaxIdle() uint64
	IdleTimeout() time.Duration
}

// 客户端池的实现。
type clientPool struct {
	// 外部定义的参数。
	maxActive   uint64
	maxIdle     uint64
	idleTimeout time.Duration
	// 内部定义的参数。
	newFunc  newFunc
	idPrefix string
	idSn     uint64
	// 内部状态。
	idleList  *list.List
	active    uint64
	destroyed bool
	// 内部并发保护。
	rwmutex sync.RWMutex
	cond    *sync.Cond
}

// 新建客户端池。
func NewPool(pconfig PoolConfig) (Pool, error) {
	if err := pconfig.Check(); err != nil {
		return nil, err
	}
	var genFunc genNewFunc
	switch pconfig.ConfigCenterType {
	case CC_TYPE_ETCD:
		genFunc = generateEtcdClientGenerator
	case CC_TYPE_CONSUL:
		genFunc = generateConsulClientGenerator
	}
	if genFunc == nil {
		return nil, errors.New("Can not generate the client generating function!")
	}
	env := pconfig.ClientConfig.Environment
	projectName := pconfig.ClientConfig.ProjectName
	idPrefix := string(pconfig.ConfigCenterType) + "-" +
		string(env) + "-" + projectName + "-"
	pool := &clientPool{
		maxActive:   pconfig.MaxActive,
		maxIdle:     pconfig.MaxIdle,
		idleTimeout: pconfig.IdleTimeout,
		newFunc:     genFunc(pconfig),
		idPrefix:    idPrefix,
		idleList:    list.New(),
		rwmutex:     sync.RWMutex{},
	}
	err := pool.initialize()
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// etcd客户端生成器的生成函数。
func generateEtcdClientGenerator(pconfig PoolConfig) newFunc {
	clientConfig := pconfig.ClientConfig
	return func(id string) (Client, error) {
		client, err := NewEtcdClient(id, clientConfig)
		if err != nil {
			errMsg := fmt.Sprintf(
				"Can not new an etcd client"+
					" (id: %s, projectName: %s, env: %s, addresses: [%s]): %s",
				id,
				pconfig.ClientConfig.ProjectName,
				pconfig.ClientConfig.Environment,
				strings.Join(pconfig.ClientConfig.Addresses, ","),
				err)
			return nil, errors.New(errMsg)
		}
		return client, nil
	}
}

func generateConsulClientGenerator(pconfig PoolConfig) newFunc {
	clientConfig := pconfig.ClientConfig
	return func(id string) (Client, error) {
		client, err := NewConsulClient(id, clientConfig)
		if err != nil {
			errMsg := fmt.Sprintf(
				"Can not new an Consul client"+
					" (id: %s, projectName: %s, env: %s, addresses: [%s]): %s",
				id,
				pconfig.ClientConfig.ProjectName,
				pconfig.ClientConfig.Environment,
				strings.Join(pconfig.ClientConfig.Addresses, ","),
				err)
			return nil, errors.New(errMsg)
		}
		return client, nil
	}
}

// 初始化池。非并发安全方法。
func (pool *clientPool) initialize() error {
	pool.clear()
	num := pool.maxIdle
	var i uint64
	for i = 1; i <= num; i++ {
		wrapper, err := pool.newClientWrapper()
		if err != nil {
			return err
		}
		pool.active++
		pool.idleList.PushBack(wrapper)
	}
	return nil
}

// 新建客户端包装器。非并发安全方法。
func (pool *clientPool) newClientWrapper() (*clientWrapper, error) {
	sn := atomic.AddUint64(&pool.idSn, 1)
	id := pool.idPrefix + strconv.FormatUint(sn, 10)
	client, err := pool.newFunc(id)
	if err != nil {
		return nil, err
	}
	wrapper := &clientWrapper{
		pool:      pool,
		client:    client,
		activated: false,
		idleTs:    time.Now().UnixNano(),
	}
	return wrapper, nil
}

// 清理池。非并发安全方法。
func (pool *clientPool) clear() {
	for {
		e := pool.idleList.Front()
		if e == nil {
			break
		}
		pool.idleList.Remove(e)
		wrapper := e.Value.(*clientWrapper)
		wrapper.close()
		if pool.active > 0 {
			pool.active--
		}
	}
}

// 释放掉（且即将关闭）一个客户端。调用该方法前必须加锁。
// 如果需要让该方法代为关闭客户端，那么就把该客户端传入。
func (pool *clientPool) releaseOne(wrapper *clientWrapper) {
	if wrapper != nil {
		wrapper.close()
	}
	pool.active--
	if pool.cond != nil {
		pool.cond.Signal()
	}
}

// 检查超时的空闲链接。已加锁。调用方不要重复加锁。
func (pool *clientPool) checkStaleClients() {
	if timeout := pool.idleTimeout; timeout > 0 {
		pool.rwmutex.Lock()
		defer pool.rwmutex.Unlock()

		diff := uint64(pool.idleList.Len()) - pool.maxIdle
		if diff <= 0 {
			return
		}
		for i := uint64(0); i < diff; i++ {
			e := pool.idleList.Front()
			if e == nil {
				break
			}
			wrapper := e.Value.(*clientWrapper)
			if !wrapper.overtime(timeout) {
				break
			}
			pool.idleList.Remove(e)
			pool.releaseOne(nil)
			pool.rwmutex.Unlock()
			wrapper.close()
			pool.rwmutex.Lock()
		}
	}
}

func (pool *clientPool) Get(blocked bool) (Client, error) {
	pool.checkStaleClients()

	pool.rwmutex.Lock()
	defer pool.rwmutex.Unlock()

	for {
		// 尝试从空闲列表中获取客户端。
		e := pool.idleList.Front()
		if e != nil {
			wrapper := pool.idleList.Remove(e).(*clientWrapper)
			wrapper.activate()
			return wrapper, nil
		}
		// 新建客户端之前先检查一下池的状态。
		if pool.destroyed {
			return nil, ErrClosedPool
		}
		// 新建一个客户端。
		maxActive := pool.maxActive
		if maxActive == 0 || pool.active < pool.maxActive {
			pool.active++
			pool.rwmutex.Unlock()
			wrapper, err := pool.newClientWrapper()
			pool.rwmutex.Lock()
			if err != nil {
				pool.releaseOne(nil)
				return nil, err
			}
			wrapper.activate()
			return wrapper, nil
		}
		// 如果不是阻塞式的获取，那么立即返回错误。
		if !blocked {
			return nil, ErrPoolExhausted
		}
		// 等待重试条件。
		if pool.cond == nil {
			pool.cond = sync.NewCond(&pool.rwmutex)
		}
		pool.cond.Wait()
	}
}

func (pool *clientPool) put(wrapper *clientWrapper, checkIdle bool) error {
	if wrapper.closed {
		pool.rwmutex.Lock()
		pool.active--
		if pool.cond != nil {
			pool.cond.Signal()
		}
		pool.rwmutex.Unlock()
		panic(errClosedClient)
	}

	pool.rwmutex.Lock()
	defer pool.rwmutex.Unlock()

	if pool.destroyed {
		wrapper.close()
		return nil
	}
	wrapper.inactivate()
	pool.idleList.PushBack(wrapper)
	if pool.cond != nil {
		pool.cond.Signal()
	}
	if checkIdle && uint64(pool.idleList.Len()) > pool.maxIdle {
		wrapper = pool.idleList.Remove(pool.idleList.Front()).(*clientWrapper)
		pool.releaseOne(wrapper)
	}
	return nil
}

func (pool *clientPool) Destroy() error {
	pool.rwmutex.Lock()
	defer pool.rwmutex.Unlock()
	if pool.destroyed {
		return nil
	}
	pool.destroyed = true
	pool.clear()
	if pool.cond != nil {
		pool.cond.Broadcast()
	}
	return nil
}

func (pool *clientPool) Destroyed() bool {
	pool.rwmutex.RLock()
	defer pool.rwmutex.RUnlock()
	return pool.destroyed
}

func (pool *clientPool) Active() uint64 {
	pool.rwmutex.RLock()
	defer pool.rwmutex.RUnlock()
	return pool.active
}

func (pool *clientPool) MaxActive() uint64 {
	return pool.maxActive
}

func (pool *clientPool) MaxIdle() uint64 {
	return pool.maxIdle
}

func (pool *clientPool) IdleTimeout() time.Duration {
	return pool.idleTimeout
}
