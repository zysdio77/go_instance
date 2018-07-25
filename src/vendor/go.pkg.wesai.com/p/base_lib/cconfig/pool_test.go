package cconfig

import (
	"testing"
	"time"
)

func getPool(
	maxActive uint64,
	maxIdle uint64,
	idleTimeout time.Duration) (Pool, error) {

	clientConfig := ClientConfig{
		ProjectName: PROJECT_NAME_FOR_CLIENT_TEST,
		Environment: ENV_DEVELOPMENT,
		Addresses:   GetCCAddressFromSys(),
		Username:    USERNAME_FOR_TEST,
		Password:    PASSWORD_FOR_TEST,
		NotGuest:    true,
	}
	poolConfig := PoolConfig{
		ConfigCenterType: CC_TYPE_ETCD,
		MaxActive:        maxActive,
		MaxIdle:          maxIdle,
		IdleTimeout:      idleTimeout,
		ClientConfig:     clientConfig,
	}
	return NewPool(poolConfig)
}

func TestNewPool(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(2)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	if pool.Destroyed() {
		t.Fatalf("Expect %v, but the actual is %v\n",
			false, pool.Destroyed())
	}
	if pool.Active() != maxIdle {
		t.Fatalf("Expect %d, but the actual is %d\n",
			maxIdle, pool.Active())
	}
	if pool.MaxActive() != maxActive {
		t.Fatalf("Expect %d, but the actual is %d\n",
			maxActive, pool.MaxActive())
	}
	if pool.MaxIdle() != maxIdle {
		t.Fatalf("Expect %d, but the actual is %d\n",
			maxIdle, pool.MaxIdle())
	}
	if pool.IdleTimeout() != idleTimeout {
		t.Fatalf("Expect %s, but the actual is %s\n",
			idleTimeout, pool.IdleTimeout())
	}
}

func TestDestroy(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(2)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	err = pool.Destroy()
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	if !pool.Destroyed() {
		t.Fatalf("Expect %v, but the actual is %v\n",
			true, pool.Destroyed())
	}
	if pool.Active() != 0 {
		t.Fatalf("Expect %d, but the actual is %d\n",
			0, pool.Active())
	}
	cpool := pool.(*clientPool)
	idleNum := cpool.idleList.Len()
	if uint64(idleNum) != 0 {
		t.Fatalf("Expect %d, but the actual is %d\n",
			0, idleNum)
	}
}

func TestIdleNumAfterNewPool(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(3)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	cpool := pool.(*clientPool)
	idleNum := cpool.idleList.Len()
	if uint64(idleNum) != maxIdle {
		t.Fatalf("Expect %d, but the actual is %d\n",
			maxIdle, idleNum)
	}
}

func TestGetWithMaxActive(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(2)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	var i uint64
	for i = 0; i < maxActive; i++ {
		client, err := pool.Get(false)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		if client == nil {
			t.Fatalf("Invalid client! (index: %d)", i)
		}
	}
	_, err = pool.Get(false)
	if err != nil {
		if err != ErrPoolExhausted {
			t.Fatalf("Occur an error: %s", err)
		}
	} else {
		t.Fatal("An exhausted pool error should occur!")
	}
}

func TestIdleNumAfterGetAndClientClose(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(3)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	cpool := pool.(*clientPool)
	var i uint64
	clients := make([]Client, 0)
	for i = 0; i < maxActive; i++ {
		client, err := pool.Get(false)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		if client == nil {
			t.Fatalf("Invalid client! (index: %d)", i)
		}
		idleNum := cpool.idleList.Len()
		if i < maxIdle {
			diff := (maxIdle - i - 1)
			if uint64(idleNum) != diff {
				t.Fatalf("Expect %d, but the actual is %d\n",
					diff, idleNum)
			}
		} else {
			expectedIdleNum := 0
			if idleNum != expectedIdleNum {
				t.Fatalf("Expect %d, but the actual is %d\n",
					expectedIdleNum, idleNum)
			}
		}
		clients = append(clients, client)
	}

	for i, client := range clients {
		client.Close()
		idleNum := cpool.idleList.Len()
		if uint64(i) < maxIdle {
			expectedIdleNum := i + 1
			if idleNum != expectedIdleNum {
				t.Fatalf("Expect %d, but the actual is %d\n",
					expectedIdleNum, idleNum)
			}
		} else {
			if uint64(idleNum) != maxIdle {
				t.Fatalf("Expect %d, but the actual is %d\n",
					maxIdle, idleNum)
			}
		}
	}
}

