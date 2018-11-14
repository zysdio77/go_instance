package main

import (
	"fmt"
	"gopkg.in/redis.v4"
)

func Sentinel() {
	//链接sentibel
	cli := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",	//主节点名字
		SentinelAddrs: []string{"104.225.154.39:26379", "104.225.154.39:26380", "104.225.154.39:26381"},	//sentinel链接地址
		DB:3,	//用哪个库
	})
	defer cli.Close()

	//插入数据
	err := cli.Set("zhang", "zhangyongsheng", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	//获取数据
	value, err := cli.Get("zhang").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)
}

func usualy() {
	//链接数据库
	cli := redis.NewClient(&redis.Options{Addr: "104.225.154.39:6379", Password: "",DB:3})
	defer cli.Close()

	//插入数据
	err := cli.Set("yong", "yong", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	//获取数据
	value, err := cli.Get("zhang").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)

	//删除数据
	number, err := cli.Del("zhang").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(number)

}

func main() {
	Sentinel()
	//usualy()

}
