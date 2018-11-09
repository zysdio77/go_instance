package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strconv"
	"os/exec"
	"go.pkg.wesai.com/p/base_lib/log"
	"syscall"
)


func ExecShellCommand(shellcommand string) error {
	//cmd := exec.Command("/bin/bash","-c","date")
	cmd := exec.Command("/bin/bash","-c",shellcommand)
	result,err :=cmd.Output()
	if err != nil {
		log.DLogger().Errorf("exec shell command %v err: %v",shellcommand,err)
		//fmt.Println(err)
		return err
	}
	fmt.Println(string(result))
	return nil
}

func Newpath(path string) (string,error) {
	files,err := ioutil.ReadDir(path)
	if err != nil{
		return "",err
	}
	filesumint := len(files)+1
	filesumstr := strconv.Itoa(filesumint)
	b:= path+"v"+filesumstr
	//fmt.Println(b)
	oldMask := syscall.Umask(0)	//设置umask为0000
	err = os.Mkdir(b,0755)
	if err != nil {
		return "",err
	}
	syscall.Umask(oldMask)	//改回原来的umask

	return b,nil
}
func Cpfile(soucepath, dispath string) error {
	cpcommand:= fmt.Sprintf("cp -ra %s/* %s",soucepath,dispath)
	err := ExecShellCommand(cpcommand)
	if err != nil {
		return err
	}
	return nil
}

func main()  {
	if len(os.Args) != 2{
		fmt.Println("Usage : you need arg!!!")
	} else {
		aa := os.Args[1]
		soucepath := "/home/ubuntu/.jenkins/workspace/casino-website/Slots_cambodia/test/"
		dispath := "/var/www/html/cambodia/"
		soucepath2 := soucepath+aa
		dispath2,err := Newpath(dispath+aa+"/")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("create successful ",dispath2)
		}

		err = Cpfile(soucepath2,dispath2)
		if err!= nil {
			fmt.Println("Cpfile err",err)
		} else {
			fmt.Println("cp file ok")
		}

	}
}
