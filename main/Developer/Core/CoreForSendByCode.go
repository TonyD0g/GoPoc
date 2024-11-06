package Core

import (
	Format "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/Log"
	"GoPoc/main/User"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

func ForSendByCode(pocOrExp string, urlsList []string, pocStruct Format.PocStruct) {
	if IsDetectionModeByBool {
		Log.Log.Println("[+] 已开启探测模式,将根据poc想要扫描的全部url(支持 Url/File)分别发起一个请求,如果不符合自定义规则,则不执行 Poc/Exp 部分")
	}

	waitGroup := &sync.WaitGroup{}
	numThreads := MaxConcurrentLevelInt // 计算要划分的小的urlsList数量
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
		panicIndex := i
		go func(subURLs []string) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("协程 %d 从报错中恢复,报错原因: %v\n\n", panicIndex, r)
				}
				waitGroup.Add(-1)
			}()

			for _, tmpUrl := range subURLs {
				parsedURL, err := url.Parse(tmpUrl)
				if err != nil {
					continue
				}
				if !IsExploitSuccessByCode(pocOrExp, tmpUrl, pocStruct) {
					continue
				}
				splitURL := strings.Split(tmpUrl, "?")
				if len(splitURL) < 2 {
					Log.Log.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + " ]\tSuccess! The target may have this vulnerability")
				} else {
					params := strings.Split(splitURL[1], "&")
					encodedParams := make([]string, len(params))
					for tmpI := range params {
						p := strings.Split(params[tmpI], "=")
						encodedParams[tmpI] = url.QueryEscape(p[0]) + "=" + url.QueryEscape(p[1])
					}
					Log.Log.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + strings.Join(encodedParams, "&") + " ]\tSuccess! The target may have this vulnerability")
				}
			}
		}(subURLsList)
	}
	waitGroup.Wait()
}

func sendDetectionPackage(hostInfo string, pocStruct Format.PocStruct) (Format.CustomResponseFormat, error) {
	config := HttpAbout.NewHttpConfig()
	config.Uri = pocStruct.Uri
	config.TimeOut = 10
	config.Method = "GET"
	config.IsRedirect = true
	return HttpAbout.SendHttpRequest(hostInfo, config)
}

