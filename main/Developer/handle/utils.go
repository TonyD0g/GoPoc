package Handle

import (
	"GoPoc/main/Log"
	"bufio"
	"os"
	"strings"
)

func HandleIni(fileName string) map[string]string {
	file, err := os.Open(fileName)
	if err != nil {
		Log.Log.Fatal("打开文件失败：%v\n", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			Log.Log.Fatal("关闭文件失败：%v\n", err)
		}
	}(file)

	// 创建一个map用于存储键值对
	config := make(map[string]string)

	// 使用Scanner按行读取文件内容
	scanner := bufio.NewScanner(file)
	var key string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if strings.HasPrefix(line, "-") {
				// 保存当前行作为键名
				key = strings.TrimPrefix(line, "-")
			} else if key != "" {
				// 存储键值对到map中
				config[key] = line
				key = ""
			}
		}
	}
	return config
}
