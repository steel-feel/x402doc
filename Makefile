PORT ?= 9000
GRPC_PORT ?= 9001

.PHONY: run build test smoke-test proto-gen lint compile

compile:
	cd spec && npx tsp compile .
	go generate ./api/...

run:
	go run ./cmd/server/. --port=$(PORT) --grpc-port=$(GRPC_PORT)

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

smoke-test:
	go test ./smoke_test.go

lint:
	golangci-lint run

