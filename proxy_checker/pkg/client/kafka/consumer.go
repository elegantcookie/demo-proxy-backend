package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"strings"
)

type IConsumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close()
}

type Client struct {
	reader *kafka.Reader
}

func (c *Client) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(context.Background())

}

func (c *Client) Close() {
	c.reader.Close()
}

func NewClient(ctx context.Context, kafkaURL, topic, groupID string) IConsumer {
	brokers := strings.Split(kafkaURL, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	client := &Client{reader: reader}
	return client
}