// IsDetectionMode 开启检测模式,不借助 fofa 进行页面规则的扫描
func IsDetectionMode(hostInfo string, pocStruct Format.PocStruct) bool {
	var bodyArray []string
	var headerArray []string
	var bodyCounter int
	var headerCounter int
	headerValuesList := ""
	resp, err := sendDetectionPackage(hostInfo, pocStruct)
	if err != nil {
		Log.Log.Fatal("[-] 解析探测语句失败! :", err)
	}
	if strings.Contains(resp.Body, "<html><head><title>Burp Suite") && strings.Contains(resp.Body, "background: #e06228; padding: 10px 15px") {
		Log.Log.Fatal("[-] 浏览器没开代理插件,比如 firefox 的 FoxyProxy 没开启代理")
		return false
	}

	if IsBuiltinFingerprint { // 是否使用内置指纹库进行匹配
		var fingerprint Format.Fingerprint
		fingerprintGlobal, _ := os.Open("./main/User/Resource/GobyFingerprint.txt") // test.txt GobyFingerprint.txt
		scanner := bufio.NewScanner(fingerprintGlobal)
		counter := 0
		fingerprintList := ""
		isHaveFingerprint := false
		for scanner.Scan() {
			jsonData := scanner.Text() // 获取当前行的内容
			err = json.Unmarshal([]byte(jsonData), &fingerprint)
			if err != nil {
				log.Fatal(err)
			}
			counter++
			transformed, keyOperatorValue := convertLogicalExpression(fingerprint.Rule)                // 精简化语句,方便进行逆波兰运算
			bodyArray, headerArray, bodyCounter, headerCounter = getKeyOperatorValue(keyOperatorValue) // 获取所有的键值对以及算数运算符
			bodyArrayByBool := make([]bool, bodyCounter)
			headerArrayByBool := make([]bool, headerCounter)
			bodyCounter = 0
			headerCounter = 0
			if len(bodyArray) == 0 && len(headerArray) == 0 {
				return false
			}

			if strings.Contains(transformed, "header") {
				headerValuesList = resp.Proto + " " + resp.Status + "\n"
				for key, values := range resp.Header {
					headerValuesList = headerValuesList + key + ": "
					for _, value := range values {
						headerValuesList = headerValuesList + value + "\n"
					}
				}
			}

			// 扫一遍,得知每个表达式是true还是false
			for _, condition := range keyOperatorValue {
				if condition.Value == "" {
					continue
				}
				switch strings.ToLower(condition.Key) {
				case "title":
				case "body":
					if !strings.Contains(condition.Operator, "!") {
						if strings.Contains(resp.Body, strings.ReplaceAll(condition.Value, "\\", "")) {
							bodyArrayByBool[bodyCounter] = true
						} else {
							bodyArrayByBool[bodyCounter] = false
						}
						bodyCounter++
					} else {
						if !strings.Contains(resp.Body, strings.ReplaceAll(condition.Value, "\\", "")) {
							bodyArrayByBool[bodyCounter] = true
						} else {
							bodyArrayByBool[bodyCounter] = false
						}
						bodyCounter++
					}
					break
				case "header":
					if !strings.Contains(condition.Operator, "!") {
						headerArrayByBool[headerCounter] = false
						if strings.Contains(headerValuesList, strings.ReplaceAll(condition.Value, "\\", "")) {
							headerArrayByBool[headerCounter] = true
						}
						headerCounter++
					} else {
						headerArrayByBool[headerCounter] = false
						if !strings.Contains(headerValuesList, strings.ReplaceAll(condition.Value, "\\", "")) {
							headerArrayByBool[headerCounter] = true
						}
						headerCounter++
					}
				}
			}
			reversePolishNotationByStr := reversePolishNotation(transformed)                     // 获取逆波兰语句
			if evaluatePostfix(reversePolishNotationByStr, bodyArrayByBool, headerArrayByBool) { // 根据逆波兰语句来计算返回整个表达式的逻辑运算结果
				fingerprintList += " " + fingerprint.Product
				isHaveFingerprint = true
				continue
			}
		}
		if isHaveFingerprint {
			Log.Log.Println(fmt.Sprintf("[+] url: %s 匹配到的指纹为: %s", hostInfo, fingerprintList))
		}
	} else { // 代码重复,但是是故意的,为了性能优化,因此不写成一个方法
		transformed, keyOperatorValue := convertLogicalExpression(pocStruct.Fofa)                  // 精简化语句,方便进行逆波兰运算
		bodyArray, headerArray, bodyCounter, headerCounter = getKeyOperatorValue(keyOperatorValue) // 获取所有的键值对以及算数运算符
		bodyArrayByBool := make([]bool, bodyCounter)
		headerArrayByBool := make([]bool, headerCounter)
		bodyCounter = 0
		headerCounter = 0

		if len(bodyArray) == 0 && len(headerArray) == 0 {
			return false
		}
		// 检测是否存在 app 语法, 有则直接返回 false
		for _, condition := range keyOperatorValue {
			if strings.ToLower(condition.Key) == "app" || strings.ToLower(condition.Key) == "product" {
				Log.Log.Fatal("[-] 由于你开启了探测模式(不借助fofa进行页面规则的扫描),因此app/product语法不能使用! 这种语法本质上由一堆基本语法所构成的,存在fofa服务器中")
			}
		}

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
		for _, condition := range keyOperatorValue {
			if condition.Value == "" {
				continue
			}
			switch strings.ToLower(condition.Key) {
			case "body":
				if !strings.Contains(condition.Operator, "!") {
					if strings.Contains(resp.Body, strings.ReplaceAll(condition.Value, "\\", "")) {
						bodyArrayByBool[bodyCounter] = true
					} else {
						bodyArrayByBool[bodyCounter] = false
					}
					bodyCounter++
				} else {
					if !strings.Contains(resp.Body, strings.ReplaceAll(condition.Value, "\\", "")) {
						bodyArrayByBool[bodyCounter] = true
					} else {
						bodyArrayByBool[bodyCounter] = false
					}
					bodyCounter++
				}
				break
			case "header":
				if !strings.Contains(condition.Operator, "!") {
					headerArrayByBool[headerCounter] = false
					if strings.Contains(headerValuesList, strings.ReplaceAll(condition.Value, "\\", "")) {
						headerArrayByBool[headerCounter] = true
					}
					headerCounter++
				} else {
					headerArrayByBool[headerCounter] = false
					if !strings.Contains(headerValuesList, strings.ReplaceAll(condition.Value, "\\", "")) {
						headerArrayByBool[headerCounter] = true
					}
					headerCounter++
				}
			}
		}
		reversePolishNotationByStr := reversePolishNotation(transformed) // 获取逆波兰语句

		if evaluatePostfix(reversePolishNotationByStr, bodyArrayByBool, headerArrayByBool) { // 根据逆波兰语句来计算返回整个表达式的逻辑运算结果
			Log.Log.Println(fmt.Sprintf("[+] url: %s 匹配到的指纹为: %s", hostInfo, pocStruct.Fofa))
			return true
		}
	}
	return false
}

func IsExploitSuccessByCode(pocOrExp, hostInfo string, pocStruct Format.PocStruct) bool {
	if IsDetectionModeByBool && !IsDetectionMode(hostInfo, pocStruct) {
		return false
	}

	if pocOrExp != "poc" {
		var expResult Format.ExpResult
		expResult.HostInfo = hostInfo
		expResult.Success = false
		expResult.Output = ""
		expResult = User.Exp(expResult)
		if expResult.Success {
			Log.Log.Println(expResult.Output)
			return true
		}
		return false
	}
	return User.Poc(hostInfo)
}
