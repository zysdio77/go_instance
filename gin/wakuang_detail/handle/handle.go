package handle

import (
	"gopkg.in/gin-gonic/gin.v1"
	"fmt"
	"strconv"
	"go_test/gin/wakuang_detail/check_mysql"
	"log"
	"go_test/shell"
	"net/http"
)

type DetailInfo struct {
	PeopleNum string   `json:"peoplenum"`
	Bettype   []int    `json:"bettype"`
	Bet       []string `json:"bet"`
	Common    string   `json:"common"`
	Super     string   `json:"super"`
	Touzhu    int      `json:"touzhu"`
	Zhichu    int      `json:"zhichu"`
	Yingli    int      `json:"yingli"`
	Payout    float64  `json:"payout"`
}

func Test(c *gin.Context)  {
	var cc = check_mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "192.168.2.186:3306", DbName: "mining"}
	var i check_mysql.Mysql
	i = &cc
	db := i.ConnDb(i.InitUserInfo())
	defer db.Close()
	fmt.Println(db)
	//sql := "show databases"
	sql := fmt.Sprintf("select count(*) from (select distinct user_id from (select * from t_dig_round where start_time >'%s 00:00:00') t) s",)
	command := check_mysql.CommandSql{sql}
	var cm check_mysql.Command
	cm = &command
	result := cm.Show(db)
	c.String(200,result[0])
}

func TodayInfo(c *gin.Context) {
	//var cc = check_mysql.DbInfo{DbUserName: "game_admin", DbPassWord: "47qu3ZKh2k6trAlc", DbHost: "guess-instance.c87tof8gczbt.ap-northeast-1.rds.amazonaws.com:3306", DbName: "mining"}
	var cc = check_mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "192.168.2.186:3306", DbName: "mining"}
	var i check_mysql.Mysql
	i = &cc
	db := i.ConnDb(i.InitUserInfo())
	defer db.Close()

	//info := &DetailInfo{}
	info := new(DetailInfo)
	//一天多少人玩
	s := "date +%Y-%m-%d"
	date := shell.ExecShell(s)
	sql := fmt.Sprintf("select count(*) from (select distinct user_id from (select * from t_dig_round where start_time >'%s 00:00:00') t) s", date)
	//sql := fmt.Sprintf("select count(*) from (select distinct user_id from (select * from t_dig_round where start_time >'2018-08-10 13:40:59') t) s")
	command := check_mysql.CommandSql{sql}
	var cm check_mysql.Command
	cm = &command
	result := cm.Show(db)
	info.PeopleNum = result[0]

	//每个BET玩了多少次
	info.Bettype = []int{1, 5, 10, 50, 100, 500}
	for _, v := range info.Bettype {
		//sql = fmt.Sprintf("select count(*) from ((select * from t_dig_round where start_time > '2018-08-10 13:40:59') t) where bet = %v", v)
		sql = fmt.Sprintf("select count(*) from ((select * from t_dig_round where start_time > '%s 00:00:00') t) where bet = %v",date,v)
		command = check_mysql.CommandSql{sql}
		cm = &command
		result := cm.Show(db)
		info.Bet = append(info.Bet, result[0])
	}
	//总产出矿晶数量
	sql = fmt.Sprintf("select sum(reward_common) from (select * from t_dig_round where start_time > '%s 00:00:00') t",date)
	//sql = "select sum(reward_common) from (select * from t_dig_round where start_time > '2018-08-10 13:40:59') t"
	command = check_mysql.CommandSql{sql}
	cm = &command
	result = cm.Show(db)
	info.Common = result[0]
	com, err := strconv.Atoi(result[0])
	if err != nil {
		log.Println(err)
		com = 0
	}


	//超级矿数量
	sql = fmt.Sprintf("select sum(reward_super) from (select * from t_dig_round where start_time > '%s 00:00:00') t",date)
	//sql = "select sum(reward_super) from (select * from t_dig_round where start_time > '2018-08-10 13:40:59') t"
	command = check_mysql.CommandSql{sql}
	cm = &command
	result= cm.Show(db)
	info.Super = result[0]
	super, err := strconv.Atoi(result[0])
	if err != nil {
		log.Printf("Atoi err: %v\n", err)
		super = 0
	}

	//总投注数量
	sql = fmt.Sprintf("select sum(bet) from (select * from t_dig_round where start_time > '%s 00:00:00') t",date)
	//sql = "select sum(bet) from (select * from t_dig_round where start_time > '2018-08-10 13:40:59') t"
	command = check_mysql.CommandSql{sql}
	cm = &command
	result = cm.Show(db)
	tz, err := strconv.Atoi(result[0])
	if err != nil {
		log.Printf("Atoi err: %v\n", err)
		tz = 0
	}
	info.Touzhu = tz

	//总支出
	zc := int((com + super) / 100)
	info.Zhichu = zc

	//盈利
	yl := tz - zc
	info.Yingli = yl
	//payout(支出/投注)
	//payout := float64(zc) / float64(tz)
	var payout float64
	if tz == 0{
		payout = 0
	}else {
		payout = float64(zc) / float64(tz)
	}
	info.Payout = payout

	c.JSON(http.StatusOK, info)
	fmt.Println(info)
	//c.String(http.StatusOK,info.PeopleNum)
}

type Inserinfo struct {
	Date string
	PeopleNum int
	Bettype []int
	Bet []int
	Com int
	Super int
	Tz int
	Zc int
	Yl int
}

func InserDate(c *gin.Context)  {
	//var cc = check_mysql.DbInfo{DbUserName: "game_admin", DbPassWord: "47qu3ZKh2k6trAlc", DbHost: "guess-instance.c87tof8gczbt.ap-northeast-1.rds.amazonaws.com:3306", DbName: "mining"}
	var cc = check_mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "192.168.2.186:3306", DbName: "test111"}
	var i check_mysql.Mysql
	i = &cc
	db := i.ConnDb(i.InitUserInfo())
	defer db.Close()

	var cmd check_mysql.CommandSqlV2
	cmd.Cmd = "insert into AAA(id,name,age) values (?,?,?)"
	cmd.Value = []string{"1","zhang,","43"}

	result := cmd.Insert(db)
	c.JSON(http.StatusOK,result)

}

func UpdateDate(c *gin.Context)  {
	cc := check_mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "192.168.2.186:3306", DbName: "test111"}
	var i check_mysql.Mysql
	i = &cc
	db := i.ConnDb(i.InitUserInfo())
}