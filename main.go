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
		Model:    "gpt-3.5-turbo",
		ApiKey:   "sk-",
		ProxyUrl: "http://127.0.0.1:7890",
	}
	var wg sync.WaitGroup
	var msgs []model.Message
	fmt.Println("请输入你要提问的内容！")
	for {
		wg.Add(1)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			question, _ := reader.ReadString('\n')
			if strings.Contains(question, "重新开始") {
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
	}

}
