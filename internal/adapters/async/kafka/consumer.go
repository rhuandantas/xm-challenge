package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rhuandantas/xm-challenge/config"
	"time"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler func([]byte) error) error
	Close()
}

type consumer struct {
	consumer *kafka.Consumer
	config   *config.Config
}

func NewConsumer(config *config.Config) (Consumer, error) {
	brokers := config.Kafka.Brokers
	groupID := config.Kafka.GroupID
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
		"security.protocol": "plaintext",
	})
	if err != nil {
		return nil, err
	}

	return &consumer{
		consumer: c,
		config:   config,
	}, nil
}

func (c *consumer) Consume(ctx context.Context, topics []string, handler func([]byte) error) error {
	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.ReadMessage(time.Second * 1)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				return err
			}

			if err := handler(msg.Value); err != nil {
				return err
			}
		}
	}
}

func (c *consumer) Close() {
	c.consumer.Close()
}
