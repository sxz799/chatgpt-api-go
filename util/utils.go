package util

import (
	"bytes"
	"chatgpt-api-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func SendChatPostMsg(msgs []model.Message, conf model.ApiConfig) (string, error) {
	if len(msgs) > conf.HistoryNumber {
		msgs = msgs[conf.HistoryNumber:]
	}

	reqData := model.Request{
		Model:    conf.Model,
		Messages: msgs,
	}

	reqBody, _ := json.Marshal(reqData)

	req, err := http.NewRequest("POST", conf.ApiServer, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.ApiKey)

	var client *http.Client

	if conf.ProxyUrl != "" {
		proxyUrl, err := url.Parse(conf.ProxyUrl)
		if err != nil {
			return "", fmt.Errorf("代理读取失败: %v", err)
		}
		client = &http.Client{
			Timeout: conf.Timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		client = &http.Client{
			Timeout: conf.Timeout,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发生了错误: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 请求失败，状态码：%d", resp.StatusCode)
	}

	var respData model.ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("解析响应数据失败: %v", err)
	}

	if len(respData.Choices) > 0 {
		return respData.Choices[0].Message.Role + ":" + respData.Choices[0].Message.Content, nil
	} else {
		return "", errors.New("API 接口访问失败")
	}
}
