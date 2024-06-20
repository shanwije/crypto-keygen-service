FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o crypto-keygen-service cmd/crypto-keygen-service/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/crypto-keygen-service .
COPY .env .env
EXPOSE 8080
CMD ["./crypto-keygen-service"]
