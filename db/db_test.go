package db_test

import (
	"testing"
	"time"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	"log"
)

type Demo struct {
	ID string `db:"id"`
	Name string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

var Db *sqlx.DB

func TestBD(t *testing.T){
	// 连接
	dbuser := "root"
	dbpwd := "somepass"
	dbhost := "127.0.0.1"
	dbname := "shopping_mall"
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbuser, dbpwd, dbhost, dbname)
	db, err := sqlx.Connect("mysql", dns)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	//事务
	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO user (id,name) VALUES ('1', 'Jack');`)
	tx.MustExec(`INSERT INTO user (id,name) VALUES ('2', 'Emily');`)
	// 提交
	//err = tx.Commit()
	// 回滚
	err = tx.Rollback()
	if err != nil {
		log.Fatalln(err)
	}


}
