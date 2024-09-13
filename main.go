package main

import (
	"flag"
	"free_proxy_pool/config"
	"free_proxy_pool/crawler"
	"free_proxy_pool/serve"
	"log"
)

func main() {
	cfgPath := flag.String("config", "./config.yaml", "配置")
	flag.Parse()

	config.Init(*cfgPath)
	info()

	go serve.Run()

	go serve.Martian()

	go crawler.Run()

	select {}
}

func info() {
	log.Printf("代理地址：%s\n", config.Cfg.Martian)
	log.Printf("服务地址：%s\n", config.Cfg.Service.Url)
}
