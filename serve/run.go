package serve

import (
	"free_proxy_pool/config"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {

	r := gin.Default()
	InitRouter(r)

	log.Println("服务地址：", config.Cfg.Service.Url)
	if err := r.Run(config.Cfg.Service.Url); err != nil {
		log.Fatalln(err)
	}
}
