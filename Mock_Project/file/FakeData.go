package file

import (
	"Mock_Project/model"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	numRecords = 900
)

func FakeAllData() {
	fakerData("0001&0", "A", "T")
	//fakerData("0002&0", "A", "T")
	//fakerData("0001&0", "A", "F")
	//fakerData("0009&0", "A", "G")
	//fakerData("0001&0", "A", "G")
	//fakerData("0008&0", "A", "G")
	//fakerData("0001&0", "A", "V")
	//fakerData("0005&0", "A", "V")
}

func fakerData(name, first, last string) {
	nameFile := first + name + last
	// Tạo file CSV
	file, err := os.Create("file/faker/" + nameFile + ".csv")
	if err != nil {
		fmt.Println("Lỗi khi tạo file:", err)
		return
	}
	defer file.Close()

	// Khởi tạo writer để ghi vào file CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Tạo dữ liệu giả lập và ghi vào file CSV
	for i := 0; i < numRecords; i++ {
		var record model.SourceObject
		record.QCD = first + name + model.StrokeCharacter + last
		record.TIME = getRandomTime()
		record.TKQKBN = first
		record.SNDC = last
		record.ZXD = strings.ReplaceAll(randomDate(), "-", "")

		// Ghi dữ liệu vào file CSV
		str := convertValues(record)
		writer.Write([]string{str})
	}
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
