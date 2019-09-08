package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Producer *kafka.Producer
}

type ProducerConfig struct {
	Host  string `json:"host"`
	NumWk int    `json:"numwk"`
	Topic string `json:"topic"`
}

func NewProducer(config ProducerConfig) (Producer, error) {
	kafkaConn := Producer{}
	var err error
	kafkaConn.Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Host,
		"group.id":          "g1",
		"auto.offset.reset": "earliest",
	})

	return kafkaConn, err
}

func (c *Producer) Publish(topic string, msg string) error {

	c.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)
	return nil
}
