package handle

import (
	"Scanner/main/format"
	"net/http"
	"reflect"
)

func HandlePackFunc(resp *http.Request, requestPackage format.RequestPackage) {
	isHasExist := make(map[string]bool)
	val := reflect.ValueOf(requestPackage.Header)
	typ := val.Type()

	defaultHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Cookie":          "",
		"Accept":          "*/*",
		"connection":      "close",
		// 添加更多的默认头部
		// ...
	}
	defaultHeadersList := map[string]string{
		"User-Agent":      "User-Agent",
		"Accept-Encoding": "Accept-Encoding",
		"Accept-Language": "Accept-Language",
		"Cookie":          "Cookie",
		"Accept":          "Accept",
		"connection":      "connection",
		// 添加更多的默认头部
		// ...
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).String()
		headerName := field.Tag
		headerName = headerName[len(headerName)-(len(headerName)-6):]
		headerName = headerName[0 : len(headerName)-1]
		isHasExist[string(headerName)] = true

		if value != "" {
			resp.Header.Add(string(headerName), value)
		} else if defaultValue, ok := defaultHeaders[string(headerName)]; ok {
			resp.Header.Add(string(headerName), defaultValue)
		}
	}

	for keyA := range isHasExist {
		if _, ok := defaultHeadersList[keyA]; !ok {
			resp.Header.Add(keyA, defaultHeaders[keyA])
		}
	}

}
