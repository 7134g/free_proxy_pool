package cell

import (
	"fmt"
	"log"
	"regexp"
)

type crawl66DaiLi struct {
	crawl
}

var (
	re66ip   = regexp.MustCompile(`<li[^>]*?min-width:\s*120px[^>]*?>([\d.]+)</li>`)
	re66port = regexp.MustCompile(`<li[^>]*?min-width:\s*60px[^>]*?>(\d+)</li>`)
)

func (c *crawl66DaiLi) name() string {
	return "crawl66DaiLi"
}

func (c *crawl66DaiLi) genSeek() {
	baseUrl := "http://www.66daili.com/?page=%d"
	for page := 1; page <= 10; page++ {
		link := fmt.Sprintf(baseUrl, page)
		c.links = append(c.links, link)
	}
}

func (c *crawl66DaiLi) parse(html []byte) ([]string, error) {
	ipList := re66ip.FindAllSubmatch(html, -1)
	portList := re66port.FindAllSubmatch(html, -1)

	urls := make([]string, 0, len(ipList))
	minLen := len(ipList)
	if len(portList) < minLen {
		minLen = len(portList)
	}
	for i := 0; i < minLen; i++ {
		link := fmt.Sprintf("http://%s:%s", ipList[i][1], portList[i][1])
		urls = append(urls, link)
	}

	log.Printf("66daili 解析到 %d 个代理", len(urls))
	return urls, nil
}
