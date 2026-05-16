package cell

import (
	"bytes"
	"fmt"
	"log"
)

type crawlGitHub struct {
	crawl
}

func (c *crawlGitHub) name() string {
	return "crawlGitHub"
}

func (c *crawlGitHub) genSeek() {
	c.links = append(c.links, "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt")
}

func (c *crawlGitHub) parse(html []byte) ([]string, error) {
	lines := bytes.Split(html, []byte("\n"))
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if reProxyLine.Match(line) {
			urls = append(urls, fmt.Sprintf("http://%s", line))
		}
	}

	log.Printf("github 解析到 %d 个代理", len(urls))
	return urls, nil
}
