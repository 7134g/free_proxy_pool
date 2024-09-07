package cell

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type crawlDaiLi66 struct {
	crawl
}

func (c *crawlDaiLi66) name() string {
	return "crawlDaiLi66"
}

func (c *crawlDaiLi66) genSeek() {
	locals := []string{
		"/areaindex_1",
		"/areaindex_2",
		"/areaindex_3",
		"/areaindex_4",
		"/areaindex_5",
		"/areaindex_6",
		"/areaindex_7",
		"/areaindex_8",
		"/areaindex_9",
		"/areaindex_10",
		"/areaindex_11",
		"/areaindex_12",
		"/areaindex_13",
		"/areaindex_14",
		"/areaindex_15",
		"/areaindex_16",
		"/areaindex_17",
		"/areaindex_18",
		"/areaindex_19",
		"/areaindex_20",
		"/areaindex_21",
		"/areaindex_22",
		"/areaindex_23",
		"/areaindex_24",
		"/areaindex_25",
		"/areaindex_26",
		"/areaindex_28",
		"/areaindex_29",
		"/areaindex_30",
		"/areaindex_31",
		"/areaindex_32",
	}
	baseUrl := "http://www.66ip.cn%s/%d.html"
	page := 1911

	for i := 0; i < 100; i++ {
		for _, local := range locals {
			link := fmt.Sprintf(baseUrl, local, page+i)
			c.links = append(c.links, link)
		}
	}
}

func (c *crawlDaiLi66) run() {
	c.defaultRun(c)
}

func (c *crawlDaiLi66) parse(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("#footer > div > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// 对于每个<div>元素，打印其内容
		if i == 0 {
			return
		}
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()
		link := fmt.Sprintf("http://%s:%s", ip, port)
		urls = append(urls, link)

	})

	return urls, nil
}
