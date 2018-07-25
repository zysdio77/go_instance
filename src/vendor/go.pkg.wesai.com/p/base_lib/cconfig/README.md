cconfig
=====================
本代码包的全称是：Centralized Configuration（client），即：集中配置客户端的Go语言版。

### 0. 版本：1.1

### 1. 用意与用途
本代码包意在方便Go程序访问配置中心。本代码包目前包含了针对Etcd的客户端以及客户端池程序。

### 2. 使用方法
在使用本代码包之前先要检查当前操作系统中是否存在下面这几个环境变量：

+ CC\_ENV：配置中心环境。可选值有：`devel`、`test`和`prod`，分别代表开发环境、测试环境和生产环境。默认值是`devel`。客户端会根据这个环境变量的值来确定其访问根目录。例如：若值为`test`，则客户端的访问根目录就是`/wepiao_root/test`。 
+ CC\_TYPE：配置中心类型。因为目前本代码包中只有对Etcd的支持，所以该环境变量的可选值和默认值都是`etcd`。所以可以不设置该环境变量。
+ CC\_ADDRS：配置中心地址。该环境变量的默认值为微票儿内部的Etcd配置中心集群的地址。强烈不建议设置该环境变量，除非出于测试或调试的目的。其值的一个示例是`http://127.0.0.1:2379,http://127.0.0.2:2379`。

客户端及客户端池的具体使用方法请见下面的示例。

### 3. 使用示例
#### 3.1 Etcd客户端的使用方法
初始化方法：

```go
package main

import (
	"fmt"
	"time"

	"go.pkg.wesai.com/p/base_lib/cconfig"
)

func main() {
	id := "" // 该客户端的ID。
	projectName := "testing"
	env := cconfig.GetEnvFromSys()           // 从环境变量获取配置中心环境。
	addr := cconfig.GetCCAddressFromSys() // 从环境变量获取配置中心地址。
	client, err := cconfig.NewEtcdClient("", projectName, env, addr...)
	if err != nil {
		fmt.Printf("Can not new an etcd client for environment '%s': %s\n",
			env, err)
		return 
	}
	defer client.Close()
	// ......
}
```

Etcd客户端的各个方法请参见本代码包的文档。

**注意：当环境变量 CONFIG\_CENTER\_ADDRS 的值为默认值时，运行上述程序的机器必须已连通腾讯云才能使该程序正常运行，否则会提示I/O超时错误。**

#### 3.2 客户端池的使用方法
初始化方法：

```go
package main

import (
	"fmt"
	"time"

	"go.pkg.wesai.com/p/base_lib/cconfig"
)

func main() {
	projectName := "testing"
	env := cconfig.GetEnvFromSys()           // 从环境变量获取配置中心环境。
	addr := cconfig.GetCCAddressFromSys() // 从环境变量获取配置中心地址。
	poolConfig := cconfig.PoolConfig{
		ConfigCenterType: cconfig.CC_TYPE_ETCD,
		ProjectName:      projectName,
		Environment:      env,
		Addresses:        addr,
		MaxActive:        3,
		MaxIdle:          1,
		IdleTimeout:      time.Second * 30,
	}
	pool, err := cconfig.NewPool(poolConfig)
	if err != nil {
		fmt.Printf("Can not new an etcd client pool for environment '%s': %s\n",
			env, err)
		return
	}
	defer pool.Destroy()
	// ......
}
```

客户端池的各个方法请参见本代码包的文档。

**注意：当环境变量 CONFIG\_CENTER\_ADDRS 的值为默认值时，运行上述程序的机器必须已连通腾讯云才能使该程序正常运行，否则会提示I/O超时错误。**

### 4. 注意事项
在运行使用了本代码包中的客户端或客户端池的程序时，需要注意当前操作系统中的相关环境变量的实际设定。

### 维护人员列表
+ 郝林（haolin@wepiao.com）