package database

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func Dbconnection() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("mysql", "root:newpassword@tcp(localhost:3306)/newdatabas")
	if err != nil {
		log.Fatal(err)
	}
	return Db, nil
}

func UserDb() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("mysql", "root:newpassword@tcp(localhost:3306)/mydatabase")
	if err != nil {
		log.Fatal(err)
	}
	return Db, nil
}
