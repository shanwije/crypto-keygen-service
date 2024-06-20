# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /build && go build -o /build/crypto-keygen-service cmd/crypto-keygen-service/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /build/crypto-keygen-service .

COPY .env .env

EXPOSE 8080

CMD ["./crypto-keygen-service"]
