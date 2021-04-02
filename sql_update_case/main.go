package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	sq "github.com/goclub/sql"
	"strconv"
)
func main(){
	accountID := sq.UUID()
	test(Data{
		Amount:10,
		AccountID:accountID,
		CostAmount:3,
	})
}
type Data struct {
	AccountID string
	Amount float64
	CostAmount float64
}
func test(data Data) {

}

/* sql */
type AccountBalanceTable struct {
	sq.DefaultLifeCycle
	sq.SoftDeleteDeletedAt
}
func (AccountBalanceTable) TableName() string {return "balance"}
type IDAccount string
type AccountBalance struct {
	AccountBalanceTable
	sq.CreatedAtUpdatedAt

	ID IDAccount
	Amount float64
}
func (AccountBalance) Column()(e struct{
	ID sq.Column
	Amount sq.Column
}){
	e.ID = "id"
	e.Amount = "amount"
	return
}

/* 乐观锁 扣费 */
func BalanceCost(ctx context.Context, accountID string, costAmount float64) (err error) {
	col := AccountBalance{}.Column()
	result, err := db.Update(ctx, sq.QB{
		Table:AccountBalanceTable{},
		Update: []sq.Update{
			sq.Update{
				Raw:sq.Raw{
					Query:string(col.Amount)+" = "+string(col.Amount)+" - "+strconv.FormatFloat(costAmount, 'E', -1, 64),
				},
			},
		},
		Where: sq.And(col.ID, sq.Equal(accountID)).
			And(col.Amount, sq.GtOrEqualFloat(costAmount)),
		Limit:1,
	})
	affectRow, err := result.RowsAffected();if err != nil {
		return
	}
	switch affectRow {
		case 0:
			err = errors.New("账户没有足够的金额");return
		case 1:
			return
		default:
			err = errors.New("affected must be 0 or 1");return
	}
}

var db *sq.Database
func init(){
	db,_, err := sq.Open("mysql", sq.MysqlDataSource{
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
}