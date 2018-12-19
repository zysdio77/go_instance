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
		soucepath := "/home/ubuntu/.jenkins/workspace/cambodia_game_website/Slots_cambodia/prepare/"
		dispath := "/data/cambodia/game-website/cambodia/"
		if aa == "all" {
			alllist := []string{"AC01","AC02","AC03","AC04","AC05","AC06","AC07","AC08","AC09","AC10","AC11","AC12","AC14","AC15","AC17","HP07","HP09"}
			for _,j := range alllist{
				soucepath2 := soucepath+j
				dispath2,err := Newpath(dispath+j+"/")
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
		} else {
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
}
