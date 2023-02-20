package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
func (s *Server) Start(port string, handler http.Handler, tls bool) error {
	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	}

	if tls {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return s.httpServer.ListenAndServeTLS(filepath.Join(pwd, "localhost.crt"), filepath.Join(pwd, "localhost.key"))
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
