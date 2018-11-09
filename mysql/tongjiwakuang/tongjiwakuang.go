package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"

)
func ConnDb() *sql.DB {
	//db,err := sql.Open("mysql","root:root@tcp(192.168.2.186:3306)/mining")
	db,err := sql.Open("mysql","game_admin:47qu3ZKh2k6trAlc@tcp(127.0.0.1:3306)/mining")
	if err != nil {
		fmt.Errorf("conn db err %v",err)
	}
	return db
}

type Cmd struct {
	cmd string
}

func (command Cmd)SelectUserId(db *sql.DB) []int {
	//command.cmd="select user_id from t_dig_round"
	row,err := db.Query(command.cmd)
	if err != nil {
		fmt.Errorf("select user id db query err : %v",err)
	}
	defer row.Close()
	var useridlist []int
	for row.Next() {
		var user_id int
		err := row.Scan(&user_id)
		if err != nil {
			fmt.Println("select user id",err)
		}
		useridlist = append(useridlist,user_id)
	}
	return useridlist
}

func (command Cmd)SelectYiLiu(db *sql.DB) int {
	row,err := db.Query(command.cmd)
	if err != nil {
		fmt.Errorf("select yiliu db query err : %v",err)
	}
	defer row.Close()
	//var useridlist []int
	var crystal int
	for row.Next() {
		//var crystal int
		err := row.Scan(&crystal)
		if err != nil {
			fmt.Println("select yiliu",err)
		}
		//useridlist = append(useridlist,crystal)
	}
	return crystal
}

func (command Cmd)Wachu(db *sql.DB) int {
	row,err := db.Query(command.cmd)
	if err != nil {
		fmt.Errorf("select Wachu db query err : %v",err)
	}
	defer row.Close()
	//var useridlist []int
	var wachu int
	for row.Next() {
		//var crystal int
		err := row.Scan(&wachu)
		if err != nil {
			fmt.Println("select Wachu",err)
		}
		//useridlist = append(useridlist,crystal)
	}
	return wachu
}
func (command Cmd)Duihuan(db *sql.DB) int {
	row,err := db.Query(command.cmd)
	if err != nil {
		fmt.Errorf("select duihuan db query err : %v",err)
	}
	defer row.Close()
	//var useridlist []int
	var duihuan int
	for row.Next() {
		//var crystal int
		err := row.Scan(&duihuan)
		if err != nil {
			fmt.Println("select duihuan ",err)
		}
		//useridlist = append(useridlist,crystal)
	}
	return duihuan
}

func main()  {
	db := ConnDb()
	defer db.Close()
	var command Cmd
	// select user_id  from t_account order by user_id，所有userID

	//select crystal from t_account where user_id = 54 遗留矿晶
	//select sum(reward_common+reward_super) from t_dig_round where user_id = 54 挖出的矿晶
	//select sum(crystal) from t_exchange where user_id = 109394 兑换的
	command.cmd = "select user_id from t_account group by user_id"
	useridlist := command.SelectUserId(db)

	fmt.Println(useridlist)
	var yiliulist []int
	for _,v := range useridlist{
		//fmt.Println(v)
		command.cmd = fmt.Sprintf("select sum(crystal) from t_account where user_id = %d",v)
		yiliu := command.SelectYiLiu(db)
		yiliulist = append(yiliulist,yiliu)
	}
	fmt.Println(yiliulist)

	var wachulist []int
	for _,v := range useridlist{
		//fmt.Println(v)
		command.cmd = fmt.Sprintf("select sum(reward_common+reward_super) as wachu from t_dig_round where user_id = %d",v)
		wachu := command.Wachu(db)
		wachulist = append(wachulist,wachu)
	}
	fmt.Println(wachulist)

	var duihuanlist []int
	for _,v := range useridlist{
		//fmt.Println(v)
		command.cmd = fmt.Sprintf("select sum(crystal) as duihuan from t_exchange where user_id = %d",v)
		duihuan := command.Duihuan(db)
		duihuanlist = append(duihuanlist,duihuan)
	}
	fmt.Println(duihuanlist)

	for i,v := range useridlist{
		fmt.Printf("user_id:%d , 遗留矿晶:%d , 挖出矿晶:%d , 兑换数量:%d\n",v,yiliulist[i],wachulist[i],duihuanlist[i])
	}

}