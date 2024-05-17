package main

import (
	"L0/internal/app/consumer"
	"L0/internal/database"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	db := database.DBConnection()
	natsStreamConnection, err := stan.Connect("test-cluster", "consumer", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	consumer.ChannelForGetJSON(natsStreamConnection, db)
	consumer.ChannelsForHandleIdDRequest(natsStreamConnection, db)

	select {}
}
