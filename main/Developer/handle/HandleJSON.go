package Handle

import (
	format2 "Scanner/main/Developer/Format"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var pocStruct format2.PocStruct

func TryToParsePocStruct(inputXml string) format2.PocStruct {
	jsonData, err := ioutil.ReadFile(inputXml)
	if err != nil {
		fmt.Println("[-] Error reading file:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(jsonData, &pocStruct)
	if err != nil {
		fmt.Println("[-] Error unmarshal Json:", err)
		os.Exit(1)
	}
	return pocStruct
}

// TraversePath Traverse the Path in the PathList and initiate a request
func TraversePath(requestPackage format2.RequestPackage, inputUrl string) []string {
	var allReqPath []string
	var i = 0

	for _, tmpPath := range requestPackage.Uri {
		if !strings.HasPrefix(tmpPath, "/") {
			tmpPath = "/" + tmpPath
		}
		allReqPath = append(allReqPath, inputUrl+tmpPath)
		i++
	}
	return allReqPath
}
