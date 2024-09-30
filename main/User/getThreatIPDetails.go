package User

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/User/Utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var Json string
var Poc func(hostInfo string) bool
var Exp func(expResult Format.ExpResult) Format.ExpResult

// todo 输出 漏洞详情翻译的内容
func init() {
	Json = `{
   "CheckIP":"false",
   "Coroutine":"200",
   "File":"",
   "Url":"https://portal.tophant.com",
   "Fofa":"",
 	"Uri":"",
   "Request":{

		"Method": "",

		"Uri": [
	   			],

		"Header":{
			"Accept-Encoding":"gzip"
		}
	},

   "Response":{

		"Operation":"",

		"Group":[

				{

                    "Regexp": "",
			        "Header":{
			            },

			        "Body":[
			            ]
			    },

          		 {
			        "Header":{
			            }
			    }
			]

}
}`

	// 获取漏洞列表
	getVulnList455445 := func(hostInfo string, page int, portalToken string) (Format.CustomResponseFormat, error) {
		config := HttpAbout.NewHttpConfig()
		Utils.FullyAutomaticFillingHeader(config, `POST /api/vip/vuln/base-list HTTP/1.1
Host: portal.tophant.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0
Accept: application/json, text/plain, */*
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Content-Type: application/json; charset=UTF-8
Portal-Token: eyJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE3MjQ4MTI0OTIsImxvZ2luX3VzZXJfa2V5IjoiMTcyMjg1NDQ2ODMzMjcwNTUifQ.WkGSxta9uWZ89Nc-I4vAQnnkeC0HDcLifROroaGrCnX9JPp_E4s5XvS2GKR64cLYSwclB0uPnaBdmpdDXRx0Cw
Content-Length: 47
Origin: https://portal.tophant.com
Referer: https://portal.tophant.com/portal/vipManage/vip/vulnManage
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
Te: trailers
Connection: close
`)
		config.Method = "POST"
		config.Body = `{
		   "lastUpdateTimeSort":0,
			"publishTimeSort":null,
			"manufacturerPublishSort":null,
			"modifyTimeSort":null,
		   "page": ` + strconv.Itoa(page) + `,
		   "pageSize": 40 
		}`
		config.Header["Portal-Token"] = portalToken
		fmt.Println("[+] 获取的是第 " + strconv.Itoa(page) + " 页, pageSize 为 40 ")
		config.Uri = "/api/vip/vuln/base-list"
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	findCsvByNvd3248329 := func(match []string) (bool, string, string) {
		config := HttpAbout.NewHttpConfig()
		config.TimeOut = 10
		re := regexp.MustCompile(`CVE-\d{4}-\d+`)
		cveMatch := re.FindString(match[1])

		resp, err := HttpAbout.SendHttpRequest(match[1], config)
		if err != nil {
			return false, "", ""
		}
		if strings.Contains(resp.Body, "Linux kernel") || !strings.Contains(resp.Body, "CVSS:3.1") { // || strings.Contains(resp.Body, "CVE ID Not Found")
			return false, "", ""
		}
		re = regexp.MustCompile(`class="tooltipCvss3CnaMetrics">(.*?)</span>`) // 没有 Nvd 链接跳转则直接下一个
		matches := re.FindStringSubmatch(resp.Body)
		if len(matches) < 2 {
			return false, "", ""
		}
		return true, matches[1], cveMatch
	}

	// 获取漏洞细节
	getVulnDetails5251552 := func(hostInfo string, ids []string, fixPath, isHavaCweNameCn, isSelectVulnName bool, portalToken string) (string, int) {
		numForVuln := 0

		config := HttpAbout.NewHttpConfig()
		config.Header["Accept"] = "application/json, text/plain, */*"
		config.Header["Accept-Encoding"] = "gzip, deflate, br"
		config.Header["Content-Type"] = "application/json; charset=UTF-8"
		config.Header["Origin"] = "https://portal.tophant.com"
		config.Header["Referer"] = "https://portal.tophant.com/portal/vipManage/vip/vulnManage"
		config.Header["Accept-Language"] = "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"
		config.Header["Portal-Token"] = portalToken
		config.Header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
		config.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0" // (非强制) 如果不写的话随机从默认URL列表上选取一个
		config.TimeOut = 10                                                                                              // (非强制) 如果不写的话默认值为 5秒

		urlList := ""
		for _, vulnId := range ids {
			realUri := "/portal/vipManage/vip/vulnManage/edit?vulnId=" + vulnId + "&redirectName=vulnManage"
			config.Uri = "/api/vip/vuln/edit/" + vulnId
			resp, err := HttpAbout.SendHttpRequest(hostInfo, config)
			if err != nil {
				log.Fatal("[-] 服务器响应错误,要么是你请求的速度太快,要么是 session 失效 !")
			}

			if isSelectVulnName {
				reForName := regexp.MustCompile(`"vulnName":"([^"]+)"`) // 漏洞标签简略则直接跳过
				matchForName := reForName.FindStringSubmatch(resp.Body)
				if strings.HasPrefix(strings.ToLower(matchForName[1]), "cve-20") {
					continue
				}
			}

			// todo bug
			fixTypePatchExists := strings.Contains(resp.Body, "\"fixType\":\"patch\"")
			if fixPath && fixTypePatchExists {
				continue
			}

			re := regexp.MustCompile(`"linkUrl":"([^"]+)"`) // 无法点击Nvd标签则直接下一个
			matches := re.FindAllStringSubmatch(resp.Body, -1)
			var nvdMatch []string
			isHaveNvd := false
			for _, match := range matches {
				if strings.Contains(match[1], "nvd.nist.gov/vuln/detail") {
					nvdMatch = match
					isHaveNvd = true
					break
				}
			}
			if !isHaveNvd {
				continue
			}

			// 如果被审核过了,直接pass掉
			if strings.Contains(resp.Body, "\"vulnAuthStatus\":\"examining\"") {
				continue
			}
			if strings.Contains(resp.Body, "\"vulnAuthStatus\":\"released\"") {
				continue
			}

			if isHavaCweNameCn {
				re = regexp.MustCompile(`"cweNameCn":"([^"]+)"`) // 没有 漏洞类型 则直接下一个
				matches := re.FindStringSubmatch(resp.Body)
				if len(matches) < 1 {
					continue
				}
			}

			isHaveNvd, nvdValue, cveValue := findCsvByNvd3248329(nvdMatch)
			if !isHaveNvd {
				continue
			}
			numForVuln++
			urlList = urlList + nvdValue + "  " + cveValue + "\n" + hostInfo + realUri + "\n"
		}

		return urlList, numForVuln
	}

	// 返回收集到的所有漏洞列表
	returnIdsList := func(resp Format.CustomResponseFormat) ([]string, bool) {
		re := regexp.MustCompile(`"id"\s*:\s*"(\d+)"`)
		matches := re.FindStringSubmatch(resp.Body)
		if len(matches) < 1 {
			return nil, false
		}
		matchesList := re.FindAllStringSubmatch(resp.Body, -1)
		var ids []string
		for _, match := range matchesList {
			if len(match) > 1 {
				// 添加匹配的 ID 到 ids 列表中
				ids = append(ids, match[1])
			}
		}
		return ids, true
	}

	// 如果使用代码模式, Poc函数为必须,其中的参数固定
	Poc = func(hostInfo string) bool {
		portalToken := `eyJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE3Mjc0MTIyNTYsImxvZ2luX3VzZXJfa2V5IjoiMTcyMjg1NDQ2ODMzMjcwNTUifQ.SwcQ3ecu7GHC3gr9vqOQWGxzC8r7mFEYyg1_Q6YtQpnwx-tWx1fHKbfNz-xHAOYEGRjgAzoj0qD2HP-jbkh-JA`
		isChoiceFixPath := true  // 是否有修复方案
		isHavaCweNameCn := true  // 是否有漏洞标签
		isSelectVulnName := true // 是否选择对漏洞名称简略的进行跳过

		resp, err := getVulnList455445(hostInfo, 1, portalToken)
		if err != nil {
			return false
		}
		if !strings.Contains(resp.Status, "200") {
			return false
		}
		re := regexp.MustCompile(`"total":"(\d+)"`)
		matches := re.FindStringSubmatch(resp.Body)
		var totalPage int
		intValue, _ := strconv.Atoi(matches[1])
		totalPage = intValue / 40
		fmt.Println("[+] 页数共有: " + strconv.Itoa(totalPage))
		for i := 1; i < totalPage; i++ {
			resp, _ := getVulnList455445(hostInfo, i, portalToken)
			ids, isHaveReturn := returnIdsList(resp)
			if !isHaveReturn {
				continue
			}

			urlList, numForVuln := getVulnDetails5251552(hostInfo, ids, isChoiceFixPath, isHavaCweNameCn, isSelectVulnName, portalToken)
			if len(urlList) == 0 {
				continue
			}
			fmt.Println("[+] (" + strconv.FormatBool(isHavaCweNameCn) + ") 漏洞标签 (" + strconv.FormatBool(isChoiceFixPath) + ") 修复方案的漏洞链接为: \n" + urlList)
			fmt.Println("目前为止输出了: " + strconv.Itoa(numForVuln) + " 个")
		}
		return true
	}

	// 如果使用代码模式, Exp函数为必须,其中的参数固定
	// Exp 你可以尝试自己写一下:
	Exp = func(expResult Format.ExpResult) Format.ExpResult {
		return expResult
	}
}
