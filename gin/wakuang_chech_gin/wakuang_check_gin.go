package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"go_test/mysql/wakuang_check"
	"net/http"
	"go_test/shell"
	"fmt"
	"strconv"
	"github.com/labstack/gommon/log"
	"go_test/mysql/caiquan_check"
)

func main()  {
	router :=gin.Default()
	router.GET("/checkwakuang", func(c *gin.Context) {
		//var cc = wakuang_check.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "192.168.2.186:3306", DbName: "mining"}
		var cc = wakuang_check.DbInfo{DbUserName: "game_admin", DbPassWord: "47qu3ZKh2k6trAlc", DbHost: "guess-instance.c87tof8gczbt.ap-northeast-1.rds.amazonaws.com:3306", DbName: "mining"}
		var i caiquan_check.Mysql
		i = &cc
		db := i.ConnDb(i.InitUserInfo())
		defer db.Close()

		//var command = mysql.CommandSql{"show tables"}
		//一天多少人玩
		s := "date +%Y-%m-%d"
		out := shell.ExecShell(s)
		sql := fmt.Sprintf("select count(*) from (select distinct user_id from (select * from t_dig_round where start_time >'%s 00:00:00') t) s",out)
		var command = caiquan_check.CommandSql{sql}
		var cm caiquan_check.Command
		cm = &command
		result :=cm.Show(db)
		c.String(http.StatusOK,"一天多少人玩：%v\n",result[0])

		//每个BET玩了多少次
		slice := [6]int{1,5,10,50,100,500}
		for _,v := range slice{
			sql = fmt.Sprintf("select count(*) from t_dig_round where bet = %v",v)
			command  = caiquan_check.CommandSql{sql}
			cm = &command
			result = cm.Show(db)
			c.String(http.StatusOK,"%v bet玩了多少次：%s\n",v,result[0])
		}

		//总产出矿晶数量
		sql = "select sum(reward_common) from t_dig_round"
		command  = caiquan_check.CommandSql{sql}
		cm = &command
		result = cm.Show(db)
		com,err := strconv.Atoi(result[0])
		if err != nil{
			log.Fatal(err)
		}
		c.String(http.StatusOK,"挖出的普通矿晶一共多少：%s\n",result[0])

		sql = "select sum(reward_super) from t_dig_round"
		command  = caiquan_check.CommandSql{sql}
		cm = &command
		result = cm.Show(db)
		super,err := strconv.Atoi(result[0])
		if err != nil{
			log.Fatal(err)
		}
		c.String(http.StatusOK,"挖出的超级矿晶一共多少：%s\n",result[0])

		//总投注数量
		sql = "select sum(bet) from t_dig_round"
		command  = caiquan_check.CommandSql{sql}
		cm = &command
		result = cm.Show(db)
		tz,err := strconv.Atoi(result[0])
		if err != nil{
			log.Fatal(err)
		}

		c.String(http.StatusOK,"总投注的GTC数量：%v\n",result[0])

		//总支出
		zc := int((com+super)/100)
		c.String(http.StatusOK,"总支出：%v\n",zc)

		//盈利
		yl :=tz-zc
		c.String(http.StatusOK,"盈利：%v\n",yl)
		//payout(支出/投注)
		payout := float64(zc)/float64(tz)
		c.String(http.StatusOK,"payout: %v\n",payout)

	})
	http.ListenAndServe(":12346", router)
}
