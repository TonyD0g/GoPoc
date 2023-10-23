package main

import (
	"Scanner/main/Developer/handle"
	"Scanner/main/Developer/http"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Get command line parameters
	args := os.Args
	// Determine whether the number of parameters is correct
	if args[3] != "-pocJson" || args[5] != "-proxy" {
		fmt.Println("[-] Input parameter error, case as follows:\n-url\nhttp://127.0.0.1/\n-pocJson\nC:\\Users\\xxx\\Desktop\\1.json\n-proxy\nhttp://127.0.0.1:8082")
		flag.Usage()
		os.Exit(1)
	}

	var userInputDetectionURL *string
	if args[1] == "-url" {
		userInputDetectionURL = flag.String("url", "", "input the url")
	} else {
		userInputDetectionURL = flag.String("file", "", "input the url file list")
	}
	maxConcurrentLevel := flag.String("maxConcurrentLevel", "2", "the max concurrency level")
	inputPocJson := flag.String("pocJson", "", "input the PocJson file")
	inputProxy := flag.String("proxy", "", "input the proxy")
	flag.Parse()

	maxConcurrentLevelInt, err := strconv.Atoi(*maxConcurrentLevel)
	if err != nil {
		fmt.Println("The maximum concurrency you entered is not a number!", err)
	}
	http.Send(handle.TryToParseJson(*inputPocJson), *userInputDetectionURL, *inputProxy, maxConcurrentLevelInt)
}
