package http

import (
	"Scanner/main/format"
	"Scanner/main/handle"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Send(requestPackage format.RequestPackage, inputProxy string) {

	// 创建一个代理函数，将请求发送到指定的本地端口
	proxy := func(_ *http.Request) (*url.URL, error) {
		if inputProxy != "" {
			return url.Parse(inputProxy)
		}
		return url.Parse("http://localhost:8080")
	}

	// 创建一个自定义的 Transport，并设置代理函数
	transport := &http.Transport{
		Proxy: proxy,
	}

	// 创建一个使用自定义 Transport 的 HTTP 客户端
	client := &http.Client{
		Transport: transport,
	}

	if !handle.CheckFileCorrectness(requestPackage) {
		return
	}

	requestBody := []byte(requestPackage.Body)

	// TODO 遍历 PathList 中的 Path,然后发起请求
	allReqPath := handle.TraversePath(requestPackage)
	for _, tmpPath := range allReqPath {
		// 创建http.Request对象
		req, err := http.NewRequest(requestPackage.Method, tmpPath, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		handle.HandlePackFunc(req, requestPackage)

		// 发送请求并获取响应结果
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("[-] error,please check io.ReadCloser")
			}
		}(resp.Body)

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 打印响应结果
		fmt.Println(string(body))
	}
}
