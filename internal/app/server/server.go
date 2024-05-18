package server

import (
	"log"
	"net/http"

	"L0/internal/api/handlers"
)

type Server struct {
	h handlers.Handler
}

func New(h handlers.Handler) Server {
	return Server{h: h}
}

func (s *Server) Run(port string) {
	http.HandleFunc("/", s.h.StartPage)
	http.HandleFunc("/data", s.h.ShowOrder)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println(err)
	}
}
