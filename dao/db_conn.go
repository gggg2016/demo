package dao

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"demo/config"
)

var host, dbname, username, password string
var port int64
var G_db *sql.DB

func init() {
	//read config file
	host = config.GetAsString("host")
	port = config.GetAsInt64("port")
	dbname = config.GetAsString("dbname")
	username = config.GetAsString("username")
	password = config.GetAsString("password")

	var err error
	url := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, username, password)
	G_db, err = sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
}
