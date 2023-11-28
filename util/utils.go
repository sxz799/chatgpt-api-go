package util

import (
	"bytes"
	"chatgpt-api-go/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendChatPostMsg(msgs []model.Message, conf model.ApiConfig) (string, error) {
	if len(msgs) > conf.HistoryNumber {
		msgs = msgs[conf.HistoryNumber:]
	}

	reqData := model.Request{
		Model:    conf.Model,
		Messages: msgs,
		Stream:   true,
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

	//用channel来接收数据
	ch := make(chan string, 10)
	go func() {
		defer close(ch)
		for {
			buf := make([]byte, 2048)
			n, err := resp.Body.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				ch <- string(buf[:n])
			}
		}
	}()
	lastStr := ""
	for msg := range ch {
		msg = strings.ReplaceAll(msg, "\n", "")
		msg = strings.ReplaceAll(msg, "data: [DONE]", "")
		if lastStr != "" {
			msg = lastStr + msg
		}
		if strings.HasPrefix(msg, "d") && strings.HasSuffix(msg, "}") {
			lastStr = ""
			ss := strings.Split(msg, "data: ")
			for _, s := range ss {
				if s == "" {
					continue
				}
				var respData model.ChatCompletionChunk
				json.Unmarshal([]byte(s), &respData)
				if len(respData.Choices) > 0 {
					fmt.Print(respData.Choices[0].Delta.Content)
				}
			}
		} else {
			lastStr = msg
			continue
		}
	}

	return "", nil
}
