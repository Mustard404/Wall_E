package ampq

import (
	"Agent/common"
	"Agent/plugins"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

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

// InitProducer 初始化MQ生产者，并创建发送和接收队列的客户端
func InitProducer() {
	// 创建发送到"port_scan_sent"队列的客户端
	sentClient, err := NewRabbitMQClient("port_scan_sent")
	if err != nil {
		log.Fatalf("[InitProducer]: %v", err)
	}
	defer sentClient.Close()

	// 创建接收"port_scan_return"队列的客户端
	returnClient, err := NewRabbitMQClient("port_scan_return")
	if err != nil {
		log.Fatalf("[InitProducer]: %v", err)
	}
	defer returnClient.Close()

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

	log.Println("Message published successfully")
}

// Consume 开始消费队列中的消息
func (c *RabbitMQClient) Consume() {
	// 消息处理函数
	msgs, err := c.ch.Consume(
		c.queue, // 队列名称
		"",      // 消费者标签
		false,   // 不自动应答
		false,   // 非独占
		false,   // 不等待服务器确认
		false,   // 不no-local
		nil,     // 额外的参数
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 打印接收到的消息内容
			// log.Printf("Received a message: %s", d.Body)
			// 将 JSON 字符串转回 Asset 切片
			var asset common.Asset
			err := json.Unmarshal(d.Body, &asset)
			if err != nil {
				log.Fatalf("Failed to unmarshal JSON to AssetInfo slice: %v", err)
			}

			Ack := plugins.PortScan(asset)
			// 确认消息（RabbitMQ的自动应答被设置为false）
			d.Ack(Ack)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever // 等待永远，除非上面的goroutine遇到错误或接收到中断信号
}

// Consumer 是一个消费者函数，用于从指定队列中消费消息
func Consumer(queueName string) {
	// 创建一个新的RabbitMQ客户端实例
	client, err := NewRabbitMQClient(queueName)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	// 开始消费队列中的消息
	client.Consume()
}
