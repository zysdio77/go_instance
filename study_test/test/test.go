package main

import (
	"log"
	"gopkg.in/resty.v1"
	"github.com/PuerkitoBio/goquery"
	"bytes"

	"fmt"
	"os"
)

func main()  {
	//resp,err := http.NewRequest("GET","http://www.zeroz.com.cn/wordpress/",nil)
	var faceImg string
	resp,err := resty.R().
		SetHeader("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36").
		Get("http://www.zeroz.com.cn/wordpress/")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp.String())
	//r:= transform.NewReader(bytes.NewReader(resp.Body()),simplifiedchinese.GBK.NewDecoder()) //gbk转utf8
	r := bytes.NewReader(resp.Body())
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.String())
	d.Find("#media_image-3").Find("img").Each(func(i int, s *goquery.Selection) {
		imgpath, exists:= s.Attr("src")
		if !exists {
			return
		}

		if i == 0 {
			faceImg = imgpath
			fmt.Println(faceImg)
		}

	})
	//if ok{
	//}



	resp,err = resty.R().Get(faceImg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Body())
	f,err := os.Create("/Users/zhangyongsheng/Desktop/1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	n,err :=f.Write(resp.Body())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(resp.Body()))
	fmt.Println(n)
	//"User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
	//defer resp.Body.Close()
	//resp.Header.Add("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	//fmt.Println(resp.StatusCode())
	//fmt.Println(resp.String())


}