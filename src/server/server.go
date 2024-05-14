package main

import (
	readDB "L0/database"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	NatsStreamConnection stan.Conn
	Page                 *template.Template
	Err                  error
	Semaphore            chan *readDB.Orders
	Cache                map[string]*readDB.Orders
}

func main() {
	natsStreamConnection, err := stan.Connect("test-cluster", "server", stan.NatsURL(stan.DefaultNatsURL))

	defer natsStreamConnection.Close()
	if err != nil {
		log.Fatal(err)
	}
	cache := make(map[string]*readDB.Orders)
	semaphore := make(chan *readDB.Orders, 1)
	_, err = natsStreamConnection.Subscribe("data", func(message *stan.Msg) {
		semaphore <- readDB.FileDeserialize(message.Data)
	})
	if err != nil {
		log.Fatal(err)
	}

	page, err := template.ParseFiles("/Users/chamomiv/go/WildBerriesTech-L0/templates/index.html")
	h := Handler{NatsStreamConnection: natsStreamConnection, Page: page, Err: err, Semaphore: semaphore, Cache: cache}
	http.HandleFunc("/", h.rootHandler)
	http.HandleFunc("/data", h.dataHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	if h.Err != nil {
		http.Error(w, h.Err.Error(), http.StatusInternalServerError)
	}
	err := h.Page.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) dataHandler(w http.ResponseWriter, r *http.Request) {
	if h.Err != nil {
		http.Error(w, h.Err.Error(), http.StatusInternalServerError)
	}
	orderID := r.FormValue("id")
	if h.Cache[orderID] == nil {
		err := h.NatsStreamConnection.Publish("id", []byte(orderID))
		if err != nil {
			log.Fatal(err)
		}
		if foundedOrder := <-h.Semaphore; foundedOrder.OrderUid != "" {
			h.Cache[orderID] = foundedOrder
		}
	}
	if h.Cache[orderID] != nil {
		err := h.Page.Execute(w, h.Cache[orderID])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_ = h.Page.Execute(w, nil)
	}
}
