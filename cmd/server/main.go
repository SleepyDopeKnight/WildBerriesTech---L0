package main

import (
	"os"

	"L0/internal/api/handlers"
	"L0/internal/app/server"
	nats "L0/internal/broker_message"
	"L0/internal/config"
)

func main() {
	config.Load(".env")

	brokerMessage := nats.New()
	handler := handlers.New(brokerMessage)
	serv := server.New(handler)
	serv.Run(os.Getenv(":SERVER_PORT"))
}
