package read_file

import (
	"bufio"
	"fmt"
	"go.pkg.wesai.com/p/base_lib/log"
	"io"
	"io/ioutil"
	"os"
)

func Readfile(infile string) (data []byte, err error) { //os用法
	//只读方式打开文件
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println("readfile os open err:", err)
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IoutilReadFile(filename string) ([]byte, error) { //iotuil的用法
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readlines(filename string) (lines []string) {
	f, err := os.Open(filename)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for true {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.DLogger().Fatal(err)
		}
		lines = append(lines, string(line))
	}
	return
}
