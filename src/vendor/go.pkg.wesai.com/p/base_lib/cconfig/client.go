package cconfig

import (
	"strings"
)

// 代表路径节点的接口。
type Node interface {
	Path() string
	IsDir() bool
	Value() string
	Children() []Node
	Recursive() bool
	Sorted() bool
}

// 代表路径节点变更的接口。
type NodeChange interface {
	Action() string
	Path() string
	IsDir() bool
	Value() string
	Prev() NodeChange
}

// 客户端接口。
type Client interface {
	ID() string
	RootPath() string
	// Node
	Exist(path string) (exist bool, dir bool, err error)
	Get(path string, recursive bool) (Node, error)
	Delete(nodePath string, recursive bool) (done bool, err error)
	// 异步的监视节点变化，并返回变化的节点。
	// errorChan通道中若有元素可接收则说明有错误发生，且监视已停止。
	Watch(path string,
		recursive bool,
		receiver chan NodeChange,
		errorChan chan error,
		stop chan bool)
	// Dir
	CreateDir(path string) error
	UpdateDir(path string) error
	// 第一个结果值代表是否新增了路径。
	SetDir(path string) (bool, error)
	// Key
	Create(path string, value string) error
	Update(path string, value string) error
	// 第一个结果值代表是否新增了路径。
	Set(path string, value string) (bool, error)
	// Connection
	Close()
}

// 客户端配置。
type ClientConfig struct {
	ProjectName string   // 当前项目的名称。该名称会作为根路径的最后一个元素。
	Environment Env      // 程序运行环境。该环境会作为根路径的中间元素。
	Addresses   []string // 配置中心的地址列表。
	Username    string   // 用户名。
	Password    string   // 密码。
	NotGuest    bool     // 是否不作为客体访问，即需要使用用户名和密码来登录。
}

// 检查自身。
func (cconfig *ClientConfig) Check() error {
	cconfig.ProjectName = strings.TrimSpace(cconfig.ProjectName)
	if cconfig.ProjectName == "" {
		return ErrEmptyProjectName
	}
	if !CheckEnvAvailability(cconfig.Environment) {
		return ErrEmptyEnv
	}
	if len(cconfig.Addresses) == 0 {
		return ErrEmptyAddresses
	}
	if cconfig.NotGuest {
		cconfig.Username = strings.TrimSpace(cconfig.Username)
		if cconfig.Username == "" {
			return ErrEmptyUsername
		}
		cconfig.Password = strings.TrimSpace(cconfig.Password)
		if cconfig.Password == "" {
			return ErrEmptyPassword
		}
	}
	return nil
}
