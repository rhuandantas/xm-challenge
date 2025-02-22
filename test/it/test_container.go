package it

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type TestContainer struct {
	Container testcontainers.Container
	Addresses []string
	URI       string
	Host      string
	Port      string
}

func SetupMySQLContainer(ctx context.Context) (*TestContainer, error) {
	mysqlPort := "3306"
	natPort := nat.Port(mysqlPort + "/tcp")

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{mysqlPort + "/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "test_db",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306"),
			wait.ForSQL(natPort, "mysql", func(host string, port nat.Port) string {
				return fmt.Sprintf("root:root@tcp(localhost:%s)/test_db", port.Port())
			}).WithStartupTimeout(time.Second*30),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(mysqlPort))
	if err != nil {
		return nil, err
	}

	return &TestContainer{
		Container: container,
		URI:       fmt.Sprintf("%s:%s", host, mappedPort.Port()),
		Host:      host,
		Port:      mappedPort.Port(),
	}, nil
}

func SetupKafkaContainer(ctx context.Context) (*TestContainer, error) {
	// Start a Kafka test container
	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		_ = fmt.Errorf("Failed to start Kafka container: %v", err)
	}
	// Get the broker address
	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		_ = fmt.Errorf("Failed to get Kafka broker address: %v", err)
	}
	fmt.Println("Kafka Broker:", brokers)

	return &TestContainer{
		Container: kafkaContainer,
		URI:       brokers[0],
		Addresses: brokers,
		Host:      "",
		Port:      "",
	}, nil
}
