package dao

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const(
	host = "localhost"
	port = 5432
	dbname = "postgres"
	username = "postgres"
	password = "postgres"
)

var G_db *sql.DB

func init(){
	var err error
	url := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",host,port,dbname,username,password)
	G_db, err = sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
}

