package main

import (
	"bytes"
	"chatgpt-api-go/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func SendChatPostMsg(msgs []model.Message, conf model.ApiConfig) string {
	reqData := model.Request{
		Model:    conf.Model,
		Messages: msgs,
	}

	// 将请求数据编码为 JSON 格式
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.ApiKey)

	// 发送请求
	proxyUrl, err := url.Parse(conf.ProxyUrl)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "发生了错误:" + err.Error()
	}
	defer resp.Body.Close()

	// 解析响应数据
	var respData model.ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		panic(err)
	}
	return respData.Choices[0].Message.Role + ":" + respData.Choices[0].Message.Content
}

func main() {
	// 准备请求数据
	conf := model.ApiConfig{
		Model:    "gpt-3.5-turbo",
		ApiKey:   "your api key",
		ProxyUrl: "http://127.0.0.1:7890",
	}
	var wg sync.WaitGroup
	var msgs []model.Message
	fmt.Println("请输入你要提问的内容！(你可以输入exitnow退出当前程序)")
	for {
		wg.Add(1)
		go func() {
			var question string
			_, err := fmt.Scanln(&question)
			if err != nil {
				fmt.Println("输入有误")
				wg.Done()
				return
			}
			if question == "exitnow" {
				os.Exit(0)
			}
			msgs = append(msgs, model.Message{
				Role:    "user",
				Content: question,
			})
			fmt.Println(SendChatPostMsg(msgs, conf))
			wg.Done()
		}()
		wg.Wait()
	}

}
