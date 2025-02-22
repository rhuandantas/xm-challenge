package kafka_test

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rhuandantas/xm-challenge/config"
	kafkaClient "github.com/rhuandantas/xm-challenge/internal/adapters/messaging/kafka"
	"time"
)

var _ = Describe("Consumer", func() {
	var (
		consumer kafkaClient.Consumer
		cfg      *config.Config
		producer *kafka.Producer
		topic    string
	)

	BeforeEach(func() {
		cfg = &config.Config{
			Kafka: config.KafkaConfig{
				Brokers: []string{"localhost:9092"},
				GroupID: "test-group",
			},
		}

		var err error
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
		})
		Expect(err).NotTo(HaveOccurred())

		topic = "test-events"
		consumer, err = kafkaClient.NewConsumer(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		consumer.Close()
		producer.Close()
	})

	It("should consume messages successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		messageReceived := make(chan []byte)
		go func() {
			err := consumer.Consume(ctx, []string{topic}, func(msg []byte) error {
				messageReceived <- msg
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
		}()

		testMessage := []byte("test message")
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          testMessage,
		}, nil)
		Expect(err).NotTo(HaveOccurred())

		var receivedMsg []byte
		Eventually(messageReceived, "5s").Should(Receive(&receivedMsg))
		Expect(receivedMsg).To(Equal(testMessage))
	})

	It("should handle context cancellation", func() {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})
		go func() {
			err := consumer.Consume(ctx, []string{topic}, func(msg []byte) error {
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			close(done)
		}()

		cancel()
		Eventually(done).Should(BeClosed())
	})
})
