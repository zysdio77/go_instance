package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"strconv"
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
	CheckServer() int
	RepairServer()
}

func (c *Cli) SshClient() (*ssh.Client, error) {
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
func (c *Cli) SshSession(client *ssh.Client) (*ssh.Session, error) {
	//NewSession为这个客户端打开一个新的会话。会话是远程的 执行一个程序

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *Cli)CheckServer() (status int) {
	client, err := c.SshClient()
	defer client.Close()
	if err != nil {
		fmt.Println(err)
	}
	session, err := c.SshSession(client)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	b,err := session.Output("netstat -autpn | grep 16776 | wc -l")
	if err != nil {
		fmt.Println(err)
	}
	//去掉\n
	s := strings.Replace(string(b),"\n","",-1)
	n,err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	if n < 1{
		fmt.Println("梯子挂了～～～～，稍等，哥看看是咋回事～～～")
		status =1
	} else {
		fmt.Println("梯子没问题，翻起来吧～～")
	}
	return status
}

func (c *Cli)RepairServer()  {
	client, err := c.SshClient()
	defer client.Close()
	if err != nil {
		fmt.Println(err)
	}
	session, err := c.SshSession(client)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	//标准输入，输出，错误，不加会没有输出
	//session.Stdout = os.Stdout
	//session.Stdin = os.Stdin
	//session.Stderr = os.Stderr
	//执行远程命令
	session.Run("/usr/bin/ssserver -c /usr/local/shadowsock/shadowsock.json -d start")

}

func main() {
	var c Cli
	c.User = "root"
	c.Password = "root"
	c.IPAddr = "127.0.0.1"
	c.Port = "22"

	status := c.CheckServer()

	if status == 1 {
		c.RepairServer()
		time.Sleep(time.Second*1)
		c.CheckServer()
	}

}
