BIN_NAME=crypto-keygen-service

.PHONY: all build test run clean

test:
	go test -v ./...