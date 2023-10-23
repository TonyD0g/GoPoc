package http

import (
	"Scanner/main/Developer/Judge"
	"Scanner/main/Developer/format"
	handle2 "Scanner/main/Developer/handle"
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

func Send(pocStruct format.PocStruct, inputUrl, inputProxy string, maxConcurrentLevel int) {
	client := SetProxy(inputProxy)
	customRequestBody := []byte(pocStruct.RequestPackage.Body)
	var urlsList []string
	if strings.HasPrefix(inputUrl, "http://") || strings.HasPrefix(inputUrl, "https://") {
		urlsList = append(urlsList, inputUrl)
	} else {
		// urlFile list
		urlFile, err := os.Open(inputUrl)
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

	var waitGroup sync.WaitGroup
	for _, tmpUrl := range urlsList {
		allRequestPath := handle2.TraversePath(pocStruct.RequestPackage, tmpUrl)
		requestCount := len(allRequestPath)
		for i := 0; i < requestCount; i++ {
			waitGroup.Add(1)
			go func(i int) {
				defer waitGroup.Done()
				if sendRequest(pocStruct, client, allRequestPath[i], customRequestBody) {
					return
				}
			}(i)

			if (i+1)%maxConcurrentLevel == 0 || i == requestCount-1 {
				waitGroup.Wait()
			}
		}
	}
}

func sendRequest(pocStruct format.PocStruct, client *http.Client, urlAndUri string, customRequestBody []byte) bool {
	_, err := url.Parse(urlAndUri)
	if err != nil {
		fmt.Println("URL parsing error:", err)
		return false
	}
	// Create an HTTP.Request object
	procedureRequest, err := http.NewRequest(pocStruct.RequestPackage.Method, urlAndUri, bytes.NewBuffer(customRequestBody))
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("[-] error, please check io.ReadCloser")
		}
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
