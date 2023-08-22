package test

import (
	"Scanner/main/format"
	"Scanner/main/handle"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

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
			"flag":            false,
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
								if keys[test] == "UserAgent" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.UserAgent == resp.Header.Get("User-Agent") {
										times++
										elements[currentTimes]["User-Agent"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
								} else if keys[test] == "AcceptEncoding" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.AcceptEncoding == resp.Header.Get("Accept-Encoding") {
										times++
										elements[currentTimes]["Accept-Encoding"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
								} else if keys[test] == "Accept" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Accept == resp.Header.Get("Accept") {

										times++
										elements[currentTimes]["Accept"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
								} else if keys[test] == "Cookie" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Cookie == resp.Header.Get("Cookie") {
										times++
										elements[currentTimes]["Cookie"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
								} else if keys[test] == "Host" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.Host == resp.Header.Get("Host") {
										times++
										elements[currentTimes]["Host"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
								} else if keys[test] == "Content-Type" {
									if pocStruct.ResponsePackage.Group[currentTimes].Contain.Header.ContentType == resp.Header.Get("Content-Type") {

										times++
										elements[currentTimes]["Content-Type"] = true
										elements[currentTimes]["flag"] = true
									} else {
										elements[currentTimes]["flag"] = false
									}
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
								elements[currentTimes]["flag"] = true
							} else {
								elements[currentTimes]["flag"] = false
							}
						}
					}
				} else {
					if fieldName == "Status" && strings.Contains(resp.Status, fieldValue.String()) {
						times++
						elements[currentTimes]["Status"] = true
						elements[currentTimes]["flag"] = true
					} else {
						elements[currentTimes]["flag"] = false
					}
				}
			}

		}

		currentTimes++
	}

	var isConform = true
	if pocStruct.ResponsePackage.Operation == "AND" {
		for i := 0; i < currentTimes; i++ {
			if elements[currentTimes]["flag"] == false {
				isConform = false
				break
			}
		}
		if isConform {
			fmt.Println("Success!")
		} else {
			fmt.Println("Error!")
		}
	} else {
		isConform = false
		for i := 0; i < currentTimes; i++ {
			if elements[currentTimes]["flag"] == true {
				isConform = true
				break
			}
		}
		if isConform {
			fmt.Println("Success!")
		} else {
			fmt.Println("Error!")
		}
	}

	return false
}