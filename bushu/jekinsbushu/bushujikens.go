package main

import (
	"os/exec"

	"go.pkg.wesai.com/p/base_lib/log"
	"fmt"
	"os"
	"crypto/md5"
	"io/ioutil"
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
func CheckSupervisor()  {

}

func StopSupervisor(supername string)  {
	shellcommand := "sudo supervisorctl stop "+ supername
	//shellcommand := "supervisorctl stop "+ supername
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("supervisorctl stop %v err :%v",supername,err)
	}
	//fmt.Printf("stop %v success\n",supername)
}

func StartSupervisor(supername string) bool {
	shellcommand := "sudo supervisorctl start "+ supername
	//shellcommand := "supervisorctl start  "+ supername
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("supervisorctl start %v err :%v",supername,err)
		return false
	}
	//fmt.Printf("start %v success\n",supername)
	return true
}

func ChechFile(filepath string) bool {
	_,err:=os.Stat(filepath)
		if os.IsNotExist(err){
			fmt.Println(err)
			log.DLogger().Errorf("check file err:%v",err)
			return false
		} else {
			return true
		}
}
func CopyFile(srcfilename,disfilename string) bool {
	shellcommand := "cp -rf "+srcfilename+" "+""+disfilename
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("copy file err:%v\n",err)
		return false
	}
	fmt.Printf("部署文件%v到%v successful\n",srcfilename,disfilename)
	return true
}
func FileRename(oldname,newname string)  {
	err := os.Rename(oldname,newname)
	if err != nil {
		//log.DLogger().Errorf("file rename err : %v",err)
		fmt.Errorf("file rename err : %v",err)
	} else {
		fmt.Println("备份可执行文件成功！！！")
	}

}
func DelHisFile(filename string)  {
	err := os.Remove(filename)
	if err != nil {
		fmt.Errorf("delete history backup file %v err : %v",filename,err)
	}
}

func Md5file(file string) string {

	f,err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf,err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	sum := md5.Sum(buf)
	s :=fmt.Sprintf("%x",sum)
	return s

}
func main()  {
	if len(os.Args) != 2{
		fmt.Println("Usage: you need a arg ")
		os.Exit(2)
	}
	/*
	supervisordisc := make(map[string]string)
	supervisordisc["chili"]="chili/chili"
	supervisordisc["sudhana"]="sudhana/sudhana"
	supervisordisc["lips"]="lips/lips"
	supervisordisc["bee"]="bee/bee"
	supervisordisc["dragonball"]="dragonball/dragonball"
	supervisordisc["dragonvstiger"]="dragon_vs_tiger/dragon_vs_tiger"
	supervisordisc["egypt"]="egypt/egypt"
	supervisordisc["fafafa"]="fafafa/fafafa"
	supervisordisc["nereids"]="nereids/nereids"
	supervisordisc["cat"]="manekineko/cat"
	supervisordisc["monkeyking"]="monkey_king/monkey_king"
	supervisordisc["pearl"]="pearl/pearl"
	supervisordisc["samba"]="samba/samba"
	supervisordisc["sweets"]="sweets/sweets"
	supervisordisc["dog"]="wangcai/dog"
*/
	publishdisc := make(map[string]string)
	publishdisc["chili"]="chili/chili"
	publishdisc["sudhana"]="sudhana/sudhana"
	publishdisc["lips"]="lips/lips"
	publishdisc["bee"]="bee/bee"
	publishdisc["dragonball"]="dragonball/dragonball"
	publishdisc["dragonvstiger"]="dragonVsTiger/dragonVsTiger"
	publishdisc["egypt"]="egypt/egypt"
	publishdisc["fafafa"]="fafafa/fafafa"
	publishdisc["nereids"]="nereids/nereids"
	publishdisc["cat"]="manekineko/manekineko"
	publishdisc["monkeyking"]="monkey_king/monkey_king"
	publishdisc["pearl"]="pearl/pearl"
	publishdisc["samba"]="samba/samba"
	publishdisc["sweets"]="sweets/sweets"
	publishdisc["dog"]="wangcai/wangcai"
	publishdisc["pet"]="pet/pet"
	publishdisc["fuguimao"]="fuguimao/fuguimao"
	publishdisc["treasure_bowl"]="treasure_bowl/treasure_bowl"

	//publishpath := "/home/ubuntu/work/publish/"
	publishpath := "/home/ubuntu/work/publish/combodia/"
	//测试
	gitrepopath := "/home/ubuntu/.jenkins/workspace/casino-machine/c-publish/"
	//正式
	//gitrepopath := "/home/ubuntu/sync_git_repo/lele_publish/"
	argname := os.Args[1]
	name :=publishdisc[argname]

	publishname := publishpath+name
	gitreponame := gitrepopath+name
	backupname := publishname+".bak"

	//fmt.Println(publishname,gitreponame,backupname)
	if ChechFile(backupname) {
		DelHisFile(backupname)
	}
	gitreponamemd5sum := Md5file(gitreponame)
	//fmt.Println(gitreponame)
	//fmt.Println(gitreponamemd5sum)
	publishnamemd5sum := Md5file(publishname)
	//fmt.Println(publishname)
	//fmt.Println(publishnamemd5sum)
	if gitreponamemd5sum == publishnamemd5sum {
		fmt.Printf("%s 和 %s 文件一样，没有改动，无需从新部署！！！\n",gitreponame,publishname)
	} else {
		StopSupervisor(argname)
		FileRename(publishname,backupname)
		c := CopyFile(gitreponame,publishname)
		if c {
			gitreponamemd5sum := Md5file(gitreponame)
			publishnamemd5sum := Md5file(publishname)
			if gitreponamemd5sum == publishnamemd5sum {
				s :=StartSupervisor(argname)
				if s {
					fmt.Println("部署成功！！！")
				} else {
					fmt.Println("重启服务失败！！！")
				}
			} else {
				fmt.Println("同步文件成功，但md5校验有问题！！！！")
			}
		} else {
			fmt.Println("同步文件失败！！！")
		}
	}

}