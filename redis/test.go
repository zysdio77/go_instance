package main

import (
	"fmt"
	"gopkg.in/redis.v4"
)

func Sentinel() {
	//链接sentibel
	cli := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"104.225.154.39:26379", "104.225.154.39:26380", "104.225.154.39:26381"},
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
	cli := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: ""})
	defer cli.Close()

	//插入数据
	err := cli.Set("yong", "yong", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	//获取数据
	value, err := cli.Get("yong").Result()
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

	usualy()
}
