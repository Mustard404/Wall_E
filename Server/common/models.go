package common

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConfig struct {
	Host     string
	Port     string
	ApiPort  string
	User     string
	Password string
	Vhost    string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

type DingTalkConfig struct {
	ClientId     string
	ClientSecret string
	Webhook      string
}

type User struct {
	ID        uint
	StaffId   string `gorm:"unique"`
	Name      string
	AssetInfo []Asset `gorm:"foreignKey:UserID"`
}

type Asset struct {
	ID         uint
	Department string
	IP         string `gorm:"unique"`
	UserID     uint
	PortInfo   []Port `gorm:"foreignKey:AssetID"`
}

type Port struct {
	ID       uint
	Port     int `gorm:"uniqueIndex:idx_unique_port_assetid"`
	Protocol string
	State    string
	Service  string
	White    bool
	AssetID  uint `gorm:"uniqueIndex:idx_unique_port_assetid"`
}

// Ports 使用 AssetID:Port 作为键
type Ports map[string]*Port

func (ps Ports) AddPort(port *Port) error {
	key := fmt.Sprintf("%d:%d", port.AssetID, port.Port)
	if _, exists := ps[key]; exists {
		return errors.New("port with IP and Port already exists")
	}
	ps[key] = port
	return nil
}

// RabbitMQClient 结构体封装了RabbitMQ的连接、通道和队列信息
type RabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}
