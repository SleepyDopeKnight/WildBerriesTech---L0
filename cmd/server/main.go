package main

import (
	"L0/internal/api/handlers"
	"L0/internal/app/server"
	nats "L0/internal/broker_message"
)

func main() {
	brokerMessage := nats.New()
	handler := handlers.New(brokerMessage)
	serv := server.New(handler)
	serv.Run(":80")
}
