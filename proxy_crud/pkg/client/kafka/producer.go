package kafka

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
)

type IProducer interface {
	Produce(ctx context.Context, key, value []byte) error
	Close()
}

type Client struct {
	writer *kafka.Writer
}

func (c *Client) Produce(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	fmt.Println(c.writer.Topic)
	fmt.Println(c.writer.Addr)
	err := c.writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.writer.Close()
}

func NewClient(ctx context.Context, kafkaURL, topic string) IProducer {
	log.Println(kafkaURL, topic)
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	client := &Client{writer: writer}
	return client
}
