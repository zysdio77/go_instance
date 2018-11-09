package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"fmt"
)

//post命令行的请求
func psot_query(c *gin.Context)  {
	//c.Query("key1")相当于url上直接带参数
	//例如：http://127.0.0.1:10004/string?key1=123&key2=234
	key1 := c.Query("key1")
	key2 := c.Query("key2")
	fmt.Println(key1,key2)
	c.String(http.StatusOK,key1)
}

//post form表单的请求
func post_form(c *gin.Context)  {
	key1 := c.PostForm("key1")
	key2 := c.PostForm("key2")
	fmt.Println(key1,key2)
	c.String(http.StatusOK,key1+key2)
	
}

//post json的请求
type Aaaa struct {
	Key1 string `json:"key1"`
	Key2 string	`json:"key2"`
}
func post_json(c *gin.Context)  {
	var key Aaaa
	//绑定json数据
	//相当于json.unmarshal的解析，从json解析到struct
	err :=c.BindJSON(&key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(key)
	c.String(http.StatusOK,key.Key1)
}

func main()  {
	router := gin.Default()
	router.POST("/string",post_json)
	router.Run(":10004")
}