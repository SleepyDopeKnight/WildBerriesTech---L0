package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	natsConnection, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err)
	}
	err = natsConnection.Publish("aboba", []byte("jotaro kujo"))
	if err != nil {
		log.Println(err)
	}
	_, err = natsConnection.Subscribe("aboba", func(message *nats.Msg) { fmt.Printf("Received a message: %s\n", string(message.Data)) })
	if err != nil {
		log.Println(err)
	}

	natsConnection.Close()
}
