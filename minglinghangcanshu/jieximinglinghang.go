package main
import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"io"
	"strconv"
)

func main()  {
	infile := flag.String("i","infile","in file ")
	outfile := flag.String("o","outfile","out file")
	flag.Parse()
	var values []int
	var err error
	if infile != nil {
		fmt.Println("infile =",*infile)
		values ,err = readfile(*infile)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(values)
		}
	}

	if outfile != nil {
		err = writefile(values,*outfile)
		if err != nil {
			fmt.Println(err)
		}else {
			fmt.Println("write file ok")
		}
	}
}
//读取文件
func readfile(infile string) (values []int,err error) {
	file,err := os.Open(infile)
	if err != nil{
		fmt.Println("readfile os open err:",err)
		return nil,err
	}
	defer file.Close()
	br :=bufio.NewReader(file)
	values = make([]int,0)
	for {
		line,isprefix,err1 :=  br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF{
				err =err1
			}
			break
		}
		if isprefix {
			fmt.Println("A too long line")
			return
		}
		str := string(line)
		value ,err1 := strconv.Atoi(str)
		if err1 != nil {
			err = err1
			return
		}
		values = append(values,value)
	}
	return
}

//写入文件

func writefile(values []int,outfile string) error{
	file,err := os.Open(outfile)
	if err != nil {
		fmt.Println("faild to create outfile",outfile)
		return err
	}
	defer file.Close()
	for _,value := range values{
		str := strconv.Itoa(value)
		file.WriteString(str+"\n")
		fmt.Println("write ",value)
	}
	return nil
}