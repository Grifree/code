package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	sq "github.com/goclub/sql"
	"log"
	"sync"
)

func main() {
	// 生成5条数据，修改3条数据（只会有3条被修改，不会因为并发导致只修改1条或2条）
	test(Data{
		TableName: "queue_pop_5_3",
		DataRows:   5,
		UpdateRows: 3,
	})
	// 生成5条数据，尝试修改10条（只有5次会反馈影响了数据）
	test(Data{
		TableName: "queue_pop_5_10",
		DataRows:   5,
		UpdateRows: 10,
	})
}
type Data struct {
	TableName string
	DataRows int
	UpdateRows int
}

func test(data Data) {
	log.Print("tableName", data.TableName)
	db,dbClose, err := sq.Open("mysql", sq.MysqlDataSource{
		User:     "root",
		Password: "somepass",
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "goclub_example",
	}.String()) ; if err != nil {
		panic(err)
	}

	defer dbClose()
	ctx := context.Background()
	_, err = db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS  ` + data.TableName + ` (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		used tinyint(4) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`, nil) ; if err != nil {
		log.Print(err); return
	}
	_, err = db.Exec(ctx, `DELETE FROM ` + data.TableName, nil) ; if err != nil {
		log.Print(err); return
	}
	for i:=0;i<data.DataRows;i++{
		_, err := db.Exec(ctx, `INSERT INTO ` + data.TableName + ` (id, used) VALUES (?, 0)`, []interface{}{i+1}) ; if err != nil {
			log.Print(err); return
		}
	}
	// 因为要使用 routine 模拟并发，所以用 WaitGroup 防止还没执行完sql程序就退出了
	wg := sync.WaitGroup{}
	for i:=0;i<data.UpdateRows;i++ {
		wg.Add(1)
		go func() {
			defer  wg.Done()
			setUsed := 1
			whereUsed := 0
			result, err := db.Exec(ctx, "UPDATE " + data.TableName +" SET used = ? WHERE used = ? LIMIT 1", []interface{}{setUsed, whereUsed}) ; if err != nil {
				log.Print(err); return
			}
			affected, err := result.RowsAffected() ; if err != nil {
				log.Print(err); return
			}
			log.Print("affected: ", affected)

		}()
	}
	wg.Wait()
}