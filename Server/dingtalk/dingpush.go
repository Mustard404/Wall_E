package dingtalk

import (
	"Server/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DingTalkMarkdownMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		IsAtAll   bool     `json:"isAtAll,omitempty"`
		AtUserIds []string `json:"atUserIds,omitempty"`
	} `json:"at,omitempty"`
}

// DingPush 使用Webhook进行消息推送
func DingPush(text string, staffId string) {

	message := DingTalkMarkdownMessage{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "端口监控告警",
			Text:  text,
		},
		At: struct {
			IsAtAll   bool     `json:"isAtAll,omitempty"`
			AtUserIds []string `json:"atUserIds,omitempty"`
		}{
			IsAtAll:   false, // 如果想要艾特所有人，将此字段设置为true
			AtUserIds: []string{staffId},
		},
	}

	messageBody, _ := json.Marshal(message)
	reader := bytes.NewReader(messageBody)

	resp, err := http.Post(common.DTConfig.Webhook, "application/json", reader)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
}
