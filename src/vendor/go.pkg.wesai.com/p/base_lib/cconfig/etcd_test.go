package cconfig

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	PROJECT_NAME_FOR_CLIENT_TEST = "cc_client_test"
	USERNAME_FOR_TEST            = "wstester"
	PASSWORD_FOR_TEST            = "567890"
)

func TestNewEtcdClient(t *testing.T) {
	envs := []Env{ENV_DEVELOPMENT, ENV_TEST, ENV_PRODUCTION}
	for _, env := range envs {
		clientConfig := ClientConfig{
			ProjectName: PROJECT_NAME_FOR_CLIENT_TEST,
			Environment: env,
			Addresses:   GetCCAddressFromSys(),
			Username:    USERNAME_FOR_TEST,
			Password:    PASSWORD_FOR_TEST,
			NotGuest:    true,
		}
		client, err := NewEtcdClient("", clientConfig)
		if err != nil {
			t.Fatalf("Can not new an etcd client for environment '%s': %s\n",
				env, err)
		}
		defer client.Close()

		// RootPath
		expectRootPath := root_base + string(env) + "/" + PROJECT_NAME_FOR_CLIENT_TEST
		actualRootPath := client.RootPath()
		if actualRootPath != expectRootPath {
			t.Fatalf("Expect '%s', but the actual is '%s'\n",
				expectRootPath, actualRootPath)
		}
		t.Logf("The etcd client (ID: %s) is ready. It's root path is '%s'.",
			client.ID(), client.RootPath())
	}
}

func getEtcdClient() (*EtcdClient, error) {
	clientConfig := ClientConfig{
		ProjectName: PROJECT_NAME_FOR_CLIENT_TEST,
		Environment: ENV_DEVELOPMENT,
		Addresses:   GetCCAddressFromSys(),
		Username:    USERNAME_FOR_TEST,
		Password:    PASSWORD_FOR_TEST,
		NotGuest:    true,
	}
	return NewEtcdClient("", clientConfig)
}

func clean(client *EtcdClient) error {
	path := "/"
	done, err := client.Delete(path, true)
	if !done {
		err = errors.New(fmt.Sprintf("Fail to delete '%s'!", path))
	}
	return err
}

func TestEtcdGet(t *testing.T) {
	client, err := getEtcdClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func() {
		err := clean(client)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Get & SetDir & Set
	nodePath0 := "/0"
	_, err = client.Get(nodePath0, false)
	if err != nil {
		if _, ok := err.(*ErrInexistentPath); !ok {
			t.Fatalf("Occur an error (path: %s): %s", nodePath0, err)
		}
	} else {
		t.Fatal("An inexistent path error should occur!")
	}
	keyPath0Xs := []string{
		nodePath0 + "/k1",
		nodePath0 + "/k2",
		nodePath0 + "/k3",
	}
	for _, keyPath0X := range keyPath0Xs {
		_, err := client.Set(keyPath0X, keyPath0X)
		if err != nil {
			t.Fatalf("Occur an error (path: %s): %s", keyPath0X, err)
		}
	}
	keyPath01 := keyPath0Xs[0]
	_, err = client.SetDir(keyPath01)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath01, err)
	}
	keyPath01Xs := []string{
		keyPath01 + "/k10",
		keyPath01 + "/k20",
		keyPath01 + "/k30",
	}
	for _, keyPath01X := range keyPath01Xs {
		_, err := client.Set(keyPath01X, keyPath01X)
		if err != nil {
			t.Fatalf("Occur an error (path: %s): %s", keyPath01X, err)
		}
	}

	// Get
	pnode, err := client.Get(nodePath0, true)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath0, err)
	}
	if pnode == nil ||
		!pnode.IsDir() ||
		!strings.HasSuffix(pnode.Path(), nodePath0) {
		t.Fatalf("Incorrect path node: %#v", pnode)
	}
	children := pnode.Children()
	childrenLen := len(children)
	expectedChildrenLen := len(keyPath0Xs)
	if childrenLen != expectedChildrenLen {
		t.Fatalf("Expect %d, but the actual is %d\n",
			expectedChildrenLen, childrenLen)
	}
	for i := 1; i < childrenLen; i++ {
		child := children[i]
		if !strings.HasSuffix(child.Path(), keyPath0Xs[i]) ||
			child.IsDir() ||
			child.Value() != keyPath0Xs[i] ||
			child.Children() != nil ||
			!child.Recursive() ||
			!child.Sorted() {
			t.Fatalf("Incorrect path node: %#v", child)
		}
	}
	child1 := children[0]
	if child1 == nil ||
		!child1.IsDir() ||
		!strings.HasSuffix(child1.Path(), keyPath0Xs[0]) {
		t.Fatalf("Incorrect path node: %#v", child1)
	}
	grandson := child1.Children()
	grandsonLen := len(grandson)
	expectedGrandsonLen := len(keyPath01Xs)
	if grandsonLen != expectedGrandsonLen {
		t.Fatalf("Expect %d, but the actual is %d\n",
			expectedGrandsonLen, grandsonLen)
	}
	for i := 1; i < grandsonLen; i++ {
		grandson := grandson[i]
		if !strings.HasSuffix(grandson.Path(), keyPath01Xs[i]) ||
			grandson.IsDir() ||
			grandson.Value() != keyPath01Xs[i] ||
			grandson.Children() != nil ||
			!grandson.Recursive() ||
			!grandson.Sorted() {
			t.Fatalf("Incorrect path node: %#v", grandson)
		}
	}
}

