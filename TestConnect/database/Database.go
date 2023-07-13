package database

import (
	"TestConnect/config"
	"context"
	"fmt"
	"log"
)

func CreateTables() {
	db := config.Connection()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := "create table if not exists users (id  int primary key auto_increment, name varchar(255) not null);"
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		fmt.Println("Failed")
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
