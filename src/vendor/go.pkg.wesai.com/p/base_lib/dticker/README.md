dticker
=====================
本代码包的全称是：Dynamic ticker，即：动态的断续器。

### 0. 版本：1.0

### 1. 用意与用途
本代码包旨在扩展Go语言标准库中的断续器，主要添加了动态更新断续间隔时间的功能。

### 2. 使用方法
动态断续器是并发安全的。

动态断续器的最小断续间隔时间是10毫秒，最小刷新断续间隔时间的间隔时间是100毫秒。

在初始化动态断续器时，使用者应该传入用于获取周期值的函数以及刷新断续间隔时间的间隔时间。具体的方法请参加下面的示例。

### 3. 使用示例
初始化及使用方法：

```go
package main

import (
	"fmt"
	"time"

	"go.pkg.wesai.com/p/base_lib/dticker"
)

func main() {
	// 累加器。
	cAcc := int64(1)
	// 单位周期：100毫秒
	minPeriod := int64(1e8)
	// 用于获取周期值的函数，会有规律的返回不同的数值。
	// 该函数需要需要实现函数类型 dticker.PeriodFetchFunc。
	pfFunc := func() int64 {
		if cAcc > 5 {
			cAcc = 1
		}
		result := minPeriod * cAcc
		cAcc++
		return result
	}
	// 刷新间隔
	refreshInterval := int64(time.Millisecond * 500)
	dTicker := dticker.NewDynamicTicker(pfFunc, refreshInterval)
	// 使用 dTicker
	for i := 0; i < 100; i++ {
		dTicker.Sign() // 等待断续器的信号
		fmt.Printf("The time(period: %d ns): %s\n",
			dTicker.Period(), time.Now().Format("2006-01-02 15:04:05.000"))
	}
	fmt.Println("End.")
}
```

### 4. 注意事项
与标准库中的断续器一样，动态断续器常常被用于对最小操作时间间隔有要求的场景中。而对于控制最大操作时间间隔，断续器是无能为力的（这种情况下应该使用定时器并配合 select 语句使用）。

另外，需要注意，断续器的实际间隔时间可能会存在误差，但通常会大于等于 -1 毫秒。

### 维护人员列表
+ 郝林（haolin@wepiao.com）