package config

import "Mock_Project/model"

func InitConfig() (*model.Server, error) {
	config := new(model.Server)
	//Config Kafka
	config.Broker = []string{"127.0.0.1:9093"}
	config.Topics = map[string]int{}
	config.MaxPartition = 20
	config.Block = 1000

	//Config End Point
	config.Endpoint = "127.0.0.1"
	config.DBName = "demo"

	//Config Database
	config.User = "root"
	config.Password = "ChuanDo@13"
	config.DriverName = "mysql"
	config.Port = 3306
	config.MaxIdleConnection = 10
	config.MaxOpenConnection = 10
	config.RetryTimes = 3000
	config.RetryWaitMs = 3000

	//Goroutine
	config.Limited = 100000
	return config, nil
}
