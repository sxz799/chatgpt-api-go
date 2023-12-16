package util

import (
	"bytes"
	"chatgpt-api-go/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Msgs []model.Message
var Conf model.ApiConfig

func SendChatPostMsg() (string, error) {
	if len(Msgs) > Conf.HistoryNumber {
		Msgs = Msgs[Conf.HistoryNumber:]
	}

	reqData := model.Request{
		Model:    Conf.Model,
		Messages: Msgs,
		Stream:   true,
	}

	reqBody, _ := json.Marshal(reqData)

	req, err := http.NewRequest("POST", Conf.ApiServer, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Conf.ApiKey)

	var client *http.Client
	if Conf.ProxyUrl != "" {
		proxyUrl, err := url.Parse(Conf.ProxyUrl)
		if err != nil {
			return "", fmt.Errorf("代理读取失败: %v", err)
		}
		client = &http.Client{
			Timeout: Conf.Timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		client = &http.Client{
			Timeout: Conf.Timeout,
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
		lastStr := ""
		for {
			buf := make([]byte, 2048)
			n, err := resp.Body.Read(buf)
			if err != nil {
				break
			}
			var str = string(buf[:n])
			if lastStr != "" {
				str = lastStr + str
				lastStr = ""
			}
			split := strings.Split(str, "data: ")
			for _, s := range split {
				s = strings.ReplaceAll(s, "\n", "")
				s = strings.ReplaceAll(s, "data: [DONE]", "")
				if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
					var chunk model.ChatCompletionChunk
					err = json.Unmarshal([]byte(s), &chunk)
					if err != nil {
						continue
					}
					if len(chunk.Choices) > 0 {
						ch <- chunk.Choices[0].Delta.Content
					}
				} else {
					lastStr = s
				}
			}

		}
	}()
	result := ""
	for msg := range ch {
		fmt.Print(msg)
		result += msg
		time.Sleep(20 * time.Millisecond)
	}
	Msgs = append(Msgs, model.Message{
		Role:    "assistant",
		Content: result,
	})

	return "", nil
}
