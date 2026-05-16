package cell

import (
	"bytes"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type crawl89ip struct {
	crawl
}

func (c *crawl89ip) name() string {
	return "crawl89ip"
}

func (c *crawl89ip) genSeek() {
	baseUrl := "https://www.89ip.cn/index_%d.html"
	for page := 1; page <= 6; page++ {
		link := fmt.Sprintf(baseUrl, page)
		c.links = append(c.links, link)
	}
}

func (c *crawl89ip) parse(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()
		link := fmt.Sprintf("http://%s:%s", ip, port)
		urls = append(urls, link)
	})

	log.Printf("89ip 解析到 %d 个代理", len(urls))
	return urls, nil
}
