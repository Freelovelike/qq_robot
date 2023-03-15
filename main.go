package main

import (
	"context"
	"strings"

	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

var message_id_list []float64

func check_message_id(message_id float64) bool {
	for _, v := range message_id_list {
		if message_id == v {
			return false
		}
	}
	return true
}

type SendQun struct {
	Group_id int64  `json:"group_id"`
	Message  string `json:"message"`
}

func xiaoxi(c *gin.Context) {
	json_parst := make(map[string]interface{})
	err := c.BindJSON(&json_parst)
	if err != nil {
		fmt.Printf("bind error:%v\n", err)
	}
	message_id, ok := json_parst["message_id"].(float64)
	if ok {
		if check_message_id(message_id) {
			message_id_list = append(message_id_list, message_id)
		} else {
			return
		}
	}

	if json_parst["message"] != nil {
		traycode := json_parst["message"].(string)
		qq_id := json_parst["user_id"].(float64)
		prefix := strings.HasPrefix(traycode, "gpt")
		fmt.Print(prefix)
		keycode2 := strings.Split(traycode, " ")[1:]
		keycode := strings.Join(keycode2, "")

		if prefix && qq_id != 1908183918 {
			fmt.Print(json_parst)
			config := openai.DefaultConfig("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")
			proxyUrl, err := url.Parse("http://127.0.0.1:7890")
			if err != nil {
				panic(err)
			}
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			config.HTTPClient = &http.Client{
				Transport: transport,
			}
			client := openai.NewClientWithConfig(config)

			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: keycode,
						},
					},
				},
			)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			var str = strings.Replace(resp.Choices[0].Message.Content, "\n", "", 2)
			// fmt.Println("http://127.0.0.1:5700/send_group_msg?group_id=831905542&message=" + str)

			res, err := http.PostForm("http://127.0.0.1:5700/send_group_msg", url.Values{"group_id": {"831905542"}, "message": {str}})
			if err != nil {
				fmt.Printf("%v", err)
			}
			if res != nil {
				// json := make(map[string]interface{})
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "hey",
		"status":  http.StatusOK,
	})

}
func main() {

	//注册路由

	r := gin.Default()
	r.POST("/sangbao", xiaoxi)
	r.Run("127.0.0.1:5701")
}
