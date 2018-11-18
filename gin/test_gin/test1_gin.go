package main

import (
	"github.com/gin-gonic/gin"
)

func main()  {

	router := gin.Default();
	router.GET("/read_cookie", func(context *gin.Context) {
		val, _ := context.Cookie("name")
		context.String(200, "Cookie:%s", val)
	})

	router.GET("/write_cookie", func(context *gin.Context) {
		context.SetCookie("name", "Shimin Li", 10, "/", "localhost", false, true)
	})

	router.GET("/clear_cookie", func(context *gin.Context) {
		context.SetCookie("name", "Shimin Li", -1, "/", "localhost", false, true)
	})

	router.Run(":8080")
}
