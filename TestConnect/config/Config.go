package config

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func Connection() *sql.DB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "ChuanDo@13",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "demo",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Initial Database Failed!!!")
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Connect Database Failed!!!")
		return nil
	}
	fmt.Println("Connect Database Success!!!")
	return db
}

func Close(db *sql.DB) {
	_ = db.Close()
	fmt.Println("Close Database Success!!!")
}
