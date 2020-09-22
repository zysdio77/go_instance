//package StressTesting
package main

import (
"fmt"
"github.com/go-resty/resty/v2"
"sync"
)

func GetRequest()  {
	client  :=resty.New()
	//resq, err := client.R().Get("http://47.252.1.155:8080/Mega-casino-temp/data/getCurrentTime")
	for i:=0;i<100;i++ {
		resq, err := client.R().Get("http://47.252.16.175:8080/Cash-hoard-slots/")
		//resq, err := client.R().Get("http://47.252.1.155:8080/Mega-casino-temp/data/getCurrentTime")
		if err != nil {
			fmt.Println(err)
			//return
		}
		fmt.Println(resq,resq.Status())
	}
	wg.Done()
}
var wg sync.WaitGroup
func main()  {
	for i:=0;i<100;i++{
		wg.Add(1)
		go GetRequest()
	}
	wg.Wait()
}