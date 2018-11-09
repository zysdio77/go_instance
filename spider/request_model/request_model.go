package main

import (
	"gopkg.in/resty.v1"
	"fmt"
)
//带参数get请求
func Get(requesturl string,params map[string]string) ([]byte,error) {
	resp,err := resty.R().SetQueryParams(params).Get(requesturl)
	if err != nil {
		return nil,err
	}
	b := resp.Body()
	return b,nil
}
//不带参数GET请求
func Get1(requesturl string) ([]byte,error) {
	resp,err := resty.R().Get(requesturl)
	if err != nil {
		return nil,err
	}
	b := resp.Body()
	return b,nil
}

func post1()  {
	resty.R().SetQueryParams().Post()
}
func main()  {
	params := make(map[string]string)
	params["menupath"]="1000000"
	b,_ := Get1("http://www.zeroz.com.cn")
	fmt.Println(string(b))
}