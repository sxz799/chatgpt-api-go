package main

import (
	"bufio"
	"chatgpt-api-go/model"
	"fmt"
	"os"
	"sync"
)

func main() {
	// 准备请求数据
	conf := model.ApiConfig{
		Model:         "gpt-3.5-turbo",
		ApiKey:        "pk-this-is-a-real-free-pool-token-for-everyone",
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
	fmt.Println(result + "(您可以输入三个空格以开始一段新的对话！)")
	if err != nil {
		return
	}
	for {
		wg.Add(1)
		fmt.Println("请继续提问: ")
		go func() {
			reader := bufio.NewReader(os.Stdin)
			question, _ := reader.ReadString('\n')
			if question == "   \n" {
				msgs = make([]model.Message, 0)
				fmt.Println("好的，现在重新开始对话！ ")
				wg.Done()
				return
			} else {
				msgs = append(msgs, model.Message{
					Role:    "user",
					Content: question,
				})
			}

			msg, tErr := model.SendChatPostMsg(msgs, conf)
			if tErr != nil {
				fmt.Println(tErr.Error())
			}
			fmt.Println(msg)
			wg.Done()
		}()
		wg.Wait()
	}

}
