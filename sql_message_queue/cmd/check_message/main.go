package main

import (
	"context"
	"github.com/grifree/code/sql_message_queue/app/connect_sql"
	"github.com/grifree/code/sql_message_queue/app/service"
	"log"
	"time"
)

func main () {
	ctx := context.Background()
	db, err := connect_sql.NewDB();if err != nil {
		panic(err)
	}
	messageService := service.Service{
		DB: db,
	}

	for {
		// 定时检查超时任务
		{
			datetime := time.Now().Add(-10 * time.Second).Format("2006-01-02 15:04:05")
			err := messageService.CheckMessageRetry(ctx, datetime, 10);if err!=nil{
				log.Print("LoopCheckMessageRetry err ", err)
			}
			err = messageService.CheckMessageFail(ctx, datetime, 10, 3);if err!=nil{
				log.Print("LoopCheckMessageRetry err ", err)
			}
		}

		time.Sleep(time.Second *2)
	}

}
