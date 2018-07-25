log
=====================
顾名思义，本代码包用于实现日志记录功能。目前封装了[logrus](https://github.com/Sirupsen/logrus)和[zap](https://github.com/uber-go/zap)，并提供了统一的访问接口。

### 1. 用意与用途
本代码包旨在隔离底层的日志记录功能，为使用不同的日志库的用户提供统一的访问接口。本代码包已对封装的日志库做了些许定制。

### 2. 使用方法
本代码包的主要访问接口请详见 logger.go 文件。下面的示例仅仅是示意。

### 3. 使用示例

```go
package main

import (
	"time"

	"go.pkg.wesai.com/p/base_lib/log"
)

func main() {
	logger := log.DLogger()
	logger.Infof("Record via logger '%s'.\n", logger.Name())
}
```

### 4. 注意事项
本代码对封装的日志库进行了定制。这些定制可能并不适合你的项目。如果确实如此，请直接使用这些或其它的日志库。本代码包的主要目的在于提高更换日志库时的便捷性，并提供使用这些日志库的最佳实践（可能不适用于你的项目）。

另外，我们虽然统一了访问接口，但是对于不同的日志库来说有些方法的实现是不同的，比如：`Fatal`和`Panic`。所以在通过本代码包使用某个日志库之前请先查看相应的子代码包中的测试源码文件。

### 维护人员列表
+ 郝林（haolin@wesai.com）
