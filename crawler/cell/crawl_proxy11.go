package cell

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type crawlProxy11 struct {
	crawl
}

func (c *crawlProxy11) name() string {
	return "crawlProxy11"
}

func (c *crawlProxy11) run() {
	/*

	    https://proxy11.com/free-proxy/speed

	   /free-proxy
	   /free-proxy/speed
	   /free-proxy/us
	   /free-proxy/anonymous
	   /free-proxy/instagram
	   /free-proxy/google
	   /free-proxy/prot-8080

	*/

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
		if err := catch(c, u); err != nil {
			log.Println("crawlIp3366 error:", err)
			continue
		}
		c.keepRunning()
	}

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
