package main

import (
	format2 "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Core"
	"GoPoc/main/Developer/Handle"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/Log"
	"GoPoc/main/User"
	"GoPoc/main/User/Utils"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func selectModule(config map[string]string, timeSet time.Duration) {
	// 两种 poc模式,第一种为json格式,第二种为代码格式
	var pocStruct format2.PocStruct
	pocModule := Core.LoadPlugin(config["vul"])
	pocStruct = Handle.TryToParsePocStruct(User.Json)
	var userInputDetectionURL string
	if pocStruct.Url != "" {
		userInputDetectionURL = pocStruct.Url
	} else if pocStruct.File != "" {
		userInputDetectionURL = pocStruct.File
	} else {
		userInputDetectionURL = ""
	}
	Log.NewLogger().Init()
	maxConcurrentLevelInt := 0
	var err error
	if pocStruct.Coroutine != "" {
		maxConcurrentLevelInt, err = strconv.Atoi(pocStruct.Coroutine)
	} else {
		maxConcurrentLevelInt, err = strconv.Atoi(config["maxConcurrentLevel"])
	}
	if err != nil {
		maxConcurrentLevelInt = 200
	}
	if maxConcurrentLevelInt < 1 {
		maxConcurrentLevelInt = 200
	}

	var urlsList []string
	detectionBurpIsOpen(config["proxy"]) // 检测是否开启了burp,如果没有开启则输出“没有开启” 且直接返回
	Log.Log.Println("[+] 加载的脚本为: " + config["vul"])
	Log.Log.Println("[+] 开启的协程数为: " + strconv.Itoa(maxConcurrentLevelInt))
	ipAddressList := getIpAddress(config["proxy"]) // 检测是否开启代理
	if len(ipAddressList) < 3 {
		ipAddressList = append(ipAddressList, "null")
	}
	Log.Log.Printf("[+] 国内发包IP地址为: %s \t国外发包IP地址为: %s \t谷歌访问测试: %s ", ipAddressList[0], ipAddressList[1], ipAddressList[2])
	Log.Log.Println("[+] 【重要,必看】请确认各项IP地址无误,如扫描国内则看\"国内发包IP地址\"是否为代理IP地址,其他以此类推,不然被溯源了!")
	Log.Log.Println("[+] 将在5秒后自动开始扫描")
	time.Sleep(timeSet * time.Second) // 休眠5秒
	Log.Log.Println("[+] 扫描开始:")
	// 发包模式1 基于 fofa 搜索
	if userInputDetectionURL == "" {
		urlsList = HttpAbout.SendForFofa(config, pocStruct)
	} else {
		// 发包模式2 基于 单个url / urlFile 文件
		urlsList = HttpAbout.SendForUrlOrFile(userInputDetectionURL)
	}
	if Core.CheckBalanced(pocStruct.Fofa) {
		Log.Log.Fatal("[-] 请检测fofa语句是否正确,着重检查括号是否正确闭合")
	}

	if pocModule == 1 {
		Core.ForSendByJson(urlsList, pocStruct, config["proxy"], maxConcurrentLevelInt)
	} else if config["mod"] == "poc" || config["mod"] == "" {
		Core.ForSendByCode("poc", urlsList, config["proxy"], maxConcurrentLevelInt, config["detectionMode"], pocStruct)
	} else {
		Core.ForSendByCode("exp", urlsList, config["proxy"], maxConcurrentLevelInt, config["detectionMode"], pocStruct)
	}
}

func parseConfigIni() map[string]string {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("使用说明:	-ini C:/config.ini\nconfig.ini内容如下:\n\n-email // fofa的email (必须)\n-key // fofa的key (必须)\n-url // 扫单个url (非必须)\n-file // 扫url文件中的每一个url (非必须)\n-vul // poc/exp文件,文件后缀为.go (必须)\n-mod // 指定poc/exp这两种模式 (必须)\n-proxy // burpsuite 代理,用于方便写poc/exp (必须)\n-maxConcurrentLevel // 最大并发量,越大扫描速度越快 (必须)\n-maxFofaSize\t   // fofa最大检索数 (必须)")
	} else if args[1] != "-ini" {
		fmt.Println("[-] 参数错误,例子:\n-email // fofa的email (必须)\n-key // fofa的key (必须)\n-url // 扫单个url (非必须)\n-file // 扫url文件中的每一个url (非必须)\n-vul // poc/exp文件,文件后缀为.go (必须)\n-mod // 指定poc/exp这两种模式 (必须)\n-proxy // burpsuite 代理,用于方便写poc/exp (必须)\n-maxConcurrentLevel // 最大并发量,越大扫描速度越快 (必须)\n-maxFofaSize\t   // fofa最大检索数 (必须)")
		os.Exit(1)
	}
	inputIniFile := flag.String("ini", ".\\config.ini", "Input the ini file")
	flag.Parse()
	config := Handle.HandleIni(*inputIniFile)

	// Determine whether the number of parameters is correct
	if !strings.Contains(config["email"], "@") || config["key"] == "" || config["maxFofaSize"] == "" {
		fmt.Println("[-] 参数错误,例子:-email // fofa的email (必须)\n-key // fofa的key (必须)\n-url // 扫单个url (非必须)\n-file // 扫url文件中的每一个url (非必须)\n-vul // poc/exp文件,文件后缀为.go (必须)\n-mod // 指定poc/exp这两种模式 (必须)\n-proxy // burpsuite 代理,用于方便写poc/exp (必须)\n-maxConcurrentLevel // 最大并发量,越大扫描速度越快 (必须)\n-maxFofaSize\t   // fofa最大检索数 (必须)")
		os.Exit(1)
	}
	return config
}

