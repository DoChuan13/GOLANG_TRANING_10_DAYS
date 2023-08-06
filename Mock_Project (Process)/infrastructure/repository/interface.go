package repository

import (
	"Mock_Project/model"
	"context"
)

type IDBRepository interface {
	InitConnection(config *model.Server, endpoint, dbName string) error
	GenerateTableAndExpFile(ctx context.Context, tableName string) error
	ImportDataFiles(ctx context.Context, tableName string) error
	ClearData(ctx context.Context, tableName string) error
	CloseAllDb() error
	//InsertData(ctx context.Context, object model.ConsumerObject, args []interface{}) error
}

type IKafkaRepository interface {
	CreateTopic(topic string, partitionNum int32) error
	SyncProducerData(topic string, partition int32, content string) error
	ASyncProducerData(topic string, partition int32, content string) error
	ConsumerData(topic string, partition int32) ([]string, error)
	InitConnection() error
	CloseTopic() error
	RemoveTopic() error
}
