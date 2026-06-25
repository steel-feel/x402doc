PORT ?= 9000

run:
	go run ./cmd/server/. --port=$(PORT)
