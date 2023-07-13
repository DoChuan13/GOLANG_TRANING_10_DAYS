package main

import (
	"TestConnect/config"
	"TestConnect/database"
)

func main() {
	config.Connection()
	database.CreateTables()
}
