package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	schema = "%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local"
	// docker-compose.ymlに設定した環境変数を取得
	username       = os.Getenv("MYSQL_USER")
	password       = os.Getenv("MYSQL_PASSWORD")
	host           = os.Getenv("MYSQL_HOST")
	dbName         = os.Getenv("MYSQL_DATABASE")
	datasourceName = fmt.Sprintf(schema, username, password, host, dbName)
	Db             *sql.DB
)

func init() {
	log.Println("database setup")
	log.Println("datasource:", datasourceName)
	connection, err := sql.Open("mysql", datasourceName)

	if err != nil {
		panic("Could not connect to the database")
	}
	Db = connection
}
