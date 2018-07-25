package cconfig

import (
	"errors"
	"fmt"
)

// 空项目名错误。
var ErrEmptyProjectName = errors.New("Empty project name!")

// 空程序运行环境错误。
var ErrEmptyEnv = errors.New("Empty environment!")

// 空用户名错误。
var ErrEmptyUsername = errors.New("Empty user name!")

// 空密码错误。
var ErrEmptyPassword = errors.New("Empty password!")

// 空地址错误。
var ErrEmptyAddresses = errors.New("Empty addresses!")

// 空路径错误。
var ErrEmptyPath = errors.New("Empty path!")

// 代表不合法的路径的错误类型。
type ErrIllegalPath struct {
	path     string
	rootPath string
}

func (err *ErrIllegalPath) Error() string {
	return fmt.Sprintf("Illegal path '%s'! (rootPath: %s)",
		err.path, err.rootPath)
}

// 代表路径不存在的错误类型。
type ErrInexistentPath struct {
	path string
}

func (err *ErrInexistentPath) Error() string {
	return fmt.Sprintf("The path '%s' does not exist! (100)", err.path)
}

// 代表路径已存在的错误类型。
type ErrExistingPath struct {
	path string
}

func (err *ErrExistingPath) Error() string {
	return fmt.Sprintf("The path '%s' already exists! (105)", err.path)
}

// 代表目录已存在的错误类型。
type ErrExistingDir struct {
	dirPath string
}

func (err *ErrExistingDir) Error() string {
	return fmt.Sprintf("The dir '%s' already exists! (102)", err.dirPath)
}

// 代表耗尽的池的错误类型。
var ErrPoolExhausted = errors.New("Connection pool exhausted!")

// 代表池已关闭的错误类型。
var ErrClosedPool = errors.New("Closed pool!")

// 代表节点变化监视已被用户停止的错误类型。
var ErrWatchStoppedByUser = errors.New("The node change watch has been stopped by user!")
