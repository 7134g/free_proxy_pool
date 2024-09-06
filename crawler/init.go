package crawler

import (
	"context"
	"free_proxy_pool/config"
	"free_proxy_pool/util/pool"
	"time"
)

type proxy struct {
	link string

	score    int // 新鲜度
	errCount int // 测试失败数
	sucCount int // 测试成功数
}

type proxyResult struct {
	link   string
	status bool

	countFail int
	createAt  time.Time
}

var (
	CacheProxyData *Store
	TaskPool       *pool.Pool
	TaskCancel     context.CancelFunc

	ProxyFinishChannel chan proxyResult // 结果队列

)

func init() {
	p, cancel := pool.NewPool(int32(config.Cfg.Concurrent), false, time.Second)
	TaskPool = p
	TaskCancel = cancel

	ProxyFinishChannel = make(chan proxyResult, 10000)

	CacheProxyData = &Store{
		body: map[string]*proxy{},
	}
}
