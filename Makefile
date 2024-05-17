.PHONY: container consumer publisher server

container:
	docker compose up -d

consumer:
	go run cmd/consumer/main.go

publisher:
	go run cmd/publisher/main.go

server:
	go run cmd/server/main.go
