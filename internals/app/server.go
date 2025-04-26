package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/istiak-004/image-extractor/internals/service"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}
	extractorService := service.NewExtractorService()
	h := NewHandler(extractorService)
	s.registerRoutes(h)

	return s

}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) registerRoutes(h *Handler) {
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/extract", h.ExtractDataFromImage).Methods("POST")
}
