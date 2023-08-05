// Package pkg for store helper func
package pkg

import (
	"Mock_Project/model"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// TransformInterfaceToString convert interface to string
func TransformInterfaceToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return model.EmptyString
	}
}

// GenerateString generate random string
func GenerateString() string {
	return uuid.New().String()
}

func LogStepProcess(start time.Time, log string) {
	text := "\n=> %s (+%s)\t\t" + log + model.NewLineCharacter
	now, diff := getLogProcessTime(start)
	fmt.Printf(text, now, diff)
}

func getLogProcessTime(start time.Time) (string, string) {
	now := time.Now()
	diff := now.Sub(start)
	out := time.Time{}.Add(diff)
	return fmt.Sprintf("%s", now.Format(time.StampMicro)),
		fmt.Sprintf("%s", out.Format(model.TimeFormatWithMicrosecond))
}
