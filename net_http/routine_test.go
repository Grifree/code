package net_http_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestGoRoutine(t *testing.T) {

	client := http.Client{}

	listCh := make(chan string)
	errCh := make(chan error)

	go func() {
		data, err := quest(client,"https://mockend.com/goclub/http/post/1");if err != nil {
			errCh <- err
			return
		}
		listCh <- data
	}()
	go func() {
		data, err := quest(client,"https://mockend.com/goclub/http/post/2");if err != nil {
			errCh <- err
			return
		}
		listCh <- data
	}()

	var list []string
	var err error
	for i:=0;i<2;i++{
		select {
			case err = <- errCh:
				log.Print(1)
			case listItem := <- listCh:
				log.Print(2)
				list = append(list, listItem)
		}
	}

	log.Print(err)
	log.Print(list)
}

func quest(client http.Client, url string) (data string, err error) {
	//return "",errors.New("test")
	resp, err := client.Get(url);if err != nil {
		return
	}
	//resp.StatusCode
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body);if err != nil {
		return
	}
	return string(content), nil
}