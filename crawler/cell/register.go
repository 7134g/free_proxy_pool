package cell

import (
	"free_proxy_pool/util"
	"log"
	"time"
)

var (
	ProxyChannel chan string         // 抓取到的代理队列
	spiderMap    map[string]struct{} // 爬虫列表
	SleepTime    time.Duration       // 抓取间隔

	linkErrorMap *util.LinkMap
)

func init() {
	ProxyChannel = make(chan string, 1000)
	spiderMap = make(map[string]struct{})
	SleepTime = time.Second * 5

	linkErrorMap = util.NewLinkMap()
}

func register(spiders ...spider) {
	for _, s := range spiders {
		_, exist := spiderMap[s.name()]
		if exist {
			continue
		}

		go runSpider(s)
	}
}

func runSpider(s spider) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("爬虫：%v 运行失败，错误信息：%v\n", s.name(), err)
		}
		delete(spiderMap, s.name())
	}()

	log.Printf("启动爬虫：%v\n", s.name())
	s.run(s)
	log.Printf("爬虫：%v 已停止运行\n", s.name())
}

// Crawler 启动爬虫
func Crawler() {

	log.Println("crawler_start......抓取爬虫")
	register(
		&crawlDaiLi66{},
		&crawlIp3366{},
		&crawlKxDaiLi{},
		&crawlProxy11{},
	)
	log.Println("crawler_stop")
}
