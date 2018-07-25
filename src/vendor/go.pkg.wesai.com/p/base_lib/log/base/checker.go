package base // import "go.pkg.wesai.com/p/base_lib/log/base"

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// debug日志开关的相关常量。
const (
	DEBUG_ENABLE_ENV_VAR_NAME_TEMPLATE = "%s_LOG_DEBUG_ENABLE" // 与debug日志开关相关的环境变量名称的模板。
	DEBUG_ENABLE_DEFAULT_ENV_VAR_NAME  = "GO_LOG_DEBUG_ENABLE" // 默认的与debug日志开关相关的环境变量名称。
)

// debug日志开关的字典。键为与项目名称对应的环境变量名称，值为对应的debug日志开关的检查器。
var debugEnableMap = map[string]*debugEnableChecker{}

// 操作上述字典的读写互斥量。
var rwm sync.RWMutex

// debug日志开关的检查器。
type debugEnableChecker struct {
	envVarName string
	value      bool
	ticker     *time.Ticker
	rwm        sync.RWMutex
	running    bool
}

// 启动周期性检查。
// 若参数interval的值为非正数，则会立即返回false。
// 参数override代表是否按照新的周期设定重启对debug日志开关的检查。
// 结果值代表本次启动是否成功。
func (checker *debugEnableChecker) start(
	interval time.Duration,
	override bool) bool {

	if interval <= 0 {
		return false
	}
	checker.rwm.Lock()
	defer checker.rwm.Unlock()
	if checker.running && !override {
		return false
	}
	if checker.running {
		checker.ticker.Stop()
		checker.ticker = nil
	}
	if checker.ticker == nil {
		checker.ticker = time.NewTicker(interval)
	}
	updateDebugEnable(checker)
	go cyclicUpdateDebugEnable(checker)
	checker.running = true
	return true
}

// 周期性的更新debug日志开关。
func cyclicUpdateDebugEnable(checker *debugEnableChecker) {
	for _ = range checker.ticker.C {
		envVarValue := os.Getenv(checker.envVarName)
		envVarValue = strings.ToLower(strings.TrimSpace(envVarValue))
		checker.rwm.Lock()
		updateDebugEnable(checker)
		checker.rwm.Unlock()
	}
}

// 更新debug日志开关。非并发安全方法！
func updateDebugEnable(checker *debugEnableChecker) {
	envVarValue := os.Getenv(checker.envVarName)
	envVarValue = strings.ToLower(strings.TrimSpace(envVarValue))
	if envVarValue == "t" ||
		envVarValue == "true" ||
		envVarValue == "1" {
		checker.value = true
	} else {
		checker.value = false
	}
}

// 停止检查。
// 结果值代表停止的成功与否。若本检查器尚未启动则会立即返回false。
func (checker *debugEnableChecker) stop() bool {
	checker.rwm.Lock()
	defer checker.rwm.Unlock()
	if !checker.running {
		return false
	}
	checker.ticker.Stop()
	checker.ticker = nil
	checker.running = false
	return true
}

// 根据系统环境变量启动/停止对指定项目的debug日志开关的周期性检查。
// 参数override代表是否按照新的周期设定检查指定项目的debug日志开关。非正数意味着会停止检查。
// 结果值代表本次启动/停止是否成功。
func CheckDebugEnable(projectName string, interval time.Duration, override bool) bool {
	envVarName := GenEnvVarName(projectName)
	rwm.Lock()
	defer rwm.Unlock()
	if _, ok := debugEnableMap[envVarName]; ok && !override {
		return false
	}
	if interval <= 0 {
		checker := debugEnableMap[envVarName]
		if checker != nil {
			checker.stop()
		}
		delete(debugEnableMap, envVarName)
		return true
	}
	var checker *debugEnableChecker
	if debugEnableMap[envVarName] == nil {
		checker = &debugEnableChecker{
			envVarName: envVarName,
		}
		debugEnableMap[envVarName] = checker
	}
	checker.start(interval, true)
	return true
}

// 获取debug日志开关。
func DebugEnable(projectName string) bool {
	envVarName := GenEnvVarName(projectName)
	rwm.RLock()
	defer rwm.RUnlock()
	if checker := debugEnableMap[envVarName]; checker != nil {
		return checker.value
	}
	return false
}

// 根据项目名称生成与debug日志开关相关的环境变量名称。
func GenEnvVarName(projectName string) string {
	projectName = strings.TrimSpace(projectName)
	var envVarName string
	if projectName == "" {
		envVarName = DEBUG_ENABLE_DEFAULT_ENV_VAR_NAME
	} else {
		envVarName = fmt.Sprintf(DEBUG_ENABLE_ENV_VAR_NAME_TEMPLATE,
			strings.ToUpper(projectName))
	}
	return envVarName
}
