package main

import (
	"L0/pkg/database"
	"L0/pkg/read_db"
	"database/sql"
	"encoding/json"
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
	ChannelForGetJSON(natsStreamConnection, db)
	ChannelsForHandleIdDRequest(natsStreamConnection, db)

	select {}
}

func ChannelForGetJSON(natsStreamConnection stan.Conn, db *sql.DB) {
	_, err := natsStreamConnection.Subscribe("orders", func(message *stan.Msg) {
		orders := readDB.readDB.FileDeserialize(message.Data)
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
