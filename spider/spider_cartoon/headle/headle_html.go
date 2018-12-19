package headle

import (
	"gopkg.in/resty.v1"
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	Host string
	Htmlbody []byte
	GetInfo
}

type GetInfo struct {
	HrefList []string
	MenuList []string
	PicList []string
}

func (i *Info)	GetHtml() (error) {	//获取网页内容
	useragent := make(map[string]string)
	useragent["User-Agent"]=" Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36"
	resp,err := resty.R().SetHeaders(useragent).Get(i.Host)
	//fmt.Println(i.Host)
	if err != nil {
		return err
	}
	i.Htmlbody=resp.Body()
	//fmt.Println(i.Htmlbody)
	return nil
}

func (i *Info)FilterHtml() error {	//过滤每章标题和链接
	r := bytes.NewReader(i.Htmlbody)
	doc , err :=goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	var hlist []string
	var piclist []string
	doc.Find("table.table").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(i,s)
		hrefaddr,ok := s.Find("td").Find("a").Attr("href")
		if ok{
			//fmt.Println(hrefaddr)
			hlist= append(hlist,hrefaddr)
		}
		picaddr , ok :=s.Find("td").Find("a").Attr("title")
		if ok{
			piclist= append(piclist,picaddr)
			//fmt.Println(picaddr)
		}
		//menuname := s.Find("a").Find()
	})
	i.GetInfo.HrefList = hlist
	i.GetInfo.MenuList = piclist
	return nil
}

func (i *Info)FilterPicHtml() error {	//过滤图片的地址
	r := bytes.NewReader(i.Htmlbody)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	var piclist []string
	//div class="list comic-imgs"
	doc.Find("div.list.comic-imgs").Find("img").Each(func(i int, s *goquery.Selection) {
		picaddr,ok := s.Attr("data-kksrc")
		if ok{
			piclist = append(piclist,picaddr)
		}
	})
	i.GetInfo.PicList =piclist
	return nil
}
