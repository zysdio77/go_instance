package main

import (
	"sync"
	"fmt"
	"time"
)


func main()  {
	var mission []int
	for i:=0;i<39;i++{
		mission = append(mission,i)
	}
	var wg sync.WaitGroup

	fmt.Println(mission)
	n:=0
	//每个协程需要执行的任务适量
	nz :=len(mission)/5
	//剩余的任务数量
	ny := len(mission)% 5
	//ny := len(mission)%5
	//启动5个协程，每启动一个，wg+1
	for i:=0;i<5;i++{
		//数据乱加个延时试试
		time.Sleep(time.Millisecond)
		wg.Add(1)
		go func() {
			//每个协程结束，wg-1
			defer wg.Done()
			//每个协程执行的任务数量nz
			for j:=0;j<nz;j++{
				fmt.Println(mission[n])

				n =n+1
			}
		}()

	}
//剩余的任务数量ny
	for k := len(mission) - ny; k < len(mission); k++ {
		fmt.Println(k,mission[k])
	}

	//阻塞在这里，保证协程运行完，wg=0时解除阻塞
	wg.Wait()
}