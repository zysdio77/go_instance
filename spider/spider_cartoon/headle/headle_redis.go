package headle


import (
	"github.com/go-redis/redis"
	"go_instance/gin/new_story_web/config"
)
//var cli *redis.Client
func Cli(conf *config.ConfigInfo) *redis.Client{	//链接redis
	cli := redis.NewClient(&redis.Options{
		Addr:conf.Redis.Addr,
		Password:conf.Redis.PassWord,
		DB:conf.Redis.Db,
	})
	return cli
}
func WriteReids(cli *redis.Client,k ,v string) error {	//写入数据
	err := cli.Set(k,v,0).Err()
	if err != nil {
		return err
	}
	return nil
}

func ReadRedis(cli *redis.Client,k string) (string,error) {//读取数据
	result ,err := cli.Get(k).Result()
	if err != nil {
		return "",err
	}
	return result,nil
}

func WriteMenuToRedis(cli *redis.Client,domain string,menu []string) error	 {//写入所有章节的标题到redis
	for _,v := range menu{
		k := domain+v
		err :=WriteReids(cli,v,k)
		if err != nil {
			return err
		}
	}
	//defer cli.Close()
	return nil
}

func HitMenu(cli *redis.Client,menu []string) ([]string,error) {	//检查章节是否在redis中已经有了
	var newmenu []string
	for _,v := range menu{
		_,err := ReadRedis(cli,v)
		if err != nil {
			if err.Error() == "redis: nil" {
				newmenu = append(newmenu,v)
			} else {
				return nil,err
			}
		}
	}
	//defer cli.Close()
	return newmenu,nil
}