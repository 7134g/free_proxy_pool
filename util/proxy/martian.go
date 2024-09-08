package proxy

import (
	"free_proxy_pool/util/cas"
	"github.com/google/martian"
	"github.com/google/martian/log"
	"github.com/google/martian/mitm"
	"github.com/google/martian/priority"
	"net"
	"net/http"
	"net/url"
)

var (
	MonitorAddress = "127.0.0.1:10888" // 监听地址
)

var (
	httpMartian *martian.Proxy // 拦截器全局对象
	certFlag    bool           // 开启自签证书验证
)

var (
	serverProxyUrlParse *url.URL // 解析代理

	serverProxyFlag     bool   // 启用代理
	serverProxy         string // 服务代理地址
	serverProxyUsername string // 用户名
	serverProxyPassword string // 密码
)

func init() {
	lock = cas.NewSpinLock()
	log.SetLevel(log.Silent)
}

func OpenCert() {
	certFlag = true
	_ = LoadCert()
}

func SetServeProxyAddress(address, username, password string) {
	if address == "" {
		serverProxyFlag = false
		return
	}
	serverProxy = address
	serverProxyUsername = username
	serverProxyPassword = password
	u, err := url.Parse(serverProxy)
	if err != nil {
		serverProxyFlag = false
		return
	}
	serverProxyUrlParse = u
	serverProxyFlag = true
}

func Martian() error {
	httpMartian = martian.NewProxy()
	if certFlag {
		mc, err := mitm.NewConfig(ca, private)
		if err != nil {
			return err
		}
		httpMartian.SetMITM(mc)
	}

	group := priority.NewGroup()
	xs := newHttpProxy()
	group.AddRequestModifier(xs, 10)
	group.AddResponseModifier(xs, 10)
	httpMartian.SetRequestModifier(group)
	httpMartian.SetResponseModifier(group)

	listener, err := net.Listen("tcp", MonitorAddress)
	if err != nil {
		return err
	}

	err = httpMartian.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

type httpProxy struct {
}

func newHttpProxy() *httpProxy {
	return &httpProxy{}
}

func (r *httpProxy) ModifyRequest(req *http.Request) error {
	if serverProxyFlag {
		httpMartian.SetDownstreamProxy(serverProxyUrlParse)
	}
	return nil
}

func (r *httpProxy) ModifyResponse(res *http.Response) error {
	switch res.StatusCode {
	case 200, 201, 202, 301, 302:
		break
	default:
		taskIncError()
	}

	return nil
}
