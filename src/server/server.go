package main

import (
	readDB "L0/database"
	"fmt"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

func main() {
	natsStreamConnection, err := stan.Connect("test-cluster", "server", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	var foundedOrder *readDB.Orders
	semaphore := make(chan struct{}, 1)
	_, err = natsStreamConnection.Subscribe("data", func(message *stan.Msg) {
		log.Printf("Received a message: %s\n", string(message.Data))
		foundedOrder, _ = readDB.FileDeserialize(message.Data)
		if foundedOrder != nil {
			semaphore <- struct{}{}
		}
	})

	page, err := template.ParseFiles("/Users/monke/go/WildBerriesTech-L0/templates/index.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = page.Execute(w, nil)
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		//page, err := template.ParseFiles("../../templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		orderID := r.FormValue("id")
		if orderID != "" {
			err = natsStreamConnection.Publish("id", []byte(orderID))

		}
		<-semaphore
		err = page.Execute(w, foundedOrder)
		fmt.Println(orderID)

	})
	http.ListenAndServe(":8080", nil)
}