func TestActiveNumAfterGetAndClientClose(t *testing.T) {
	maxActive := uint64(12)
	maxIdle := uint64(5)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	cpool := pool.(*clientPool)
	var i uint64
	clients := make([]Client, 0)
	for i = 0; i < maxActive; i++ {
		client, err := pool.Get(false)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		if client == nil {
			t.Fatalf("Invalid client! (index: %d)", i)
		}
		if i < maxIdle {
			if cpool.active != maxIdle {
				t.Fatalf("Expect %d, but the actual is %d\n",
					maxIdle, cpool.active)
			}
		} else {
			expectedActiveNum := uint64(i + 1)
			if cpool.active != expectedActiveNum {
				t.Fatalf("Expect %d, but the actual is %d\n",
					expectedActiveNum, cpool.active)
			}
		}
		clients = append(clients, client)
	}

	for i, client := range clients {
		client.Close()
		if uint64(i+1) < maxIdle {
			if cpool.active != maxActive {
				t.Fatalf("Expect %d, but the actual is %d\n",
					maxActive, cpool.active)
			}
		} else {
			expectedActiveNum := maxActive - (uint64(i+1) - maxIdle)
			if cpool.active != expectedActiveNum {
				t.Errorf("Expect %d, but the actual is %d\n",
					expectedActiveNum, cpool.active)
			}
		}
	}
}

func TestIdleTimeout(t *testing.T) {
	maxActive := uint64(7)
	maxIdle := uint64(2)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	var i uint64
	clients := make([]Client, 0)
	for i = 0; i < maxActive; i++ {
		client, err := pool.Get(false)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		if client == nil {
			t.Fatalf("Invalid client! (index: %d)", i)
		}
		clients = append(clients, client)
	}
	cpool := pool.(*clientPool)
	for _, client := range clients {
		cpool.put(client.(*clientWrapper), false)
	}

	time.Sleep(idleTimeout) // 等待空闲链接超时。
	client, err := pool.Get(false)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	if client == nil {
		t.Fatalf("Invalid client! (index: %d)", i)
	}
	if cpool.active != maxIdle {
		t.Fatalf("Expect %d, but the actual is %d\n",
			maxIdle, cpool.active)
	}
}

func TestBlockedGet(t *testing.T) {
	maxActive := uint64(5)
	maxIdle := uint64(3)
	idleTimeout := time.Second
	pool, err := getPool(maxActive, maxIdle, idleTimeout)
	if err != nil {
		t.Fatalf("Occur an error: %s", err)
	}
	defer func() {
		pool.Destroy()
	}()

	var i uint64
	var lastClient Client
	for i = 0; i < maxActive; i++ {
		client, err := pool.Get(false)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		if client == nil {
			t.Fatalf("Invalid client! (index: %d)", i)
		}
		lastClient = client
	}

	clientCh := make(chan Client, 1)
	go func() {
		client, err := pool.Get(true)
		if err != nil {
			t.Fatalf("Occur an error: %s", err)
		}
		clientCh <- client
	}()
	time.Sleep(time.Millisecond * 10)
	lastClient.Close()
	client := <-clientCh
	if client == nil {
		t.Fatalf("Invalid client! (index: %d)", i)
	}
	if client.ID() != lastClient.ID() {
		t.Fatalf("Expect %d, but the actual is %d\n",
			lastClient.ID(), client.ID())
	}
}
