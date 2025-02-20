FROM golang:1.23-alpine as base
LABEL authors="rhuan"
RUN apk update

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api

FROM alpine as binary
COPY --from=base /app/api .
EXPOSE 3000

CMD ["./api"]
