package crawler

import (
	"free_proxy_pool/config"
	"free_proxy_pool/crawler/cell"
	"free_proxy_pool/util/pool"
	"log"
	"net/http"
	"net/url"
	"time"
)

func TestProxy(proxy string) *pool.Task {
	t := &pool.Task{}
	t.Param = []interface{}{proxy}
	t.TaskFunc = func(i []interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("testProxy error:", err)
			}
		}()
		proxy := i[0].(string)

		for _, link := range config.Cfg.TestUrls {
			status := runTestProxy(link, proxy)
			addProxyFinish(proxy, status)
		}
	}
	return t
}

func addProxyFinish(link string, status bool) {
	ProxyFinishChannel <- proxyResult{link: link, status: status, createAt: time.Now()}
}

// TestStoreProxy 测试此时库里的代理是否有效
func TestStoreProxy() {
	log.Println("tester_start......测试代理")
	list := CacheProxyData.slice
	for _, p := range list {
		cell.ProxyChannel <- p.Link
	}
	log.Println("tester_stop")
}

func runTestProxy(link, proxy string) bool {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	u, err := url.Parse(proxy)
	client := &http.Client{}
	client.Timeout = time.Second * 10
	client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
	resp, err := client.Do(req)
	if err != nil {
		//log.Println(err)
		return false
	}

	if resp != nil {
		defer resp.Body.Close()
	} else {
		//log.Println("response is nil")
		return false
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true
	default:
		//log.Printf("response status code is %d\n", resp.StatusCode)
		return false
	}
}
