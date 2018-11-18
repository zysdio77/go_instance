package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go.pkg.wesai.com/p/base_lib/log"
	"database/sql"
	"fmt"
)

var (db *sql.DB
	err error
	)
func query_sql() {
	rows, err := db.Query("show databases")
	defer rows.Close()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	var databases string
	for rows.Next() {
		err := rows.Scan(&databases)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		fmt.Println(databases)
	}
}

func main() {
	connectinfo := "root:123456@tcp(192.168.1.104:3306)/mysql?charset=utf8"
	db,err = sql.Open("mysql", connectinfo)
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()
	//defer db.Close()
	db.SetMaxOpenConns(140)	//最大连接数
	db.SetMaxIdleConns(30)	//最大空闲连接数
	//db.Ping()

	query_sql()
}
