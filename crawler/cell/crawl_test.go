package cell

import (
	"testing"
)

func TestCrawler_crawlDaiLi66(t *testing.T) {
	c := crawlDaiLi66{}
	c.run()
}

func TestCrawler_crawlIp3366(t *testing.T) {
	c := crawlIp3366{}
	c.run()
}

func TestCrawler_crawlKxDaiLi(t *testing.T) {
	c := crawlKxDaiLi{}
	c.run()
}

func TestCrawler_crawlProxy11i(t *testing.T) {
	c := crawlProxy11{}
	c.run()
}
