package main

import (
	"flag"
	"free_proxy_pool/config"
	"free_proxy_pool/crawler"
	"free_proxy_pool/serve"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfgPath := flag.String("config", "./config.yaml", "配置")
	flag.Parse()

	config.Init(*cfgPath)
	info()

	go serve.Run()

	go serve.Martian()

	go crawler.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	config.CloseRedis()
}

func info() {
	log.Printf("代理地址：%s\n", config.Cfg.Martian.Url)
	log.Printf("服务地址：%s\n", config.Cfg.Service.Url)
}
