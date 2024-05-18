package server

import (
	"L0/internal/api/handlers"
	"log"
	"net/http"
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
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println(err)
	}
}
