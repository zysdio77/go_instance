package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"strings"
	"time"
)

type Cli struct {
	User     string
	Password string
	IPAddr   string
	Port     string
}
type CliMethod interface {
	SshClient() (*ssh.Client, error)
	SshSession(client *ssh.Client) (*ssh.Session, error)
}

func (c Cli) SshClient() (*ssh.Client, error) {
	addr := c.IPAddr + ":" + c.Port
	//client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
	//Dial启动到给定SSH服务器的客户机连接
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		//需要验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

//会话表示到远程命令或shell的连接
func (c Cli) SshSession(client *ssh.Client) (*ssh.Session, error) {
	//NewSession为这个客户端打开一个新的会话。会话是远程的 执行一个程序
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func main() {
	var c Cli
	c.User = "root"
	c.Password = "123"
	c.IPAddr = "127.0.0.1"
	c.Port = "22"

	var cm CliMethod
	cm = &c
	client, err := cm.SshClient()
	defer client.Close()
	if err != nil {
		fmt.Println(err)
	}
	session, err := cm.SshSession(client)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	///////////////
	//输出在远程主机上运行cmd并返回其标准输出。
	b, err := session.Output("netstat -autpn | grep 16776 | wc -l")
	if err != nil {
		fmt.Println(err)
	}
	//去掉\n
	s := strings.Replace(string(b), "\n", "", -1)
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	if n < 1 {
		fmt.Println("梯子挂了～～～～，稍等，让剑哥看看是咋回事～～～")
	} else {
		fmt.Println("梯子还好，翻起来吧～～")
	}
	time.Sleep(time.Second * 1)

}
