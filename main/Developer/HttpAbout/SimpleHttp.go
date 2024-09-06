package HttpAbout

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Fofa"
	"GoPoc/main/Log"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func SendForFofa(config map[string]string, pocStruct Format.PocStruct) []string {
	maxFofaSizeInt, err := strconv.Atoi(config["maxFofaSize"])
	if err != nil {
		Log.Log.Fatal("maxFofaSize 并不是一个有效数字\n")
	}

	var urlsList []string
	var queryResponse Fofa.QueryResponse
	err = json.Unmarshal(Fofa.SearchReturnByte(config, pocStruct, maxFofaSizeInt), &queryResponse)
	if err != nil {
		Log.Log.Fatal("Failed to parse JSON:", err)
	}

	for _, tmpOutcome := range queryResponse.Results {
		if !strings.HasPrefix(tmpOutcome[1].(string), tmpOutcome[0].(string)) {
			urlsList = append(urlsList, tmpOutcome[0].(string)+"://"+tmpOutcome[1].(string))
		} else {
			urlsList = append(urlsList, tmpOutcome[1].(string))
		}
	}
	Log.Log.Println(fmt.Printf("[+] 查询 fofa 语句为: %v 该fofa语句查询到: %v 条,你最大想搜索 %v 条\n", queryResponse.Query, queryResponse.Size, config["maxFofaSize"]))
	return urlsList
}

func SendForUrlOrFile(userInputDetectionURL string) []string {
	var urlsList []string
	if strings.HasPrefix(strings.ToLower(userInputDetectionURL), "http://") || strings.HasPrefix(strings.ToLower(userInputDetectionURL), "https://") {
		urlsList = append(urlsList, userInputDetectionURL)
	} else {
		// urlFile list
		urlFile, err := os.Open(userInputDetectionURL)
		if err != nil {
			Log.Log.Fatal("can't open file:", err)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				Log.Log.Fatal("can't close file:", err)
			}
		}(urlFile)

		reader := bufio.NewReader(urlFile)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && line == "" {
				break
			}
			urlsList = append(urlsList, strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", ""))
			line = ""
		}
	}
	return urlsList
}
