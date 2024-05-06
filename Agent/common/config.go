package common

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
)

var AmpqConfig *AMQPConfig
var MSConfig *MasscanConfig

// Config loads configuration from config.ini file into corresponding struct instances
func Config() {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Print(err)
	}

	amqpSection := cfg.Section("AMQP")
	AmpqConfig = &AMQPConfig{
		Host:     amqpSection.Key("host").String(),
		Port:     amqpSection.Key("port").String(),
		User:     amqpSection.Key("user").String(),
		Password: amqpSection.Key("password").String(),
		Vhost:    amqpSection.Key("vhost").String(),
	}

	masscanSection := cfg.Section("Masscan")
	rateValue, err := masscanSection.Key("rate").Int()
	if err != nil {
		log.Printf("无法获取 Masscan 的 Rate 值: %v", err)
	}

	MSConfig = &MasscanConfig{
		Rate: rateValue,
	}

}
