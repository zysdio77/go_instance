package main

import (
	"flag"
	"go.pkg.wesai.com/p/base_lib/log"
	"go_instance/spider/spider_cartoon/headle"

)

func main() {
	var fileaddr string
	flag.StringVar(&fileaddr,
		"config",
		"/Users/zhangyongsheng/data/src/go_instance/gin/new_story_web/config.toml",
		"Please flag a config file")
	flag.Parse()

	conf, err := headle.ReadConf(fileaddr)
	if err != nil {
		log.DLogger().Fatal(err)
	}

	var info *headle.Info
	info.Host = conf.SrcDomain
	info.GetHtml()
}