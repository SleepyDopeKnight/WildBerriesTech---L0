package main

import (
	readDB "L0/database"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	fileData, _ := readDB.FileOpen("../../models/model.json")

	natsStreamConnection, err := stan.Connect("test-cluster", "publisher", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	err = natsStreamConnection.Publish("orders", []byte(fileData))
	if err != nil {
		log.Fatal(err)
	}
}
