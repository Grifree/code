package connect_sql

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	sq "github.com/goclub/sql"
)

func NewDB() (db *sq.Database, err error) {
	db,_, err = sq.Open("mysql", sq.MysqlDataSource{
		User:     "root",
		Password: "somepass",
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "goclub_example",
	}.String()) ; if err != nil {
		panic(err)
	}
	err = db.Ping(context.TODO()) ; if err != nil {
		panic(err)
	}
	return
}