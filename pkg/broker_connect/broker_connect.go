package nats

import (
	"github.com/nats-io/stan.go"
	"log"
)

func Connect(clusterID, clientID string) stan.Conn {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	return nc
}
