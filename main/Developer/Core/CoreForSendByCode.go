package Core

import (
	Format "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Http"
	"GoPoc/main/User"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func ForSendByCode(pocOrExp string, urlsList []string, inputProxy string, maxConcurrentLevel int) {
	client := Http.SetProxy(inputProxy)
	waitGroup := &sync.WaitGroup{}

	// 计算要划分的小的urlsList数量
	numThreads := maxConcurrentLevel
	if numThreads > len(urlsList) {
		numThreads = len(urlsList)
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

				if IsExploitSuccessByCode(pocOrExp, tmpUrl, client) {
					if splitURL := strings.Split(tmpUrl, "?"); len(splitURL) >= 2 {
						params := strings.Split(splitURL[1], "&")
						encodedParams := make([]string, len(params))
						for tmpI := range params {
							p := strings.Split(params[tmpI], "=")
							encodedParams[tmpI] = url.QueryEscape(p[0]) + "=" + url.QueryEscape(p[1])
						}
						fmt.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + strings.Join(encodedParams, "&") + " ]\tSuccess! The target may have this vulnerability")
					} else {
						fmt.Println("[+] [ " + parsedURL.Scheme + "://" + parsedURL.Host + "/" + " ]\tSuccess! The target may have this vulnerability")
					}
				}
			}
		}(subURLsList)
	}
	waitGroup.Wait()
}

func IsExploitSuccessByCode(pocOrExp, hostInfo string, client *http.Client) bool {
	if pocOrExp == "poc" {
		return User.Poc(hostInfo, client)
	} else {
		var expResult Format.ExpResult
		expResult.HostInfo = hostInfo
		expResult.Success = false
		expResult.Output = ""
		expResult = User.Exp(expResult, client)
		if expResult.Success {
			fmt.Println(expResult.Output)
			return true
		}
		return false
	}
}
