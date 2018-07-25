package cconfig

import (
	"errors"
	"sync"
	"time"
)

var (
	errRepeatedActivation = errors.New("Can not activate the client has been activated!")
	errInactivatedClient  = errors.New("Can not use the inactivated client!")
	errClosedClient       = errors.New("Can not use the closed client!")
)

// 客户端包装器。
type clientWrapper struct {
	pool      *clientPool
	client    Client
	idleTs    int64
	activated bool
	closed    bool
	rwm       sync.RWMutex
}

func (wrapper *clientWrapper) overtime(idleTimeout time.Duration) bool {
	wrapper.rwm.RLock()
	defer wrapper.rwm.RUnlock()
	var overtime bool
	if !wrapper.activated &&
		time.Now().UnixNano()-wrapper.idleTs >= idleTimeout.Nanoseconds() {
		overtime = true
	}
	return overtime
}

func (wrapper *clientWrapper) check() error {
	wrapper.rwm.RLock()
	defer wrapper.rwm.RUnlock()
	if wrapper.pool.destroyed {
		return ErrClosedPool
	}
	// if wrapper.closed {
	// 	return errClosedClient
	// }
	return nil
}

func (wrapper *clientWrapper) activate() {
	wrapper.rwm.Lock()
	defer wrapper.rwm.Unlock()
	if wrapper.closed {
		panic(errClosedClient)
	}
	if wrapper.activated {
		panic(errRepeatedActivation)
	}
	wrapper.activated = true
	wrapper.idleTs = 0
}

func (wrapper *clientWrapper) inactivate() {
	wrapper.rwm.Lock()
	defer wrapper.rwm.Unlock()
	if wrapper.closed || !wrapper.activated {
		return
	}
	wrapper.activated = false
	wrapper.idleTs = time.Now().UnixNano()
}

func (wrapper *clientWrapper) close() {
	wrapper.rwm.Lock()
	defer wrapper.rwm.Unlock()
	if !wrapper.closed {
		wrapper.closed = true
		wrapper.client.Close()
	}
}

// ** Implement methods of interface 'Client' - Begin ** //

func (wrapper *clientWrapper) ID() string {
	return wrapper.client.ID()
}

func (wrapper *clientWrapper) RootPath() string {
	return wrapper.client.RootPath()
}

func (wrapper *clientWrapper) Exist(path string) (exist bool, dir bool, err error) {
	if err = wrapper.check(); err != nil {
		return
	}
	return wrapper.client.Exist(path)
}

func (wrapper *clientWrapper) Get(path string, recursive bool) (Node, error) {
	if err := wrapper.check(); err != nil {
		return nil, err
	}
	return wrapper.client.Get(path, recursive)
}

func (wrapper *clientWrapper) Delete(nodePath string, recursive bool) (done bool, err error) {
	if err = wrapper.check(); err != nil {
		return
	}
	return wrapper.client.Delete(nodePath, recursive)
}

func (wrapper *clientWrapper) Watch(
	path string,
	recursive bool,
	receiver chan NodeChange,
	errorChan chan error,
	stop chan bool) {

	var err error
	defer func() {
		if err != nil {
			errorChan <- err
			close(errorChan)
			close(receiver)
		}
	}()
	if err = wrapper.check(); err != nil {
		return
	}
	if stop == nil {
		err = errors.New("Invalid stop channel!")
		return
	}
	wrapper.client.Watch(path, recursive, receiver, errorChan, stop)
}

func (wrapper *clientWrapper) CreateDir(path string) error {
	if err := wrapper.check(); err != nil {
		return err
	}
	return wrapper.client.CreateDir(path)
}

func (wrapper *clientWrapper) UpdateDir(path string) error {
	if err := wrapper.check(); err != nil {
		return err
	}
	return wrapper.client.UpdateDir(path)
}

func (wrapper *clientWrapper) SetDir(path string) (bool, error) {
	if err := wrapper.check(); err != nil {
		return false, err
	}
	return wrapper.client.SetDir(path)
}

func (wrapper *clientWrapper) Create(path string, value string) error {
	if err := wrapper.check(); err != nil {
		return err
	}
	return wrapper.client.Create(path, value)
}

func (wrapper *clientWrapper) Update(path string, value string) error {
	if err := wrapper.check(); err != nil {
		return err
	}
	return wrapper.client.Update(path, value)
}

func (wrapper *clientWrapper) Set(path string, value string) (bool, error) {
	if err := wrapper.check(); err != nil {
		return false, err
	}
	return wrapper.client.Set(path, value)
}

func (wrapper *clientWrapper) Close() {
	if wrapper.pool.destroyed {
		wrapper.close()
	} else {
		wrapper.pool.put(wrapper, true)
	}
}

// ** Implement methods of interface 'Client' - End ** //
