package file

import (
	"Mock_Project/model"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Code struct {
	name  string
	first string
	last  string
}

const (
	numRecords = 100
)

func FakeAllData() error {
	code1 := Code{name: "A0001&0/T", first: "A", last: "T"}
	code2 := Code{name: "A0002&0/T", first: "A", last: "T"}
	code3 := Code{name: "A0001&0/F", first: "A", last: "F"}
	code4 := Code{name: "A0009&0/G", first: "A", last: "G"}
	code5 := Code{name: "A0001&0/G", first: "A", last: "G"}
	code6 := Code{name: "A0008&0/G", first: "A", last: "G"}
	code7 := Code{name: "A0001&0/V", first: "A", last: "V"}
	code8 := Code{name: "A0005&0/V", first: "A", last: "V"}

	codeList := []Code{code1, code2, code3, code4, code5, code6, code7, code8}
	return fakerData(codeList)
}

func fakerData(codeList []Code) error {
	path := "file/ListValue.csv"
	// Tạo file CSV
	file := OpenFile(path)
	defer CloseFile(file)

	for _, code := range codeList {
		for i := 0; i < numRecords; i++ {
			var record model.SourceObject
			record.QCD = code.name
			record.TIME = getRandomTime()
			record.TKQKBN = code.first
			record.SNDC = code.last
			record.ZXD = strings.ReplaceAll(randomDate(), "-", "")

			// Ghi dữ liệu vào file CSV
			str := convertValues(record)
			err := WriteFile(file, str)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func convertValues(source model.SourceObject) string {
	result := ""
	val := reflect.ValueOf(source)
	for j := 0; j < val.NumField(); j++ {
		temp := val.Field(j).Interface()
		switch v := temp.(type) {
		case string:
			result += v
		case int:
			result += strconv.Itoa(v)
		}

		if j < val.NumField()-1 {
			result += model.CommaCharacter
		}
	}
	return result
}

func randomDate() string {
	now := time.Now()
	randDay := rand.Intn(365)
	randTime := now.AddDate(0, 0, -randDay)
	return randTime.Format(model.DateFormatWithStroke)
}

func getRandomTime() string {
	// Tạo seed mới cho random
	rand.Seed(time.Now().UnixNano())

	randomHour := rand.Intn(24)
	randomMinute := rand.Intn(60)

	return fmt.Sprintf("%02d:%02d", randomHour, randomMinute)
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

func WriteFile(file *os.File, content string) error {
	_, err := fmt.Fprintln(file, content)
	if err != nil {
		return err
	}
	return nil
}

func CloseFile(f *os.File) {
	_ = f.Close()
}
