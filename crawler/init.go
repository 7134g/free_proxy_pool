package crawler

import (
	"context"
	"free_proxy_pool/config"
	"free_proxy_pool/util/pool"
	"sync/atomic"
	"time"
)

type proxy struct {
	Link  string // 代理地址
	Score int    // 新鲜度

	errCount int // 测试失败数
	sucCount int // 测试成功数
}

func newProxy(link string) *proxy {
	return &proxy{
		Link:     link,
		Score:    config.Cfg.FlashScore,
		errCount: 0,
		sucCount: 0,
	}
}

type proxyResult struct {
	link   string
	status bool
}

var (
	CacheProxyData     *Store
	TaskPool           *pool.Pool
	TaskCancel         context.CancelFunc
	ProxyFinishChannel chan proxyResult // 结果队列
	TesterRunning      atomic.Bool      // 是否处于测试爬虫运行中
)

func init() {
	p, cancel := pool.NewPool(200, false, time.Second)
	TaskPool = p
	TaskCancel = cancel

	ProxyFinishChannel = make(chan proxyResult, 10000)

	CacheProxyData = &Store{
		body: map[string]*proxy{},
	}
}
