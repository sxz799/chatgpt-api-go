package main

import (
	"bufio"
	"chatgpt-api-go/model"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	// 准备请求数据
	conf := model.ApiConfig{
		Model:         "gpt-3.5-turbo",
		ApiKey:        "1pk-this-is-a-real-free-pool-token-for-everyone",
		ApiServer:     "https://ai.fakeopen.com/v1/chat/completions",
		ProxyUrl:      "",
		HistoryNumber: 6,
	}
	var wg sync.WaitGroup
	var msgs []model.Message
	msgs = append(msgs, model.Message{
		Role:    "user",
		Content: "你好！",
	})
	result, err := model.SendChatPostMsg(msgs, conf)
	fmt.Println(result)
	if err != nil {
		return
	}
	for {
		wg.Add(1)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			question, _ := reader.ReadString('\n')
			if strings.EqualFold(question, "重新开始") {
				msgs = make([]model.Message, 0)
			}
			msgs = append(msgs, model.Message{
				Role:    "user",
				Content: question,
			})
			fmt.Println(model.SendChatPostMsg(msgs, conf))
			wg.Done()
		}()
		wg.Wait()
		fmt.Println("请继续提问: ")
	}

}
