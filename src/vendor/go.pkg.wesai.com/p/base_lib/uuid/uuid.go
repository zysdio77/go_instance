package uuid

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// 生成全局唯一的ID。
// 参数sys用于指定是否利用操作系统的伪随机数生成机制（以下简称随机机制）。这种随机机制仅在类Unix系统下存在。
// 若sys的值为false或sys的值为true但随机机制无效，则会利用系统时间戳和math/rand包内的随机数生成函数组合生成UUID。
func UUID(sys bool) string {
	b := make([]byte, 16)
	if sys {
		f, err := os.Open("/dev/urandom")
		if err == nil {
			f.Read(b)
			f.Close()
		}
	}
	if b[0] == 0 {
		putUint64(b, uint64(time.Now().UnixNano()), uint64(rand.Int63()))
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// putUint64 用于向指定的切片b中先后放入无符号整数v1和v2的大端形式字节序列。
// 注意！参数b的容量必须大于等于16！
func putUint64(b []byte, v1 uint64, v2 uint64) {
	b[0] = byte(v1 >> 56)
	b[1] = byte(v1 >> 48)
	b[2] = byte(v1 >> 40)
	b[3] = byte(v1 >> 32)
	b[4] = byte(v1 >> 24)
	b[5] = byte(v1 >> 16)
	b[6] = byte(v1 >> 8)
	b[7] = byte(v1)
	b[8] = byte(v2 >> 56)
	b[9] = byte(v2 >> 48)
	b[10] = byte(v2 >> 40)
	b[11] = byte(v2 >> 32)
	b[12] = byte(v2 >> 24)
	b[13] = byte(v2 >> 16)
	b[14] = byte(v2 >> 8)
	b[15] = byte(v2)
}
