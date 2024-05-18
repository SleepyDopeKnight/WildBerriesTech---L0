package nats

import (
	"L0/internal/database/models"
	"L0/internal/serialization"
	nats "L0/pkg/broker_connect"
	"L0/pkg/html"

	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

type BrokerMessage struct {
	nc        stan.Conn
	semaphore chan *models.Orders
}

func New() BrokerMessage {
	bm := BrokerMessage{
		nc:        nats.Connect("test-cluster", "server"),
		semaphore: make(chan *models.Orders, 1),
	}
	bm.subscribe()
	return bm
}

func (b *BrokerMessage) subscribe() {
	_, err := b.nc.Subscribe("data", func(message *stan.Msg) {
		foundedOrder := serialization.FileDeserialize(message.Data)
		if foundedOrder != nil {
			b.semaphore <- foundedOrder
		} else {
			b.semaphore <- &models.Orders{}
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func (b *BrokerMessage) GetOrder(orderID string, w http.ResponseWriter) *models.Orders {
	err := b.nc.Publish("id", []byte(orderID))
	if err != nil {
		log.Println(err)
		html.ParseTemplate(w, "./assets/errors/500.html", nil)
	}
	return <-b.semaphore
}
