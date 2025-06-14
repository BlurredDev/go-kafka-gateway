# Stage 1: Build the Go binary
FROM golang:1.24.1 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gateway ./cmd/gateway

# Stage 2: Minimal runtime using UBI 9
FROM registry.access.redhat.com/ubi9/ubi-minimal

WORKDIR /srv/gateway
COPY --from=builder /app/gateway .

EXPOSE 8080
ENTRYPOINT ["./gateway"]