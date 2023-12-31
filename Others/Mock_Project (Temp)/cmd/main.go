package main

import (
	"Mock_Project/api"
	configs "Mock_Project/config"
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
	config, err := configs.InitConfig()
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
