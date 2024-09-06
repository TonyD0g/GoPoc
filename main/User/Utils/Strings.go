package Utils

import (
	"GoPoc/main/Developer/HttpAbout"
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// RandomStringByModule 根据模式来返回对应的字符串. 模式1为所有字符(即包含数字和英文字符),模式2为纯数字,模式3为纯英文字符
func RandomStringByModule(size, module int) string {
	var allChar string
	if module == 1 {
		allChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	} else if module == 2 {
		allChar = "1234567890"
	} else {
		allChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		buffer.WriteByte(allChar[rand.Intn(len(allChar))])
	}
	return buffer.String()
}

// SplitStr 返回被 两两分割 的数组
func SplitStr(key string) []int {
	str := []rune(key)
	result := make([]int, len(str)/2)
	l, i, m := 2, 0, 0

	for i+m < len(str) {
		if i%l == 0 && i != 0 {
			str = append(str[:i+m], append([]rune{','}, str[i+m:]...)...)
			m++
		}
		i++
	}

	se := string(str)
	a := strings.Split(se, ",")
	for i := 0; i < len(result); i++ {
		j := a[i]
		if j == "" {
			break
		}
		val, _ := strconv.Atoi(j)
		result[i] = val
	}

	return result
}

// ParseHeaders 解析输入字符串，并返回键值对的映射
func ParseHeaders(input string) {
	// 创建一个正则表达式，用于匹配键值对
	re := regexp.MustCompile(`(?m)^\s*(\S+):\s*(.*)$`)
	matches := re.FindAllStringSubmatch(input, -1)

	// 创建一个map来存储结果
	headers := make(map[string]string)

	// 遍历匹配结果，并将其存储到map中
	for _, match := range matches {
		if len(match) > 2 {
			headers[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}
	fmt.Println("[+] 解析 header 字符串,输出的结果为:")
	// 打印结果
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
	}
}

// FullyAutomaticFillingHeader 输入请求体，全自动构造 header 【但可能存在Bug,如果遇到bug了,那还是人为构造吧】
func FullyAutomaticFillingHeader(config HttpAbout.SetHttpConfig, input string) HttpAbout.SetHttpConfig {
	// 创建一个正则表达式，用于匹配键值对
	re := regexp.MustCompile(`(?m)^\s*(\S+):\s*(.*)$`) // todo 这个正则表达式可能会存在问题
	matches := re.FindAllStringSubmatch(input, -1)

	// 创建一个map来存储结果
	headers := make(map[string]string)

	// 遍历匹配结果，并将其存储到map中
	for _, match := range matches {
		if len(match) > 2 {
			headers[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}
	for key, value := range headers {
		if key == "Host" || key == "User-Agent" {
			continue
		}
		config.Header[key] = value
	}
	return config
}
