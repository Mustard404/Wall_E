package main

import (
	"Agent/ampq"
	"Agent/common"
)

func main() {
	common.Banner()
	common.Config()
	ampq.InitProducer()
	ampq.Consumer("port_scan_sent")
}
