package Fofa

import (
	"GoPoc/main/Developer/AllFormat"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func SearchReturnByte(config map[string]string, pocStruct Format.PocStruct, maxFofaSizeInt int) []byte {
	fofaClient := NewFofaClient([]byte(config["email"]), []byte(config["key"]))
	if fofaClient == nil {
		fmt.Printf("fofa 查询失败\n")
		os.Exit(1)
	}
	var (
		query        = []byte(nil)
		fields       = []byte("protocol,host,ip")
		lastQueryUrl = []byte(nil)
	)

	query = []byte(pocStruct.Fofa)
	lastQueryUrl = []byte(base64.StdEncoding.EncodeToString(query))
	lastQueryUrl = bytes.Join([][]byte{[]byte("https://fofa.info/api/v1/search/all?"),
		[]byte("email="), fofaClient.Email,
		[]byte("&key="), fofaClient.Key,
		[]byte("&qbase64="), lastQueryUrl,
		[]byte("&fields="), fields,
		[]byte("&page="), []byte(strconv.Itoa(1)),
		[]byte("&size="), []byte(strconv.Itoa(maxFofaSizeInt)),
	}, []byte(""))
	fmt.Printf("%s\n", lastQueryUrl) // fofa 查询语句
	content, err := fofaClient.Get(string(lastQueryUrl))

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}
	if strings.Contains(string(content), "查询语法错误") {
		fmt.Printf("%s\n", "fofa 查询语法错误")
		os.Exit(1)
	}

	return content
}
