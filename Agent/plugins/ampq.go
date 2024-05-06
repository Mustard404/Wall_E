package plugins

import (
	"Agent/common"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// RabbitMQClient 结构体封装了RabbitMQ的连接、通道和队列信息
type RabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Close 关闭RabbitMQClient实例，包括关闭通道和连接
func (c *RabbitMQClient) Close() error {
	if err := c.ch.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	return c.conn.Close()
}

// NewRabbitMQClient 创建一个新的RabbitMQ客户端实例
func NewRabbitMQClient(queueName string) (*RabbitMQClient, error) {
	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial("amqp://" + common.AmpqConfig.User + ":" + common.AmpqConfig.Password + "@" + common.AmpqConfig.Host +
		":" + common.AmpqConfig.Port + "/" + common.AmpqConfig.Vhost)
	if err != nil {
		return nil, err
	}

	// 在连接上创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // 如果创建通道失败，关闭连接
		return nil, err
	}

	// 声明一个队列，如果队列不存在则创建它
	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()   // 如果声明队列失败，关闭通道
		conn.Close() // 关闭连接
		return nil, err
	}

	// 返回RabbitMQClient实例
	return &RabbitMQClient{
		conn:  conn,
		ch:    ch,
		queue: queue.Name, // 队列名称
	}, nil
}

// Publish 发布消息到RabbitMQ队列
func (c *RabbitMQClient) Publish(message []byte) error {
	return c.ch.PublishWithContext(
		context.Background(),
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

// Producer MQ生产者
func Producer(queueName string, message []byte) {
	// 创建一个新的RabbitMQ客户端实例
	client, err := NewRabbitMQClient(queueName)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	// 发布消息到队列
	err = client.Publish(message)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}
