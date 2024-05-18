package nats

import (
	"log"

	"github.com/nats-io/stan.go"
)

func Connect(clusterID, clientID string) stan.Conn {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Println(err)
	}

	return nc
}
