package main

import (
	"gopkg.in/redis.v4"
	"fmt"
)

func main()  {
	//链接数据库
	cli := redis.NewClient(&redis.Options{Addr:"192.168.2.237:6379",Password:""})

	//获取数据
	value,err := cli.Get("xxx").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)

	//插入数据
	err = cli.Set("zhang", "zhangyongsheng", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

}