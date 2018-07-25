package dticker // import "go.pkg.wesai.com/p/base_lib/dticker"

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	MIN_PERIOD           = 1e7 // 最小周期，10毫秒。
	TIMEOUT_DURATION     = 1e7 // 超时时间，10毫秒。
	MIN_REFRESH_INTERVAL = 1e8 // 最小刷新间隔，100毫秒。

)

// 用于获取周期值的函数的类型。
type PeriodFetchFunc func() int64

// 动态断续器。
type DynamicTicker struct {
	refreshInterval time.Duration
	pfFunc          PeriodFetchFunc
	period          int64
	tickerValue     *atomic.Value
}

func NewDynamicTicker(pfFunc PeriodFetchFunc, refreshInterval int64) *DynamicTicker {
	if pfFunc == nil {
		panic(errors.New("Invalid period fetcher!"))
	}
	if refreshInterval <= MIN_REFRESH_INTERVAL {
		refreshInterval = MIN_REFRESH_INTERVAL
	}
	dTicker := &DynamicTicker{
		refreshInterval: time.Duration(refreshInterval),
		pfFunc:          pfFunc,
		tickerValue:     &atomic.Value{},
	}
	dTicker.asyncFresh()
	return dTicker
}

// 获取断续器的信号，信号来临之前该方法会阻塞。
func (dTicker *DynamicTicker) Sign() {
	timeout := time.NewTimer(TIMEOUT_DURATION)
	for {
		ticker := dTicker.tickerValue.Load().(*time.Ticker)
		select {
		case <-ticker.C:
			timeout.Stop()
			return
		case <-timeout.C:
			timeout.Reset(TIMEOUT_DURATION)
			continue
		}
	}
}

// 获取断续器的当前的周期值，单位：纳秒。
func (dTicker *DynamicTicker) Period() int64 {
	return atomic.LoadInt64(&dTicker.period)
}

// 异步刷新断续器的周期。
func (dTicker *DynamicTicker) asyncFresh() {
	dTicker.fresh()
	go func() {
		ch := time.Tick(dTicker.refreshInterval)
		for _ = range ch {
			dTicker.fresh()
		}
	}()
}

func (dTicker *DynamicTicker) fresh() {
	newPeriod := dTicker.pfFunc()
	if newPeriod < MIN_PERIOD {
		newPeriod = MIN_PERIOD
	}
	if newPeriod != atomic.LoadInt64(&dTicker.period) {
		oldTicker := dTicker.tickerValue.Load()
		newTicker := time.NewTicker(time.Duration(newPeriod))
		dTicker.tickerValue.Store(newTicker)
		atomic.StoreInt64(&dTicker.period, newPeriod)
		if oldTicker != nil {
			oldTicker.(*time.Ticker).Stop()
		}
	}
}
