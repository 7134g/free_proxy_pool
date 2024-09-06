package cell

import (
	"errors"
	"free_proxy_pool/util/xhttp"
	"log"
	"time"
)

var (
	ProxyChannel chan string // 抓取到的代理队列
	spiderMap    map[string]struct{}
	SleepTime    time.Duration
)

func init() {
	ProxyChannel = make(chan string, 10000)
	spiderMap = make(map[string]struct{})
	SleepTime = time.Second * 5
}

type spider interface {
	getStatus() bool
	stop()

	name() string
	run()
	parse(html []byte) ([]string, error)
}

type crawl struct {
	status bool
}

func (c *crawl) getStatus() bool {
	return c.status
}

func (c *crawl) keepRunning() {
	c.status = false
}

func (c *crawl) stop() {
	c.status = true
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
	s.run()
	log.Printf("爬虫：%v 已停止运行\n", s.name())
}

func addUrl(link string) {
	ProxyChannel <- link
}

func catch(s spider, link string) error {
	if s.getStatus() {
		return nil
	}

	log.Printf("正在抓取：%v\n", link)
	dat, err := xhttp.Get(link)
	if err != nil {
		return err
	}

	if len(dat) == 0 {
		return errors.New("crawl error")
	}

	urls, err := s.parse(dat)
	if err != nil {
		return err
	}

	for _, url := range urls {
		addUrl(url)
	}

	if len(urls) == 0 {
		s.stop()
	}

	time.Sleep(SleepTime)

	return nil
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
