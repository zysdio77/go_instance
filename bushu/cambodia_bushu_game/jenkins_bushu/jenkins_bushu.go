package main

import (
"os/exec"

"crypto/md5"
"fmt"
"go.pkg.wesai.com/p/base_lib/log"
"io/ioutil"
"os"
)

func ExecShellCommand(shellcommand string) error {
	//cmd := exec.Command("/bin/bash","-c","date")
	cmd := exec.Command("/bin/bash", "-c", shellcommand)
	result, err := cmd.Output()
	if err != nil {
		log.DLogger().Errorf("exec shell command %v err: %v", shellcommand, err)
		//fmt.Println(err)
		return err
	}
	fmt.Println(string(result))
	return nil
}
func CheckSupervisor() {

}

func StopSupervisor(supername string) {
	shellcommand := "sudo supervisorctl stop " + supername
	//shellcommand := "supervisorctl stop "+ supername
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("supervisorctl stop %v err :%v", supername, err)
	}
	//fmt.Printf("stop %v success\n",supername)
}

func StartSupervisor(supername string) bool {
	shellcommand := "sudo supervisorctl start " + supername
	//shellcommand := "supervisorctl start  "+ supername
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("supervisorctl start %v err :%v", supername, err)
		return false
	}
	//fmt.Printf("start %v success\n",supername)
	return true
}

func ChechFile(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		//fmt.Println(err)
		//log.DLogger().Errorf("check file err:%v", err)
		return false
	}
	return false
}

func CopyFile(srcfilename, disfilename string) bool {
	shellcommand := "cp -rf " + srcfilename + " " + "" + disfilename
	err := ExecShellCommand(shellcommand)
	if err != nil {
		fmt.Errorf("copy file err:%v\n", err)
		return false
	}
	fmt.Printf("部署文件%v到%v successful\n", srcfilename, disfilename)
	return true
}
func FileRename(oldname, newname string) {
	err := os.Rename(oldname, newname)
	if err != nil {
		//log.DLogger().Errorf("file rename err : %v",err)
		fmt.Errorf("file rename err : %v", err)
	} else {
		fmt.Println("备份可执行文件成功！！！")
	}

}
func DelHisFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Errorf("delete history backup file %v err : %v", filename, err)
	}
}

