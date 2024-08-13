package main

import (
	format2 "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Core"
	"GoPoc/main/Developer/Handle"
	"GoPoc/main/Developer/Http"
	"GoPoc/main/Log"
	"GoPoc/main/User"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("            ______          \n            | ___ \\         \n  __ _  ___ | |_/ /__   ___ \n / _` |/ _ \\|  __/ _ \\ / __|\n| (_| | (_) | | | (_) | (__ \n \\__, |\\___/\\_|  \\___/ \\___|\n  __/ |                     \n |___/                      ")
	fmt.Println("基于 Json 、自定义Go脚本、fofa的快速验证扫描引擎，可用于快速验证目标是否存在该漏洞。\nVersion 1.5.4")
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

	var userInputDetectionURL string
	if config["url"] != "" {
		userInputDetectionURL = config["url"]
	} else if config["file"] != "" {
		userInputDetectionURL = config["file"]
	} else {
		userInputDetectionURL = ""
	}
	Log.NewLogger().Init()
	maxConcurrentLevelInt, err := strconv.Atoi(config["maxConcurrentLevel"])
	if err != nil {
		Log.Log.Fatal("The maximum concurrency you entered is not a number!", err)
	}
	// 两种 poc模式,第一种为json格式,第二种为代码格式
	var pocStruct format2.PocStruct
	pocModule := Core.LoadPlugin(config["vul"])
	pocStruct = Handle.TryToParsePocStruct(User.Json)

	var urlsList []string
	Log.Log.Println("[+] 扫描开始,记得挂全局socks代理以及burpsuite,burpsuite代理端口取决于config文件所设置的! burpsuite是必须要开启的,否则程序不执行")
	// 发包模式1 基于 fofa 搜索
	if userInputDetectionURL == "" {
		urlsList = Http.SendForFofa(config, pocStruct)
	} else {
		// 发包模式2 基于 单个url / urlFile 文件
		urlsList = Http.SendForUrlOrFile(userInputDetectionURL)
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
	Log.Log.Println("\n[+] 扫描结束,如果什么都没有则说明没有扫出来")
}
