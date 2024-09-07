package xhttp

import (
	"free_proxy_pool/util/cas"
	"io"
	"log"
	"net/http"
)

var (
	lock       = cas.NewSpinLock()
	errCount   int
	localProxy string
)

func SetLocalProxy(link string) {
	localProxy = link
}

func IncHttpErrorCount() {
	lock.Lock()
	defer lock.Unlock()
	errCount++
}

func UpdateLocalProxy() {
	if errCount < 100 {
		// 错误数大于50
		return
	}

	resp, err := http.Get("http://127.0.0.1:5555/max")
	if err != nil {
		log.Fatalln(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	errCount = 0
	log.Println("当前使用代理：", localProxy)
	localProxy = string(b)
}
