package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/istiak-004/image-extractor/internals/service"
	"github.com/rs/cors"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}
	extractorService := service.NewExtractorService()
	creatorService := service.NewPNGCreatorService()

	h := NewHandler(extractorService, creatorService)
	s.registerRoutes(h)

	return s

}

func (s *Server) Start(addr string) error {
	// add headers
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://json-extraction-challenge.intellixio.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(s.router)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) registerRoutes(h *Handler) {
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/extract", h.ExtractDataFromImage).Methods("POST")
	api.HandleFunc("/create", h.PNGCreatorHandler).Methods("POST")
}
