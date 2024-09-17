package serve

import (
	"fmt"
	"free_proxy_pool/config"
	"free_proxy_pool/crawler"
	"github.com/gin-gonic/gin"
	"log"
)

func home(ctx *gin.Context) {
	html := `
<!DOCTYPE html>

<head>
<meta charset="UTF-8" />
</head>

<h2>
<a href="https://github.com/7134g">欢迎来到我的代理池，点我即刻跳转github</a>
</h2>

<br/>

<h3>没错就是广告</h3>

`
	ctx.Data(200, "text/html", []byte(html))
}

func proxyUpdateConfig(ctx *gin.Context) {
	log.Println("更新配置文件")
	config.Init(config.ConfigPath)
	ctx.String(200, "ok")
}

func proxyMax(ctx *gin.Context) {
	ctx.String(200, crawler.CacheProxyData.GetOnce(0))
}

func proxyList(ctx *gin.Context) {
	list := crawler.CacheProxyData.GetMaxList()

	ctx.JSON(200, list)
}

func proxyRandom(ctx *gin.Context) {
	ctx.String(200, crawler.CacheProxyData.Random())
}

func proxyUseless(ctx *gin.Context) {
	link := ctx.GetString("url")
	crawler.CacheProxyData.Del(link)
	ctx.String(200, "ok")
}

func proxyCount(ctx *gin.Context) {
	ctx.String(200, fmt.Sprintf("%d", crawler.CacheProxyData.GetCount()))
}
