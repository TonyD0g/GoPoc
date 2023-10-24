package Http

import (
	"Scanner/main/Developer/Fofa"
	"Scanner/main/Developer/Format"
	handle2 "Scanner/main/Developer/Handle"
	"Scanner/main/Developer/Judge"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func Send(pocStruct Format.PocStruct, queryResponse Fofa.QueryResponse, userInputDetectionURL *string, inputProxy string, maxConcurrentLevel int) {
	client := SetProxy(inputProxy)
	customRequestBody := []byte(pocStruct.RequestPackage.Body)
	var urlsList []string
	if userInputDetectionURL != nil {
		if strings.HasPrefix(strings.ToLower(*userInputDetectionURL), "http://") || strings.HasPrefix(strings.ToLower(*userInputDetectionURL), "https://") {
			urlsList = append(urlsList, *userInputDetectionURL)
		} else {
			// urlFile list
			urlFile, err := os.Open(*userInputDetectionURL)
			if err != nil {
				fmt.Println("can't open file:", err)
				return
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					fmt.Println("can't close file:", err)
					return
				}
			}(urlFile)

			reader := bufio.NewReader(urlFile)
			for {
				line, err := reader.ReadString('\n')
				if err != nil && line == "" {
					break
				}
				urlsList = append(urlsList, strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", ""))
				line = ""
			}
		}
	} else {
		for _, tmpOutcome := range queryResponse.Results {
			if !strings.HasPrefix(tmpOutcome[1].(string), tmpOutcome[0].(string)) {
				urlsList = append(urlsList, tmpOutcome[0].(string)+"://"+tmpOutcome[1].(string))
			} else {
				urlsList = append(urlsList, tmpOutcome[1].(string))
			}
		}
	}

	var waitGroup sync.WaitGroup
	var done = make(chan struct{})
	processedURLs := make(map[string]struct{}) // 用于存储已处理的URL

	for _, tmpUrl := range urlsList {
		allRequestPath := handle2.TraversePath(pocStruct.RequestPackage, tmpUrl)
		requestCount := len(allRequestPath)
		sem := make(chan struct{}, maxConcurrentLevel) // 控制并发度的信号量

		for i := 0; i < requestCount; i++ {
			tmpUrlForAllRequestPath := allRequestPath[i]

			// 检查URL是否已经处理过，如果处理过则跳过
			if _, exists := processedURLs[tmpUrlForAllRequestPath]; exists {
				continue
			}

			sem <- struct{}{} // 占用一个并发度信号
			waitGroup.Add(1)
			processedURLs[tmpUrlForAllRequestPath] = struct{}{} // 标记URL已处理

			go func(url string) {
				defer func() {
					<-sem // 释放一个并发度信号
					waitGroup.Done()
				}()

				if !sendRequest(pocStruct, client, url, customRequestBody) {
					return
				}
			}(tmpUrlForAllRequestPath)

			if (i+1)%maxConcurrentLevel == 0 || i == requestCount-1 {
				go func() {
					waitGroup.Wait()
				}()
			}
		}
	}

	// 在主循环外部等待所有请求完成
	go func() {
		waitGroup.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}
}

func sendRequest(pocStruct Format.PocStruct, client *http.Client, urlAndUri string, customRequestBody []byte) bool {
	_, err := url.Parse(urlAndUri)
	if err != nil {
		return false
	}
	// Create an HTTP.Request object
	procedureRequest, err := http.NewRequest(pocStruct.RequestPackage.Method, urlAndUri, bytes.NewBuffer(customRequestBody))
	if err != nil {
		return false
	}
	handle2.ProcessPackages(procedureRequest, pocStruct)

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	procedureRequest = procedureRequest.WithContext(ctx)

	// Send request and obtain response results
	procedureResponse, err := client.Do(procedureRequest)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(procedureResponse.Body)

	if Judge.IsExploitSuccess(pocStruct, procedureResponse, customRequestBody) {
		splitURL := strings.Split(urlAndUri, "?")
		baseURL := splitURL[0] + "?"
		params := strings.Split(splitURL[1], "&")
		encodedParams := make([]string, len(params))
		for i := range params {
			p := strings.Split(params[i], "=")
			encodedParams[i] = url.QueryEscape(p[0]) + "=" + url.QueryEscape(p[1])
		}

		fmt.Println("[+] [ " + baseURL + strings.Join(encodedParams, "&") + " ]\tSuccess! The target may have this vulnerability")
		return true
	}
	return false
}
