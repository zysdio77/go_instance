package cconfig

import (
	"errors"
	"fmt"
)

// 根路径基础。
const root_base = "/wesai_root/"

// 获取根路径。
func GetRootPath(env Env, projectName string) (string, error) {
	switch env {
	case ENV_DEVELOPMENT, ENV_TEST, ENV_PRODUCTION:
		return root_base + string(env) + "/" + projectName, nil
	default:
		errMsg := fmt.Sprintf("Unsupported env '%s'!\n", env)
		return "", errors.New(errMsg)
	}
}
