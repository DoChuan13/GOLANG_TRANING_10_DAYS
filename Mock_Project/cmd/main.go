package main

import (
	"Mock_Project/api"
	config2 "Mock_Project/config"
	"Mock_Project/file"
	"Mock_Project/infrastructure"
	"context"
	"fmt"
	"runtime"
)

func init() {
	initLog()
}

func initLog() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	err := file.FakeAllData()
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := config2.InitConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	infra, err := infrastructure.Init(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	server := api.New(infra, config)
	defer func() {
		server.Close()
	}()
	err = server.Start(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

}
