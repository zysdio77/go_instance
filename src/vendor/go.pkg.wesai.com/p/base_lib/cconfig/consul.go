package cconfig

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	_ "go.pkg.wesai.com/p/base_lib/log"
)

const (
	SMS_KV_CCONFIG_TOKEN string = "SMS_KV_TOKEN"
)

type ConsulNode struct {
	path      string
	isDir     bool
	value     string
	children  []Node
	recursive bool
	sorted    bool
}

func (en ConsulNode) Path() string {
	return en.path
}

func (en ConsulNode) IsDir() bool {
	return en.isDir
}

func (en ConsulNode) Value() string {
	return en.value
}

func (en ConsulNode) Children() []Node {
	return en.children
}

func (en ConsulNode) Recursive() bool {
	return en.recursive
}

func (en ConsulNode) Sorted() bool {
	return en.sorted
}

// etcd的节点变更。
type ConsulNodeChange struct {
	action string
	path   string
	isDir  bool
	value  string
	prev   *EtcdNodeChange
}

func (enc ConsulNodeChange) Action() string {
	return enc.action
}

func (enc ConsulNodeChange) Path() string {
	return enc.path
}

func (enc ConsulNodeChange) IsDir() bool {
	return enc.isDir
}

func (enc ConsulNodeChange) Value() string {
	return enc.value
}

func (enc ConsulNodeChange) Prev() NodeChange {
	return enc.prev
}

type ConsulClient struct { //implements  interface Client
	id         string
	consulAPI  *api.Client
	env        Env
	rootPath   string
	dataCenter string
	token      string
}

func NewConsulClient(id string, config ClientConfig) (cclient *ConsulClient, err error) {

	if err := config.Check(); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && cclient != nil {
			cclient.Close()
			cclient = nil
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

	consulConfig := &api.Config{
		Address:    config.Addresses[0],
		Scheme:     "http",
		Datacenter: config.Username,
		HttpAuth: &api.HttpBasicAuth{
			Username: "",
			Password: "",
		},
		Token: config.Password,
	}

	if consulConfig.Datacenter == "" {
		consulConfig.Datacenter = "dc1"
	}
	if consulConfig.Token == "" {
		consulConfig.Token = os.Getenv(SMS_KV_CCONFIG_TOKEN)
	}

	fmt.Println("Address = " + consulConfig.Address)
	fmt.Println("Scheme = " + consulConfig.Scheme)
	fmt.Println("Datacenter = " + consulConfig.Datacenter)
	fmt.Println("Token = " + consulConfig.Token)

	rawClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cclient = &ConsulClient{
		id:         id,
		consulAPI:  rawClient,
		env:        env,
		rootPath:   rootPath,
		dataCenter: config.Username,
		token:      config.Password,
	}

	var exist bool
	if exist, _, err = cclient.Exist(rootPath); err == nil && !exist {

		if err = cclient.CreateDir(rootPath); err != nil {
			return
		}
	} else {
		if err != nil {
			return
		}
	}

	return cclient, nil
}

func (cclient *ConsulClient) ID() string {
	return cclient.id
}

func (cclient *ConsulClient) RootPath() string {
	return cclient.rootPath
}

func (cclient *ConsulClient) Exist(nodePath string) (exist bool, dir bool, err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}
	var node Node
	node, err = cclient.Get(nodePath, false)
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

func (cclient *ConsulClient) CreateDir(nodePath string) (err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	if strings.HasPrefix(nodePath, "/") {
		nodePath = strings.TrimLeft(nodePath, "/")
	}
	if !strings.HasSuffix(nodePath, "/") { //目录必须带"/"
		nodePath += "/"
	}
	kvp := &api.KVPair{
		Key: nodePath,
		//Value: nil,
	}

	ops := &api.WriteOptions{
		Datacenter: cclient.dataCenter,
		Token:      cclient.token,
	}

	var consulKV *api.KV = cclient.consulAPI.KV()
	_, err = consulKV.Put(kvp, ops)
	return
}

func (cclient *ConsulClient) UpdateDir(nodePath string) (err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	kvp := &api.KVPair{
		Key:   nodePath,
		Value: nil,
	}

	ops := &api.WriteOptions{
		Datacenter: cclient.dataCenter,
		Token:      cclient.token,
	}

	var consulKV *api.KV = cclient.consulAPI.KV()
	_, err = consulKV.Put(kvp, ops)
	return

}

func (cclient *ConsulClient) SetDir(nodePath string) (incremental bool, err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	var consulKV *api.KV = cclient.consulAPI.KV()

	queryOpt := &api.QueryOptions{
		Datacenter:        cclient.dataCenter,
		AllowStale:        true,
		RequireConsistent: false, //to decide
		WaitIndex:         0,
		WaitTime:          time.Duration(0),
		Token:             cclient.token,
	}

	var kbP *api.KVPair
	///var queryMt *api.QueryMeta
	kbP, _, err = consulKV.Get(nodePath, queryOpt)

	if err != nil {
		return false, err
	}

	if kbP == nil {
		return false, cclient.CreateDir(nodePath)
	}
	if strings.HasSuffix(kbP.Key, "/") {
		return false, nil ///该目录已经存在
	} else {
		wopts := &api.WriteOptions{
			Datacenter: cclient.dataCenter,
			Token:      cclient.token,
		}
		_, err = consulKV.Delete(nodePath, wopts)
		if err != nil {
			return false, err
		}
		return false, cclient.CreateDir(nodePath)
	}
}

func (cclient *ConsulClient) Create(nodePath string, value string) (err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	var consulKV *api.KV = cclient.consulAPI.KV()
	kvP := &api.KVPair{
		Key:   nodePath,
		Value: []byte(value),
	}

	wopts := &api.WriteOptions{
		Datacenter: cclient.dataCenter,
		Token:      cclient.token,
	}
	_, err = consulKV.Put(kvP, wopts)
	return
}

func (cclient *ConsulClient) Update(nodePath string, value string) (err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	var consulKV *api.KV = cclient.consulAPI.KV()
	kvP := &api.KVPair{
		Key:   nodePath,
		Value: []byte(value),
	}

	wopts := &api.WriteOptions{
		Datacenter: cclient.dataCenter,
		Token:      cclient.token,
	}
	_, err = consulKV.Put(kvP, wopts)
	return
}

func (cclient *ConsulClient) Set(nodePath string, value string) (incremental bool, err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return false, err
	}

	var consulKV *api.KV = cclient.consulAPI.KV()

	queryOpt := &api.QueryOptions{
		Datacenter:        cclient.dataCenter,
		AllowStale:        true,
		RequireConsistent: false, //to decide
		WaitIndex:         0,
		WaitTime:          time.Duration(0),
		Token:             cclient.token,
	}

	var kbP *api.KVPair
	///var queryMt *api.QueryMeta
	kbP, _, err = consulKV.Get(nodePath, queryOpt)

	if kbP != nil {
		if strings.HasSuffix(kbP.Key, "/") {
			return false, fmt.Errorf("%s is a directory", nodePath)
		} else {
			err = cclient.Update(nodePath, value)
			if err != nil {
				return false, err
			} else {
				return true, nil
			}
		}
	}

	return false, cclient.Create(nodePath, value)
}

