package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"fmt"
	"net/http"
)
//返回sting
func response_string(c *gin.Context)  {
	params := c.Query("params")
	fmt.Println(params)
	c.String(http.StatusOK,params)
}
//返回json
func response_json(c *gin.Context)  {
	params := c.Query("params")
	j:= map[string]string{"params":params}
	c.JSON(http.StatusOK,j)
}
//返回html
func response_html(c *gin.Context)  {
	//第二个参数是router.LoadHTMLFiles(htmlfile)加载的文件名字
	//第三个参数是需要渲染的键值对，这里没有渲染所以是空，response_html2位渲染的例子
	c.HTML(http.StatusOK,"test.html",gin.H{})
}
//返回html，有渲染
func response_html2(c *gin.Context)  {
	//第二个参数是router.LoadHTMLFiles(htmlfile)加载的文件名字
	//第三个参数是需要渲染的键值对
	value := "valuse"
	c.HTML(http.StatusOK,"test.html",gin.H{
		"key":value,
	})
}

func main()  {
	htmlfile := "/Users/zhangyongsheng/data/src/go_instance/templates_file/templates_html/test.html"
	router := gin.Default()

	router.LoadHTMLFiles(htmlfile)
	router.GET("/string",response_string)
	router.GET("/json",response_json)
	router.GET("/html",response_html)
	router.GET("/html2",response_html2)
	router.Run(":10005")
}