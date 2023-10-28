package Http

import (
	"Scanner/main/Developer/Fofa"
	"Scanner/main/Developer/Format"
	handle2 "Scanner/main/Developer/Handle"
	"Scanner/main/Developer/Judge"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SendForFofa(config map[string]string, pocStruct Format.PocStruct) []string {
	maxFofaSizeInt, err := strconv.Atoi(config["maxFofaSize"])
	if err != nil {
		fmt.Printf("maxFofaSize 并不是一个有效数字\n")
		os.Exit(1)
	}

	var urlsList []string
	var queryResponse Fofa.QueryResponse
	err = json.Unmarshal(Fofa.SearchReturnByte(config, pocStruct, maxFofaSizeInt), &queryResponse)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		os.Exit(1)
	}

	for _, tmpOutcome := range queryResponse.Results {
		if !strings.HasPrefix(tmpOutcome[1].(string), tmpOutcome[0].(string)) {
			urlsList = append(urlsList, tmpOutcome[0].(string)+"://"+tmpOutcome[1].(string))
		} else {
			urlsList = append(urlsList, tmpOutcome[1].(string))
		}
	}
	fmt.Printf("[+] 此 fofa 语句: %v 查询到: %v 条", queryResponse.Query, queryResponse.Size)
	return urlsList
}

func SendForUrlOrFile(userInputDetectionURL string) []string {
	var urlsList []string
	if strings.HasPrefix(strings.ToLower(userInputDetectionURL), "http://") || strings.HasPrefix(strings.ToLower(userInputDetectionURL), "https://") {
		urlsList = append(urlsList, userInputDetectionURL)
	} else {
		// urlFile list
		urlFile, err := os.Open(userInputDetectionURL)
		if err != nil {
			fmt.Println("can't open file:", err)
			os.Exit(1)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				fmt.Println("can't close file:", err)
				os.Exit(1)
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
	return urlsList
}

func CoreForSend(urlsList []string, pocStruct Format.PocStruct, inputProxy string, maxConcurrentLevel int) {
	client := SetProxy(inputProxy)
	customRequestBody := []byte(pocStruct.RequestPackage.Body)
	waitGroup := &sync.WaitGroup{}
	processedURLs := make(map[string]struct{}) // 用于存储已处理的URL

	// 计算要划分的小的urlsList数量
	numThreads := maxConcurrentLevel
	if numThreads > len(urlsList) {
		numThreads = len(urlsList)
	}
	urlsPerThread := len(urlsList) / numThreads
	var processedURLsMutex sync.Mutex
	for i := 0; i < numThreads; i++ {
		start := i * urlsPerThread
		end := start + urlsPerThread
		if i == numThreads-1 {
			end = len(urlsList)
		}

		subURLsList := urlsList[start:end]
		waitGroup.Add(1)
		go func(subURLs []string) {
			defer waitGroup.Add(-1)

			for _, tmpUrl := range subURLs {
				allRequestPath := handle2.TraversePath(pocStruct.RequestPackage, tmpUrl)
				requestCount := len(allRequestPath)

				for tmpI := 0; tmpI < requestCount; tmpI++ {
					tmpUrlForAllRequestPath := allRequestPath[tmpI]

					processedURLsMutex.Lock() // 使用互斥锁保护对 processedURLs 的访问
					// 检查URL是否已经处理过，如果处理过则跳过
					if _, exists := processedURLs[tmpUrlForAllRequestPath]; exists {
						processedURLsMutex.Unlock()
						continue
					}
					processedURLsMutex.Unlock()                         // 解锁 processedURLs
					processedURLs[tmpUrlForAllRequestPath] = struct{}{} // 标记URL已处理

					parsedURL, err := url.Parse(tmpUrlForAllRequestPath)
					if err != nil {
						continue
					}
					// Create an HTTP.Request object
					procedureRequest, err := http.NewRequest(pocStruct.RequestPackage.Method, tmpUrlForAllRequestPath, bytes.NewBuffer(customRequestBody))
					if err != nil {
						continue
					}
					handle2.ProcessPackages(procedureRequest, pocStruct)

					// Set a timeout for the request
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
					defer cancel()
					procedureRequest = procedureRequest.WithContext(ctx)

					// Send request and obtain response results
					procedureResponse, err := client.Do(procedureRequest)
					if err != nil {
						continue
					}
					defer func(Body io.ReadCloser) {
						if err = Body.Close(); err != nil {
							return
						}
					}(procedureResponse.Body)

					if Judge.IsExploitSuccess(pocStruct, procedureResponse, customRequestBody) {
						if splitURL := strings.Split(tmpUrlForAllRequestPath, "?"); len(splitURL) >= 2 {
							params := strings.Split(splitURL[1], "&")
							encodedParams := make([]string, len(params))
							for tmpI := range params {
								p := strings.Split(params[tmpI], "=")
								encodedParams[tmpI] = url.QueryEscape(p[0]) + "=" + url.QueryEscape(p[1])
							}
							fmt.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + strings.Join(encodedParams, "&") + " ]\tSuccess! The target may have this vulnerability")
						} else {
							fmt.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + " ]\tSuccess! The target may have this vulnerability")
						}
					}

				}

			}
		}(subURLsList)
	}
	waitGroup.Wait()
}
