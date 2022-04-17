package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AXlIS/go-shortener/internal/config"
	s "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	ShortUrl string `json:"url"`
	Message  string `json:"message"`
}

type NotFoundResponse struct {
	Message string `json:"message"`
}

// APIServer ...
type APIServer struct {
	config  *config.Config
	storage s.Storage
}

// New ...
func New(config *config.Config, store s.Storage) *APIServer {
	return &APIServer{
		config:  config,
		storage: store,
	}
}

// Start ...
func (s *APIServer) Start() error {
	server := &http.Server{
		Addr: s.config.Port,
	}

	http.HandleFunc("/", s.APIHandlerUrl())
	fmt.Println("Start")
	return server.ListenAndServe()
}

func (s *APIServer) APIHandlerUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			url := string(body)
			shortUrl := utils.GenerateShortUrl(url)
			s.AddValue(shortUrl, url)

			r := Response{
				ShortUrl: "http://" + r.Host + "/" + shortUrl,
				Message:  "Short Url was created",
			}

			resp, err := json.Marshal(r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(resp)
			return

		case "GET":
			key := strings.Split(r.URL.RequestURI(), "/")[1]
			url, err := s.GetValue(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return

		default:
			body := NotFoundResponse{
				Message: "Resource Not Found",
			}
			resp, err := json.Marshal(body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}

	}
}

func (s *APIServer) AddValue(key, value string) {
	s.storage[key] = value
}

func (s *APIServer) GetValue(key string) (string, error) {
	if value, found := s.storage[key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}
