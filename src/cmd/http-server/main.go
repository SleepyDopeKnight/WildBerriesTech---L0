package main

import (
	"L0/pkg/read_db"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	NatsStreamConnection stan.Conn
	Semaphore            chan *readDB.Orders
	Cache                map[string]*readDB.Orders
}

func main() {
	cache := make(map[string]*readDB.Orders)

	h := Handler{Cache: cache}
	h.connection()
	h.subscribe()
	http.HandleFunc("/", h.rootHandler)
	http.HandleFunc("/data", h.dataHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) connection() {
	h.Semaphore = make(chan *readDB.Orders, 1)
	var err error
	h.NatsStreamConnection, err = stan.Connect("test-cluster", "server", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) subscribe() {
	_, err := h.NatsStreamConnection.Subscribe("data", func(message *stan.Msg) {
		foundedOrder := readDB.FileDeserialize(message.Data)
		if foundedOrder != nil {
			h.Semaphore <- foundedOrder
		} else {
			h.Semaphore <- &readDB.Orders{}
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	h.hmtlParse(w, "/Users/chamomiv/go/WildBerriesTech-L0/templates/main/index.html", nil)
}

func (h *Handler) dataHandler(w http.ResponseWriter, r *http.Request) {
	h.rootHandler(w, r)
	orderID := r.FormValue("id")
	if h.Cache[orderID] == nil {
		h.publishId(orderID, w)
	}
	if h.Cache[orderID] != nil {
		h.hmtlParse(w, "/Users/chamomiv/go/WildBerriesTech-L0/templates/main/order_data.html", h.Cache[orderID])
	} else {
		h.hmtlParse(w, "/Users/chamomiv/go/WildBerriesTech-L0/templates/errors/404.html", nil)
	}
}

func (h *Handler) publishId(orderID string, w http.ResponseWriter) {
	err := h.NatsStreamConnection.Publish("id", []byte(orderID))
	if err != nil {
		log.Println(err)
		h.hmtlParse(w, "/Users/chamomiv/go/WildBerriesTech-L0/templates/errors/500.html", nil)
	}
	if foundedOrder := <-h.Semaphore; foundedOrder.OrderUid != "" {
		h.Cache[orderID] = foundedOrder
	}
}

func (h *Handler) hmtlParse(w http.ResponseWriter, htmlFile string, data any) {
	page, err := template.ParseFiles(htmlFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = page.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
