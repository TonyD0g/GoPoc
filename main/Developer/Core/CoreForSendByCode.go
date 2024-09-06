package Core

import (
	Format "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/Log"
	"GoPoc/main/User"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func ForSendByCode(pocOrExp string, urlsList []string, inputProxy string, maxConcurrentLevel int, isDetectionMode string, pocStruct Format.PocStruct) {
	client := HttpAbout.SetProxy(inputProxy)
	waitGroup := &sync.WaitGroup{}

	// 计算要划分的小的urlsList数量
	numThreads := maxConcurrentLevel
	if numThreads > len(urlsList) {
		numThreads = len(urlsList)
	}
	if numThreads == 0 {
		numThreads = 1
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
				parsedURL, err := url.Parse(tmpUrl)
				if err != nil {
					continue
				}

				if IsExploitSuccessByCode(pocOrExp, tmpUrl, client, isDetectionMode, pocStruct) {
					if splitURL := strings.Split(tmpUrl, "?"); len(splitURL) >= 2 {
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
			}
		}(subURLsList)
	}
	waitGroup.Wait()
}

// IsDetectionMode 开启检测模式
func IsDetectionMode(hostInfo string, client *http.Client, pocStruct Format.PocStruct) bool {
	Log.Log.Println("[+] 已开启探测模式,将根据规则中的Uri发起一个请求,如果不符合自定义规则,则不执行 Poc/Exp 部分")
	conditions, bodyArray, headerArray, bodyCounter, headerCounter := getKeyOperatorValue(pocStruct.Fofa) // 获取所有的键值对以及算数运算符

	// 检测是否存在 app 语法, 有则直接返回 false
	for _, condition := range conditions {
		if strings.ToLower(condition.Key) == "app" {
			Log.Log.Fatal("[-] 由于你开启了探测模式(不借助fofa进行页面规则的扫描),因此app语法不能使用! app 语法本质上由一堆基本语法所构成的,存在fofa服务器中 ")
		}
	}

	// 不借助 fofa 进行页面规则的扫描
	config := HttpAbout.NewHttpConfig()
	config.Uri = pocStruct.Uri
	config.TimeOut = 10
	config.Method = "GET"
	config.Client = client
	resp, err := HttpAbout.SendHttpRequest(hostInfo, config)
	if err != nil {
		Log.Log.Fatal("[-] 解析探测语句失败! :", err)
	}
	if strings.Contains(resp.Body, "<html><head><title>Burp Suite") && strings.Contains(resp.Body, "background: #e06228; padding: 10px 15px") {
		Log.Log.Fatal("[-] 浏览器没开代理插件,比如 firefox 的 FoxyProxy 没开启代理")
		return false
	}

	if len(bodyArray) == 0 || len(headerArray) == 0 {
		return false
	}
	bodyArrayByBool := make([]bool, bodyCounter)
	headerArrayByBool := make([]bool, headerCounter)
	bodyCounter = 0
	headerCounter = 0
	headerValuesList := ""
	if strings.Contains(pocStruct.Fofa, "header") {
		headerValuesList = resp.Proto + " " + resp.Status + "\n"
		for key, values := range resp.Header {
			headerValuesList = headerValuesList + key + ": "
			for _, value := range values {
				headerValuesList = headerValuesList + value + "\n"
			}
		}
	}

	// 扫一遍,得知每个表达式是true还是false
	for _, condition := range conditions {
		if strings.ToLower(condition.Key) == "body" {
			if !strings.Contains(condition.Operator, "!") && strings.Contains(resp.Body, condition.Value) {
				bodyArrayByBool[bodyCounter] = true
				bodyCounter++
			} else if !strings.Contains(resp.Body, condition.Value) {
				bodyArrayByBool[bodyCounter] = true
				bodyCounter++
			}
		} else {
			if !strings.Contains(condition.Operator, "!") && strings.Contains(headerValuesList, condition.Value) {
				headerArrayByBool[headerCounter] = true
				headerCounter++
			} else if !strings.Contains(headerValuesList, condition.Value) {
				headerArrayByBool[headerCounter] = true
				headerCounter++
			}
		}
	}

	transformed := convertLogicalExpression(pocStruct.Fofa)                                // 精简化语句,方便进行逆波兰运算
	reversePolishNotationByStr := reversePolishNotation(transformed)                       // 获取逆波兰语句
	return evaluatePostfix(reversePolishNotationByStr, bodyArrayByBool, headerArrayByBool) // 根据逆波兰语句来计算返回整个表达式的逻辑运算结果
}

func IsExploitSuccessByCode(pocOrExp, hostInfo string, client *http.Client, isDetectionMode string, pocStruct Format.PocStruct) bool {
	if isDetectionMode == "true" {
		if !IsDetectionMode(hostInfo, client, pocStruct) {
			return false
		}
	}

	if pocOrExp != "poc" {
		var expResult Format.ExpResult
		expResult.HostInfo = hostInfo
		expResult.Success = false
		expResult.Output = ""
		expResult = User.Exp(expResult, client)
		if expResult.Success {
			Log.Log.Println(expResult.Output)
			return true
		}
		return false
	}
	return User.Poc(hostInfo, client)
}
