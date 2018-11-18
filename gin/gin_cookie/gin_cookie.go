package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strings"
)

func ReadCookie(c *gin.Context)  {
	str,err := c.Cookie("name")
	if err != nil {
		fmt.Println(err)
	}
	slist := strings.Split(str,",")
	fmt.Println(slist)
	c.String(http.StatusOK,"Cookie:%s",str)
}
func WriteCookie(c *gin.Context)  {
	c.SetCookie("name","zhang,32",0,"/","localhost",false,true)
}

func CleanCookie(c *gin.Context)  {
	c.SetCookie("name", "Shimin Li", -1, "/", "localhost", false, true)
}


func main()  {

	router := gin.Default();
	router.GET("/read_cookie", ReadCookie)

	router.GET("/write_cookie", WriteCookie)

	router.GET("/clear_cookie", CleanCookie)

	router.Run(":8080")
}

//http://localhost:8080/read_cookie
//
//http://localhost:8080/write_cookie
//
//http://localhost:8080/clear_cookie