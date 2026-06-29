PORT ?= 9000
GRPC_PORT ?= 9001

.PHONY: run build test smoke-test proto-gen lint compile

compile:
	cd spec && npx tsp compile .
	mkdir -p proto/document/v1
	mv proto/document/v1.proto proto/document/v1/v1.proto 2>/dev/null || true
	PATH=$$PATH:/Users/steel-feel/go/bin protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mproto/document/v1/v1.proto=github.com/steel-feel/prac/proto/document/v1 --go-grpc_out=. --go-grpc_opt=paths=source_relative --go-grpc_opt=Mproto/document/v1/v1.proto=github.com/steel-feel/prac/proto/document/v1 proto/document/v1/v1.proto
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

