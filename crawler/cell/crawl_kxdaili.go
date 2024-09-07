package cell

import (
	"bytes"
	"fmt"
	"free_proxy_pool/util"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type crawlKxDaiLi struct {
	crawl
}

func (c *crawlKxDaiLi) genSeek() {
	urls := []string{
		"http://www.kxdaili.com/dailiip/1/%d.html",
		"http://www.kxdaili.com/dailiip/2/%d.html",
	}

	for _, u := range urls {
		for page := 1; page < 50; page++ {
			link := fmt.Sprintf(u, page)
			c.links = append(c.links, link)
		}
	}
}

func (c *crawlKxDaiLi) name() string {
	return "crawlKxDaiLi"
}

func (c *crawlKxDaiLi) run() {
	c.defaultRun(c)
}

func (c *crawlKxDaiLi) parse(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("table.active > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// 对于每个<div>元素，打印其内容
		//if i == 0 {
		//	return
		//}
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()
		scheme := s.Find("td:nth-child(4)").Text()
		scheme = util.FixScheme(scheme)
		link := fmt.Sprintf("%s://%s:%s", strings.ToLower(scheme), ip, port)
		urls = append(urls, link)
	})

	return urls, nil
}
