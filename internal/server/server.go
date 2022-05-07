package server

import (
	"net/http"
)

type Response struct {
	ShortURL string `json:"url"`
	Message  string `json:"message"`
}

type NotFoundResponse struct {
	Message string `json:"message"`
}

// Server ...
type Server struct {
	httpServer *http.Server
}

// Start ...
func (s *Server) Start(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	}

	return s.httpServer.ListenAndServe()
}
