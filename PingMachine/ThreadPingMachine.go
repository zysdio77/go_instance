package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func CheckMachine(shell string) int {
	_, err := exec.Command("/bin/bash", "-c", shell).Output()
	if err != nil {
		return 0
	} else {
		return 1
	}
}
func threadip(i int, c chan int) {
	ipaddr := "192.168.1." + strconv.Itoa(i)
	shell := "ping -c 2 " + ipaddr
	r := CheckMachine(shell)
	c <- r
	if r == 1 {
		fmt.Println(i)
	}
}
func main() {
	a := time.Now().Unix()
	//fmt.Println(a)
	c := make(chan int, 10)
	for i := 1; i < 255; i++ {
		go threadip(i, c)
	}
	for i := 1; i < 255; i++ {
		<-c
	}
	b := time.Now().Unix()
	//fmt.Println(b)
	fmt.Println("共用时：", b-a, "秒")
}
