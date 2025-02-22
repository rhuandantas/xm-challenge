package kafka

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rhuandantas/xm-challenge/config"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Produce(ctx context.Context, topic string, key string, value interface{}) error
	Close()
}

type producer struct {
	producer *kafka.Producer
}

func NewProducer(config *config.Config) (Producer, error) {
	brokers := config.Kafka.Brokers
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"acks":              "all",
		"security.protocol": "plaintext",
	})
	if err != nil {
		return nil, err
	}

	return &producer{
		producer: p,
	}, nil
}

func (p *producer) Produce(ctx context.Context, topic string, key string, value interface{}) error {
	log.Info().Msgf("publishing message to %s", topic)
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          jsonValue,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func (p *producer) Close() {
	p.producer.Close()
}
