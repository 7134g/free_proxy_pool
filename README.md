# 免费代理池

## 安装方式
go1.18 环境以上

```shell
go mod tidy
go build .
```

## 运行方式
`./free_proxy_pool`

## 配置文件
`config.yaml` 具体含义看里面的注释

## 使用方式
- 一、通过`get`请求

  http://127.0.0.1:5555/max  或者 http://127.0.0.1:5555/random   
  获取代理

```go
func UpdateLocalProxy() {
  if errCount < 50 {
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

```


- 二、通过代理端口方式  
  由`free_proxy_pool`转发你的代理请求
  默认端口为`:10888`

```go
func httpProxy() http.RoundTripper {
	proxyUrl, err := url.Parse(`http://127.0.0.1:10888`)
	if err != nil {
		log.Fatalln(err)
	}
	fc := http.ProxyURL(proxyUrl)

	ht := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           fc,
	}

	return ht
}

...

  client := &http.Client{
    Transport: httpProxy(config.Configs.Proxy),
    Timeout:   time.Second * 10,
  }
```


## todo
1. 将抓取的代理存入redis中
2. docker 化
