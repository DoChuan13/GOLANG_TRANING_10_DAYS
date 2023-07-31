package config

import "Mock_Project/model"

func InitConfig() (*model.Server, error) {
	config := new(model.Server)
	config.Broker = []string{"da"}
	config.User = "root"
	config.Password = "123456123"
	config.DriverName = "mysql"
	config.Port = 3306
	config.MaxIdleConnection = 10
	config.MaxOpenConnection = 10
	config.RetryTimes = 3000
	config.RetryWaitMs = 3000
	return config, nil
}
