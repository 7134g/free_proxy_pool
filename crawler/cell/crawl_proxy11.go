package cell

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type crawlProxy11 struct {
	crawl
}

func (c *crawlProxy11) name() string {
	return "crawlProxy11"
}

func (c *crawlProxy11) genSeek() {
	urls := []string{
		"https://proxy11.com/free-proxy",
		"https://proxy11.com/free-proxy/speed",
		"https://proxy11.com/free-proxy/us",
		"https://proxy11.com/free-proxy/anonymous",
		"https://proxy11.com/free-proxy/instagram",
		"https://proxy11.com/free-proxy/google",
		"https://proxy11.com/free-proxy/prot-8080",
	}

	for _, u := range urls {
		c.links = append(c.links, u)
	}
}

func (c *crawlProxy11) run() {
	c.defaultRun(c)
}

func (c *crawlProxy11) parse(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("div.row > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// 对于每个<div>元素，打印其内容
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()
		link := fmt.Sprintf("%s://%s:%s", "http", ip, port)
		urls = append(urls, link)
	})

	return urls, nil
}
