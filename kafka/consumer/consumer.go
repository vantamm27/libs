package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Consumer *kafka.Consumer
}

type ConsumerConfig struct {
	Host  string `json:"host"`
	NumWk int    `json:"numwk"`
	Topic string `json:"topic"`
}

func NewConsumer(config ConsumerConfig) (Consumer, error) {
	kafkaConn := Consumer{}
	var err error
	kafkaConn.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Host,
		"group.id":          "g1",
		"auto.offset.reset": "earliest",
	})
	return kafkaConn, err
}

func (c *Consumer) Subscribe(topic string) error {
	return c.Consumer.Subscribe(topic, nil)
}
