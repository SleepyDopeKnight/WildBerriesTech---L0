package main

import (
	"L0/internal/serialization"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	filesData := serialization.OpenOrdersJSON("/Users/chamomiv/go/WildBerriesTech-L0/schema/")

	natsStreamConnection, err := stan.Connect("test-cluster", "publisher", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Println(err)
	}

	for _, fileData := range filesData {
		err = natsStreamConnection.Publish("orders", []byte(fileData))
		if err != nil {
			log.Println(err)
		}
	}
}
