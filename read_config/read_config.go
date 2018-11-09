package read_config

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"flag"
	"fmt"
)

type ConfigInfo struct {
	DBConfig MysqlConfig `toml:"mysql"`		//config文件中的[mysql]标签
	Templateconfig TemplateConfig `toml:"template"`	//config文件中的[tempalte]标签
}


type MysqlConfig struct {	//[mysql]标签下的配置
	Host 		string `toml:"host"`
	Port        int    `toml:"port"`
	DBName      string `toml:"db"`
	User        string `toml:"user"`
	Pwd         string `toml:"password"`
}

type TemplateConfig struct {	//[template]标签下的配置
	Dir string `toml:"dir"`
}


func ParseConfig(filename string) (*ConfigInfo,error) {
	var config ConfigInfo
	//ReadFile读取文件名指定的文件并返回内容
	data,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,err
	}

	//解析toml格式的配置文件，把解码内容存入config指针
	_,err =toml.Decode(string(data),&config)
	if err != nil {
		return nil,err
	}
	return &config,nil
}

func main()  {

	var filename string
	//stringvar方法第一个参数必须的是string类型的指针，相当于把配合文件的路径赋值给了filename
	//第二个参数是参数关键字，命令行中的选项就是-config="filename"，
	//第三个参数是默认filename的路径
	//第四个参数是方法提示
	flag.StringVar(&filename,"config","/Users/zhangyongsheng/data/src/go_instance/read_config/config.toml","-config=configfile")

	//解析命令行参数
	flag.Parse()

	conf,err := ParseConfig(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
}