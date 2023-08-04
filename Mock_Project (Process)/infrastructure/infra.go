package infrastructure

import (
	"Mock_Project/infrastructure/db"
	"Mock_Project/infrastructure/kafka"
	"Mock_Project/model"
	"Mock_Project/pkg/logger"
)

type Infra struct {
	DBHandler    db.IDBHandler
	KafkaHandler kafka.IKafkaHandler
}

func Init(cfg *model.Server) (*Infra, error) {
	kafkaDBHandler, err := db.NewDBHandler(cfg)
	if err != nil {
		return nil, err
	}

	kafkaHandler, err := kafka.NewKafkaHandler(cfg)
	if err != nil {
		logger.AppLogKafka(
			logger.ErrorCodeInitS3ConnectionFail, logger.ErrorLevelLog,
		)
		return nil, err
	}
	return &Infra{
		DBHandler:    kafkaDBHandler,
		KafkaHandler: kafkaHandler,
	}, nil
}

// Close closes resources gracefully
func (f *Infra) Close() {
	if f.DBHandler != nil {
		_ = f.DBHandler.Close()
	}
}
