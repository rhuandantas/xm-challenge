run:
	go run cmd/main.go

build:
	docker-compose build

up:
	docker-compose up -d

install:
	go get github.com/onsi/ginkgo/v2@latest
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

mock:
	mockgen -source=internal/adapters/repository/company.go -package=mock_mysql -destination=test/mock/repository/company.go
	mockgen -source=internal/adapters/messaging/kafka/producer.go -package=mock_kafka -destination=test/mock/kafka/producer.go

unit-test: install
	go test -v ./...

run-linter:
	golangci-lint run