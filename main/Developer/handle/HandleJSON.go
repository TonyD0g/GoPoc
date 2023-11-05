package Handle

import (
	format2 "GoPoc/main/Developer/AllFormat"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var pocStruct format2.PocStruct

func TryToParsePocStruct(jsonData string) format2.PocStruct {
	err := json.Unmarshal([]byte(jsonData), &pocStruct)
	if err != nil {
		fmt.Println("[-] Error unmarshal Json:", err)
		os.Exit(1)
	}
	return pocStruct
}

// TraversePath 遍历PathList中的Path,并添加到allReqPath,用于处理poc中同时出现多个uri的情况
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
