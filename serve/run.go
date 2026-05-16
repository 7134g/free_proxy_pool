package serve

import (
	"free_proxy_pool/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	InitRouter(r)

	log.Println("启动服务===========>：", config.Cfg.Service.Url)
	if err := r.Run(config.Cfg.Service.Url); err != nil {
		log.Fatalln(err)
	}
}
