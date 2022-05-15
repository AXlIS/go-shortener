package server

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Want struct {
	code     int
	response string
}

func TestServer_CreateJSONShorten(t *testing.T) {
	tests := []struct {
		name    string
		request string
		body    string
		storage *storage.Storage
		want    Want
	}{
		{
			name:    "POST 201 OK test",
			request: "/api/shorten",
			body:    "{\"url\":\"https://www.yandex.ru/\"}",
			storage: storage.NewStorage(),
			want: Want{
				code:     201,
				response: "{\"result\":\"http://localhost:8080/KRJARhJf5S\"}",
			},
		},
	}

	store := storage.NewStorage()
	services := service.NewService(store)
	conf := config.NewConfig("http://localhost:8080")
	handlers := handler.NewHandler(services, conf)

	router := handlers.InitRoutes()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, tt.request, reader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, string(body))
			assert.NoError(t, err)
		})
	}
}

func TestServer_GetShorten(t *testing.T) {
	tests := []struct {
		name    string
		request string
		storage *storage.Storage
		want    Want
	}{
		{
			name:    "GET 200 url test",
			request: "/VzGUU3fuyV",
			storage: &storage.Storage{
				List: map[string]string{
					"VzGUU3fuyV": "https://www.yandex.ru/",
				},
			},
			want: Want{
				code: 307,
			},
		},
	}

	for _, tt := range tests {

		services := service.NewService(tt.storage)
		conf := config.NewConfig("http://localhost:8080")
		handlers := handler.NewHandler(services, conf)

		router := handlers.InitRoutes()

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func TestServer_CreateShorten(t *testing.T) {
	tests := []struct {
		name    string
		request string
		body    string
		storage *storage.Storage
		want    Want
	}{
		{
			name:    "POST 201 OK test",
			request: "/",
			storage: &storage.Storage{},
			body:    "https://www.yandex.ru/",
			want: Want{
				code:     201,
				response: "http://localhost:8080/KRJARhJf5S",
			},
		},
	}

	store := storage.NewStorage()
	services := service.NewService(store)
	conf := config.NewConfig("http://localhost:8080")
	handlers := handler.NewHandler(services, conf)

	router := handlers.InitRoutes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, tt.request, reader)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, string(body))
			assert.NoError(t, err)
		})
	}
}
