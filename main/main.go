package main

import (
	"Scanner/main/handle"
	"Scanner/main/http"
	"flag"
	"fmt"
	"os"
)

func main() {

	// Get command line parameters
	args := os.Args

	// Determine whether the number of parameters is correct
	if len(args) != 7 && args[1] != "-url" && args[3] != "-pocJson" {
		fmt.Println("[-] Incorrect number of parameters")
		flag.Usage()
		os.Exit(1)
	}

	inputUrl := flag.String("url", "", "input the url")
	inputPocJson := flag.String("pocJson", "", "input the PocJson file")
	inputProxy := flag.String("proxy", "", "input the proxy")

	flag.Parse()

	inputUrlValue := *inputUrl
	inputPocJsonValue := *inputPocJson
	inputProxyValue := *inputProxy

	var pocStruct = handle.HandleJsonFunc(inputPocJsonValue)

	http.Send(pocStruct, inputUrlValue, inputProxyValue)

}
