package main

import (
	"TestConnect/config"
	"fmt"
)

func main() {
	//config.Connection()
	//database.CreateTables()

	fmt.Println("input Number")
	demo := config.InputFloat()
	fmt.Println("Result =>", demo)

}
