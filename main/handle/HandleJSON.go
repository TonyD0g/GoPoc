package handle

import (
	"Scanner/main/format"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var pocStruct format.PocStruct
var InterfaceData map[string]interface{}

func HandleJsonFunc(inputXml string) format.PocStruct {
	jsonData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		fmt.Println("[-] Error reading file:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(jsonData, &InterfaceData)
	err = json.Unmarshal(jsonData, &pocStruct)
	if err != nil {
		fmt.Println("[-] Error unmarshal Json:", err)
		os.Exit(1)
	}
	return pocStruct
}

// TraversePath Traverse the Path in the PathList and initiate a request
func TraversePath(requestPackage format.RequestPackage) []string {
	var allReqPath []string
	var i = 0

	for _, tmpPath := range requestPackage.Uri {
		if !strings.HasPrefix(tmpPath, "/") {
			tmpPath = "/" + tmpPath
		}
		allReqPath = append(allReqPath, requestPackage.Url+tmpPath)
		i++
	}
	return allReqPath
}

// 递归遍历map对象，将指定结构下的键添加到keys变量中
func GetKeys(obj map[string]interface{}, topKey string, targetKey string, keys *[]string) {
	if obj[topKey] != nil {
		if subObj, ok := obj[topKey].(map[string]interface{}); ok {
			if subObj[targetKey] != nil {
				for key := range subObj[targetKey].(map[string]interface{}) {
					*keys = append(*keys, key)
				}
			}
		}
	} else {
		for _, value := range obj {
			switch valueType := value.(type) {
			case map[string]interface{}:
				GetKeys(valueType, topKey, targetKey, keys)
			case []interface{}:
				for _, subValue := range valueType {
					if subMap, ok := subValue.(map[string]interface{}); ok {
						GetKeys(subMap, topKey, targetKey, keys)
					}
				}
			}
		}
	}
}

// 递归遍历map对象，将指定结构下的键值对添加到result变量中
func GetValue(obj map[string]interface{}, topKey string, targetKey string, result *map[string]interface{}) {
	if obj[topKey] != nil {
		if subObj, ok := obj[topKey].(map[string]interface{}); ok {
			if subObj[targetKey] != nil {
				for key, value := range subObj[targetKey].(map[string]interface{}) {
					(*result)[key] = value
				}
			}
		}
	} else {
		for _, value := range obj {
			switch valueType := value.(type) {
			case map[string]interface{}:
				GetValue(valueType, topKey, targetKey, result)
			case []interface{}:
				for _, subValue := range valueType {
					if subMap, ok := subValue.(map[string]interface{}); ok {
						GetValue(subMap, topKey, targetKey, result)
					}
				}
			}
		}
	}
}
