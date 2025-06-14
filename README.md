# Go Kafka Gateway

[![Go Report Card](https://goreportcard.com/badge/github.com/BlurredDev/go-kafka-gateway)](https://goreportcard.com/report/github.com/BlurredDev/go-kafka-gateway)
[![Build Status](https://github.com/BlurredDev/go-kafka-gateway/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/BlurredDev/go-kafka-gateway/actions/workflows/test.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](./coverage.html)
[![License](https://img.shields.io/github/license/BlurredDev/go-kafka-gateway)](https://github.com/BlurredDev/go-kafka-gateway/blob/main/LICENSE)

---

## 🚀 Purpose

**Go Kafka Gateway** is a high-performance HTTP gateway written in Go that accepts arbitrary JSON payloads and forwards them directly to a Kafka topic. It is designed to simplify ingestion pipelines by offloading validation and processing to downstream systems.

---

## ⚙️ Features

- ✅ Accepts dynamic JSON with no schema requirements  
- ✅ Sends raw messages directly to Kafka  
- ✅ Built in Go for performance and simplicity  
- ✅ Stateless, production-ready microservice  
- ✅ Easily containerized (UBI 9 base for enterprise use)  
- ✅ Automatic fallback to Dead Letter Queue (DLQ) on message delivery failure  

---

## 🐳 Build & Run with Docker


### 🧰 Dockerfile (UBI 9-based)
```dockerfile
# Stage 1: Build the Go binary
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway ./cmd/gateway

# Stage 2: Minimal runtime using UBI 9
FROM registry.access.redhat.com/ubi9/ubi-minimal

WORKDIR /srv/gateway
COPY --from=builder /app/gateway .

EXPOSE 8080
ENTRYPOINT ["./gateway"]
```
⸻

###  🔨 Build the image

```bash
docker build -t BlurredDev/go-kafka-gateway .
```
⸻
### ▶️ Run the container
```bash

docker run -p 8080:8080 \
  -e KAFKA_BROKER=kafka.example.com:9092 \
  -e KAFKA_TOPIC=my-topic \
  -e KAFKA_DLQ_TOPIC=my-dlq-topic \
  BlurredDev/go-kafka-gateway
```

⸻

### 📬 Example Request
```bash
curl -X POST http://localhost:8080/publish \
  -H "Content-Type: application/json" \
  -d '{"event":"signup", "user":{"id":1,"name":"Alice"}}'
```

⸻

### DLQ (Dead Letter Queue) Support

The gateway supports a Dead Letter Queue (DLQ) to improve reliability. If a message fails to be delivered to the primary Kafka topic (e.g., due to broker issues or serialization errors), the message will automatically be sent to the configured DLQ topic. This ensures no data is lost and allows downstream systems to handle or inspect problematic messages separately.

A message will be routed to the DLQ if the attempt to publish to the main Kafka topic fails — typically due to:
- Kafka broker/network errors
- Partition unavailability
- Message serialization failures

🔎 DLQ logic is implemented in [`internal/kafka/producer.go`](./internal/kafka/producer.go)

💡 If `KAFKA_DLQ_TOPIC` is not set, failed messages will be dropped and logged instead of retried.

---

🔧 Environment Variables

| Variable       | Description                              | Default     |
| -------------- |:-------------------------------------:|:-----------:|
| KAFKA_BROKER   | Kafka bootstrap server address          | kafka:9092  |
| KAFKA_TOPIC    | Kafka topic to publish messages to      | (required)  |
| KAFKA_DLQ_TOPIC| Kafka Dead Letter Queue topic for fallback messages | (optional) |
| HTTP_ADDR      | HTTP server binding address              | :8080       |


⸻

👥 Contributing

We welcome PRs! Here’s how:
	1.	Fork this repo
	2.	Create a feature branch (git checkout -b my-feature)
	3.	Commit your changes
	4.	Push to your fork and submit a Pull Request

⸻

📄 License

This project is licensed under the MIT License.

⸻

💬 Questions?
	•	Open an issue
	•	Start a discussion

----
