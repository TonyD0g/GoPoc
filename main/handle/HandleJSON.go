package handle

import (
	"Scanner/main/format"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func HandleJsonFunc(inputXml string) format.RequestPackage {
	jsonData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		fmt.Println("[-] Error reading file:", err)
		os.Exit(1)
	}

	var requestPackage format.RequestPackage
	err = json.Unmarshal(jsonData, &requestPackage)
	if err != nil {
		fmt.Println("[-] Error unmarshal Json:", err)
		os.Exit(1)
	}
	return requestPackage
}

func TraversePath(requestPackage format.RequestPackage) []string {
	var allReqPath []string
	var i = 0
	pathFieldValue := reflect.ValueOf(requestPackage.PathList).FieldByName("Path")
	var pathList []string
	for i := 0; i < pathFieldValue.Len(); i++ {
		pathList = append(pathList, pathFieldValue.Index(i).String())
	}

	for _, tmpPath := range pathList {
		if !strings.HasPrefix(tmpPath, "/") {
			tmpPath = "/" + tmpPath
		}
		i++
		allReqPath[i] = requestPackage.Url + tmpPath
	}
	return allReqPath
}

func CheckFileCorrectness(requestPackage format.RequestPackage) bool {
	if requestPackage.Method != "GET" && requestPackage.Url != "POST" {
		fmt.Print("[-] error! You must provide the correct method")
		return false
	}
	if requestPackage.Url == "" {
		fmt.Print("[-] error! You must provide the url")
		return false
	}
	return true
}
