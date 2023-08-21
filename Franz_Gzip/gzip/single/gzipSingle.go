package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var path = "gzip/single/test.txt.gz"

func main() {
	writeFIle()
	readFile()
}

func writeFIle() {
	//file, err := os.Create(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := gzip.NewWriter(file)
	defer writer.Close()
	//return
	for i := 0; i < 10; i++ {
		writer.Write([]byte("Gophers rule!\n"))
	}
}

func readFile() {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	strContent := string(content)
	strSlice := strings.Split(strContent, "\n")
	for row, value := range strSlice {
		time.Sleep(time.Second)
		fmt.Printf("Row %d with value: %s\n", row, value)
	}
}
