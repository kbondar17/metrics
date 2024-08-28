server:
	go run ./cmd/server -k=secret! -d="host=localhost port=6432 dbname=yandex user=postgres password=postgres sslmode=disable"

server-mem:
	go run ./cmd/server -a=localhost:37557


db:
	docker-compose up -d


migr:
	goose -dir internal/database/postgres/migrations postgres "host=localhost port=6432 dbname=yandex user=postgres password=postgres sslmode=disable" up

agent:
	go run -race ./cmd/agent -r=4 -p=2 -k=secret!   -d="host=localhost port=6432 dbname=yandex user=postgres password=postgres sslmode=disable"
 	

temp:
	go run ./cmd/test

ping:
	curl localhost:8080/ping

test:
	go test ./...   

auto:
	/bin/bash test.sh

vet:
	go vet -vettool=statictest ./...

build:
	go build -o ./cmd/server/server ./cmd/server
	go build -o ./cmd/agent/agent ./cmd/agent




	