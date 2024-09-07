package serve

import (
	"free_proxy_pool/config"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"syscall"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = os.NewFile(uintptr(syscall.Stdin), os.DevNull)

	r := gin.Default()
	InitRouter(r)

	log.Println("服务地址：", config.Cfg.Service.Url)
	if err := r.Run(config.Cfg.Service.Url); err != nil {
		log.Fatalln(err)
	}
}
