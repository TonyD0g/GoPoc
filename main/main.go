package main

import (
	"Scanner/main/Developer/Fofa"
	"Scanner/main/Developer/Handle"
	"Scanner/main/Developer/Http"
	"Scanner/main/Developer/Input"
	"encoding/json"
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
	if !strings.Contains(config["email"], "@") || config["key"] == "" {
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
	var queryResponse Fofa.QueryResponse
	if userInputDetectionURL == "" {
		err = json.Unmarshal(Fofa.Query(config["email"], config["key"], pocStruct.Fofa, 6), &queryResponse)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
			os.Exit(1)
		}
	}

	fmt.Println("[+] 扫描开始,记得挂全局socks代理! :")
	Http.Send(pocStruct, queryResponse, userInputDetectionURL, config["proxy"], maxConcurrentLevelInt)
	fmt.Println("\n[+] 扫描结束")
}
