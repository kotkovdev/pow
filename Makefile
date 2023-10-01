.phony: server
server: 
	go run cmd/server/main.go

.phony: client
client:
	go run cmd/client/main.go

.phony: test
test:
	go test ./... -race

.phony: lint
lint: 
	golangci-lint run ./...