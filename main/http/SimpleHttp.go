package http

import (
	"Scanner/main/format"
	"Scanner/main/handle"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func Send(requestPackage format.RequestPackage) {
	// 创建http.Client对象
	client := &http.Client{}

	if !handle.CheckFileCorrectness(requestPackage) {
		return
	}

	// 创建http.Request对象
	req, err := http.NewRequest(requestPackage.Method, requestPackage.Url, nil)
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
