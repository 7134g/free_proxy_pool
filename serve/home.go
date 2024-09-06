package serve

import (
	"github.com/gin-gonic/gin"
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