func detectionBurpIsOpen(proxyAddress string) {
	matches := regexp.MustCompile(`:(\d+)`).FindStringSubmatch(proxyAddress)
	if len(matches) < 1 {
		Log.Log.Println("[-] 输入的 config 没有找到对应的 -proxy 参数!")
		os.Exit(1)
	}
	conn, err := net.DialTimeout("tcp", strings.Replace(strings.ToLower(proxyAddress), "http://", "", -1), 10*time.Second)
	if err != nil {
		Log.Log.Println(`[-] Burpsuite是必开项,但你未开启! 你应该在此处: "Burpsuite --> Proxy --> Options --> listeners" 添加或开启一个端口为: ` + matches[1]) // matches[1] 是捕获的端口号部分
		os.Exit(1)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("[-] 关闭探测 Burpsuite 连接出错")
			os.Exit(1)
		}
	}(conn) // 关闭连接
}

func getIpAddress(proxy string) []string {
	var ipAddressList []string
	config := HttpAbout.NewHttpConfig()
	Utils.FullyAutomaticFillingHeader(config, `GET / HTTP/1.1
Host: www.ip111.cn
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Priority: u=0, i
Te: trailers
Connection: close`)
	config.TimeOut = 60
	client := HttpAbout.SetProxy(proxy)
	config.Client = client
	resp, err := HttpAbout.SendHttpRequest(`Http://www.ip111.cn`, config)
	if err != nil {
		log.Fatal("[-] 获取ip地址失误,请检查你的网络是否正常或者 Http://www.ip111.cn 这个链接是否挂了")
	}

	matches := regexp.MustCompile(`<p>\s*(.*?)\s*</p>`).FindStringSubmatch(resp.Body)

	if len(matches) > 1 { // 国内测试
		ipAddressList = append(ipAddressList, matches[1])
	}

	Utils.FullyAutomaticFillingHeader(config, `GET /ip.php HTTP/1.1
Host: sspanel.net
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Referer: https://www.ip111.cn/
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: iframe
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: cross-site
Priority: u=4
Te: trailers
Connection: close
`)
	resp, err = HttpAbout.SendHttpRequest(`Http://sspanel.net/ip.php`, config)
	if err != nil {
		fmt.Println("[-] 获取ip地址失误,请检查你的网络是否正常或者 Http://sspanel.net/ip.php 这个链接是否挂了")
	} else {
		matches = regexp.MustCompile(`(\d{1,3}(?:\.\d{1,3}){3})\s+(\p{Han}+)\s+(\p{Han}+)`).FindStringSubmatch(resp.Body)
	}

	if len(matches) > 1 { // 国外测试
		ipAddressList = append(ipAddressList, matches[0])
	}
	Utils.FullyAutomaticFillingHeader(config, `GET /ip.php HTTP/1.1
Host: us.ip111.cn
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Referer: https://www.ip111.cn/
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: iframe
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: same-site
Priority: u=4
Te: trailers
Connection: close
`)

	resp, err = HttpAbout.SendHttpRequest(`Http://us.ip111.cn/ip.php`, config)
	if err != nil {
		log.Fatal("[-] 获取ip地址失误,请检查你的网络是否正常或者 Http://us.ip111.cn/ip.php 这个链接是否挂了")
	}
	matches = regexp.MustCompile(`(\d{1,3}(?:\.\d{1,3}){3})\s+(\p{Han}+)\s+(\p{Han}+)`).FindStringSubmatch(resp.Body)

	if len(matches) > 1 { // 从谷歌测试
		ipAddressList = append(ipAddressList, matches[0])
	}

	return ipAddressList
}

func main() {
	fmt.Println("            ______          \n            | ___ \\         \n  __ _  ___ | |_/ /__   ___ \n / _` |/ _ \\|  __/ _ \\ / __|\n| (_| | (_) | | | (_) | (__ \n \\__, |\\___/\\_|  \\___/ \\___|\n  __/ |                     \n |___/                      ")
	fmt.Println("基于 Json 、自定义Go脚本、fofa的快速验证扫描引擎，可用于快速验证目标是否存在该漏洞或者帮助你优化工作流程	-- TonyD0g")
	fmt.Println("Version 1.5.6")
	config := parseConfigIni()
	selectModule(config, 5)
	Log.Log.Println("\n[+] 扫描结束,如果什么输出链接则说明没有扫出来")
}
