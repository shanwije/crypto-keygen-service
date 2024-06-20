BIN_NAME=crypto-keygen-service

.PHONY: all build test run clean

run: build
	@env $(shell cat .env | xargs) ./$(BIN_NAME)

test:
	go test -v ./...