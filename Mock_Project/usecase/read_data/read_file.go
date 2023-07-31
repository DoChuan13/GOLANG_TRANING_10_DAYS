package read_data

import (
	"Mock_Project/model"
	"os"
	"strings"
)

type Service struct {
}

// NewService service constructor
func NewService() IFile {
	return &Service{}
}

func (s Service) ReadFileProcess(path string) ([]string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(string(file), model.EnterLineLF)
	var result []string
	lastRow := rows[len(rows)-1]
	if lastRow == model.EmptyString {
		result = rows[:len(rows)-1]
	} else {
		result = rows
	}
	return result, nil
}
