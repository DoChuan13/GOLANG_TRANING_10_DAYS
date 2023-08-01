package repository

import (
	"Mock_Project/model"
	"context"
)

type IDBRepository interface {
	InitConnection(config *model.Server, endpoint, dbName string) error
	CreateNewTable(ctx context.Context, object model.ObjectProcess) error
	InsertData(ctx context.Context, objectProcess model.ObjectProcess, args []interface{}) error
	Close() error
}

type IKafkaRepository interface {
	ProducerData(broker []string, topic string, partition int32, content string) error
	ConsumerData(broker []string, topic string, partition int32) ([]interface{}, error)
	ClearData(server *model.Server) error
}