func TestEtcdOpDir(t *testing.T) {
	client, err := getEtcdClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func() {
		err := clean(client)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Exist & CreatDir & Get
	nodePath1 := "/1"
	exist, _, err := client.Exist(nodePath1)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
	}
	if exist {
		t.Fatalf("The path '%s' already exists!", nodePath1)
	}
	err = client.CreateDir(nodePath1)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
	}
	exist, dir, err := client.Exist(nodePath1)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
	}
	if !exist {
		t.Fatalf("The path '%s' does not exist!", nodePath1)
	}
	if !dir {
		t.Fatalf("The path '%s' is not a dir! But it should be.", nodePath1)
	}
	err = client.CreateDir(nodePath1)
	if err != nil {
		if _, ok := err.(*ErrExistingPath); !ok {
			t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
		}
	} else {
		t.Fatal("An existing path error should occur!")
	}
	pnode, err := client.Get(nodePath1, false)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
	}
	if pnode == nil ||
		!pnode.IsDir() ||
		!strings.HasSuffix(pnode.Path(), nodePath1) {
		t.Fatalf("Incorrect path node: %#v", pnode)
	}

	// UpdateDir & SetDir
	err = client.UpdateDir(nodePath1)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePath1, err)
	}
	keyPath1 := "/1/k1"
	incremental, err := client.SetDir(keyPath1)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath1, err)
	}
	expectedIncremental := true
	if incremental != expectedIncremental {
		t.Fatalf("Expect %v, but the actual is %v\n",
			expectedIncremental, incremental)
	}
	_, err = client.SetDir(keyPath1)
	if err != nil {
		if _, ok := err.(*ErrExistingDir); !ok {
			t.Fatalf("Occur an error (path: %s): %s", keyPath1, err)
		}
	} else {
		t.Fatal("An existing dir error should occur!")
	}
}

