package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/resty.v1"
	"sync"
	"strconv"
	"io/ioutil"
)

//爬去的第一个页面
func FirstPage(firstpage string) ([]byte, error) {
	resp, err := resty.R().Get(firstpage)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

//把需要的元素找出来
func clean_flag(htmlbyte []byte) ([]string, error) {
	newreader := bytes.NewReader(htmlbyte)
	doc, err := goquery.NewDocumentFromReader(newreader)
	if err != nil {
		return nil, err
	}
	var imglist []string
	//<ul class="book-list">每一个<li>标签
	doc.Find("ul.book-list").Find("li").Each(func(i int, s *goquery.Selection) {
		//每个<li>标签下的<img src="需要的属性">
		src, ok := s.Find("img").Attr("src")
		if ok {
			imglist = append(imglist, src)
		}
		//fmt.Println(src)
	})
	return imglist, nil
}

//携程并发爬去数据，并写入文件
func main() {
	var wg sync.WaitGroup
	firstpage := "https://www.golang123.com/book"
	b, err := FirstPage(firstpage)
	if err != nil {
		fmt.Println(err)
	}
	imglist, err := clean_flag(b)
	if err != nil {
		fmt.Println(err)
	}

	//n是一共要执行的次数，在携程中会自增
	n := 0
	//每个携程执行的任务数量
	nz := len(imglist) / 10
	//需要单独执行的剩余的任务
	ny := len(imglist) % 10
	//启动10个携程
	for i := 0; i < 10; i++ {
		//每启动一个携程，wg+1
		wg.Add(1)
		go func() {
			//每个携程退出后wg-1
			defer wg.Done()
			//每个携程执行nz个任务
			for j:=0;j<nz;j++  {
				resp,err:= resty.R().Get(imglist[n])
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(n,imglist[n])
				filename := "/Users/zhangyongsheng/Desktop/"+strconv.Itoa(n)+".jpg"
				ioutil.WriteFile(filename,resp.Body(),0755)
				if err != nil {
					fmt.Println(err)
				}
				n = n+1
			}
		}()

	}
	//剩余任务在这里执行，总任务-剩余是起始点，
	for k:=len(imglist)-ny;k<len(imglist);k++{
		//fmt.Println(k)
		resp,err:= resty.R().Get(imglist[k])
		if err != nil {
			fmt.Println(err)
		}
		filename := "/Users/zhangyongsheng/Desktop/"+strconv.Itoa(k)+".jpg"
		ioutil.WriteFile(filename,resp.Body(),0755)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(k,imglist[k])
	}

	//阻塞在这里，知道wg=0时退出
	wg.Wait()
}
