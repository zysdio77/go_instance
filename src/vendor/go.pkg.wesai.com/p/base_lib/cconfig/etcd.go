package cconfig

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"go.pkg.wesai.com/p/base_lib/log"
)

var logger = log.DLogger()

// etcd的节点。
type EtcdNode struct {
	path      string
	isDir     bool
	value     string
	children  []Node
	recursive bool
	sorted    bool
}

func (en *EtcdNode) Path() string {
	return en.path
}

func (en *EtcdNode) IsDir() bool {
	return en.isDir
}

func (en *EtcdNode) Value() string {
	return en.value
}

func (en *EtcdNode) Children() []Node {
	return en.children
}

func (en *EtcdNode) Recursive() bool {
	return en.recursive
}

func (en *EtcdNode) Sorted() bool {
	return en.sorted
}

// etcd的节点变更。
type EtcdNodeChange struct {
	action string
	path   string
	isDir  bool
	value  string
	prev   *EtcdNodeChange
}

func (enc *EtcdNodeChange) Action() string {
	return enc.action
}

func (enc *EtcdNodeChange) Path() string {
	return enc.path
}

func (enc *EtcdNodeChange) IsDir() bool {
	return enc.isDir
}

func (enc *EtcdNodeChange) Value() string {
	return enc.value
}

func (enc *EtcdNodeChange) Prev() NodeChange {
	return enc.prev
}

type EtcdClient struct {
	id           string
	etcdAPI      client.KeysAPI
	env          Env
	rootPath     string
	innerContext context.Context
	cancelFunc   context.CancelFunc
}

func NewEtcdClient(id string, config ClientConfig) (eclient *EtcdClient, err error) {

	if err := config.Check(); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && eclient != nil {
			eclient.Close()
			eclient = nil
		}
	}()

	id = strings.TrimSpace(id)
	if id == "" {
		id = config.ProjectName
	}
	env := config.Environment
	rootPath, err := GetRootPath(env, config.ProjectName)
	if err != nil {
		return
	}

	var username, password string
	if config.NotGuest {
		username = config.Username
		password = config.Password
	}
	etcdConfig := client.Config{
		Endpoints:               config.Addresses,
		Transport:               client.DefaultTransport,
		Username:                username,
		Password:                password,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("Can not new a etcd client: %s", err)
	}
	kapi := client.NewKeysAPI(etcdClient)
	innerContext, cancelFunc := context.WithCancel(context.Background())

	eclient = &EtcdClient{
		id:           id,
		etcdAPI:      kapi,
		env:          env,
		rootPath:     rootPath,
		innerContext: innerContext,
		cancelFunc:   cancelFunc,
	}
	var exist bool
	exist, _, err = eclient.Exist(rootPath)
	if exist, _, err = eclient.Exist(rootPath); err == nil && !exist {
		if err = eclient.CreateDir(rootPath); err != nil {
			return
		}
	} else {
		if err != nil {
			return
		}
	}
	return eclient, nil
}

// 浅表复制etcd的Node并生成公用的Node实例。
func copyEtcdNodeShallow(node *client.Node, sorted, recursive bool) *EtcdNode {

	if node == nil {
		return nil
	}
	return &EtcdNode{
		path:      node.Key,
		isDir:     node.Dir,
		value:     node.Value,
		recursive: recursive,
		sorted:    sorted,
	}
}

// 复制etcd的Node并生成公用的Node实例。
func copyEtcdNode(etcdNode *client.Node, sorted, recursive bool) *EtcdNode {

	if etcdNode == nil {
		return nil
	}
	node := copyEtcdNodeShallow(etcdNode, sorted, recursive)
	var children []Node
	subNodes := etcdNode.Nodes
	if subNodes != nil {
		subNodesLen := len(subNodes)
		children = make([]Node, subNodesLen)
		for i := 0; i < subNodesLen; i++ {
			children[i] = copyEtcdNode(subNodes[i], sorted, recursive)
		}
	}
	node.children = children
	return node
}

