package repository

import (
	"Mock_Project/model"
	"context"
)

type IDBRepository interface {
	ImportData(
		ctx context.Context, object model.ObjectProcess, tableName, endpoint, dbName string, args []interface{},
	) error
	Close() error
}

type IKafkaRepository interface {
	ProducerData(broker []string, topic string, partition int32, content string) error
	ConsumerData(broker []string, topic string, partition int32) ([]interface{}, error)
	ClearData(server *model.Server) error
}
