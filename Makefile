server:
	go run ./cmd/server

agent:
	go run ./cmd/agent --freq 3

test:
	go run ./cmd/test


test-all:
	go test ./...   


	
