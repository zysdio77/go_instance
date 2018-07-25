uuid
=====================
本代码包提供了可以生成全局唯一ID的函数。

### 1. 用意与用途
在分布式环境下，需要UUID的场景很多。这里，我们提供一种通用的方法，并保证在多机环境下UUID的唯一性。

### 2. 使用方法
本代码包目前只有一个函数`UUID`。调用该函数时需要提供一个bool类型的参数`sys`。参数`sys`用于指定是否利用操作系统的伪随机数生成机制（以下简称随机机制）。这种随机机制仅在类Unix系统下存在。若`sys`的值为`false`或`sys`的值为`true`但随机机制无效，则会利用系统时间戳和`math/rand`包内的随机数生成函数组合生成UUID。

### 3. 使用示例

```go
package main

import (
	"fmt"

	"go.pkg.wesai.com/p/base_lib/uuid"
)

func main() {
	fmt.Println(uuid(true))
}
```

### 4. 注意事项
通常情况下，函数`UUID`的参数`sys`应该被给定为`true`。

### 维护人员列表
+ 郝林（haolin@wesai.com）
