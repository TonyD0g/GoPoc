package Http

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Fofa"
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
		fmt.Printf("maxFofaSize 并不是一个有效数字\n")
		os.Exit(1)
	}

	var urlsList []string
	var queryResponse Fofa.QueryResponse
	err = json.Unmarshal(Fofa.SearchReturnByte(config, pocStruct, maxFofaSizeInt), &queryResponse)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		os.Exit(1)
	}

	for _, tmpOutcome := range queryResponse.Results {
		if !strings.HasPrefix(tmpOutcome[1].(string), tmpOutcome[0].(string)) {
			urlsList = append(urlsList, tmpOutcome[0].(string)+"://"+tmpOutcome[1].(string))
		} else {
			urlsList = append(urlsList, tmpOutcome[1].(string))
		}
	}
	fmt.Printf("[+] 此 fofa 语句: %v 查询到: %v 条,你想搜索 %v 条\n", queryResponse.Query, queryResponse.Size, config["maxFofaSize"])
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
			fmt.Println("can't open file:", err)
			os.Exit(1)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				fmt.Println("can't close file:", err)
				os.Exit(1)
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
