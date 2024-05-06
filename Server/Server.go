package main

import (
	"Server/common"
	"Server/cron"
	"Server/dingtalk"
)

func main() {
	common.Banner()
	common.Config()
	common.InitDB()
	common.InitProducer()
	go cron.Monitor()
	go common.Consumer("port_scan_return")
	dingtalk.DingStream()
}
