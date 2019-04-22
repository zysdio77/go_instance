package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.pkg.wesai.com/p/base_lib/log"
	"go_instance/read_config/toml_config"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"strconv"
)

func ConnectDb(filename string) (*sql.DB, error) {
	config, err := toml_config.ParseConfig(filename)
	if err != nil {
		log.DLogger().Fatal(err)
		return nil, err
	}
	dd := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config.DBConfig.User, config.DBConfig.Pwd, config.DBConfig.Host, config.DBConfig.Port, config.DBConfig.DBName)
	//db,err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	//fmt.Println(dd)
	db, err := sql.Open("mysql", dd)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Info struct {
	Id       string
	Name     string
	Author   string
	Menupath string
	Menuname string
	Content  string
}

func SelectData(db *sql.DB, table_name, id string) Info {
	defer db.Close()
	ss := fmt.Sprintf("select * from %s where id = '%s'", table_name, id)
	row := db.QueryRow(ss)
	var info Info
	err := row.Scan(&info.Id, &info.Name, &info.Author, &info.Menupath, &info.Menuname, &info.Content)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	//fmt.Printf("%s",info)
	return info
}
func SelectMenuPath(db *sql.DB, table_name string) (string, []string, []string, map[string]string) {
	menudct := make(map[string]string)
	var menupathlist []string
	var menunamelist []string
	//var menupath string
	var id string
	var name string
	var menuname string
	defer db.Close()
	sqlcmd := fmt.Sprintf("select id,name,menuname from %s", table_name)
	row, err := db.Query(sqlcmd)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	for row.Next() {
		err := row.Scan(&id, &name, &menuname)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		//fmt.Println(menupath,menuname)
		//ss := strings.Split(menupath,".")
		menudct[id] = menuname
		menupathlist = append(menupathlist, id)
		menunamelist = append(menunamelist, menuname)
	}
	//fmt.Println(menudct,name)
	return name, menupathlist, menunamelist, menudct

}
func index(c *gin.Context) {
	db, err := ConnectDb(configFile)
	if err != nil {
		fmt.Println(err)
	}
	table_name := c.Query("table_name")
	name, _, _, menudct := SelectMenuPath(db, table_name)
	//fmt.Println(table_name)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":        name,
		"table_name":   table_name,
		//"menupathlist": menupathlist,
		//"menunamelist":menunamelist,
		"menudct": menudct,
	})

}
func content(c *gin.Context) {
	db, err := ConnectDb(configFile)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	var s string
	s = c.Query("menupath")
	table_name := c.Query("table_name")
	pagenumb, err := strconv.Atoi(s)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	beforpagenumb := pagenumb - 1
	afterpagenumb := pagenumb + 1
	fmt.Println()
	info := SelectData(db, table_name, s)
	c.HTML(http.StatusOK, "content.html", gin.H{
		"title":         info.Menuname,
		"table_name":    table_name,
		"content":       info.Content,
		"beforpagenumb": beforpagenumb,
		"afterpagenumb": afterpagenumb,
	})

}

var configFile string

func main() {
	//name_dic := make(map[string]string)
	//name_dic["story_mjgs"] = "苗疆蛊事"
	//name_dic["story_mshy"] = "茅山后裔"
	//fmt.Println(name_dic)

	flag.StringVar(&configFile, "config", `/Users/zhangyongsheng/data/src/go_instance/read_config/config.toml`, "config file path")
	flag.Parse()

	config, err := toml_config.ParseConfig(configFile)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	router := gin.Default()
	htmlfile := fmt.Sprintf("%s*", config.Templateconfig.Dir)

	router.LoadHTMLGlob(htmlfile)
	//router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/menu", index)
	router.GET("/content", content)
	router.Run(":10003")
}
