package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type ApiConfig struct {
	Model    string
	ApiKey   string
	ProxyUrl string
}

func SendChatPostMsg(msgs []Message, conf ApiConfig) string {
	reqData := Request{
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
	var respData ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		panic(err)
	}
	return respData.Choices[0].Message.Role + ":" + respData.Choices[0].Message.Content
}
