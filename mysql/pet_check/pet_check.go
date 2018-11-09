package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)


func ConnectDb(connectinfo string)  *sql.DB {
	db,err := sql.Open("mysql",connectinfo)
	if err != nil {
		fmt.Println(err)
	}
	return db

}

func ExecSql(db *sql.DB,s string) string  {
	row := db.QueryRow(s)
	var num string
	err := row.Scan(&num)
	if err != nil {
		fmt.Println("Execsql err:",err)
	}
	return num
}

func CheckDb(c *gin.Context)  {
	connectinfo := "game_admin:47qu3ZKh2k6trAlc@tcp(guess-instance.c87tof8gczbt.ap-northeast-1.rds.amazonaws.com:3306)/AC20"
	db := ConnectDb(connectinfo)
	defer db.Close()
	gtc_total_bet := "select sum(total_bet) from t_game_round where currency = 'gtc0'"
	gtc_total_bet_result := ExecSql(db,gtc_total_bet)
	gtc_total_reward := "select sum(total_reward) from t_game_round where currency = 'gtc0'"
	gtc_total_reward_result := ExecSql(db,gtc_total_reward)
	trx_total_bet  := "select sum(total_bet) from t_game_round where currency = 'trx'"
	trx_total_bet_result := ExecSql(db,trx_total_bet)
	trx_total_reward := "select sum(total_reward) from t_game_round where currency = 'trx'"
	trx_total_reward_result := ExecSql(db,trx_total_reward)

	result := fmt.Sprintf(" gtc总押注：%v\n gtc总产出：%v\n trx总押注：%v\n trx总产出：%v",gtc_total_bet_result,gtc_total_reward_result,trx_total_bet_result,trx_total_reward_result)
	c.String(http.StatusOK,result)
}

func main()  {
	router := gin.Default()
	router.GET("/check_pet",CheckDb)
	router.Run(":12347")
}