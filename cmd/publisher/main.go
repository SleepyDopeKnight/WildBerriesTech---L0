package main

import (
	nats "L0/pkg/broker_connect"
	"L0/pkg/reader_json"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	filesData := reader_json.Open("./schema/")
	nc := nats.Connect("test-cluster", "publisher")

	for _, fileData := range filesData {
		err := nc.Publish("orders", []byte(fileData))
		if err != nil {
			log.Println(err)
		}
	}
}
