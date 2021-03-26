package net_http_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var client = &http.Client{} //Client的Transport字段一般会含有内部状态（缓存TCP连接），因此Client类型值应尽量被重用而不是每次需要都创建新的
const urlGetStr = "http://unpkg.com/goclub@0.0.1/package.json"

// GET
func TestRequest_get(t *testing.T)  {
	resp, err := http.Get(urlGetStr);if err !=nil {panic(err)}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body);if err !=nil {panic(err)}
	log.Print(string(body))
}
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
	log.Print(resp.Header)
}

// POST
func TestRequest_post_json(t *testing.T)  {
	// 请求数据
	req := struct{
		ID string `json:"jsonID"`
	}{
		ID:"id_1",
	}
	reqByte, err:= json.Marshal(&req);if err !=nil {panic(err)}
	reqBody := bytes.NewReader(reqByte)
	resp, err := http.Post("http://127.0.0.1:1219/post_json", "application/json", reqBody);if err !=nil {panic(err)}
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
func TestRequest_post_form(t *testing.T)  {
	// 请求数据
	reqData := url.Values{}
	reqData.Set("name", "free")
	reqBody := strings.NewReader(reqData.Encode())

	//resp, err := http.Post("http://127.0.0.1:1219/post_json", "application/x-www-form-urlencoded", reqBody);if err !=nil {panic(err)}
	//defer resp.Body.Close()
	httpReq, err := http.NewRequest("POST", "http://127.0.0.1:1219/post_json", reqBody);if err !=nil {panic(err)}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(httpReq);if err !=nil {panic(err)}
	defer resp.Body.Close()
}
func TestRequest_post_multform(t *testing.T)  {
	reqData := &bytes.Buffer{}
	w := multipart.NewWriter(reqData)
	err:=w.WriteField("ttt","multform");if err !=nil {panic(err)}
	err=w.Close();if err !=nil {panic(err)}
	resp, err := http.Post("http://127.0.0.1:1219/post_json", "multipart/form-data;", reqData);if err !=nil {panic(err)}
	defer resp.Body.Close()
}
// go响应/发送图片
// go接受/发送文件
// cookie操作