package serve

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.GET("/", home)
	r.GET("/max", proxyMax)
	r.GET("/random", proxyRandom)
	r.GET("/useless", proxyUseless)
	r.GET("/count", proxyCount)
}
