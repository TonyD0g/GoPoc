package http

import (
	"Scanner/main/format"
	"Scanner/main/handle"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var requestBody []byte

func Send(pocStruct format.PocStruct, inputProxy string) {

	client := SetProxy(inputProxy)

	if !handle.CheckFileCorrectness(pocStruct.RequestPackage) {
		return
	}

	requestBody := []byte(pocStruct.RequestPackage.Body)

	allReqPath := handle.TraversePath(pocStruct.RequestPackage)
	for _, tmpPath := range allReqPath {
		// 创建http.Request对象
		req, err := http.NewRequest(pocStruct.RequestPackage.Method, tmpPath, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println(err)
			return
		}
		handle.HandlePackFunc(req, pocStruct.RequestPackage)

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

		// TODO 处理response
		isUtilizeSuccess(pocStruct, resp)

		// 读取响应内容
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//
		//// 打印响应结果
		//fmt.Println(string(body))
	}
}

func isUtilizeSuccess(pocStruct format.PocStruct, resp *http.Response) bool {
	var currentTimes = 0
	var times = 0

	// 定义一个[]string类型的变量存储指定结构的所有键
	var keys []string
	// 指定要获取的最上层键
	topKey := "Contain"
	// 指定要获取的上层结构的键
	targetKey := "Header"

	var elements []map[string]bool

	if pocStruct.ResponsePackage.Operation == "AND" {

	} else {

	}

	for _, currentGroup := range pocStruct.ResponsePackage.Group {
		currentGroupReflect := reflect.ValueOf(currentGroup)
		typ := currentGroupReflect.Type()

		elements = append(elements, map[string]bool{
			"User-Agent":      false,
			"Accept-Encoding": false,
			"Accept-Language": false,
			"Cookie":          false,
			"Accept":          false,
			"connection":      false,
			"Content-Type":    false,
			"Body":            false,
			"Status":          false,
		})

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldName := field.Name

			fieldValue := currentGroupReflect.FieldByName(fieldName)
			fieldValueInterface := fieldValue.Interface()

			if fieldValue.String() != "" {
				if containValue, ok := fieldValueInterface.(format.Contain); ok {
					v := reflect.ValueOf(containValue)
					t := v.Type()
					for tmp := 0; tmp < t.NumField(); tmp++ {
						if t.Field(tmp).Name == "Header" {
							// 处理Header部分
							// 递归遍历map对象，将指定结构下的键添加到keys变量中
							handle.GetKeys(handle.InterfaceData, topKey, targetKey, &keys)

							for test := 0; test < len(keys); test++ {
								if keys[test] == "UserAgent" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.UserAgent == resp.Header.Get("User-Agent") {
									times++
									elements[currentTimes]["User-Agent"] = true
								} else if keys[test] == "AcceptEncoding" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.AcceptEncoding == resp.Header.Get("Accept-Encoding") {
									times++
									elements[currentTimes]["Accept-Encoding"] = true
								} else if keys[test] == "Accept" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Accept == resp.Header.Get("Accept") {
									times++
									elements[currentTimes]["Accept"] = true
								} else if keys[test] == "Cookie" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Cookie == resp.Header.Get("Cookie") {
									times++
									elements[currentTimes]["Cookie"] = true
								} else if keys[test] == "Host" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Host == resp.Header.Get("Host") {
									times++
									elements[currentTimes]["Host"] = true
								} else if keys[test] == "Content-Type" && pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.ContentType == resp.Header.Get("Content-Type") {
									times++
									elements[currentTimes]["Content-Type"] = true
								}
							}

						} else {
							body, err := ioutil.ReadAll(resp.Body)
							if len(body) != 0 {
								requestBody = body
							}
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							if ContainsAny(string(requestBody), pocStruct.ResponsePackage.Group[currentTimes].Contain.Body) {
								elements[currentTimes]["Body"] = true
							}
						}
					}
				} else {
					if fieldName == "Status" && strings.Contains(resp.Status, fieldValue.String()) {
						times++
						elements[currentTimes]["Status"] = true
					}
				}
			}

		}

		currentTimes++
	}

	return false
}

func ContainsAny(target string, slice []string) bool {
	for _, s := range slice {
		if strings.Contains(target, s) {
			return true
		}
	}
	return false
}
