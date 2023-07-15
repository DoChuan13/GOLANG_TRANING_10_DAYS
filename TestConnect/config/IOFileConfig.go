package config

import (
	"fmt"
	"os"
)

var FilePath = "Practice/ReadWrite/test.csv"

func ReadFile(filePath string) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Read file failed==>", err)
		return
	}
	fmt.Println("Content File==>")
	fmt.Println(string(f))
}

func OverWriteFile(filePath string, content []string) {
	for _, value := range content {
		err := os.WriteFile(filePath, []byte(value), 0777)
		if err != nil {
			fmt.Println("Write file failed==>", err)
			return
		}
	}
}

func OpenFile(filePath string) *os.File {
	_, info := os.Stat(filePath)
	if info != nil {
		fmt.Println(info)
		fmt.Println("New file is created!!!")
		file, _ := os.Create(filePath)
		return file
	}
	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	return file
}

func WriteFile(file *os.File, content []string) {
	for _, value := range content {
		_, err := fmt.Fprintln(file, value)
		if err != nil {
			fmt.Println("Write files failed", err)
		}
	}
}

func CloseFile(f *os.File) {
	_ = f.Close()
}