func Md5file(file string) string {

	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	sum := md5.Sum(buf)
	s := fmt.Sprintf("%x", sum)
	return s

}
func main() {
	if len(os.Args) != 2 {
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
	publishdisc["chili"] = "chili/chili"
	publishdisc["sudhana"] = "sudhana/sudhana"
	publishdisc["lips"] = "lips/lips"
	publishdisc["bee"] = "bee/bee"
	publishdisc["dragonball"] = "dragonball/dragonball"
	publishdisc["dragonvstiger"] = "dragonVsTiger/dragonVsTiger"
	publishdisc["egypt"] = "egypt/egypt"
	publishdisc["fafafa"] = "fafafa/fafafa"
	publishdisc["nereids"] = "nereids/nereids"
	publishdisc["manekineko"] = "manekineko/manekineko"
	publishdisc["monkeyking"] = "monkey_king/monkey_king"
	publishdisc["pearl"] = "pearl/pearl"
	publishdisc["samba"] = "samba/samba"
	publishdisc["sweets"] = "sweets/sweets"
	publishdisc["wangcai"] = "wangcai/wangcai"
	//publishdisc["pet"] = "pet/pet"
	publishdisc["fuguimao"] = "fuguimao/fuguimao"
	publishdisc["treasure_bowl"] = "treasure_bowl/treasure_bowl"
	publishdisc["night"] = "night/night"

	//正式
	//publishpath := "/home/ubuntu/work/cambodia/"
	//gitrepopath := "/home/ubuntu/sync_git_repo/cambodia/"
	//stage
	publishpath := "/home/ubuntu/work/publish/combodia/"
	gitrepopath := "/home/ubuntu/.jenkins/workspace/casino-machine/c-publish/"
	argname := os.Args[1]

	if argname == "all" {
		for i,j := range publishdisc {
			//name := publishdisc[i]
			name := j
			publishname := publishpath + name
			fmt.Println(publishname)
			gitreponame := gitrepopath + name
			backupname := publishname + ".bak"
			jsonfile := name + "_config.json"
			publishjsonfile := publishpath + jsonfile
			gitrepojsonfile := gitrepopath + jsonfile
			if ChechFile(backupname) {
				DelHisFile(backupname)
			}
			gitreponamemd5sum := Md5file(gitreponame)
			gitrepojsonfilemd5sum := Md5file(gitrepojsonfile)
			publishjsonfilemd5sum := Md5file(publishjsonfile)
			//fmt.Println(gitreponame)
			//fmt.Println(gitreponamemd5sum)
			publishnamemd5sum := Md5file(publishname)
			//fmt.Println(publishname)
			//fmt.Println(publishnamemd5sum)
			gamecode := 0
			jsoncode := 0
			if gitreponamemd5sum == publishnamemd5sum {
				//fmt.Printf("%s 和 %s 文件一样，没有改动，无需从新部署！！！\n", gitreponame, publishname)
				fmt.Printf("%s 执行文件，没有改动，无需从新部署！！！\n", i)
			} else {
				gamecode = 1
			}
			if gitrepojsonfilemd5sum == publishjsonfilemd5sum {
				//fmt.Printf("%s 和 %s json文件一样，没有改动，无需从新部署！！！\n", gitrepojsonfile, publishjsonfile)
				fmt.Printf("%s json文件一样，没有改动，无需从新部署！！！\n", i)
			} else {
				jsoncode = 1
			}
			if gamecode == 1 || jsoncode == 1 {
				StopSupervisor(i)
				if gamecode == 1 {
					FileRename(publishname, backupname)
					c := CopyFile(gitreponame, publishname)
					if c {
						gitreponamemd5sum := Md5file(gitreponame)
						publishnamemd5sum := Md5file(publishname)
						if gitreponamemd5sum == publishnamemd5sum {
							gamecode = 2
						} else {
							fmt.Println("同步文件成功，但md5校验有问题！！！！")
						}
					}
				}
				if jsoncode == 1 {
					c := CopyFile(gitrepojsonfile, publishjsonfile)
					if c {
						gitrepojsonfilemd5sum := Md5file(gitrepojsonfile)
						publishjsonfilemd5sum := Md5file(publishjsonfile)
						if gitrepojsonfilemd5sum == publishjsonfilemd5sum {
							jsoncode = 2
						} else {
							fmt.Println("同步文件成功，但md5校验有问题！！！！")
						}
					}
				}
			}
			if gamecode == 2 || jsoncode ==2 {
				s := StartSupervisor(i)
				if s {
					fmt.Println("部署成功！！！")
					fmt.Println("####################")
				} else {
					fmt.Println("重启服务失败！！！")
					fmt.Println("####################")
				}
			}
		}
	} else {
	name := publishdisc[argname]

	publishname := publishpath + name
	gitreponame := gitrepopath + name
	backupname := publishname + ".bak"

	jsonfile := name + "_config.json"
	publishjsonfile := publishpath + jsonfile
	gitrepojsonfile := gitrepopath + jsonfile
	//fmt.Println(publishname,gitreponame,backupname)
	if ChechFile(backupname) {
		DelHisFile(backupname)
	}
	gitreponamemd5sum := Md5file(gitreponame)
	gitrepojsonfilemd5sum := Md5file(gitrepojsonfile)
	publishjsonfilemd5sum := Md5file(publishjsonfile)
	//fmt.Println(gitreponame)
	//fmt.Println(gitreponamemd5sum)
	publishnamemd5sum := Md5file(publishname)
	//fmt.Println(publishname)
	//fmt.Println(publishnamemd5sum)
	gamecode := 0
	jsoncode := 0

	if gitreponamemd5sum == publishnamemd5sum {
		//fmt.Printf("%s 和 %s 文件一样，没有改动，无需从新部署！！！\n", gitreponame, publishname)
		fmt.Printf("%s 执行文件，没有改动，无需从新部署！！！\n", argname)
	} else {
		gamecode = 1
	}
	if gitrepojsonfilemd5sum == publishjsonfilemd5sum {
		fmt.Printf("%s json文件一样，没有改动，无需从新部署！！！\n", argname)
	} else {
		jsoncode = 1
	}
	if gamecode == 1 || jsoncode == 1 {
		StopSupervisor(argname)
		if gamecode == 1 {
			FileRename(publishname, backupname)
			c := CopyFile(gitreponame, publishname)
			if c {
				gitreponamemd5sum := Md5file(gitreponame)
				publishnamemd5sum := Md5file(publishname)
				if gitreponamemd5sum == publishnamemd5sum {
					gamecode = 2
				} else {
					fmt.Println("同步文件成功，但md5校验有问题！！！！")
				}
			}
		}
		if jsoncode == 1 {
			c := CopyFile(gitrepojsonfile, publishjsonfile)
			if c {
				gitrepojsonfilemd5sum := Md5file(gitrepojsonfile)
				publishjsonfilemd5sum := Md5file(publishjsonfile)
				if gitrepojsonfilemd5sum == publishjsonfilemd5sum {
					jsoncode = 2
				} else {
					fmt.Println("同步文件成功，但md5校验有问题！！！！")
				}
			}
		}
	}
	if gamecode == 2 || jsoncode ==2 {
		s := StartSupervisor(argname)
		if s {
			fmt.Println("部署成功！！！")
			fmt.Println("####################")
		} else {
			fmt.Println("重启服务失败！！！")
			fmt.Println("####################")
		}
	}
}
}