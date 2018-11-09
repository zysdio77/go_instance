package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/resty.v1"
	"log"
	"os"
)

func GetPicturePath(b []byte) ([]string, error) {
	var dates []string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b)) //创建doc
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	doc.Find("img").Each(func(i int, s *goquery.Selection) { //img标签下，所有src的属性
		data, ok := s.Attr("src")
		if ok {
			dates = append(dates, data) //吧地址添加到切片
		}
	})
	return dates, nil
}

func GetHtml() []byte {
	resp, err := resty.R().Get("http://www.zeroz.com.cn/wordpress/")
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body()
}
func ReadFile(srcfile string) ([]byte, error) {
	resp, err := resty.R().Get(srcfile)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func WriteFile(disfile string, data []byte) error {
	file, err := os.Create(disfile)
	if err != nil {

	}
	defer file.Close()
	n, err := file.Write(data)
	if err != nil {
		return err
	}
	if len(data) != n {
		return err
	}
	return nil
}
func main() {
	r := GetHtml()
	list, err := GetPicturePath(r)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range list {
		fmt.Println(k, v)
		data, err := ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		disfile := fmt.Sprintf("/Users/zhangyongsheng/Desktop/%d.jpg", k)
		err = WriteFile(disfile, data)
		if err != nil {
			log.Fatal(err)
		}
	}

}
