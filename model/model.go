package model

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Model         string
	ApiKey        string
	ApiServer     string
	ProxyUrl      string
	HistoryNumber int
}

func SendChatPostMsg(msgs []Message, conf ApiConfig) (string, error) {
	if len(msgs) > conf.HistoryNumber {
		msgs = msgs[conf.HistoryNumber:]
	}
	reqData := Request{
		Model:    conf.Model,
		Messages: msgs,
	}

	// 将请求数据编码为 JSON 格式
	reqBody, _ := json.Marshal(reqData)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", conf.ApiServer, bytes.NewBuffer(reqBody))
	if err != nil {
		return "创建 HTTP 请求失败!", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.ApiKey)

	var client *http.Client
	// 发送请求
	if conf.ProxyUrl != "" {
		proxyUrl, err := url.Parse(conf.ProxyUrl)
		if err != nil {
			return "代理读取失败!", err
		}
		client = &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		client = &http.Client{
			Timeout: 60 * time.Second,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "发生了错误:", err
	}
	defer resp.Body.Close()

	// 解析响应数据
	var respData ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "解析响应数据失败:", err
	}
	if len(respData.Choices) > 0 {
		return respData.Choices[0].Message.Role + ":" + respData.Choices[0].Message.Content, nil
	} else {
		return "api接口访问失败", errors.New("api接口访问失败")
	}

}
