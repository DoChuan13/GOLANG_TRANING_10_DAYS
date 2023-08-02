package read_data

import (
	"Mock_Project/model"
	"fmt"
	"os"
	"strings"
)

type Service struct {
	parentPath string
	filePath   string
	name       string
}

// NewService service constructor
func NewService(path, name string) IFile {
	return &Service{
		parentPath: path,
		filePath:   path + "/" + name,
		name:       name,
	}
}

func (s Service) ReadFileProcess() ([]string, error) {
	file, err := os.ReadFile(s.filePath)
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

func (s Service) InsertCurrentFiles(rows *[]string) error {
	filePath := s.parentPath + model.StrokeCharacter + s.name
	file, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer closeFile(file)
	for _, row := range *rows {
		err := writeFile(file, row)
		if err != nil {
			return err
		}
	}
	return nil
}

func openFile(filePath string) (*os.File, error) {
	_, info := os.Stat(filePath)
	if info != nil {
		fmt.Println(info)
		fmt.Println("New file is created!!!")
		err := os.Mkdir(filePath, 0755)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	return file, nil
}

func writeFile(file *os.File, content string) error {
	_, err := fmt.Fprintln(file, content)
	if err != nil {
		fmt.Println("Write files failed")
		return err
	}
	return nil
}

func closeFile(f *os.File) {
	_ = f.Close()
}
