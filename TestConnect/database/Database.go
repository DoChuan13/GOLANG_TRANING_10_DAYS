package database

import (
	"TestConnect/model"
	"context"
	"database/sql"
	"fmt"
	"log"
)

func CreateTables(db *sql.DB) {
	//db := config.Connection()
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

func SaveDataDB(db *sql.DB, user *model.User) {
	//db := config.Connection()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := "insert into users (id,name) values (?,?)"
	_, err = tx.ExecContext(ctx, query, user.Id, user.Name)
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