func TestEtcdOpKeyValue(t *testing.T) {
	client, err := getEtcdClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func() {
		err := clean(client)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Exist & Create & Update & Get & SetDir
	keyPath2 := "/1/k2"
	value2 := "v2"
	exist, _, err := client.Exist(keyPath2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	if exist {
		t.Fatalf("The key path '%s' already exists!", keyPath2)
	}
	err = client.Create(keyPath2, value2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	exist, dir, err := client.Exist(keyPath2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	if !exist {
		t.Fatalf("The key path '%s' does not exist!", keyPath2)
	}
	if dir {
		t.Fatalf("The key path '%s' is a dir! But it should not be.", keyPath2)
	}
	err = client.Create(keyPath2, value2)
	if err != nil {
		if _, ok := err.(*ErrExistingPath); !ok {
			t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
		}
	} else {
		t.Fatal("An existing path error should occur!")
	}
	err = client.Update(keyPath2, value2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	err = client.Update(keyPath2+"x", value2)
	if err != nil {
		if _, ok := err.(*ErrInexistentPath); !ok {
			t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
		}
	} else {
		t.Fatal("An inexistent path error should occur!")
	}
	pnode, err := client.Get(keyPath2, false)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	if pnode == nil ||
		pnode.IsDir() ||
		!strings.HasSuffix(pnode.Path(), keyPath2) {
		t.Fatalf("Incorrect path node: %#v", pnode)
	}
	incremental, err := client.SetDir(keyPath2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
	}
	expectedIncremental := false
	if incremental != expectedIncremental {
		t.Fatalf("Expect %v, but the actual is %v\n",
			expectedIncremental, incremental)
	}

	// Set
	_, err = client.Set(keyPath2, value2)
	if err != nil {
		if _, ok := err.(*ErrExistingDir); !ok {
			t.Fatalf("Occur an error (path: %s): %s", keyPath2, err)
		}
	} else {
		t.Fatal("An existing dir error should occur!")
	}
	keyPath3 := "/1/k3"
	value3 := "v3"
	incremental, err = client.Set(keyPath3, value3)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath3, err)
	}
	expectedIncremental = true
	if incremental != expectedIncremental {
		t.Fatalf("Expect %v, but the actual is %v\n",
			expectedIncremental, incremental)
	}
	incremental, err = client.Set(keyPath3, value3)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", keyPath3, err)
	}
	expectedIncremental = false
	if incremental != expectedIncremental {
		t.Fatalf("Expect %v, but the actual is %v\n",
			expectedIncremental, incremental)
	}
}

func TestEtcdClientWatch(t *testing.T) {
	client, err := getEtcdClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func() {
		err := clean(client)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	nodePathW2 := "/w/1"
	keyPathWs := []string{
		nodePathW2 + "/k1",
		nodePathW2 + "/k2",
		nodePathW2 + "/k3",
	}
	receiver := make(chan NodeChange, 1)
	errorCh := make(chan error, 1)
	stopCh := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for err := range errorCh {
			if err != nil && err != ErrWatchStoppedByUser {
				t.Fatalf("Occur an error (path: %s): %s", nodePathW2, err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		nodeChange := <-receiver
		if nodeChange == nil ||
			nodeChange.Action() != "set" ||
			!strings.HasSuffix(nodeChange.Path(), nodePathW2) ||
			!nodeChange.IsDir() {

			t.Fatalf("Incorrect node change (key): %#v", nodeChange)
		}
		t.Logf("Received a node chang (key): %#v\n", nodeChange)
		var index int
		for nodeChange = range receiver {
			suffix := keyPathWs[index]
			if nodeChange == nil ||
				nodeChange.Action() != "set" ||
				!strings.HasSuffix(nodeChange.Path(), suffix) ||
				nodeChange.IsDir() {

				t.Fatalf("Incorrect node chang (key): %#v", nodeChange)
			}
			t.Logf("Received a node chang (key): %#v\n", nodeChange)
			index++
		}
	}()
	client.Watch(nodePathW2, true, receiver, errorCh, stopCh)
	time.Sleep(time.Millisecond * 10)
	_, err = client.SetDir(nodePathW2)
	if err != nil {
		t.Fatalf("Occur an error (path: %s): %s", nodePathW2, err)
	}
	for _, keyPathW := range keyPathWs {
		_, err = client.Set(keyPathW, keyPathW)
		if err != nil {
			t.Fatalf("Occur an error (path: %s): %s", keyPathW, err)
		}
	}
	time.Sleep(time.Millisecond * 10)
	stopCh <- true
	wg.Wait()
}
