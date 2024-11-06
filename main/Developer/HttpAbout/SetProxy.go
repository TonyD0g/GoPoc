package HttpAbout

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

func SetProxy(inputProxy string, followRedirects bool) *http.Client {
	if inputProxy == "" {
		fmt.Println("[-] 你忘记输入代理地址了,自动为你配置为: http://127.0.0.1:8080")
		inputProxy = "http://127.0.0.1:8080"
	}

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
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !followRedirects {
				return http.ErrUseLastResponse // 如果不追随重定向，则返回上一个response
			}
			return nil // 如果跟随重定向，则允许继续
		},
	}
	return client
}
