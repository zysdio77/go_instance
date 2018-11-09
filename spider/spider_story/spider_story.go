package main

import (
	"gopkg.in/resty.v1"
	"golang.org/x/text/transform"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"github.com/labstack/gommon/log"
	"github.com/PuerkitoBio/goquery"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

type Introduce struct {
	Author string
	Name string
	MenuPath	[]string
	MenuName []string
	Content string
}

func FirstPage(url_addr string) ([]byte,error){
	resp,err := resty.R().Get(url_addr)
	if err != nil {
		return nil, err	
	}
	//readre := transform.NewReader(bytes.NewReader(resp.Body()),simplifiedchinese.GBK.NewDecoder()) //gbk转utf8
	//fmt.Println(readre)
	return resp.Body(),nil
}

func GbkToUtf8(r []byte) io.Reader  {
	//r是传过来的resp.body()内容
	readre := transform.NewReader(bytes.NewReader(r),simplifiedchinese.GBK.NewDecoder()) //gbk转utf8
	return readre
}

func (introduce *Introduce)GetIntroduce(utf8_html io.Reader) (*Introduce,error) {
	//goquery，从reader读出文档
	doc,err := goquery.NewDocumentFromReader(utf8_html)
	if err != nil {
		return introduce,err
	}
	//过滤标签内容，查找出<div class="jieshao">下的<div class="rt">下的<h1>下的文本内容
	introduce.Name = doc.Find("div.jieshao").Find("div.rt").Find("h1").Text()
	//fmt.Println(introduce.Name)
	//找出<div class="jieshao">下的<div class="rt">下的<div class=msg"">下的<em>下的第一条(如果有多个内容只要第一个，索引从0开始)内容
	introduce.Author = doc.Find("div.jieshao").Find("div.rt").Find("div.msg").Find("em").Eq(0).Text()
	//fmt.Println(introduce.Author)
	//循环遍历<div class="mulu">下的所有<li>标签
	doc.Find("div.mulu").Find("li").Each(func(i int, s *goquery.Selection) {
		//每条<li>下的<a>里href的属性 <a href = "属性">
		ss,ok := s.Find("a").Attr("href")
		//每一条<a>标签的文本内容 <a>文本内容</a>
		neme :=s.Find("a").Text()
		if ok{
			introduce.MenuPath = append(introduce.MenuPath,ss)
		}
		introduce.MenuName = append(introduce.MenuName,neme)
	})
	//fmt.Println(introduce.Menu[1])
	return introduce,nil
}
func (introduce *Introduce)GetContent(host string) error {
		resp,err := resty.R().Get(host)
		r := GbkToUtf8(resp.Body())
		if err != nil {
			return err
		}
		doc,err := goquery.NewDocumentFromReader(r)
		if err != nil{
			return err
		}
		//s:= doc.Find("div.novel").Find("div.yd_text2").Text()
		//找出<div class = "novel">下<div class="yd_text2">下的文本内容
		introduce.Content= doc.Find("div.novel").Find("div.yd_text2").Text()
		//fmt.Println("getcontent",introduce.Content)
	return nil
}
func WriteFile(words,filename string) error {
	//f,err :=os.Create(filename)
	//打开文件，读写，如果没有则创建，追加的方式，权限644
	f,err := os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	defer f.Close()
	if err != nil{
		return err
	}
	//写入文字
	_,err = f.WriteString(words)
	if err != nil{
		return err
	}

	return nil

}
func ConnectDb() (*sql.DB,error) {
	//链接数据库
	db,err := sql.Open("mysql","root:root@tcp(104.225.154.39:3306)/test?charset=utf8")
	if err != nil {
		return nil,err
	}
	return db,nil
}
func (introduce *Introduce )WriteSql(db *sql.DB,k int)  {
	//for k,_ := range introduce.MenuPath{
		//fmt.Println(introduce.Name,introduce.Author,introduce.MenuPath[k],introduce.MenuName[k],introduce.Content)
		//db.Exec("insert into story (name,author,menupath,menuname,content) values (?,?,?,?,?)",introduce.Name,introduce.Author,introduce.MenuPath[k],introduce.MenuName[k],introduce.Content)

	//准备执行sql
	s,err :=db.Prepare("insert into story_01 (name,author,menupath,menuname,content) values (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	//执行sql，把valuse传进去
	r,err :=s.Exec(introduce.Name,introduce.Author,introduce.MenuPath[k],introduce.MenuName[k],introduce.Content)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(r.RowsAffected())
	//}
}

func main()  {
	var introduce Introduce
	fmt.Println(time.Now())
	url_addr := "https://www.88dush.com/xiaoshuo/0/801/"
	r,err := FirstPage(url_addr)
	if err != nil {
		log.Fatal(err)
	}
	resp := GbkToUtf8(r)
	//fmt.Println(string(b))
	info ,err := introduce.GetIntroduce(resp)
	if err != nil {
		log.Fatal(err)
	}
	//introduce.WriteFile("/Users/zhangyongsheng/Desktop/q.txt")
	//introduce.WriteMysql()
	//localpath := "/Users/zhangyongsheng/Desktop/"+introduce.Name+".txt"
	//localpath := introduce.Name+".txt"
	db,err :=ConnectDb()
	defer db.Close()
	for k,v := range introduce.MenuPath{
		host := url_addr+v
		//host := "https://www.88dush.com/xiaoshuo/0/801/247999.html"
		err = introduce.GetContent(host)
		if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(info.Content)
		info.WriteSql(db,k)
		//WriteFile(introduce.MenuName[k],localpath)
		//WriteFile(s,localpath)
	}
	fmt.Println(time.Now())
}
