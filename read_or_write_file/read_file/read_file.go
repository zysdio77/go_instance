package read_file

import (
	"os"
	"fmt"
	"io/ioutil"
)

func Readfile(infile string) (data []byte,err error) {	//os用法
	//只读方式打开文件
	file,err := os.Open(infile)
	if err != nil{
		fmt.Println("readfile os open err:",err)
		return nil,err
	}
	defer file.Close()
	b,err :=ioutil.ReadAll(file)
	if err != nil {
		return nil,err
	}
	return b,nil
}

func IoutilReadFile(filename string) ([]byte,error) { 	//iotuil的用法
	b,err := ioutil.ReadFile(filename)
	if err != nil{
		return nil,err
	}
	return b,nil
}