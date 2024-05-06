package cron

import (
	"Server/common"
	"Server/dingtalk"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Monitor 定时器与MQ监控
func Monitor() {
	ticker := time.NewTicker(time.Hour * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if mqMessageCountAll() {
				log.Println("[Info]:消息队列无任务，开始消息推送和新一轮扫描")
				// 遍历用户查看是否存在非白名单端口
				for _, user := range common.SelectAllUser() {
					sentMsg, accMsg, errMsg := common.SelectNotWhitePort(user)
					replyMsg := fmt.Sprintf("\n> | **部门** | **IP** | **端口** | **协议** | **状态** | **服务** | **白名单** | **结果** |")
					replyMsg += fmt.Sprintf("\n> |:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|:-----:|")
					replyMsg += accMsg
					replyMsg += errMsg
					replyMsg += fmt.Sprintf("\r\n")
					replyMsg += fmt.Sprintf("\n> **处置人员**: @%s", user.StaffId)
					if sentMsg {
						dingtalk.DingPush(replyMsg, user.StaffId)
					}
				}
				// 开启新一轮扫描
				for _, asset := range common.SelectAllAsset() {
					jsonAsset, _ := json.Marshal(asset)
					//log.Println(jsonAsset)
					common.Producer("port_scan_sent", jsonAsset)
				}
			} else {
				log.Println("[Info]:当前任务未完成，无法开启下一轮扫描！")
			}
		}
	}
}

// mqMessageCountAll MQ中的消息总量
func mqMessageCountAll() bool {
	sentCount, err := common.MQMessageCount("port_scan_sent")
	returnCount, err := common.MQMessageCount("port_scan_return")
	if err != nil {
		fmt.Println("Error getting message count:", err)
	}
	return sentCount+returnCount == 0
}
