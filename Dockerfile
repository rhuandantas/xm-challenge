FROM golang:1.23-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Remove precompiled librdkafka
RUN rm -rf /go/pkg/mod/github.com/confluentinc/confluent-kafka-go*/kafka/librdkafka_vendor/ && \
    go mod download && \
    go mod tidy

# Copy source code
COPY . .

# Build the application with ARM64 architecture
RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig \
    go build -tags dynamic -ldflags="-w -s" -o main .

# Final stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 3001

CMD ["./main"]