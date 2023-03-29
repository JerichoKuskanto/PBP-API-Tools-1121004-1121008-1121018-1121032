package controller

import (
	"database/sql"
	"log"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pbp_tugas_eksplorasi2?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal()
	}
	return db
}
