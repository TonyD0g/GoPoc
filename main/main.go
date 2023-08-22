package main

import (
	"Scanner/main/handle"
	"Scanner/main/http"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 获取命令行参数
	args := os.Args

	// 判断参数数量是否正确
	if len(args) != 3 && args[1] != "-pocJson" {
		fmt.Println("[-] Incorrect number of parameters")
		flag.Usage()
		os.Exit(1)
	}

	inputPocJson := flag.String("pocJson", "", "input the PocJson file")
	inputProxy := flag.String("proxy", "", "input the proxy")

	flag.Parse()

	inputPocJsonValue := *inputPocJson
	inputProxyValue := *inputProxy

	var pocStruct = handle.HandleJsonFunc(inputPocJsonValue)

	http.Send(pocStruct, inputProxyValue)

}