func (cclient *ConsulClient) Watch(
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

	path, err = cclient.ensurePath(path)
	if err != nil {
		return
	}

	var flagExit int32 = 0
	errChanInner := make(chan error)

	go func() {
		var lastIndex uint64 = 0

		var consulKV *api.KV = cclient.consulAPI.KV()
		queryOpt := &api.QueryOptions{
			Datacenter:        cclient.dataCenter,
			AllowStale:        true,
			RequireConsistent: false, //to decide
			WaitIndex:         lastIndex,
			WaitTime:          time.Duration(0),
			Token:             cclient.token,
		}

		var kvPairs api.KVPairs
		var queryMeta *api.QueryMeta
		kvPairs, queryMeta, err2 := consulKV.List(path, queryOpt)
		if err != nil {
			errChanInner <- err2
			return
		}

		if len(kvPairs) < 1 {
			errChanInner <- fmt.Errorf("Invalid key")
			return
		}

		lastIndex = queryMeta.LastIndex
		for atomic.LoadInt32(&flagExit) == 0 {
			queryOpt.WaitIndex = lastIndex
			queryOpt.WaitTime = time.Minute //todo time

			kvPairs, queryMeta, err2 = consulKV.List(path, queryOpt)

			if err2 != nil {
				errChanInner <- err2
				return
			}

			nodechange := ConsulNodeChange{}
			if queryMeta.LastIndex > lastIndex {
				nodechange.action = "update"
				nodechange.isDir = strings.HasSuffix(path, "/")
				nodechange.path = path
				nodechange.prev = nil
				nodechange.value = ""

				changeChan <- nodechange
				lastIndex = queryMeta.LastIndex
			}
		}
		errChanInner <- errors.New("External stop")
	}()

	for {
		select {
		case <-stopChan:
			atomic.StoreInt32(&flagExit, 1)
		case errInner := <-errChanInner:
			errorChan <- errInner
			break
		}
	}
}

func (cclient *ConsulClient) Delete(nodePath string, recursive bool) (done bool, err error) {
	nodePath, err = cclient.ensurePath(nodePath)
	if err != nil {
		return
	}

	var consulKV *api.KV = cclient.consulAPI.KV()

	wropts := &api.WriteOptions{
		Datacenter: cclient.dataCenter,
		Token:      cclient.token,
	}

	_, err = consulKV.DeleteTree(nodePath, wropts)

	if err != nil {
		done = false
	} else {
		done = true
	}
	return done, err
}

