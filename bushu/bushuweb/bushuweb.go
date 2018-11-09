package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"os/exec"
	"net/http"
)

func main()  {
	router :=gin.Default()
	router.GET("/bushuwebfile", func(c *gin.Context) {
		err := Ecexsyncfile()
		if err != nil {
			c.String(http.StatusOK,"同步文件失败！！！！")
		}else {
			c.String(http.StatusOK,"同步文件成功！！！！")
		}
		err1 := Ecexcpfile()
		if err1 != nil {
			c.String(http.StatusOK,"部署前端文件失败！！！！")
		}else {
			c.String(http.StatusOK,"部署前端文件完成！！！！")
		}
	})

	router.Run(":12345")
}

func Ecexsyncfile() error{
	shell := "sudo rsync -a --delete ruser@192.168.2.90::webfile /home/rddd1212/bushu/syncfile/"
	cmd := exec.Command("/bin/bash","-c",shell)
	err := cmd.Run()
	if err != nil{
		return err
	}
	return nil

}

func Ecexcpfile() error{
	shell := "sudo cp -rf /home/rddd1212/bushu/syncfile/* /home/rddd1212/install/nginx/html/"
	cmd := exec.Command("/bin/bash","-c",shell)
	err := cmd.Run()
	if err != nil{
		return err
	}
	return nil

}