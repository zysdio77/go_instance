package caiquan_check

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func checherr(err error) {
	if err != nil {
		fmt.Println("Error is ", err)
		os.Exit(1)
	}
}

type Mysql interface {
	InitUserInfo() string
	ConnDb(string) *sql.DB
}
type Command interface {
	Show(db *sql.DB) []string
}

//定义dbinfo结构体存放登陆数据库是需要的信息
type DbInfo struct {
	DbUserName string
	DbPassWord string
	DbHost     string
	DbName     string
}

//定义commandsql结构体储存要执行的sql
type CommandSql struct {
	Cmd string
}

// 执行sql操作，打印执行的结果
func (cmondsql CommandSql) Show(db *sql.DB) []string{
	rows, err := db.Query(cmondsql.Cmd)
	defer rows.Close()
	checherr(err)
	var databases []string
	for rows.Next() {
		var database string
		err = rows.Scan(&database)
		checherr(err)
		databases = append(databases,database)
		//fmt.Println(database)

	}
	return databases
}

//初始化链接数据库是的字符串，sqlOpen(的第二个参数)
// sql.Open("mysql", "用户名:密码@tcp(IP:端口)/数据库?charset=utf8")
func (info DbInfo) InitUserInfo() string {
	var buffer bytes.Buffer
	buffer.WriteString(info.DbUserName)
	buffer.WriteString(":")
	buffer.WriteString(info.DbPassWord)
	buffer.WriteString("@(")
	buffer.WriteString(info.DbHost)
	buffer.WriteString(")/")
	buffer.WriteString(info.DbName)
	buffer.WriteString("?charset=utf8")
	return buffer.String()
}

func (info DbInfo) ConnDb(infostring string) *sql.DB {
	db, err := sql.Open("mysql", infostring)
	checherr(err)
	return db
}