func (cclient *ConsulClient) Get(nodePath string, recursive bool) (Node, error) {

	nodePath, err := cclient.ensurePath(nodePath)
	nodePath = strings.TrimRight(nodePath, "/")
	if err != nil {
		return nil, err
	}

	var consulKV *api.KV = cclient.consulAPI.KV()

	queryOpt := &api.QueryOptions{
		Datacenter:        cclient.dataCenter,
		AllowStale:        true,
		RequireConsistent: false, //to decide
		WaitIndex:         0,
		WaitTime:          time.Duration(0),
		Token:             cclient.token,
	}

	var kvPairs api.KVPairs
	kvPairs, _, err = consulKV.List(nodePath, queryOpt)
	if err != nil {
		return nil, err
	}

	if len(kvPairs) < 1 {
		return nil, fmt.Errorf("%s : list return is empty", nodePath)
	}
	retNode, err := convToNodeFromKV(kvPairs, nodePath)
	if err != nil {
		return nil, err
	}
	return retNode, nil
}

func (cclient *ConsulClient) Close() {
	return
}

func convToNodeFromKV(kvps api.KVPairs, nodePath string) (ConsulNode, error) {

	for _, kvP := range kvps {
		if !strings.HasPrefix(kvP.Key, "/") {
			kvP.Key = "/" + kvP.Key
		}
	}

	consulNode := ConsulNode{
		path:      nodePath,
		isDir:     true,
		value:     "",
		children:  nil,
		recursive: true,
		sorted:    true,
	}
	for _, kvP := range kvps {
		if !strings.HasPrefix(kvP.Key, consulNode.path) {
			logger.Errorf("List Return's data wrong\n")
			continue
		}
		_ = recursiveMerge(kvP, &consulNode)
	}

	return consulNode, nil
}

func recursiveMerge(kvp *api.KVPair, cNode *ConsulNode) error {
	if strings.HasPrefix(kvp.Key, cNode.path) {
		var keys string
		if strings.HasSuffix(cNode.path, "/") {
			keys = kvp.Key[(len(cNode.path) - 1):]
		} else {
			keys = kvp.Key[len(cNode.path):]
		}

		if len(keys) == 0 {
			return nil
		}

		count := strings.Count(keys, "/") - 1

		if count == 0 {
			tpConsulNode := ConsulNode{
				path:      kvp.Key,
				isDir:     false,
				value:     string(kvp.Value),
				children:  nil,
				recursive: true,
				sorted:    true,
			}
			addChildren(&tpConsulNode, cNode)
			return nil
		} else if count == 1 && strings.HasSuffix(keys, "/") {
			tpConsulNode := ConsulNode{
				path:      kvp.Key,
				isDir:     true,
				value:     "",
				children:  nil,
				recursive: true,
				sorted:    true,
			}
			addChildren(&tpConsulNode, cNode)
			return nil
		} else {
			keys = strings.TrimLeft(keys, "/")
			elems := strings.Split(keys, "/")
			var tppath string
			if strings.HasSuffix(cNode.path, "/") {
				tppath = cNode.path + elems[0]
			} else {
				tppath = cNode.path + "/" + elems[0]
			}

			tpConsulNode := getChildren(tppath, cNode)
			if tpConsulNode == nil {
				tpConsulNode = &ConsulNode{
					path:      tppath,
					isDir:     true,
					value:     "",
					children:  nil,
					recursive: true,
					sorted:    true,
				}
			}
			errRet := recursiveMerge(kvp, tpConsulNode)
			addChildren(tpConsulNode, cNode)
			return errRet
		}
	}
	return nil
}

func getChildren(childKey string, cNode *ConsulNode) *ConsulNode {
	if cNode.children == nil {
		return nil
	}
	for index, child := range cNode.children {
		if childKey == child.Path() {
			tpNode := cNode.children[index]
			if index < len(cNode.children)-1 {
				cNode.children = append(cNode.children[0:index], cNode.children[(index+1):]...)
			} else {
				cNode.children = cNode.children[0:index]
			}

			if rNode, ok := tpNode.(ConsulNode); ok {
				return &rNode
			}
		}
	}
	return nil
}

func addChildren(cChil *ConsulNode, cNode *ConsulNode) {
	if cNode.children == nil {
		cslNode := []Node{}
		cNode.children = append(cslNode, *cChil)
	} else {
		found := false
		for _, nodechildren := range cNode.children {

			if nodechildren.Path() == cChil.path {
				found = true
				break
			}
		}

		if found == false {
			cNode.children = append(cNode.children, *cChil)
		}
	}
}

func (eclient *ConsulClient) ensurePath(nodePath string) (string, error) {
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
