package main

import (
	"os"

	"L0/internal/app/consumer"
	"L0/internal/config"
	"L0/internal/database"
	nats "L0/pkg/broker_connect"

	_ "github.com/lib/pq"
)

func main() {
	config.Load(".env")
	db := database.DBConnection(os.Getenv("DB_DSN"))
	nc := nats.Connect("test-cluster", "consumer")
	consumer.ChannelForGetJSON(nc, db)
	consumer.ChannelsForHandleIdDRequest(nc, db)

	select {}
}
