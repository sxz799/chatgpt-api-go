package main

import (
	"bufio"
	"chatgpt-api-go/model"
	"chatgpt-api-go/util"
	"fmt"
	"os"
	"time"
)

var (
	chatModel     = "gpt-4"
	apiKey        = "sk-xxxxxxxxxxxxxxxxxxxxxxxxxx"
	apiServer     = "https://api.openai.com/v1/chat/completions"
	proxyUrl      = ""
	historyNumber = 6
	timeout       = 60 * time.Second
)

func main() {
	util.Conf = model.ApiConfig{
		Model:         chatModel,
		ApiKey:        apiKey,
		ApiServer:     apiServer,
		ProxyUrl:      proxyUrl,
		HistoryNumber: historyNumber,
		Timeout:       timeout,
	}

	fmt.Println("您可以输入三个空格以开始新的对话！")

	for {
		fmt.Print("You:")
		reader := bufio.NewReader(os.Stdin)
		question, _ := reader.ReadString('\n')

		if question == "   \n" {
			fmt.Println("好的，现在重新开始对话！")
			util.Msgs = []model.Message{}
		} else {
			util.Msgs = append(util.Msgs, model.Message{
				Role:    "user",
				Content: question,
			})
			result, err := util.SendChatPostMsg()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(result)
		}
	}
}
