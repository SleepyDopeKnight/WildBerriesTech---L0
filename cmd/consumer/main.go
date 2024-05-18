package main

import (
	"L0/internal/app/consumer"
	"L0/internal/database"
	nats "L0/pkg/broker_connect"
	_ "github.com/lib/pq"
)

func main() {
	db := database.DBConnection()
	nc := nats.Connect("test-cluster", "consumer")
	consumer.ChannelForGetJSON(nc, db)
	consumer.ChannelsForHandleIdDRequest(nc, db)

	select {}
}
