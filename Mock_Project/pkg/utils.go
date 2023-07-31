// Package pkg for store helper func
package pkg

import (
	"Mock_Project/model"
	"github.com/google/uuid"
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
