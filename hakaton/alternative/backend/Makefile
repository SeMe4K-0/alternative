.PHONY: docker-up docker-down install test test-coverage test-verbose

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

install:
	go mod download
	go mod tidy

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
