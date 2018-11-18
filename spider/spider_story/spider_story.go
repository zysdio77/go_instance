package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"go.pkg.wesai.com/p/base_lib/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/resty.v1"
	"io"
	"os"
	"time"
)

type Introduce struct {
	Author   string
	Name     string
	MenuPath []string
	MenuName []string
	Content  string
}

func FirstPage(url_addr string) ([]byte, error) {
	resp, err := resty.R().Get(url_addr)
	if err != nil {
		return nil, err
	}
	//readre := transform.NewReader(bytes.NewReader(resp.Body()),simplifiedchinese.GBK.NewDecoder()) //gbk转utf8
	//fmt.Println(readre)
	return resp.Body(), nil
}

func GbkToUtf8(r []byte) io.Reader {
	//r是传过来的resp.body()内容
	readre := transform.NewReader(bytes.NewReader(r), simplifiedchinese.GBK.NewDecoder()) //gbk转utf8
	return readre
}

func (introduce *Introduce) GetIntroduce(utf8_html io.Reader) (*Introduce, error) {
	//goquery，从reader读出文档
	doc, err := goquery.NewDocumentFromReader(utf8_html)
	if err != nil {
		return introduce, err
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
		ss, ok := s.Find("a").Attr("href")
		//每一条<a>标签的文本内容 <a>文本内容</a>
		neme := s.Find("a").Text()
		if ok {
			introduce.MenuPath = append(introduce.MenuPath, ss)
		}
		introduce.MenuName = append(introduce.MenuName, neme)
	})
	//fmt.Println(introduce.Menu[1])
	return introduce, nil
}
func (introduce *Introduce) GetContent(host string) error {
	resp, err := resty.R().Get(host)
	r := GbkToUtf8(resp.Body())
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	//s:= doc.Find("div.novel").Find("div.yd_text2").Text()
	//找出<div class = "novel">下<div class="yd_text2">下的文本内容
	introduce.Content = doc.Find("div.novel").Find("div.yd_text2").Text()
	//fmt.Println("getcontent",introduce.Content)
	return nil
}
func WriteFile(words, filename string) error {
	//f,err :=os.Create(filename)
	//打开文件，读写，如果没有则创建，追加的方式，权限644
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	//写入文字
	_, err = f.WriteString(words)
	if err != nil {
		return err
	}

	return nil

}
func ConnectDb() (*sql.DB, error) {
	//链接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (introduce *Introduce) WriteSql(db *sql.DB, k int,table_name string) {
	sqlcmd := fmt.Sprintf("insert into %s (name,author,menupath,menuname,content) values (?,?,?,?,?)",table_name)
	//准备执行sql
	s, err := db.Prepare(sqlcmd)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer s.Close()
	//执行sql，把valuse传进去
	r, err := s.Exec(introduce.Name, introduce.Author, introduce.MenuPath[k], introduce.MenuName[k], introduce.Content)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(r.RowsAffected())
	//}
}
func drop_table(table_name string)  {
	db,err := ConnectDb()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer db.Close()
	sqlcmd := fmt.Sprintf("DROP TABLE IF EXISTS `%s`",table_name)
	stmt,err := db.Prepare(sqlcmd)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer stmt.Close()
	r, err :=stmt.Exec()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(r.RowsAffected())
}

func create_table(table_name string) {
	db, err := ConnectDb()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer db.Close()
	sqlcmd := fmt.Sprintf("create table %s (id int auto_increment primary key,name varchar(100),author varchar(100),menupath varchar(100),menuname varchar(100),content longtext)", table_name)
	stmt, err := db.Prepare(sqlcmd)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer stmt.Close()
	r, err := stmt.Exec()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(r.RowsAffected())
}

func alter_table(table_name string) {
	db, err := ConnectDb()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer db.Close()
	sqlcmd := fmt.Sprintf("alter table %s AUTO_INCREMENT=%d", table_name, 1000000)
	stmt, err := db.Prepare(sqlcmd)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer stmt.Close()
	r, err := stmt.Exec()
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(r.RowsAffected())
}

func do_work(url_addr ,table_name string) {
	start_time := time.Now().Unix()
	var introduce Introduce
	r, err := FirstPage(url_addr)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	resp := GbkToUtf8(r)
	//fmt.Println(string(b))
	info, err := introduce.GetIntroduce(resp)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(info)

	//introduce.WriteFile("/Users/zhangyongsheng/Desktop/q.txt")
	//introduce.WriteMysql()
	//localpath := "/Users/zhangyongsheng/Desktop/"+introduce.Name+".txt"
	//localpath := introduce.Name+".txt"

	db, err := ConnectDb()
	defer db.Close()
	for k, v := range introduce.MenuPath {
		host := url_addr + v
		err = introduce.GetContent(host)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		if err != nil {
			log.DLogger().Fatal(err)
		}
		//fmt.Println(info.Content)
		info.WriteSql(db, k,table_name)
		//WriteFile(introduce.MenuName[k],localpath)
		//WriteFile(s,localpath)
	}
	end_time := time.Now().Unix()
	used_second := end_time - start_time
	display_minit := used_second / 60
	display_second := used_second % 60
	fmt.Printf("用时%d分钟：%d秒", display_minit, display_second)

}
func test(url_addr string) {
	start_time := time.Now().Unix()
	var introduce Introduce
	r, err := FirstPage(url_addr)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	resp := GbkToUtf8(r)
	//fmt.Println(string(b))
	info, err := introduce.GetIntroduce(resp)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	fmt.Println(info)

	//introduce.WriteFile("/Users/zhangyongsheng/Desktop/q.txt")
	//introduce.WriteMysql()
	//localpath := "/Users/zhangyongsheng/Desktop/"+introduce.Name+".txt"
	//localpath := introduce.Name+".txt"

	for k, v := range introduce.MenuPath {
		host := url_addr + v
		//host := "https://www.88dush.com/xiaoshuo/0/801/247999.html"
		err = introduce.GetContent(host)
		if err != nil {
			log.DLogger().Fatal(err)
		}
		if err != nil {
			log.DLogger().Fatal(err)
		}
		//fmt.Println(info.Content)
			fmt.Println(introduce.MenuName[k],introduce.Content)
		//WriteFile(introduce.MenuName[k],localpath)
		//WriteFile(s,localpath)
	}
	end_time := time.Now().Unix()
	used_second := end_time - start_time
	display_minit := used_second / 60
	display_second := used_second % 60
	fmt.Printf("用时%d分钟：%d秒", display_minit, display_second)

}

func main() {
	tablename := "story_jyfhwz"
	drop_table(tablename)
	create_table(tablename)
	alter_table(tablename)
	//url_addr := "https://www.88dush.com/xiaoshuo/0/801/" //茅山后裔
	//url_addr := "https://www.88dush.com/xiaoshuo/70/70239/"	//苗疆蛊事
	url_addr := "https://www.88dush.com/xiaoshuo/3/3896/"	//捉蛊记
	do_work(url_addr,tablename)
	//test(url_addr)
}
