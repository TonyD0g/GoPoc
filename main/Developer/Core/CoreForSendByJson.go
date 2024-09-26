package Core

import (
	Format "GoPoc/main/Developer/AllFormat"
	handle2 "GoPoc/main/Developer/Handle"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/Log"
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func ForSendByJson(urlsList []string, pocStruct Format.PocStruct, maxConcurrentLevel int) {
	client := HttpAbout.SetProxy(HttpAbout.InputProxy, false)
	customRequestBody := []byte(pocStruct.RequestPackage.Body)
	waitGroup := &sync.WaitGroup{}

	// 计算要划分的小的urlsList数量
	numThreads := maxConcurrentLevel
	if numThreads > len(urlsList) {
		numThreads = len(urlsList)
	}
	urlsPerThread := len(urlsList) / numThreads
	for i := 0; i < numThreads; i++ {
		start := i * urlsPerThread
		end := start + urlsPerThread
		if i == numThreads-1 {
			end = len(urlsList)
		}

		subURLsList := urlsList[start:end]
		waitGroup.Add(1)
		go func(subURLs []string) {
			defer func() {
				waitGroup.Add(-1)
			}()

			for _, tmpUrl := range subURLs {
				allRequestPath := handle2.TraversePath(pocStruct.RequestPackage, tmpUrl)
				requestCount := len(allRequestPath)

				for tmpI := 0; tmpI < requestCount; tmpI++ {
					tmpUrlForAllRequestPath := allRequestPath[tmpI]
					parsedURL, err := url.Parse(tmpUrlForAllRequestPath)
					if err != nil {
						continue
					}
					// Create an HTTP.Request object
					procedureRequest, err := http.NewRequest(pocStruct.RequestPackage.Method, tmpUrlForAllRequestPath, bytes.NewBuffer(customRequestBody))
					if err != nil {
						continue
					}
					// Set a timeout for the request
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					handle2.ProcessPackagesForJson(procedureRequest, pocStruct)
					procedureRequest = procedureRequest.WithContext(ctx)

					// Send request and obtain response results
					procedureResponse, err := client.Do(procedureRequest)
					if err != nil {
						cancel()
						continue
					}

					if IsExploitSuccessByJson(pocStruct, procedureResponse, customRequestBody) {
						if splitURL := strings.Split(tmpUrlForAllRequestPath, "?"); len(splitURL) >= 2 {
							params := strings.Split(splitURL[1], "&")
							encodedParams := make([]string, len(params))
							for tmpI := range params {
								p := strings.Split(params[tmpI], "=")
								encodedParams[tmpI] = url.QueryEscape(p[0]) + "=" + url.QueryEscape(p[1])
							}
							Log.Log.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + strings.Join(encodedParams, "&") + " ]\tSuccess! The target may have this vulnerability")
						} else {
							Log.Log.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + " ]\tSuccess! The target may have this vulnerability")
						}
					}
					cancel()
				}

			}
		}(subURLsList)
	}
	waitGroup.Wait()
}
