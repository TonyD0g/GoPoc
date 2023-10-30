package Http

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

func SetProxy(inputProxy string) *http.Client {
	// 创建一个代理函数，将请求发送到指定的本地端口
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(inputProxy)
	}

	// 创建一个自定义的 Transport，并设置代理函数
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           proxy,
	}

	// 创建一个使用自定义 Transport 的 HTTP 客户端
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}
