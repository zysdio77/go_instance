package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"go.pkg.wesai.com/p/base_lib/log"
	"fmt"
)
func (s Sqlcmd)Select(db *sql.DB) {
	row, err := db.Query(s.Command) //"select content from test.story_01 limit 10"
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer row.Close()
	var count string
	for row.Next() {
		err = row.Scan(&count)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		fmt.Println(count)
	}
}

type Sqlcmd struct {
	Command string
}
type SqlcmdInter interface {
	Select(db *sql.DB)
	Insert(db *sql.DB)
}

func (s Sqlcmd)Insert(db *sql.DB)  {
	stmt,err := db.Prepare(s.Command)
	defer stmt.Close()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	result,err :=stmt.Exec()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(result.RowsAffected())
}
func main()  {
	connectinfo :="root:123456@tcp(192.168.1.30:3306)/mysql?charset=utf8"
	db,err := sql.Open("mysql",connectinfo)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer db.Close()

	var s Sqlcmd

	s.Command="show databases"

	var ss SqlcmdInter
	ss = &s
	ss.Insert(db)
	ss.Select(db)


}