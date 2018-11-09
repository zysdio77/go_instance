package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/gommon/log"
)

func ConnDb() *sql.DB {
	db,err := sql.Open("mysql","root:root@tcp(104.225.154.39:3306)/actest")
	if err != nil {
		log.Errorf("ConnDb err : %v",err)
	}
	return db
}

func CommandSql(db *sql.DB,commandsql string) {
	db.Exec(commandsql)
}

func main()  {
	db := ConnDb()
	commandsql := "show databases"
	CommandSql(db,commandsql)
	defer db.Close()
}