package db_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //初始化一个mysql驱动，必须
	"github.com/jmoiron/sqlx" // 语法介绍 https://www.jianshu.com/p/01634def632c
	"testing"
	"time"
)

type User struct {
	ID string `db:"id"`
	Name string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

var Db *sqlx.DB

func init() {
	// 连接
	dbuser := "root"
	dbpwd := "somepass"
	dbhost := "127.0.0.1"
	dbname := "shopping_mall"
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbuser, dbpwd, dbhost, dbname)
	/* way1 打开db不连接
	db := sqlx.Open("mysql", dns)
	db = sqlx.NewDb(db = sqlx.Open("mysql", dns))
	err = db.Ping()  //强制连接并测试它是否工作
	*/
	/* way2 打开一个DB并同时连接
	db, err := sqlx.Connect("mysql", dns)
	*/
	// way3 打开、连接数据库的同时panic error
	db := sqlx.MustConnect("mysql", dns)

	Db = db
}

/*
func TestAddBDBasic(t *testing.T){
	defer Db.Close()

	//事务
	tx := Db.MustBegin()
	// 插数据
	// MustExec 等同 原生Exec,并错误时painc
	tx.MustExec(`INSERT INTO user (id,name) VALUES ('0004', 'Jack');`)
	//tx.MustExec(`INSERT INTO demo (id,name) VALUES (?, ?);`, "2", "Emily") // ? 是占位符 可以防止SQL注入攻击
	// 回滚
	//tx.Rollback()
	// 提交
	tx.Commit()
}
*/

/*
func TestAddBDBasicJudge(t *testing.T){
	defer Db.Close()

	//事务
	tx, err := Db.Begin()
	if err != nil {
		return
	}

	// 插数据
	row, err := Db.Exec(`INSERT INTO user (id,name) VALUES ('0004', 'Jack');`)
	if err != nil {
		fmt.Println("exec failed,", err)
		//回滚事务
		tx.Rollback()
		return
	}
	fmt.Println("row ,", row)

	// TODO id 获取不到
	//id, err := row.LastInsertId()
	//fmt.Println("insert succ:", id)

	//提交事务
	err = tx.Commit()
	if err != nil {
		fmt.Println("Commit failed,", err)
		return
	}
}
*/

/*
func TestAddBDStruct(t *testing.T){
	defer Db.Close()

	//事务
	tx, err := Db.Beginx()
	if err != nil {
		fmt.Println("Beginx failed,", err)
		return
	}

	// 插数据
	user := User{
		ID: "0002",
		Name: "name",
	}
	// TODO tx.NamedExec
	_, err = tx.NamedQuery(`INSERT INTO user (id, name) VALUES (:id, :name)`, user)
	if err != nil {
		fmt.Println("namequery failed,", err)
		//回滚事务
		err = tx.Rollback()
		if err != nil {
			fmt.Println("Rollback failed,", err)
			return
		}
		return
	}

	// 提交
	err = tx.Commit()
	if err != nil {
		fmt.Println("Commit failed,", err)
		return
	}
}
*/

// TODO TestAddBDMultiStruct

func TestQueryBDBasic(t *testing.T){
	defer Db.Close()

	rows, err := Db.Query("SELECT id, name FROM user")
	if err != nil {
		fmt.Println("Query failed,", err)
	}

	for rows.Next() {
		var id string
		var name string
		//var createdAt time.Time
		var deletedAt sql.NullTime
		err = rows.Scan(&id, &name, &deletedAt)
		if err != nil {
			fmt.Println("rows failed,", err)
		}else{
			fmt.Printf("\n%v, %v, %v ", id, name, deletedAt)
		}
	}
}
/*
func TestQueryxBDBasic(t *testing.T){
	defer Db.Close()

	// Queryx(...) (*sqlx.Rows, error) 基本等同 Query(...) (*sql.Rows, error), but return an sqlx.Rows
	rows, err := Db.Queryx("SELECT id, name, deleted_at FROM user")
	if err != nil {
		fmt.Println("Queryx failed,", err)
	}


	for rows.Next() {
		var user User
		err = rows.StructScan(&user)
		if err != nil {
			fmt.Println("rows failed,", err)
		}else{
			//fmt.Printf("rows: \n%+v\n", user) // {Key:键值 Value:数据}
			bytes, _ := json.Marshal(user)
			fmt.Println(string(bytes)) // {"Key":"键值","Value":"数据"}
		}
	}
}

*/