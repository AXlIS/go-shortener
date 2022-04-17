package server

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itchyny/base58-go"
	"io"
	"log"
	"math/big"
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
	config *Config
	Data   map[string]string
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		Data:   make(map[string]string),
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
			shortUrl := GenerateShortUrl(url)
			s.AddValue(shortUrl, url)

			r := Response{
				ShortUrl: "http://" + r.Host + "/" + shortUrl,
				Message:  "Short Url was created",
			}

			response, err := json.Marshal(r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(response)

		case "GET":
			key := strings.Split(r.URL.RequestURI(), "/")[1]
			url, err := s.GetValue(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)

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
	s.Data[key] = value
}

func (s *APIServer) GetValue(key string) (string, error) {
	if value, found := s.Data[key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}

func GenerateShortUrl(url string) string {
	urlHashBytes := sha256Of(url)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:10]
}

func sha256Of(initialString string) []byte {
	encoder := sha256.New()
	encoder.Write([]byte(initialString))
	return encoder.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return string(encoded)
}
