package api

import (
	"Mock_Project/infrastructure"
	"Mock_Project/model"
	"Mock_Project/pkg/logger"
	"context"
)

type Server struct {
	infra *infrastructure.Infra
	cfg   *model.Server
}

// New create new server instance
func New(infra *infrastructure.Infra, cfg *model.Server) *Server {
	return &Server{
		infra: infra,
		cfg:   cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	//cfg := model.KafkaSystem{
	//    Broker:    []string{"0.0.0.0:9093"},
	//    Topics:     "SZN-TSE1",
	//    Partition: 0,
	//}
	//kafka, err := kafka2.NewKafkaHandler(&cfg)
	//if err != nil {
	//    return
	//}
	////_ = kafka.ProducerData(cfg.Broker, cfg.Topics, cfg.Partition, "Hello World 20")
	//value, _ := kafka.ConsumerData(cfg.Broker, cfg.Topics, cfg.Partition, parseConsumer)
	//fmt.Println(value)
	return nil
}

// Close closes server and related resources
func (s *Server) Close() {
	s.infra.Close()
	logger.AppLog(logger.InfoAppNewRecoverDataFinished, logger.InfoLevelLog)
}
