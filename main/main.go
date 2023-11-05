package main

import (
	format2 "GoPoc/main/Developer/AllFormat"
	Handle "GoPoc/main/Developer/Handle"
	"GoPoc/main/Developer/Http"
	"GoPoc/main/Developer/Input"
	"GoPoc/main/Developer/LoadingGo"
	"GoPoc/main/User"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("            ______          \n            | ___ \\         \n  __ _  ___ | |_/ /__   ___ \n / _` |/ _ \\|  __/ _ \\ / __|\n| (_| | (_) | | | (_) | (__ \n \\__, |\\___/\\_|  \\___/ \\___|\n  __/ |                     \n |___/                      ")
	fmt.Println("Version 1.3")
	args := os.Args
	if len(args) == 1 {
		fmt.Println("使用说明:	-ini C:/config.ini")
	} else if args[1] != "-ini" {
		os.Exit(1)
	}
	inputIniFile := flag.String("ini", ".\\config.ini", "Input the ini file")
	flag.Parse()
	config := Input.HandleIni(*inputIniFile)

	// Determine whether the number of parameters is correct
	if !strings.Contains(config["email"], "@") || config["key"] == "" || config["maxFofaSize"] == "" {
		fmt.Println("[-] 参数错误,例子:\n-email\nxxxxx@qq.com\n-key\nxxxxx\n-url\nHttp://127.0.0.1/\n-pocJson\nC:\\Users\\xxx\\Desktop\\1.json\n-proxy\nHttp://127.0.0.1:8082")
		flag.Usage()
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

	maxConcurrentLevelInt, err := strconv.Atoi(config["maxConcurrentLevel"])
	if err != nil {
		fmt.Println("The maximum concurrency you entered is not a number!", err)
	}
	// 两种 poc模式,第一种为json格式,第二种为代码格式
	var pocStruct format2.PocStruct
	pocModule := LoadingGo.LoadPlugin(config["poc"])
	if pocModule == 1 {
		pocStruct = Handle.TryToParsePocStruct(User.Json)
	}

	var urlsList []string
	fmt.Println("[+] 扫描开始,记得挂全局socks代理! :")
	// 发包模式1 基于 fofa 搜索
	if userInputDetectionURL == "" {
		urlsList = Http.SendForFofa(config, pocStruct)
	} else {
		// 发包模式2 基于 单个url / urlFile 文件
		urlsList = Http.SendForUrlOrFile(userInputDetectionURL)
	}
	if pocModule == 1 {
		Http.CoreForSendByJson(urlsList, pocStruct, config["proxy"], maxConcurrentLevelInt)
	} else {
		Http.CoreForSendByCode("poc", urlsList, config["proxy"], maxConcurrentLevelInt)
	}
	fmt.Println("\n[+] 扫描结束")
}
