package cell

import (
	"errors"
	"free_proxy_pool/util/xhttp"
	"log"
	"time"
)

type spider interface {
	name() string
	genSeek()
	run()
	parse(html []byte) ([]string, error)
}

type crawl struct {
	links []string
}

func (c *crawl) defaultRun(s spider) {
	s.genSeek()
	for _, link := range c.links {
		if err := catch(s, link); err != nil {
			log.Printf("%s error: %v\n", s.name(), err)
			continue
		}
	}
}

func addUrl(link string) {
	ProxyChannel <- link
}

func catch(s spider, link string) error {
	if !linkErrorMap.Check(link) {
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

	if len(urls) == 0 {
		linkErrorMap.Add(link)
	}

	log.Printf("抓取：%s -> %v 个代理地址\n", s.name(), len(urls))
	for _, url := range urls {
		addUrl(url)
	}

	time.Sleep(SleepTime)

	return nil
}
