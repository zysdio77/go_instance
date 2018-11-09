package caiquan_chech_gin

import (
	"go_test/mysql/caiquan_check"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
)

func main()  {
	router	 := gin.Default()
	/*var c = mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "104.225.154.39:3306", DbName: "mysql"}
	var i mysql.Mysql
	i = &c

	db := i.ConnDb(i.InitUserInfo())
	defer db.Close()

	var command = mysql.CommandSql{"show tables"}
	var cm mysql.Command
	cm = &command
	result :=cm.Show(db)

	fmt.Println(result)*/
	router.GET("/caiquan", func(c *gin.Context) {
		//var cc = mysql.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "104.225.154.39:3306", DbName: "mysql"}
		var cc = caiquan_check.DbInfo{DbUserName: "root", DbPassWord: "root", DbHost: "104.225.154.39:3306", DbName: "mining"}
		//var i mysql.Mysql
		var i caiquan_check.Mysql
		i = &cc

		db := i.ConnDb(i.InitUserInfo())
		defer db.Close()

		//var command = mysql.CommandSql{"show tables"}
		var command = caiquan_check.CommandSql{"show tables"}
		var cm caiquan_check.Command
		cm = &command
		result :=cm.Show(db)
		c.String(http.StatusOK,"50匹配%s\n",result)

		command  = caiquan_check.CommandSql{"select count(*) as value from user"}
		cm = &command
		result = cm.Show(db)
		c.String(http.StatusOK,"一共%s\n",result)
	})
	http.ListenAndServe(":18080", router)
}