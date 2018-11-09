package check_mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"bytes"
	//"log"
	"go.pkg.wesai.com/p/base_lib/log"
)
func checherr(err error) {
	//logger := log.Logger{}
	if err != nil {
		//logger.Println(err)
		log.DLogger().Errorf("%v",err)
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
type CommandSqlV2 struct {
	Cmd string
	Value []string
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
func (cmondsql CommandSqlV2) Insert(db *sql.DB) error{
	//stmt,err :=db.Prepare("insert into user(name,age) values(?,?)")
	stmt,err :=db.Prepare(cmondsql.Cmd)
	defer stmt.Close()
	if err != nil {
		log.DLogger().Errorf("db prepare err : %v",err)
		return err
	}
	//result,err := stmt.Exec("zhang",21)
	_,err = stmt.Exec(cmondsql.Value[0],cmondsql.Value[1],cmondsql.Value[2])
	if err != nil {
		log.DLogger().Errorf("stmt exec err : %v",err)
		return err
	}
	//id,err := result.LastInsertId()
	//affect ,err := result.RowsAffected()
	//lipce :=[]int64{id,affect}

	return nil

}

func (cmondsql CommandSqlV2)Update(db *sql.DB) error {
	stmt,err := db.Prepare(cmondsql.Cmd)
	if err != nil {
		log.DLogger().Errorf("update db prepare err : %v",err)
		return err
	}
	_,err = stmt.Exec(cmondsql.Value[0],cmondsql.Value[1],cmondsql.Value[2])
	if err != nil{
		log.DLogger().Errorf("update stmt exec err : %v",err)
		return err
	}
	return nil
}
func Delete()  {

}