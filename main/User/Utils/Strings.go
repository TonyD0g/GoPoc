package Utils

import (
	"bytes"
	"math/rand"
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
