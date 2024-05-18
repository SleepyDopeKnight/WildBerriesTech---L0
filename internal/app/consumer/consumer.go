package consumer

import (
	"L0/internal/database"
	"L0/internal/serialization"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

func ChannelForGetJSON(nc stan.Conn, db *sql.DB) {
	_, err := nc.Subscribe("orders", func(message *stan.Msg) {
		orders := serialization.FileDeserialize(message.Data)
		if orders != nil {
			database.FillDatabase(orders, db)
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func ChannelsForHandleIdDRequest(nc stan.Conn, db *sql.DB) {
	_, err := nc.Subscribe("id", func(message *stan.Msg) {
		err := db.Ping()
		if err == nil {
			publicationOrder(message, nc, db)
		} else {
			publicationOrder(message, nc, nil)
		}
	})
	if err != nil {
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
	err := nc.Publish("data", order)
	if err != nil {
		log.Println(err)
	}
}