// 确保路径合法。
func (eclient *EtcdClient) ensurePath(nodePath string) (string, error) {
	if len(nodePath) == 0 {
		return "", ErrEmptyPath
	}
	if strings.HasPrefix(nodePath, root_base) {
		if !strings.HasPrefix(nodePath, eclient.rootPath) {
			return "", &ErrIllegalPath{nodePath, eclient.rootPath}
		}
		return nodePath, nil
	}
	nodePath = path.Clean(nodePath)
	if !strings.HasPrefix(nodePath, "/") {
		nodePath = "/" + nodePath
	}
	nodePath = eclient.rootPath + nodePath
	return nodePath, nil
}

func (eclient *EtcdClient) ID() string {
	return eclient.id
}

func (eclient *EtcdClient) RootPath() string {
	return eclient.rootPath
}

func (eclient *EtcdClient) Exist(nodePath string) (exist bool, dir bool, err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	var node Node
	node, err = eclient.Get(nodePath, false)
	if err != nil {
		if _, ok := err.(*ErrInexistentPath); ok {
			err = nil
		}
		return
	}
	if node == nil {
		return
	}
	exist = true
	dir = node.IsDir()
	return
}

func (eclient *EtcdClient) Get(nodePath string, recursive bool) (Node, error) {
	nodePath, err := eclient.ensurePath(nodePath)
	if err != nil {
		return nil, err
	}
	sorted := recursive
	opts := &client.GetOptions{
		Recursive: recursive,
		Sort:      sorted,
	}
	var resp *client.Response
	resp, err = eclient.etcdAPI.Get(eclient.innerContext, nodePath, opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeKeyNotFound {
				err = &ErrInexistentPath{nodePath}
			}
		}
		return nil, err
	} else {
		return copyEtcdNode(resp.Node, sorted, recursive), nil
	}
}

// 删除一个已存在的路径。
// 当该路径不存在时，该方法会返回ErrInexistentPath类型的错误。
func (eclient *EtcdClient) Delete(nodePath string, recursive bool) (done bool, err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.DeleteOptions{
		Recursive: recursive,
	}
	_, err = eclient.etcdAPI.Delete(eclient.innerContext, nodePath, opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeKeyNotFound {
				err = &ErrInexistentPath{nodePath}
			}
		}
	} else {
		done = true
	}
	return
}

func (eclient *EtcdClient) Watch(
	path string,
	recursive bool,
	changeChan chan NodeChange,
	errorChan chan error,
	stopChan chan bool) {

	var err error
	closeChannels := func() {
		close(changeChan)
		close(errorChan)
	}
	defer func() {
		if p := recover(); p != nil {
			var ok bool
			err, ok = p.(error)
			if !ok {
				err = errors.New(fmt.Sprintf("%+v", p))
			}
		}
		if err != nil {
			errorChan <- err
			closeChannels()
		}
	}()

	path, err = eclient.ensurePath(path)
	if err != nil {
		return
	}

	respChan := make(chan *client.Response, 10)
	go eclient.watch(path, recursive, respChan, errorChan, false, stopChan)
	go func() {
		defer closeChannels()
		for resp := range respChan {
			nodeChange := eclient.parseWatcherResp(resp)
			if nodeChange == nil {
				continue
			}
			changeChan <- nodeChange
		}
	}()
}

// 解析观察器的响应、
func (eclient *EtcdClient) parseWatcherResp(resp *client.Response) NodeChange {
	if resp == nil {
		return nil
	}
	node := resp.Node
	prevNode := resp.PrevNode
	nodeChange := &EtcdNodeChange{
		action: resp.Action,
	}
	if node != nil {
		nodeChange.path = node.Key
		nodeChange.isDir = node.Dir
		nodeChange.value = node.Value
	}
	if prevNode != nil {
		prev := &EtcdNodeChange{
			action: resp.Action,
			path:   prevNode.Key,
			isDir:  prevNode.Dir,
			value:  prevNode.Value,
		}
		nodeChange.prev = prev
	}
	return nodeChange
}

