package main

import (
	"context"
	sq "github.com/goclub/sql"
	"github.com/grifree/code/sql_message_queue/app/connect_sql"
	pd "github.com/grifree/code/sql_message_queue/app/persisence_data"
	"github.com/grifree/code/sql_message_queue/app/service"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
	
)
func main () {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		reqByte ,err := httputil.DumpRequest(request, false)
		log.Print(string(reqByte), err)
		result := []byte("ok")
		db, err := connect_sql.NewDB();if err != nil {
			panic(err)
		}
		messageService := service.Service{
			DB: db,
		}
		// 每次访问时发布一条消息
		err = messageService.PublishMessage(context.TODO(), pd.Message{
			UserID:sq.UUID(),
			Amount:float64(time.Now().Second()),
		}); if err != nil {
			result = []byte(err.Error())
		}
		_, err = writer.Write(result) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500)
		}
	})
	log.Print(http.ListenAndServe(":1111", nil))
}