package common

import (
	"fmt"
	"github.com/go-ini/ini"
)

var AmpqConfig *AMQPConfig
var DBConfig *DatabaseConfig
var DTConfig *DingTalkConfig

// Config 从 config.ini 文件中加载配置并返回相应的结构体实例
func Config() {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Print(err)
	}

	amqpSection := cfg.Section("AMQP")
	AmpqConfig = &AMQPConfig{
		Host:     amqpSection.Key("host").String(),
		Port:     amqpSection.Key("port").String(),
		ApiPort:  amqpSection.Key("api_port").String(),
		User:     amqpSection.Key("user").String(),
		Password: amqpSection.Key("password").String(),
		Vhost:    amqpSection.Key("vhost").String(),
	}

	databaseSection := cfg.Section("Database")
	DBConfig = &DatabaseConfig{
		Host:     databaseSection.Key("host").String(),
		Port:     databaseSection.Key("port").String(),
		DBName:   databaseSection.Key("dbname").String(),
		User:     databaseSection.Key("user").String(),
		Password: databaseSection.Key("password").String(),
	}

	dingtalkSection := cfg.Section("DingTalk")
	DTConfig = &DingTalkConfig{
		ClientId:     dingtalkSection.Key("clientId").String(),
		ClientSecret: dingtalkSection.Key("clientSecret").String(),
		Webhook:      dingtalkSection.Key("Webhook").String(),
	}
}
