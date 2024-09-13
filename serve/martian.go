package serve

import (
	"free_proxy_pool/config"
	"free_proxy_pool/crawler"
	"free_proxy_pool/util/proxy"
	"log"
	"time"
)

func Martian() {
	if config.Cfg.Martian == "" {
		return
	}
	go monitor()

	log.Println("开启代理服务==================>", config.Cfg.Martian)
	proxy.MonitorAddress = config.Cfg.Martian
	if err := proxy.Martian(); err != nil {
		log.Fatalln(err)
	}
}

func monitor() {
	ticker := time.NewTicker(time.Second)
	checkMartianProxyTicker := time.NewTicker(time.Second * 5)
	var martianProxyStatus bool
	now := time.Now()

	for {
		select {
		case <-checkMartianProxyTicker.C:
			if now.Sub(proxy.RunningTime) <= time.Second*5 {
				continue
			}

			if time.Now().Sub(proxy.RunningTime) <= time.Minute {
				if martianProxyStatus {
					continue
				}
				log.Println("==================================")
				log.Println("正在使用代理池的端口代理模式")
				martianProxyStatus = true
			} else {
				martianProxyStatus = false
			}
		case <-ticker.C:
			if !proxy.TaskCheckError() {
				continue
			}

			var newProxy string
			switch config.Cfg.MartianMode {
			case "random":
				newProxy = crawler.CacheProxyData.Random()
			case "max":
				newProxy = crawler.CacheProxyData.GetOnce(0)
			default:
				newProxy = crawler.CacheProxyData.GetOnce(0)
			}

			proxy.SetServeProxyAddress(newProxy, "", "")
		}
	}
}
