package cell

import (
	"bytes"
	"fmt"
	"free_proxy_pool/util"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type crawlIp3366 struct {
	crawl
}

func (c *crawlIp3366) name() string {
	return "crawlIp3366"
}

func (c *crawlIp3366) genSeek() {
	baseUrl := "http://www.ip3366.net/?stype=1&page=%d"
	for page := 1; page < 100; page++ {
		link := fmt.Sprintf(baseUrl, page)
		c.links = append(c.links, link)
	}
}

func (c *crawlIp3366) parse(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("#list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
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
