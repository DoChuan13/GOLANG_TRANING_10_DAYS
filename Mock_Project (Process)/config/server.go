package config

import "Mock_Project/model"

func InitConfig() (*model.Server, error) {
	config := new(model.Server)
	//Config Kafka
	config.Broker = []string{"127.0.0.1:9093"}
	config.Topics = map[string]int{}
	config.MaxPartition = 1
	config.Block = 1000

	//Config End Point
	config.SqlPath = "/docker-entrypoint-initdb.d/temp/"
	//config.SqlPath = "/Users/Chuan/Personal/Docker/Bitnami/sqldump/temp/"
	config.LocalPath = "/Users/Chuan/Personal/Docker/Bitnami/sqldump/temp/"
	config.Endpoint = "127.0.0.1"
	config.DBName = "demo"
	config.DBLtd = 200

	//Config Database
	config.User = "root"
	config.Password = "ChuanDo@13"
	config.DriverName = "mysql"
	config.Port = 3306
	config.MaxIdleConnection = 100
	config.MaxOpenConnection = 100
	config.RetryTimes = 5000
	config.RetryWaitMs = 5000

	//Goroutine
	config.ConsumerLtd = 200
	config.ProducerLtd = 200000
	return config, nil
}
