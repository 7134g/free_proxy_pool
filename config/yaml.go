package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Redis struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Key      string `yaml:"key"`
}

type Service struct {
	Url string `yaml:"url"`
}

type setting struct {
	Redis       Redis    `yaml:"redis"`        // redis 配置
	Service     Service  `yaml:"service"`      // 服务器地址
	TestTime    string   `yaml:"test_time"`    // 测试周期
	CrawlerTime string   `yaml:"crawler_time"` // 抓取周期
	PoolCap     int      `yaml:"pool_cap"`     // 代理池中最大存在的代理数
	TestUrls    []string `yaml:"test_urls"`    // 测试链接
	FlashScore  int      `yaml:"flash_score"`  // 新鲜度
}

var Cfg setting

func Init(p string) {
	f, err := os.Open(p)
	if err != nil {
		log.Fatalln(err)
	}

	decode := yaml.NewDecoder(f)
	if err := decode.Decode(&Cfg); err != nil {
		log.Fatalln(err)
	}

	InitRedis()
}

func InitRedis() {
	// todo
}
