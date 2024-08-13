package Http

import (
	"GoPoc/main/Developer/AllFormat"
	"context"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func SendHttpRequest(hostInfo string, config SetHttpConfig) (Format.CustomResponseFormat, error) {
	var customResponse Format.CustomResponseFormat
	// Create an HTTP.Request object
	procedureRequest, err := http.NewRequest(config.Method, hostInfo+config.Uri, strings.NewReader(config.Body))
	if err != nil {
		return customResponse, err
	}

	// Set a timeout for the request
	if config.TimeOut <= 0 {
		config.TimeOut = 5
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.TimeOut)

	ProcessPackagesForCode(procedureRequest, config)
	procedureRequest = procedureRequest.WithContext(ctx)
	// Send request and obtain response results
	procedureResponse, err := config.Client.Do(procedureRequest)
	if err != nil {
		cancel()
		return customResponse, err
	}

	// 对自定义Response进行赋值
	bodyOfExecutionResults, err := io.ReadAll(procedureResponse.Body)
	if config.Method == "GET" && err != nil {
		cancel()
		return customResponse, err
	} else if config.Method == "POST" && err == nil {
		cancel()
		return customResponse, err
	}

	customResponse.Body = string(bodyOfExecutionResults)
	customResponse.Request = procedureResponse.Request
	customResponse.Header = procedureResponse.Header
	customResponse.Status = procedureResponse.Status
	customResponse.RawBody = procedureResponse.Body
	customResponse.Close = procedureResponse.Close
	customResponse.ContentLength = procedureResponse.ContentLength
	customResponse.Proto = procedureResponse.Proto
	customResponse.ProtoMajor = procedureResponse.ProtoMajor
	customResponse.ProtoMinor = procedureResponse.ProtoMinor
	customResponse.TLS = procedureResponse.TLS
	customResponse.StatusCode = procedureResponse.StatusCode
	customResponse.Trailer = procedureResponse.Trailer
	customResponse.TransferEncoding = procedureResponse.TransferEncoding
	customResponse.Uncompressed = procedureResponse.Uncompressed

	cancel()
	return customResponse, nil
}

// ProcessPackagesForCode 被用来处理包,包含请求包和响应包. 用处给未写明的header头字段赋予初值
func ProcessPackagesForCode(procedureRequest *http.Request, config SetHttpConfig) {

	var randomUserAgent string
	if config.Header["User-Agent"] != "" {
		randomUserAgent = config.Header["User-Agent"]
	} else {
		// 随机 UA
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
			"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; Hot Lingo 2.0)",
			"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3451.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.155 Safari/537.36 OPR/31.0.1889.174",
			"Mozilla/5.0 (X11; Linux i686) AppleWebKit/535.21 (KHTML, like Gecko) Chrome/19.0.1041.0 Safari/535.21",
			"Mozilla/5.0 (Macintosh; U; PPC Mac OS X; ja-jp) AppleWebKit/418.9.1 (KHTML, like Gecko) Safari/419.3",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
		}
		rand.Seed(time.Now().UnixNano())
		randomUserAgent = userAgents[rand.Intn(len(userAgents))]
	}

	// 对 header 头赋予默认值
	defaultHeaders := map[string]string{
		"User-Agent":      randomUserAgent,
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept":          "*/*",
		"Connection":      "close",
		"Content-Type":    "text/html; charset=UTF-8",
		// 你可以添加更多默认的键值对，当 poc 没写时可以赋值给它
	}

	// 使用 map 来标记是否存在
	isHasExist := make(map[string]bool)

	// 情况1: header字段已设置值则使用设置值,不赋予默认值
	for headerName, headerValue := range config.Header {
		procedureRequest.Header.Add(headerName, headerValue)
	}

	// 情况2: header字段未设置值则赋予默认值.
	for headerName, headerValue := range defaultHeaders {
		if !isHasExist[headerName] {
			procedureRequest.Header.Add(headerName, headerValue)
		}
	}
}
