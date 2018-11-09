package write_file

import (
	"os"
	"fmt"
	"io/ioutil"
	"go_instance/read_or_write_file/read_file"
)

func IoutilWriteFile(filename string,words_byte []byte) error {
	//创建文件并写入，不能追加，每次都是从头开始写
	err := ioutil.WriteFile(filename,words_byte,0755)
	if err != nil {
		return err
	}
	return nil
}

func OsWriteFile(words,filename string) error {
	//os.Create方法是创建新文件，如果文件存在会从新创建
	//f,err :=os.Create(filename)

	//os.OpenFile方法打开文件配置度高
	//第一个参数是需要打开的文件名字
	//第二个参数是需要的方式，多个权限用|分割，下列权限是读写方式，文件不存在则创建，追加写入
	//第三个参数是文件的权限
	f,err := os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	defer f.Close()
	if err != nil{
		return err
	}
	//写入文件,写的是字符切片(可以写音乐，视频，图片)，返回写入了多少个字节数，
	//n,err :=f.Write([]byte(words))
	//写入字符串到文件,返回写入了多少个字节数
	n,err := f.WriteString(words)
	if err != nil{
		return err
	}
	fmt.Println(n)

	return nil

}

func main()  {
	data ,err :=read_file.IoutilReadFile("/Users/zhangyongsheng/Desktop/801s.jpg")
	//fmt.Println(data)
	if err != nil {
		fmt.Println(err)
	}
	filename := "/Users/zhangyongsheng/Desktop/3.txt"
	//words := "heelo world"
	//OsWriteFile(words,"2.txt")
	//words_byte := []byte(words)
	IoutilWriteFile(filename,data)
}