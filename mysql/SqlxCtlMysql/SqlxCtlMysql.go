package SqlxCtlMysql
import (
_ "github.com/go-sql-driver/mysql"

"fmt"
"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var err error

type UserInfo struct {
	Id uint `db:"id"`
	Name string `db:"name"`
	Gender string `db:"gender"`
	Hobby string `db:"hobby"`

}
func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test1?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//fmt.Printf("connect DB failed, err:%v\n", err)
		return err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return nil
}
// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select * from user_info2"
	var us []UserInfo
	//查询多行数据用select
	err := db.Select(&us, sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", us)
	for _,j := range us {
		fmt.Printf("id:%v,name:%v,gender:%v,hobby:%v\n",j.Id,j.Name,j.Gender,j.Hobby)

	}

}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select * from user_info2 where id=?"
	var u UserInfo
	//查询一行数据用get
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Println(u.Id,u.Name,u.Gender,u.Hobby)
	//fmt.Printf("id:%d name:%s age:%d\n", , u.Name, u.Age)
}
//插入数据
func insertRow()  {
	sqlStr := "insert into user_info2(name,gender,hobby) values (?,?,?)"
	result,err := db.Exec(sqlStr," 老张","男","旅游")
	if err != nil {
		fmt.Println(err)
		return
	}
	insertId ,err := result.LastInsertId()//自增ID的id号
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("insert data success,id:%v\n",insertId)
}

//更改数据
func updateRow(){
	sqlStr := "update user_info2 set name = ? where id = ?"
	result,err  := db.Exec(sqlStr,"李鹏",3)
	if err != nil {
		fmt.Println(err)
		return
	}
	affecteRows,err := result.RowsAffected()//影响的行数
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("update data success,affected rows:%d\n",affecteRows)
}

//删除一行
func deleteRow()  {
	sqlStr := "delete from user_info2 where id =?"
	result,err := db.Exec(sqlStr,14)
	if err != nil {
		fmt.Println(err)
		return
	}
	affectdRows,err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete data success,affected rows:%d\n",affectdRows)
}

func main() {
	err = initDB()
	if err != nil {
		fmt.Println(err)
	}
	queryRowDemo()
	queryMultiRowDemo()
	insertRow()
	updateRow()
	deleteRow()
}