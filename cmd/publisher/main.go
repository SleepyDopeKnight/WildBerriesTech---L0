package main

import (
	"log"

	nats "L0/pkg/broker_connect"
	"L0/pkg/reader_json"

	_ "github.com/lib/pq"
)

func main() {
	filesData := reader_json.Open("./schema/")
	nc := nats.Connect("test-cluster", "publisher")

	for _, fileData := range filesData {
		if err := nc.Publish("orders", []byte(fileData)); err != nil {
			log.Println(err)
		}
	}
}
