package repository

import (
	"Mock_Project/model"
	"context"
)

type IDBRepository interface {
	InitConnection(config *model.Server, endpoint, dbName string) error
	CreateNewTable(ctx context.Context, object model.ConsumerObject) error
	ImportDataFiles(file string, ctx context.Context, object model.ConsumerObject) error
	ExportDataFiles(file string, ctx context.Context, object model.ConsumerObject) error
	ClearData(ctx context.Context, object model.ConsumerObject) error
	Close() error
	//InsertData(ctx context.Context, object model.ConsumerObject, args []interface{}) error
}

type IKafkaRepository interface {
	CreateTopic(topic string, partitionNum int32) error
	ProducerData(topic string, partition int32, content string) error
	ConsumerData(topic string, partition int32) ([]string, error)
	InitConnection(topic string) error
	CloseTopic() error
	RemoveTopic() error
}
