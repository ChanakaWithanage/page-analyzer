APP := page-analyzer
BIN := bin/$(APP)

.PHONY: run build test tidy

run:
	go run ./cmd/web

build:
	GOOS=linux GOARCH=amd64 go build -o $(BIN) ./cmd/web

test:
	go test ./...

tidy:
	go mod tidy
