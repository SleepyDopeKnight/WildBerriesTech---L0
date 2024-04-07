package main

import (
	readDB "L0/database"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	fileData, _ := readDB.FileOpen("/Users/chamomiv/go/WildBerriesTech-L0/src/model.json")

	natsStreamConnection, err := stan.Connect("test-cluster", "publisher", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Println(err)
	}
	_, err = natsStreamConnection.Subscribe("jojo", func(message *stan.Msg) { log.Printf("Received a message: %s\n", string(message.Data)) })
	if err != nil {
		log.Println(err)
	}
	err = natsStreamConnection.Publish("jojo", []byte(fileData))
	if err != nil {
		log.Println(err)
	}
	//orders, _ := readDB.FileDeserialize(fileData)
	//err = natsStreamConnection.Close()
	//if err != nil {
	//	log.Println(err)
	//}
	select {}
}
