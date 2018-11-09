package main

import (
	"go_test/gin/wakuang_detail/handle"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
)

func main()  {
	router :=gin.Default()
	router.GET("/checkwakuang",handle.TodayInfo)
	//router.GET("/checkwakuang/allinfo",handle.AllInfo)
	router.GET("/insert",handle.InserDate)
	router.GET("/update",handle.InserDate)
	router.GET("/delete",handle.InserDate)
	http.ListenAndServe(":12346", router)


}
