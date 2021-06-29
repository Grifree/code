package net_http_test

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGoRoutine_Retry(t *testing.T) {
	c := 0
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {c++}()
		time.Sleep(time.Second)
		switch c {
		case 0:
			writer.WriteHeader(500)
		case 1:
			writer.WriteHeader(501)
		case 2:
			writer.WriteHeader(200)
		case 3:
			writer.WriteHeader(200)
		default:
			writer.WriteHeader(200)
		}
		_, _ = writer.Write([]byte("abc"))
	})
	go func() {
		log.Print(http.ListenAndServe(":1111", nil))
	}()
	//ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()


	data, err := retry("http://127.0.0.1:1111/test")
	log.Print(string(data))
	log.Print(err)

}

func send(url string) (data []byte, err error) {

	resp, err := http.Get(url);if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("status fail")
	}
	return ioutil.ReadAll(resp.Body)
}

func retry(url string) (data []byte, err error) {
	timer := time.NewTimer(time.Millisecond*1500)
	for retryCount:=0;retryCount < 3;retryCount++ {
		select {
			case <-timer.C:
				return nil, errors.New("timer over")
		default:
			log.Print(url, " ",retryCount)
			data, err = send(url);if err != nil {
				continue
			}
			return
		}
	}
	return
}

func TestRetry_time(t *testing.T){
	timer := time.NewTimer(time.Second)
	log.Print(1)
	log.Print(<- timer.C)
	log.Print(2)
}

type RecoverError struct {
	Err error
	RecoverValue interface{}
}
func Go(routine func () (err error)) (reCh chan RecoverError) {
	reCh = make(chan RecoverError)
	go func() {
		re := RecoverError{}
		defer func() {
			re.RecoverValue = recover()
			reCh <- re
		}()
		re.Err = routine()
	}()
	return reCh
}
func TestRetry_GoRoutine (t *testing.T) {
	dataCh := make(chan []byte)
	reCh := Go(func() (err error) {
		var data []byte
		data, err = send("");if err != nil {
			return  err
		}
		dataCh <- data
		return nil
	})
	select {
	case re := <-reCh:
		if re.Err != nil {

		}
		if re.RecoverValue != nil {

		}
	case data := <- dataCh:
		log.Print(string(data))
	}

}