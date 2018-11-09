package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"net/http"
	"flag"
	"go_instance/read_config"
	"go.pkg.wesai.com/p/base_lib/log"
	"strconv"
)




func ConnectDb(filename string) (*sql.DB,error) {
	config,err :=read_config.ParseConfig(filename)
	if err != nil {
		log.DLogger().Fatal(err)
		return nil,err
	}
	dd := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",config.DBConfig.User,config.DBConfig.Pwd,config.DBConfig.Host,config.DBConfig.Port,config.DBConfig.DBName)
	//db,err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	//fmt.Println(dd)
	db,err := sql.Open("mysql",dd)
	if err != nil {
		return nil,err
	}
	return db,nil
}
type Info struct {
	Id string
	Name string
	Author string
	Menupath string
	Menuname string
	Content string
}
func SelectData(db *sql.DB,id string) Info {
	defer db.Close()
	ss :=fmt.Sprintf("select * from story_01 where id = '%s'",id)
	row := db.QueryRow(ss)
	var info Info
	err := row.Scan(&info.Id,&info.Name,&info.Author,&info.Menupath,&info.Menuname,&info.Content)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	//fmt.Printf("%s",info)
	return info
}
func SelectMenuPath(db *sql.DB) ([]string,[]string,map[string]string) {
	menudct := make(map[string]string)
	var menupathlist []string
	var menunamelist []string
	//var menupath string
	var id string
	var menuname string
	defer db.Close()
	row,err	 := db.Query("select id,menuname from story_01")
	if err != nil{
		log.DLogger().Fatal(err)
	}
	for row.Next(){
		err :=row.Scan(&id,&menuname)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		//fmt.Println(menupath,menuname)
		//ss := strings.Split(menupath,".")
		menudct[id]=menuname
		menupathlist = append(menupathlist,id)
		menunamelist = append(menunamelist,menuname)
	}
	//fmt.Println(menudct)
	return menupathlist,menunamelist,menudct

}
func index(c *gin.Context)  {
	db,err := ConnectDb(configFile)
	if err != nil {
		fmt.Println(err)
	}
	_,_,menudct := SelectMenuPath(db)
	//fmt.Println(menudict)
	c.HTML(http.StatusOK,"index.html",gin.H{
		"title":"茅山后裔",
		//"menupathlist":menupathlist,
		//"menunamelist":menunamelist,
		"menudct":menudct,
	})
}
func content(c *gin.Context)  {
	db,err := ConnectDb(configFile)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	var s string
	s = c.Query("menupath")
	pagenumb,err:= strconv.Atoi(s)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	beforpagenumb:=pagenumb-1
	afterpagenumb := pagenumb+1
	fmt.Println()
	info :=SelectData(db,s)
	c.HTML(http.StatusOK,"content.html",gin.H{
		"title":info.Menuname,
		"content":info.Content,
		"beforpagenumb":beforpagenumb,
		"afterpagenumb":afterpagenumb,
	})

}
var configFile string
func main()  {

	//db,err := ConnectDb()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//SelectData(db)
	flag.StringVar(&configFile, "config", `/Users/zhangyongsheng/data/src/go_instance/read_config/config.toml`, "config file path")
	flag.Parse()

	config,err :=read_config.ParseConfig(configFile)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	router := gin.Default()
	htmlfile := fmt.Sprintf("%s*",config.Templateconfig.Dir)

	router.LoadHTMLGlob(htmlfile)
	//router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/menu",index)
	router.GET("/content",content)
	router.Run(":10003")
}
