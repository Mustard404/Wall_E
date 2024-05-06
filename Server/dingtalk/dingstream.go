package dingtalk

import (
	"Server/common"
	"context"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"strings"
)

func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	// 钉钉身份处理
	user, err := common.SelectUser(data.SenderStaffId, data.SenderNick)
	if err != nil {
		return nil, err
	}

	// 指令处理
	menu := 0
	replyMsg, accMsgs, errMsgs := "", "", ""

	for i, line := range strings.Split(data.Text.Content, "\n") {
		accMsg, errMsg := "", ""
		if i == 0 {
			switch line {
			case "添加资产":
				menu = 1
			case "删除资产":
				menu = 2
			case "更新资产":
				menu = 3
			case "查询资产":
				menu = 4
				accMsg, errMsg = handleSelectAsset(user)
			case "删除端口":
				menu = 5
			case "添加白名单":
				menu = 6
			case "查询端口":
				menu = 7
				accMsg, errMsg = handleSelectPort(user)
			}
		} else {
			switch menu {
			case 1:
				accMsg, errMsg = handleAddAsset(user, line)
			case 2:
				accMsg, errMsg = handleDelAsset(user, line)
			case 3:
				accMsg, errMsg = handleUpAsset(user, line)
			case 5:
				accMsg, errMsg = handleDelPort(user, line)
			case 6:
				accMsg, errMsg = handleWhitePort(user, line)
			}
		}
		accMsgs += accMsg
		errMsgs += errMsg
	}
	switch menu {
	case 0:
		replyMsg += fmt.Sprintf(
			"\n > ## 帮助文档" +
				"\n > - 添加资产" +
				"\n > ```" +
				"\n > 添加资产" +
				"\n > 部门-192.168.1.1" +
				"\n > ```" +
				"\n > - 删除资产" +
				"\n > ```" +
				"\n > 删除资产" +
				"\n > 192.168.1.1" +
				"\n > ```" +
				"\n > - 更新资产" +
				"\n > ```" +
				"\n > 更新资产" +
				"\n > 部门修改-192.168.1.1" +
				"\n > ```" +
				"\n > - 查询资产" +
				"\n > ```" +
				"\n > 查询资产" +
				"\n > ```" +
				"\n > - 删除端口" +
				"\n > ```" +
				"\n > 删除端口" +
				"\n > 192.168.1.1-80" +
				"\n > ```" +
				"\n > - 添加白名单" +
				"\n > ```" +
				"\n > 添加白名单" +
				"\n > 192.168.1.1-80" +
				"\n > ```" +
				"\n > - 查询端口" +
				"\n > ```" +
				"\n > 查询端口" +
				"\n > ```",
		)
	case 1:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|")
	case 2:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|")
	case 3:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|")
	case 4:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|")
	case 5:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **端口** | **协议** | **状态** | **服务** | **白名单** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|")
	case 6:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **端口** | **协议** | **状态** | **服务** | **白名单** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|")
	case 7:
		replyMsg += fmt.Sprintf("\n> | **部门** | **IP** | **端口** | **协议** | **状态** | **服务** | **白名单** | **结果** |")
		replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|")
	}

	replyMsg += accMsgs
	replyMsg += errMsgs
	replier := chatbot.NewChatbotReplier()
	if err := replier.SimpleReplyMarkdown(ctx, data.SessionWebhook, []byte("Wall-E"), []byte(replyMsg)); err != nil {
		return nil, err
	}
	return []byte(""), nil
}

func DingStream() {

	logger.SetLogger(logger.NewStdTestLogger())

	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(
		common.DTConfig.ClientId,
		common.DTConfig.ClientSecret,
	)))
	cli.RegisterChatBotCallbackRouter(OnChatBotMessageReceived)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
