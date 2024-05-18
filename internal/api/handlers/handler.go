package handlers

import (
	"net/http"

	nats "L0/internal/broker_message"
	"L0/internal/database/models"
	"L0/pkg/html"
)

type Handler struct {
	cache map[string]*models.Orders
	bm    nats.BrokerMessage
}

func New(bm nats.BrokerMessage) Handler {
	return Handler{cache: map[string]*models.Orders{}, bm: bm}
}

func (h *Handler) StartPage(w http.ResponseWriter, _ *http.Request) {
	html.ParseTemplate(w, "./assets/main/index.html", nil)
}

func (h *Handler) ShowOrder(w http.ResponseWriter, r *http.Request) {
	h.StartPage(w, r)
	orderID := r.FormValue("id")
	if h.cache[orderID] == nil {
		if foundedOrder := h.bm.GetOrder(orderID, w); foundedOrder.OrderUid != "" {
			h.cache[orderID] = foundedOrder
		}
	}
	if h.cache[orderID] != nil {
		html.ParseTemplate(w, "./assets/main/order_data.html", h.cache[orderID])
	} else {
		html.ParseTemplate(w, "./assets/errors/404.html", nil)
	}
}
