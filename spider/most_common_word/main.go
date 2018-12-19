package main

import (
	"gopkg.in/resty.v1"
	"go.pkg.wesai.com/p/base_lib/log"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"bytes"
)

func FirstPage(firstaddr string) ([]byte,error) {
	resp,err := resty.R().Get(firstaddr)
	if err != nil {
		log.DLogger().Errorf("FirstPage err: %v",err)
		return nil,err
	}
	//resp.Body()
	//fmt.Println(string(resp.Body()))
	return resp.Body(),nil
}

func HeadleFirstPage(resp []byte)  {
	doc ,err :=goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		log.DLogger().Errorf("HeadleFirstPage err: %v",err)
	}
	doc.Find()
}

func main()  {
	firstaddr := "http://www.51voa.com/"
	resp,err := FirstPage(firstaddr)
	if err != nil {
		log.DLogger().Fatal(err)
	}


}
