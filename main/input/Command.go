package input

import (
	"Scanner/main/handle"
	"Scanner/main/http"
	"flag"
	"fmt"
	"os"
)

func HandleUserInput() {
	// 获取命令行参数
	args := os.Args

	// 判断参数数量是否正确
	if len(args) != 1 {
		fmt.Println("[-] Incorrect number of parameters")
		flag.Usage()
		os.Exit(1)
	}

	inputPocXml := flag.String("pocXml", "", "input the method")

	flag.Parse()

	inputPocXmlValue := *inputPocXml

	var requestPackage = handle.HandleXMLFunc(inputPocXmlValue)

	http.Send(requestPackage)

}
