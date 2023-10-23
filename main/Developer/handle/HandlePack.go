package handle

import (
	"Scanner/main/Developer/format"
	"net/http"
	"reflect"
)

// ProcessPackages Used to process packages
func ProcessPackages(procedureResponse *http.Request, pocStruct format.PocStruct) {
	isHasExist := make(map[string]bool)
	reflectValue := reflect.ValueOf(pocStruct.RequestPackage.Header)
	reflectValueType := reflectValue.Type()

	// todo 增加随机 UA
	defaultHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Cookie":          "",
		"Accept":          "*/*",
		"connection":      "close",
		"Content-Type":    "text/html; charset=UTF-8",
		// you can add more default headers
		// ...
	}
	defaultHeadersList := map[string]string{
		"User-Agent":      "User-Agent",
		"Accept-Encoding": "Accept-Encoding",
		"Accept-Language": "Accept-Language",
		"Cookie":          "Cookie",
		"Accept":          "Accept",
		"connection":      "connection",
		"Content-Type":    "text/html; charset=UTF-8",
		// you can add more default headers
		// ...
	}

	for index := 0; index < reflectValueType.NumField(); index++ {
		field := reflectValueType.Field(index)
		value := reflectValue.Field(index).String()
		headerName := field.Tag
		headerName = headerName[len(headerName)-(len(headerName)-6):]
		headerName = headerName[0 : len(headerName)-1]
		isHasExist[string(headerName)] = true

		if value != "" {
			procedureResponse.Header.Add(string(headerName), value)
		} else if defaultValue, ok := defaultHeaders[string(headerName)]; ok {
			procedureResponse.Header.Add(string(headerName), defaultValue)
		}
	}

	for index := range isHasExist {
		if _, ok := defaultHeadersList[index]; !ok {
			procedureResponse.Header.Add(index, defaultHeaders[index])
		}
	}
}
