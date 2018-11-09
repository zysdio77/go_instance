package main

import (
	"crypto/md5"
	"os"
	"fmt"
	"io/ioutil"
)

func main()  {

	f,err := os.Open("/Users/zhangyongsheng/Downloads/test111")
	if err != nil {
		fmt.Println(err)
	}
	buf,err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	sum := md5.Sum(buf)

	s :=fmt.Sprintf("%x",sum)
	fmt.Printf("%T",s)
}
