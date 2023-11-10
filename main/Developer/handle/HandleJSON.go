package Handle

import (
	format2 "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Log"
	"encoding/json"
	"strings"
)

var pocStruct format2.PocStruct

func TryToParsePocStruct(jsonData string) format2.PocStruct {
	err := json.Unmarshal([]byte(jsonData), &pocStruct)
	if err != nil {
		Log.Log.Fatal("[-] Error unmarshal Json:", err)
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
