package cconfig // import "go.pkg.wesai.com/p/base_lib/cconfig"

import (
	"net/url"
	"os"
	"strings"
)

// 微影内部的etcd集群的地址。
const (
	etcd_addresses_default = "http://10.2.2.9:2379,http://10.2.2.8:2379,http://10.2.2.7:2379"
)

// 配置中心类型。
type ConfigCenterType string

// 预设的客户端类型。
const (
	CC_TYPE_ETCD   = "etcd"
	CC_TYPE_CONSUL = "consul"
	CC_TYPE_ZK     = "zookeeper"
)

// 可用的配置中心类型的字典。
var availableConfigCenterTypeMap = map[ConfigCenterType]bool{
	CC_TYPE_ETCD:   true,
	CC_TYPE_CONSUL: true,
}

// 检查配置中心类型的可用性。
func CheckTypeAvailability(ccType ConfigCenterType) bool {
	return availableConfigCenterTypeMap[ccType]
}

// 程序运行环境。
type Env string

// 预定义的程序运行环境常量。
const (
	ENV_DEVELOPMENT Env = "devel" // 开发环境。
	ENV_TEST        Env = "test"  // 测试环境。
	ENV_PRODUCTION  Env = "prod"  // 生产环境。
)

// 相关环境变量的名称。
const (
	ENV_VAR_NAME_CC_TYPE  = "CC_TYPE"  // 配置中心的类型。
	ENV_VAR_NAME_CC_ADDRS = "CC_ADDRS" // 配置中心的地址列表，多个地址之间应该以“,”分隔。
	ENV_VAR_NAME_CC_ENV   = "CC_ENV"   // 程序运行环境。
)

// 从系统环境变量获取配置中心的类型。默认为：CC_TYPE_ETCD。
func GetCCenterTypeFromSys() ConfigCenterType {
	varValue := os.Getenv(ENV_VAR_NAME_CC_TYPE)
	ccenterType := ConfigCenterType(strings.ToLower(varValue))
	switch ccenterType {
	case CC_TYPE_ETCD, CC_TYPE_CONSUL, CC_TYPE_ZK:
	default:
		ccenterType = CC_TYPE_ETCD
	}
	return ccenterType
}

// 从系统环境变量获取配置中心的地址列表。
func GetCCAddressFromSys() []string {
	envVarValue := os.Getenv(ENV_VAR_NAME_CC_ADDRS)
	parts := strings.Split(envVarValue, ",")
	addrs := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		ccenterUrl, err := url.Parse(part)
		if err != nil {
			continue
		}
		addrs = append(addrs, ccenterUrl.String())
	}
	if len(addrs) == 0 {
		addrs = strings.Split(etcd_addresses_default, ",")
	}
	return addrs
}

// 从系统环境变量获取程序运行环境的值。默认：ENV_DEVELOPMENT。
func GetEnvFromSys() Env {
	varValue := os.Getenv(ENV_VAR_NAME_CC_ENV)
	env := Env(strings.ToLower(varValue))
	switch env {
	case ENV_DEVELOPMENT, ENV_TEST, ENV_PRODUCTION:
	default:
		env = ENV_DEVELOPMENT
	}
	return env
}

// 检查程序运行环境标记的可用性。
func CheckEnvAvailability(env Env) bool {
	if env != ENV_DEVELOPMENT &&
		env != ENV_TEST &&
		env != ENV_PRODUCTION {
		return false
	}
	return true
}
