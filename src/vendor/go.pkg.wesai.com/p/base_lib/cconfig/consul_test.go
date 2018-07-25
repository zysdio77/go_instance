package cconfig

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewConsulClient(t *testing.T) {
	envs := []Env{ENV_DEVELOPMENT, ENV_TEST, ENV_PRODUCTION}
	for _, env := range envs {
		clientConfig := ClientConfig{
			ProjectName: "msg_gateway",
			Environment: env,
			Addresses:   GetCCAddressFromSys(),
			Username:    "dc1",
			Password:    "9f0fead1-9cbe-4512-b0d8-caf64857f2b2",
			NotGuest:    true,
		}
		client, err := NewConsulClient("", clientConfig)
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

func Test_watch(t *testing.T) {
	clientConfig := ClientConfig{
		ProjectName: "msg_gateway",
		Environment: "devel",
		Addresses:   GetCCAddressFromSys(),
		Username:    "dc1",
		Password:    "9f0fead1-9cbe-4512-b0d8-caf64857f2b2",
		NotGuest:    true,
	}

	client, err := NewConsulClient("", clientConfig)
	if err != nil {
		t.Fatalf("Can not new an consul client for environment '%s': %s\n",
			"devel", err)
	}
	defer client.Close()

	nodec := make(chan NodeChange, 10)
	errC := make(chan error, 10)
	stopC := make(chan bool, 10)

	client.Watch("/wesai_root/devel/msg_gateway", true, nodec, errC, stopC)

	for {
		select {
		case ndCh := <-nodec:
			t.Logf("Action : %s  Path: %s\n", ndCh.Action(), ndCh.Path())

		}
	}

}

func Test_Get(t *testing.T) {
	clientConfig := ClientConfig{
		ProjectName: "msg_gateway",
		Environment: "devel",
		Addresses:   GetCCAddressFromSys(),
		Username:    "dc1",
		Password:    "9f0fead1-9cbe-4512-b0d8-caf64857f2b2",
		NotGuest:    true,
	}

	client, err := NewConsulClient("", clientConfig)
	if err != nil {
		t.Fatalf("Can not new an consul client for environment '%s': %s\n",
			"devel", err)
	}
	defer client.Close()

	node, err := client.Get("/wesai_root/devel/msg_gateway", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(node.Path())
	t.Log(node.Value())

}