// 同步的单次或循环观察指定路径的变化。该方法返回时会关闭响应通道。
func (eclient *EtcdClient) watch(
	path string,
	recursive bool,
	respChan chan<- *client.Response,
	errorChan chan<- error,
	loop bool,
	stopChan <-chan bool) {

	defer close(respChan)
	opts := &client.WatcherOptions{
		Recursive: recursive,
	}
	watcher := eclient.etcdAPI.Watcher(path, opts)
	watcherContext, cancelFunc := context.WithCancel(eclient.innerContext)
	go func() {
		for {
			select {
			case <-stopChan:
				cancelFunc()
				break
			default:
			}
			time.Sleep(time.Millisecond)
		}
	}()
	for {
		resp, err := watcher.Next(watcherContext)
		if err != nil {
			if err == context.Canceled {
				err = ErrWatchStoppedByUser
			}
			errorChan <- err
			break
		}
		respChan <- resp
		if !loop {
			break
		}
	}
}

// 创建一个目录路径。
// 当该路径已存在时，该方法会返回ErrExistingPath类型的错误。
func (eclient *EtcdClient) CreateDir(nodePath string) (err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		Dir:       true,
	}
	_, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, "", opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeNodeExist {
				err = &ErrExistingPath{nodePath}
			}
		}
	}
	return
}

// 更新一个已存在的目录路径。
// 若该路径为目录路径，则只会更新该键路径的变更索引。
// 若该路径为键路径，则会删除该键路径的值，但并不会改变该路径的类型。
func (eclient *EtcdClient) UpdateDir(nodePath string) (err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevExist,
		Dir:       true,
	}
	_, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, "", opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeKeyNotFound {
				err = &ErrInexistentPath{nodePath}
			}
		}
	}
	return
}

// 创建一个新的目录，或者把一个键路径变成目录路径（键路径对应的值会被删除）。
// 当该目录路径已存在时，该方法会返回ErrExistingDir类型的错误。
func (eclient *EtcdClient) SetDir(nodePath string) (incremental bool, err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevIgnore,
		Dir:       true,
	}
	var resp *client.Response
	resp, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, "", opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeNotFile {
				err = &ErrExistingDir{nodePath}
			}
		}
		return
	}
	if resp.PrevNode == nil {
		incremental = true
	}
	return
}

// 创建一个键路径并设置其值。
// 当该路径已存在时，该方法会返回ErrExistingPath类型的错误。
func (eclient *EtcdClient) Create(nodePath string, value string) (err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		Dir:       false,
	}
	_, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, value, opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeNodeExist {
				err = &ErrExistingPath{nodePath}
			}
		}
	}
	return
}

// 更新一个已存在的键路径的值。
// 当该路径不存在时，该方法会返回ErrInexistentPath类型的错误。
// 当该路径代表了一个已存在的目录路径时，该方法会返回ErrExistingDir类型的错误。
func (eclient *EtcdClient) Update(nodePath string, value string) (err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevExist,
		Dir:       false,
	}
	_, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, value, opts)
	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeKeyNotFound {
				err = &ErrInexistentPath{nodePath}
			} else if eerr.Code == client.ErrorCodeNotFile {
				err = &ErrExistingDir{nodePath}
			}
		}
	}
	return
}

// 创建一个键路径并设置其值，或者更新一个已存在的键路径的值。
// 当该路径代表一个已存在的目录路径时，该方法会返回ErrExistingDir类型的错误。
func (eclient *EtcdClient) Set(nodePath string, value string) (incremental bool, err error) {
	nodePath, err = eclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	opts := &client.SetOptions{
		PrevExist: client.PrevIgnore,
		Dir:       false,
	}
	var resp *client.Response
	resp, err = eclient.etcdAPI.Set(eclient.innerContext, nodePath, value, opts)

	if err != nil {
		if eerr, ok := err.(client.Error); ok {
			if eerr.Code == client.ErrorCodeNotFile {
				err = &ErrExistingDir{nodePath}
			}
		}
		return
	}
	if resp.PrevNode == nil {
		incremental = true
	}
	return
}

func (eclient *EtcdClient) Close() {
	if eclient.cancelFunc != nil {
		eclient.cancelFunc()
	}
}
