package dticker

import (
	"testing"
	"time"
)

func TestDTicker(t *testing.T) {
	cAcc := int64(1)
	// 最小周期，10毫秒
	minPeriod := int64(1e7)
	pfFunc := func() int64 {
		if cAcc > 5 {
			cAcc = 1
		}
		result := minPeriod * cAcc
		cAcc++
		return result
	}
	// pfFuncWithoutAcc := func() int64 {
	// 	return unit * cAcc
	// }
	refreshInterval := int64(time.Millisecond * 500)
	dTicker := NewDynamicTicker(pfFunc, refreshInterval)
	if dTicker == nil {
		t.Fatal("Invalid Dynamic Ticker!")
	}
	// 常规信号发送的最大误差，5毫秒
	maxDiffNormal := int64(1e6 * 5)
	// 周期变更后的首次信号发送的最大误差，60毫秒
	maxDiffFirst := int64(1e7 * 6)
	var maxDiff int64
	var count int
	currentPeriod := minPeriod
	// 经实验，误差的值并不是发散的
	for {
		if count == 200 { // 经实验，1000也可通过测试
			break
		}
		begin := time.Now().UnixNano()
		dTicker.Sign()
		end := time.Now().UnixNano()
		actualPeriod := end - begin
		frequency := dTicker.Period()
		if currentPeriod != frequency {
			maxDiff = maxDiffFirst
			currentPeriod = frequency
		} else {
			maxDiff = maxDiffNormal
		}
		diff := actualPeriod - frequency
		if diff > maxDiff || diff < -maxDiff {
			t.Errorf("Timeout! E: %d, A: %d, D: %d, DD: %d! (maxD=%d, C=%d)\n",
				frequency, actualPeriod, diff, diff-maxDiff, maxDiff, count)
			// } else {
			// 	t.Logf("E: %d, A: %d, D: %d, DD: %d! (maxD=%d, C=%d)\n",
			// 		frequency, actualPeriod, diff, diff-maxDiff, maxDiff, count)
		}
		count++
	}
}
