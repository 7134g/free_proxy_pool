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
			xhttp.SetLocalProxy(CacheProxyData.GetMax(0))

		case link := <-cell.ProxyChannel:
			if err := TaskPool.Submit(TestProxy(link)); err != nil {
				log.Println(err)
			}

		case result := <-ProxyFinishChannel:
			// 新鲜度
			if result.status {
				CacheProxyData.inc(result.link)
				continue
			}
			if exist := CacheProxyData.dnc(result.link); exist {
				continue
			}

			// 未添加过
			ctLessFiveMinute := result.createAt.Add(-time.Minute * 5)
			if result.createAt.After(ctLessFiveMinute) && result.countFail < 5 {
				result.countFail++
				if err := TaskPool.Submit(TestProxy(result.link)); err != nil {
					log.Println(err)
				}
			}
		}
	}

}
