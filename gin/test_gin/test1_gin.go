package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"encoding/json"
	"fmt"
)

type Person struct {
	Username string
	Password string
}

func aa(c *gin.Context)  {
	//userneme :=c.Query("username")
	//c.String(200,userneme)
	//uu,err := c.GetCookie(userneme)
	//if err != nil {
	//	fmt.Println(err)
	//}
	c.Writer.WriteString("aaaaa")
	p := Person{"zhang","123"}
	//c.String(200,uu)
	c.JSON(200,&p)
	j,err :=json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(j))
}
func main ()  {
	router := gin.Default()
	router.GET("/a",aa)
	router.Run("0.0.0.0:9999")
}