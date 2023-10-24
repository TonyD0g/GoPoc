package main

import (
	"Scanner/main/Developer/Fofa"
	"Scanner/main/Developer/Handle"
	"Scanner/main/Developer/Http"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Get command line parameters
	args := os.Args
	// Determine whether the number of parameters is correct
	if args[1] != "-email" || args[3] != "-key" {
		fmt.Println("[-] Input parameter error, case as follows:\n-email\nxxxxx@qq.com\n-key\nxxxxx\n-url\nHttp://127.0.0.1/\n-pocJson\nC:\\Users\\xxx\\Desktop\\1.json\n-proxy\nHttp://127.0.0.1:8082")
		flag.Usage()
		os.Exit(1)
	}

	var userInputDetectionURL *string
	if args[5] == "-url" {
		userInputDetectionURL = flag.String("url", "", "Input the url")
	} else if args[5] == "-file" {
		userInputDetectionURL = flag.String("file", "", "Input the url file list")
	} else {
		userInputDetectionURL = nil
	}
	maxConcurrentLevel := flag.String("maxConcurrentLevel", "2", "the max concurrency level")
	inputPocJson := flag.String("pocJson", "", "Input the PocJson file")
	inputProxy := flag.String("proxy", "http://127.0.0.1:8082", "Input the proxy")
	inputFofaEmail := flag.String("email", "", "Input the fofa email")
	inputFofaKey := flag.String("key", "", "Input the fofa key")
	flag.Parse()

	maxConcurrentLevelInt, err := strconv.Atoi(*maxConcurrentLevel)
	if err != nil {
		fmt.Println("The maximum concurrency you entered is not a number!", err)
	}
	pocStruct := Handle.TryToParsePocStruct(*inputPocJson)
	var queryResponse Fofa.QueryResponse
	if userInputDetectionURL == nil {
		err = json.Unmarshal(Fofa.Query(*inputFofaEmail, *inputFofaKey, pocStruct.Fofa, 3), &queryResponse)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
			os.Exit(1)
		}
	}

	fmt.Println("[+] 扫描开始,记得挂全局socks代理! :")
	Http.Send(pocStruct, queryResponse, userInputDetectionURL, *inputProxy, maxConcurrentLevelInt)
}
