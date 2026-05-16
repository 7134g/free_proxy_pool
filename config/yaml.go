package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Redis struct {
	Enabled  bool   `yaml:"enabled"`
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Key      string `yaml:"key"`
}

type Service struct {
	Url string `yaml:"url"`
}

type Martian struct {
	Url           string `yaml:"url"`
	Mode          string `yaml:"mode"`            // 代理模式类型
	ErrorMaxCount int    `yaml:"error_max_count"` // 错误最大值
}

type setting struct {
	Redis       Redis    `yaml:"redis"`        // redis 配置
	Service     Service  `yaml:"service"`      // 服务器地址
	Martian     Martian  `yaml:"martian"`      // 代理服务
	TestTime    string   `yaml:"test_time"`    // 测试周期
	CrawlerTime string   `yaml:"crawler_time"` // 抓取周期
	PoolCap     int      `yaml:"pool_cap"`     // 代理池中最大存在的代理数
	TestUrls    []string `yaml:"test_urls"`    // 测试链接
	FlashScore  int      `yaml:"flash_score"`  // 新鲜度
}

var (
	Cfg         setting
	ConfigPath  string
	RedisClient *redis.Client
)

func Init(p string) {
	if ConfigPath == "" {
		ConfigPath = p
	}
	f, err := os.Open(p)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	decode := yaml.NewDecoder(f)
	if err := decode.Decode(&Cfg); err != nil {
		log.Fatalln(err)
	}

	InitRedis()
}

func InitRedis() {
	if !Cfg.Redis.Enabled {
		log.Println("Redis is disabled, using in-memory storage only")
		return
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Url,
		Password: Cfg.Redis.Password,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis connection failed: %v, falling back to in-memory storage only\n", err)
		RedisClient = nil
		return
	}

	log.Println("Redis connected successfully")
}

func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Redis close failed: %v\n", err)
		}
	}
}
