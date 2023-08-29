package http

import (
	"Scanner/main/Judge"
	"Scanner/main/format"
	"Scanner/main/handle"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// TODO 加入  -r 参数
func Send(pocStruct format.PocStruct, inputUrl string, inputProxy string) {

	client := SetProxy(inputProxy)
	handle.CheckFileCorrectness(pocStruct.RequestPackage)

	customRequestBody := []byte(pocStruct.RequestPackage.Body)

	allReqPath := handle.TraversePath(pocStruct.RequestPackage, inputUrl)
	var waitGroupList sync.WaitGroup
	for _, tmpPath := range allReqPath {
		waitGroupList.Add(1)
		go func(path string, requestBody []byte) {
			// Create an HTTP. Request object
			procedureRequest, err := http.NewRequest(pocStruct.RequestPackage.Method, path, bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Println(err)
				waitGroupList.Done()
				return
			}
			handle.HandlePackFunc(procedureRequest, pocStruct.RequestPackage)

			// Send request and obtain response results
			procedureResponse, err := client.Do(procedureRequest)
			if err != nil {
				fmt.Println(err)
				waitGroupList.Done()
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println("[-] error,please check io.ReadCloser")
				}
			}(procedureResponse.Body)

			if Judge.IsExploitSuccess(pocStruct, procedureResponse, requestBody) {
				fmt.Println("[+] " + path + "\tSuccess! The target may have this vulnerability, please use burpsuite for further investigation")
			}

			waitGroupList.Done()
		}(tmpPath, customRequestBody)
	}

	waitGroupList.Wait()
}
