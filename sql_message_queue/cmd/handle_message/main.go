package main

import (
	"context"
	sq "github.com/goclub/sql"
	"github.com/grifree/code/sql_message_queue/app/connect_sql"
	"github.com/grifree/code/sql_message_queue/app/service"
	"log"
	"math/rand"
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
		defer func() {
			time.Sleep(2*time.Second)
		}()
		ownerKey := sq.UUID()
		// 消费信息
		{
			message, hasMessage, err := messageService.ConsumeMessage(ctx, ownerKey);if err != nil {
				log.Print("ConsumeMessage ", err)
				continue
			}
			if hasMessage == false {
				continue
			}
			{
				// 模拟异步处理消息
				log.Print("message ", message)
				time.Sleep(2*time.Second)
			}
			// 随机模拟部分消息完成, 部分超时
			mockNetInterrupt := false
			if rand.Int()%2 == 0 {
				log.Print("网络中断了")
				mockNetInterrupt = true
			}
			if mockNetInterrupt == false {
				log.Print("消息完成了")
				err = messageService.DoneMessage(ctx, ownerKey);if err != nil {
					log.Print(err)
					continue
				}
			}
		}
	}
}
