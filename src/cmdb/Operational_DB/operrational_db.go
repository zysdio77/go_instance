package Operational_DB

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)

func connct_db() (*sql.DB,error) {
	db,err := sql.Open("mysql","root:root@tcp(104.225.154.39:3306)/AC01?charset=utf8")
	defer db.Close()
	if err!= nil {
		return nil,fmt.Errorf("connnet DB err %v",err)
	}
	return db,nil
}

func add_slq(db *sql.DB) error{

	//插入数据
	stmt, err := db.Prepare("insert  aaa set user_id=?,user_name=?") //准备插入数据，
	defer stmt.Close()
	if err != nil {
		fmt.Println("perpare err:", err)
		return fmt.Errorf("add_sql perpare err %v",err)
	}
	res, err := stmt.Exec("12", "bob") //插入数据，args是要插入的值
	if err != nil {
		return fmt.Errorf("add_sql exec err %v",err)
	}
	id, err := res.LastInsertId() //最后插入的ID
	if err != nil {
		fmt.Errorf("add_sql lastinsertid err %v",err)
	}
	fmt.Println(id)
	return nil
}
func select_sql (db *sql.DB) error	{
	//查询数据
	rows,err := db.Query("select * from aaa;;")
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("select_sql query err %v",err)
	}
	for rows.Next(){
		var id int
		var name string
		err := rows.Scan(&id,&name)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("select_sql rows scan err %v",err)
		}
		fmt.Println(id,name)
	}

	return nil
}