.PHONY: lint test build run mocks

lint:
	golangci-lint run --fix --verbose --config .golangci.yml

test:
	go test ./...

build:
	go build -o server cmd/server/main.go

run: build
	go run cmd/server/main.go

mocks:
	@mockery

