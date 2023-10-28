package main

import (
	"Scanner/main/Developer/Handle"
	"Scanner/main/Developer/Http"
	"Scanner/main/Developer/Input"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	if args[1] != "-ini" {
		flag.Usage()
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
	pocStruct := Handle.TryToParsePocStruct(config["pocJson"])

	var urlsList []string
	fmt.Println("[+] 扫描开始,记得挂全局socks代理! :")
	// 模式1 基于 fofa 搜索
	if userInputDetectionURL == "" {
		urlsList = Http.SendForFofa(config, pocStruct)
	} else {
		// 模式2 基于 单个url / urlFile 文件
		urlsList = Http.SendForUrlOrFile(userInputDetectionURL)
	}

	Http.CoreForSend(urlsList, pocStruct, config["proxy"], maxConcurrentLevelInt)
	fmt.Println("\n[+] 扫描结束")
}
