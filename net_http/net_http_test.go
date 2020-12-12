package net_http_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var client = &http.Client{}
const urlGetStr = "http://unpkg.com/goclub@0.0.1/package.json"

// GET
func TestRequest_getBody(t *testing.T)  {
	resp, err := http.Get(urlGetStr);if err !=nil {panic(err)}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body);if err !=nil {panic(err)}
	log.Print(string(body))
}
/*Client的Transport字段一般会含有内部状态（缓存TCP连接），因此Client类型值应尽量被重用而不是每次需要都创建新的*/
func TestRequest_client_get(t *testing.T){
	resp, err := client.Get(urlGetStr);if err !=nil {panic(err)}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body);if err !=nil {panic(err)}
	log.Print(string(body))
}
func TestRequest_client_do(t *testing.T){
	u,err := url.Parse(urlGetStr);if err !=nil {panic(err)}
	var req = &http.Request{
		Method:"GET",
		URL: u,
	}
	resp, err := client.Do(req);if err !=nil {panic(err)}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body);if err !=nil {panic(err)}
	log.Print(string(body))
	//log.Print(resp.Header)
}

// POST
func TestRequest_postBodyJson(t *testing.T)  {
	// 请求数据
	req := struct{
		ID string `json:"jsonID" form:"formID"`
	}{
		ID:"id_1",
	}
	reqByte, err:= json.Marshal(&req);if err !=nil {panic(err)}
	reqBody := bytes.NewReader(reqByte)
	resp, err := http.Post("http://127.0.0.1:1219/post_json", "application/json", reqBody);if err !=nil {panic(err)}
	//resp, err := http.Post("http://127.0.0.1:1219/post_json", "multipart/form-data;", reqBody);if err !=nil {panic(err)}
	//resp, err := http.Post("http://127.0.0.1:1219/post_json", "application/x-www-form-urlencoded;", reqBody);if err !=nil {panic(err)}
	defer resp.Body.Close()

	// 读值
	body, err := ioutil.ReadAll(resp.Body);if err !=nil {panic(err)}
	log.Print("body: ",string(body))
	// 绑定结构体
	var res struct{
		Type string `json:"type"`
		Data struct{
			Name string `json:"name"`
			Age int `json:"age"`
		} `json:"data"`
	}
	err =json.Unmarshal(body, &res);if err !=nil {panic(err)}
	log.Print("res: ",res)
}

func TestRequest_postBodyForm(t *testing.T)  {
	// 请求数据
	v := url.Values{}
	v.Set("name", "free")
	log.Print(v.Encode())
	reqBody := strings.NewReader(v.Encode())

	resp, err := http.Post("http://127.0.0.1:1219/post_json", "application/x-www-form-urlencoded", reqBody);if err !=nil {panic(err)}
	defer resp.Body.Close()
}
// go响应/发送图片
// go接受/发送文件
// cookie操作