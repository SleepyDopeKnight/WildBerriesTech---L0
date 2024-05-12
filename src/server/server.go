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
	_, err = natsStreamConnection.Subscribe("data", func(message *stan.Msg) {
		log.Printf("Received a message: %s\n", string(message.Data))
		foundedOrder, _ = readDB.FileDeserialize(message.Data)
	})

	page, err := template.ParseFiles("/Users/chamomiv/go/WildBerriesTech-L0/templates/index.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = page.Execute(w, nil)
	})

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		//page, err := template.ParseFiles("../../templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		orderID := r.FormValue("id")
		if orderID != "" {
			err = natsStreamConnection.Publish("id", []byte(orderID))

		}
		if foundedOrder != nil {
			err = page.Execute(w, foundedOrder)
			//} else {
			//	err = page.Execute(w, nil)
		}

		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//}
		//
		//if err != nil {
		//	log.Fatal(err)
		//}
		fmt.Println(orderID)

	})
	http.ListenAndServe(":8080", nil)
}
