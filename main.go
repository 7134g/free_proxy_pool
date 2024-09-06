package main

import (
	"flag"
	"free_proxy_pool/config"
	"free_proxy_pool/crawler"
	"free_proxy_pool/serve"
)

func main() {
	cfgPath := flag.String("config", "./config.yaml", "配置")
	flag.Parse()

	config.Init(*cfgPath)

	go serve.Run()

	go crawler.Run()

	select {}
}
