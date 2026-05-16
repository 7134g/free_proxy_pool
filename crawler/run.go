package crawler

import (
	"free_proxy_pool/config"
	"free_proxy_pool/crawler/cell"
	"free_proxy_pool/util/xhttp"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func Run() {
	recoverFromRedis()

	c := cron.New()
	go monitor()
	// 启动脚本时候立马启动，后面启动定时任务
	cell.Crawler()

	if _, err := c.AddFunc(config.Cfg.CrawlerTime, cell.Crawler); err != nil {
		log.Fatal(err)
	}
	if _, err := c.AddFunc(config.Cfg.TestTime, TestStoreProxy); err != nil {
		log.Fatal(err)
	}

	c.Run()
}

func monitor() {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			log.Println("proxy pool size:", CacheProxyData.GetCount())
			CacheProxyData.sort()
			xhttp.SetLocalProxy(CacheProxyData.GetOnce(0))

		case link := <-cell.ProxyChannel:
			if err := TaskPool.Submit(checkProxy(link)); err != nil {
				log.Println(err)
			}

		case result := <-ProxyFinishChannel:
			if result.status {
				CacheProxyData.inc(result.link)
				continue
			}
			if exist := CacheProxyData.dnc(result.link); exist {
				continue
			}
		}
	}

}
