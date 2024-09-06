package serve

import (
	"fmt"
	"free_proxy_pool/crawler"
	"github.com/gin-gonic/gin"
	"math/rand"
)

func proxyMax(ctx *gin.Context) {
	ctx.String(200, crawler.CacheProxyData.GetOne(0))
}

func proxyList(ctx *gin.Context) {
	list := crawler.CacheProxyData.GetMaxList()

	ctx.JSON(200, list)
}

func proxyRandom(ctx *gin.Context) {
	count := crawler.CacheProxyData.GetCount()
	index := rand.Intn(count)

	ctx.String(200, crawler.CacheProxyData.GetOne(index))
}

func proxyUseless(ctx *gin.Context) {
	link := ctx.GetString("url")
	crawler.CacheProxyData.Del(link)
	ctx.String(200, "ok")
}

func proxyCount(ctx *gin.Context) {
	ctx.String(200, fmt.Sprintf("%d", crawler.CacheProxyData.GetCount()))
}
