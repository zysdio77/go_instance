package headle

import (
	"gopkg.in/resty.v1"
	"io/ioutil"
	"strconv"
	"fmt"
	"os"
)

func Recover(i []string) []string{	//反序切片
	alist :=[]string{}
	n := len(i)
	for k,_ := range i{
		//fmt.Println(n)
		alist = append(alist,i[n-k-1])
	}
	return alist

}
func ReadWebFile(picaddr string) ([]byte,error) {	//读取网页文件内容
	resp,err := resty.R().Get(picaddr)
	if err != nil {
		return nil,err
	}
	return  resp.Body(),nil

}
func WriteLocalFile(saveaddr string,filedata []byte) error {	//写入本地文件
	err := ioutil.WriteFile(saveaddr,filedata,0644)
	if err != nil {
		return err
	}
	return nil
}

func SavePic(savepath string,piclist []string) error {	//批量读取网页文件并写入本地

	for k,v := range piclist {
		d,err := ReadWebFile(v)
		saveaddr :=savepath+strconv.Itoa(k)
		fmt.Println(saveaddr)
		if err != nil {
			return err
		}
		err = WriteLocalFile(saveaddr,d)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExistDir(picpath string) (bool,error) {
	_,err :=os.Stat(picpath)
	if err != nil {
		if os.IsExist(err) {
			return true,nil
		} else {
			return false,err
		}
	}
	return false,nil
}

func CreateDir(picpath string) error  {
	err := os.Mkdir(picpath,0755)
	if err != nil {
		return err
	}
	return  nil
}

func Dirlist(dirname string) ([]string ,error){
	v,err :=ioutil.ReadDir(dirname)
	if err != nil {
		return nil,err
	}
	var jj []string
	for i,j := range v{
		if i != 0 {
			jj = append(jj,j.Name())
		}
	}
	return jj,nil
}