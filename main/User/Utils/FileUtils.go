package Utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// WriteLinesToFile [写文件] 写内容进文件中
func WriteLinesToFile(filePath string, lines []string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file) // 确保在函数退出时关闭文件

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// 确保所有缓冲区的数据都写入文件
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}

// ReadFileToStringArray [读文件] 按行读取文件,返回值为字符串数组
func ReadFileToStringArray(filePath string) []string {
	var lines []string
	file, err := os.Open(filePath) // 替换为你的文件路径
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file) // 确保在函数结束时关闭文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // 将每一行添加到切片
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	return lines
}

// AppendToFile [写文件] 使用追加的方式将内容写入进文件
func AppendToFile(filePath string, content string) {
	// 打开文件，使用追加访问模式及创建文件的权限
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// 将字符串写入文件
	if _, err = file.WriteString(content + "\n"); err != nil { // 可以根据需要添加换行符
		fmt.Printf("failed to write to file: %v\n", err)
	}
}
