package consumer

import (
	"L0/internal/database"
	"L0/internal/serialization"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

func ChannelForGetJSON(natsStreamConnection stan.Conn, db *sql.DB) {
	_, err := natsStreamConnection.Subscribe("orders", func(message *stan.Msg) {
		orders := serialization.FileDeserialize(message.Data)
		if orders != nil {
			database.FillDatabase(orders, db)
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func ChannelsForHandleIdDRequest(natsStreamConnection stan.Conn, db *sql.DB) {
	_, err := natsStreamConnection.Subscribe("id", func(message *stan.Msg) {
		err := db.Ping()
		if err == nil {
			publicationOrder(message, natsStreamConnection, db)
		} else {
			publicationOrder(message, natsStreamConnection, nil)
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func publicationOrder(message *stan.Msg, natsStreamConnection stan.Conn, db *sql.DB) {
	if db != nil {
		wantedOrder := database.FindOrder(message, db)
		outgoingOrder, err := json.Marshal(wantedOrder)
		if err != nil {
			log.Println(err)
		}
		publish(natsStreamConnection, outgoingOrder)
	} else {
		publish(natsStreamConnection, nil)
	}
}

func publish(natsStreamConnection stan.Conn, order []byte) {
	err := natsStreamConnection.Publish("data", order)
	if err != nil {
		log.Println(err)
	}
}
