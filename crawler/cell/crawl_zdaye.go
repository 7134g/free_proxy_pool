package cell

import (
	"fmt"
	"free_proxy_pool/util/xhttp"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type crawlZdaye struct {
	crawl
}

var reZdayeIP = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+):(\d+)`)

func (c *crawlZdaye) name() string {
	return "crawlZdaye"
}

func (c *crawlZdaye) genSeek() {
	baseUrl := "https://www.zdaye.com/dayProxy/%d.html"
	for page := 1; page <= 4; page++ {
		link := fmt.Sprintf(baseUrl, page)
		c.links = append(c.links, link)
	}
}

func (c *crawlZdaye) parse(data []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("#J_posts_list .thread_item div div p a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		detailUrl := "https://www.zdaye.com" + href
		log.Printf("zdaye 详情页: %s", detailUrl)

		detailData, err := xhttp.GetHeader(detailUrl, map[string]string{
			"Cookie": "__root_domain_v=.zdaye.com; _qddaz=QD.0.0.0.0; _qdda=3-0.0; _qddab=3-0.0; Hm_lvt_80f407a85cf0bc32ab5f9cc91c15f88b=0; Hm_lpvt_80f407a85cf0bc32ab5f9cc91c15f88b=0; acw_sc__v2=0",
		})
		if err != nil {
			log.Printf("zdaye 详情页获取失败: %v", err)
			return
		}

		detailDoc, err := goquery.NewDocumentFromReader(strings.NewReader(string(detailData)))
		if err != nil {
			return
		}

		detailDoc.Find(".cont br").Each(func(j int, br *goquery.Selection) {
			node := br.Get(0)
			if node != nil && node.NextSibling != nil && node.NextSibling.Type == html.TextNode {
				line := strings.TrimSpace(node.NextSibling.Data)
				match := reZdayeIP.FindStringSubmatch(line)
				if len(match) >= 3 {
					link := fmt.Sprintf("http://%s:%s", match[1], match[2])
					urls = append(urls, link)
				}
			}
		})
	})

	log.Printf("zdaye 解析到 %d 个代理", len(urls))
	return urls, nil
}
