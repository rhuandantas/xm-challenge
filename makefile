run:
	go run cmd/main.go

build:
	docker-compose build

up:
	docker-compose up -d

install:
	go get github.com/onsi/ginkgo/v2@latest
	go install github.com/onsi/ginkgo/v2/ginkgo

mock:
	mockgen -source=internal/adapters/repository/cache/memcache.go -package=mock_cache -destination=test/mock/repository/cache/memcache.go

unit-test: install
	go test -v ./...	