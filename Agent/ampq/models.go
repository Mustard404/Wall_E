package ampq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient 结构体封装了RabbitMQ的连接、通道和队列信息
type RabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}
