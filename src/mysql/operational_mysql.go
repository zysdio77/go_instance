package mysql

import (
"database/sql"
"fmt"
_"github.com/go-sql-driver/mysql"
)

func main() {
	//db, err := sql.Open("mysql", "root:root@tcp(104.225.154.39:3306)/AC01?charset=utf8")
	db, err := sql.Open("mysql", "BJPemberton_pro:bjCocaCola_PR@@tcp(10.11.12.134:3306)/actest?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Println("open err:",err)
	}

	//插入数据
	stmt,err := db.Prepare("insert  aaa set user_id=?,user_name=?")
	if err != nil {
		fmt.Println("perpare err:",err)
	}
	res,err := stmt.Exec("12","bob")
	if err != nil {
		fmt.Println("exrc err:",err)
	}
	id ,err :=res.LastInsertId()
	if err !=nil {
		fmt.Println("lastinsertid err:",err)
	}
	fmt.Println(id)
	//查询数据
	rows,err := db.Query("select * from aaa;;")
	defer rows.Close()
	if err != nil {
		fmt.Println("query err:",err)
	}
	for rows.Next(){
		var id int
		var name string
		err := rows.Scan(&id,&name)
		if err != nil {
			fmt.Println(err)

		}
		fmt.Println(id,name)
	}


}