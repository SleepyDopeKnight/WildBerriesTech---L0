package consumer

import (
	"database/sql"
	"encoding/json"
	"log"

	"L0/internal/database"
	"L0/internal/serialization"

	"github.com/nats-io/stan.go"
)

func ChannelForGetJSON(nc stan.Conn, db *sql.DB) {
	if _, err := nc.Subscribe("orders", func(message *stan.Msg) {
		if orders := serialization.FileDeserialize(message.Data); orders != nil {
			database.FillDatabase(orders, db)
		}
	}); err != nil {
		log.Println(err)
	}
}

func ChannelsForHandleIdDRequest(nc stan.Conn, db *sql.DB) {
	if _, err := nc.Subscribe("id", func(message *stan.Msg) {
		if err := db.Ping(); err == nil {
			publicationOrder(message, nc, db)
		} else {
			publicationOrder(message, nc, nil)
		}
	}); err != nil {
		log.Println(err)
	}
}

func publicationOrder(message *stan.Msg, nc stan.Conn, db *sql.DB) {
	if db != nil {
		wantedOrder := database.FindOrder(message, db)

		outgoingOrder, err := json.Marshal(wantedOrder)
		if err != nil {
			log.Println(err)
		}

		publish(nc, outgoingOrder)
	} else {
		publish(nc, nil)
	}
}

func publish(nc stan.Conn, order []byte) {
	if err := nc.Publish("data", order); err != nil {
		log.Println(err)
	}
}
