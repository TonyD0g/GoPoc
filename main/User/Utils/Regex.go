package Utils

import (
	"regexp"
	"strings"
)

// GetXXXContent 获取字符串中括号内/外的内容
func GetXXXContent(_type string, str string) string {
	var regEx string
	switch _type {
	case "rightBrackets":
		regEx = "(?<=\\))(.+)"
	case "brackets":
		regEx = "(?<=\\()(.+?)(?=\\))"
	case "middleBrackets":
		regEx = "\\[(.*?)]"
	case "brace":
		regEx = "\\{(.+?)\\}"
	}

	reg := regexp.MustCompile(regEx)
	matches := reg.FindAllStringSubmatch(str, -1)

	var stringBuilder strings.Builder
	for _, match := range matches {
		stringBuilder.WriteString(match[1])
	}
	return stringBuilder.String()
}

// RegexSpecialChair 正则过滤特殊字符(1.空格 2.特殊字符[中英文])
func RegexSpecialChair(str string) string {
	regEx := "[\\s`~!@#$%^&*()+=|{}':;',\\[\\].<>/?~！@#￥%……&*（）——+|{}【】‘；：”“’。，、？]"
	reg := regexp.MustCompile(regEx)
	return reg.ReplaceAllString(str, "")
}
