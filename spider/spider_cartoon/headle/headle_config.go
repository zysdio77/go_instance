package headle

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

type ConfigInfo struct {
	SrcDomain string `toml:srcdomain`
	SrcUrl string `toml:srcurl`
	TargetPath string `toml:targetpath`
	Redis RedisInfo `toml:redis`


}

type RedisInfo struct {
	Addr string `toml:addr`
	PassWord string `toml:password`
	Db int `toml:db`
}


func ReadConf(conffile string) (*ConfigInfo,error) {
	b,err := ioutil.ReadFile(conffile)
	if err != nil {
		return nil,err
	}

	var config ConfigInfo
	_,err = toml.Decode(string(b),&config)
	if err != nil {
		return nil,err
	}
	return &config,nil
}

