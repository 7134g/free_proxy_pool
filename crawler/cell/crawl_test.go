package cell

import (
	"testing"
)

func TestCrawler_crawlDaiLi66(t *testing.T) {
	c := crawlDaiLi66{}
	c.run(&c)
}

func TestCrawler_crawlIp3366(t *testing.T) {
	c := crawlIp3366{}
	c.run(&c)
}

func TestCrawler_crawlKxDaiLi(t *testing.T) {
	c := crawlKxDaiLi{}
	c.run(&c)
}

func TestCrawler_crawlProxy11i(t *testing.T) {
	c := crawlProxy11{}
	c.run(&c)
}

func TestCrawler_crawl89ip(t *testing.T) {
	c := crawl89ip{}
	c.run(&c)
}

func TestCrawler_crawl66DaiLi(t *testing.T) {
	c := crawl66DaiLi{}
	c.run(&c)
}

func TestCrawler_crawlProxyScrape(t *testing.T) {
	c := crawlProxyScrape{}
	c.run(&c)
}

func TestCrawler_crawlGitHub(t *testing.T) {
	c := crawlGitHub{}
	c.run(&c)
}

func TestCrawler_crawlZdaye(t *testing.T) {
	c := crawlZdaye{}
	c.run(&c)
}
