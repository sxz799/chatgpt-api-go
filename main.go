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
	chatModel     = "gpt-3.5-turbo"
	apiKey        = "fk-AzZ6z5NrR94_PGP-Mzttd3r7JMedvwc9ZhYmy_hROfM"
	apiServer     = "https://ai.fakeopen.com/v1/chat/completions"
	proxyUrl      = ""
	historyNumber = 6
	timeout       = 60 * time.Second
)

func main() {
	conf := model.ApiConfig{
		Model:         chatModel,
		ApiKey:        apiKey,
		ApiServer:     apiServer,
		ProxyUrl:      proxyUrl,
		HistoryNumber: historyNumber,
		Timeout:       timeout,
	}

	var msgs []model.Message
	fmt.Println("您可以输入三个空格以开始新的对话！")

	for {
		fmt.Println("请提问:")
		reader := bufio.NewReader(os.Stdin)
		question, _ := reader.ReadString('\n')

		if question == "   \n" {
			fmt.Println("好的，现在重新开始对话！")
			msgs = nil
		} else {
			msgs = append(msgs, model.Message{
				Role:    "user",
				Content: question,
			})
			result, err := util.SendChatPostMsg(msgs, conf)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(result)
		}
	}
}
