# Go Kafka Gateway

![Build](https://img.shields.io/github/actions/workflow/status/BlurredDev/go-kafka-gateway/ci.yml?branch=main)
![Docker Image](https://img.shields.io/docker/pulls/BlurredDev/go-kafka-gateway)
![License](https://img.shields.io/github/license/BlurredDev/go-kafka-gateway)
![Go Version](https://img.shields.io/github/go-mod/go-version/BlurredDev/go-kafka-gateway)
![Last Commit](https://img.shields.io/github/last-commit/BlurredDev/go-kafka-gateway)

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

🔧 Environment Variables

|Variable	|Description	Default|
| ------------- |:-------------:|
|KAFKA_BROKER|	Kafka bootstrap server address	kafka:9092|
|KAFKA_TOPIC|	Kafka topic to publish 
|KAFKA_DQL_TOPIC |Kafka secondary topic to preserve message 
|HTTP_ADDR	|HTTP server binding address	:8080


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
