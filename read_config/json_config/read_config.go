package json_config

import (
	"io/ioutil"
	"go.pkg.wesai.com/p/base_lib/log"
	"encoding/json"
)

type Config struct {

}


func readfile(configfile string) ([]byte,error) {
	b,err :=ioutil.ReadFile(configfile)
	if err != nil {
		log.DLogger().Error(err)
		return nil,err
	}
	return b,nil
}
func headerjson(b []byte)  {
	json.Unmarshal(b,&j)
}

func main()  {
	
}
