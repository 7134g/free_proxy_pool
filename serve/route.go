package serve

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.GET("/", home)
	r.GET("/config", proxyUpdateConfig)
	r.GET("/max", proxyMax)
	r.GET("/list", proxyList)
	r.GET("/random", proxyRandom)
	r.GET("/useless", proxyUseless)
	r.GET("/count", proxyCount)
}
