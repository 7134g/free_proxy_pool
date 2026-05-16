package cell

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
)

type crawlProxyScrape struct {
	crawl
}

var reProxyLine = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+:\d+`)

func (c *crawlProxyScrape) name() string {
	return "crawlProxyScrape"
}

func (c *crawlProxyScrape) genSeek() {
	c.links = append(c.links, "https://api.proxyscrape.com/v2/?request=displayproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all")
}

func (c *crawlProxyScrape) parse(html []byte) ([]string, error) {
	lines := bytes.Split(html, []byte("\n"))
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if reProxyLine.Match(line) {
			urls = append(urls, fmt.Sprintf("http://%s", line))
		}
	}

	log.Printf("proxyscrape 解析到 %d 个代理", len(urls))
	return urls, nil
}
