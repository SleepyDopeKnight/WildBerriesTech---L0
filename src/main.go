package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

func main() {
	natsStreamConnection, err := stan.Connect("test-cluster", "publisher", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Println(err)
	}
	_, err = natsStreamConnection.Subscribe("aboba", func(message *stan.Msg) { log.Printf("Received a message: %s\n", string(message.Data)) })
	if err != nil {
		log.Println(err)
	}
	err = natsStreamConnection.Publish("aboba", []byte("jotaro kujo"))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
	err = natsStreamConnection.Publish("aboba", []byte("skibidi"))
	if err != nil {
		log.Println(err)
	}
	//err = natsStreamConnection.Close()
	//if err != nil {
	//	log.Println(err)
	//}
	select {}
}
