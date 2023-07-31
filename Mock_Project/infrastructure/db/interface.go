package db

import (
	"Mock_Project/model"
	"context"
)

type IDBHandler interface {
	InitConnection(config *model.KafkaDB, endpoint, dbName string) error
	Exec(ctx context.Context, endpoint, dbName, sql string, args []interface{}) error
	Close() error
}
