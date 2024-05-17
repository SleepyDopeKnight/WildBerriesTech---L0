package main

import (
	"L0/pkg/read_db"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	filesData := readDB.FileOpen("/Users/chamomiv/go/WildBerriesTech-L0/models/")

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
