package db

import (
	"Mock_Project/model"
	"context"
)

type IDBHandler interface {
	InitConnection(config *model.Server, endpoint, dbName string) error
	Exec(ctx context.Context, endpoint, dbName, sql string, args []interface{}) error
	CloseAllDb() error
}
