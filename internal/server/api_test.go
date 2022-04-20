package server

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var conf = config.NewConfig()

type Want struct {
	code     int
	response string
}

func TestAPIServer_APIHandlerPost(t *testing.T) {
	tests := []struct {
		name    string
		request string
		body    string
		storage storage.Storage
		want    Want
	}{
		{
			name:    "POST 201 OK test",
			request: "/",
			body:    "https://www.yandex.ru/",
			storage: map[string]string{},
			want: Want{
				code:     201,
				response: "http://localhost:8080/KRJARhJf5S",
			},
		},
	}

	for _, tt := range tests {

		server := New(conf, tt.storage)
		router := server.SetupRouter()

		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, tt.request, reader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, resp.StatusCode, tt.want.code)
			assert.Equal(t, strings.Replace(tt.want.response, "localhost:8080", req.Host, 1), string(body))
			assert.NoError(t, err)
		})
	}
}

func TestAPIServer_APIHandlerGet(t *testing.T) {
	tests := []struct {
		name    string
		request string
		storage storage.Storage
		want    Want
	}{
		{
			name:    "GET 200 url test",
			request: "/VzGUU3fuyV",
			storage: map[string]string{
				"VzGUU3fuyV": "https://www.yandex.ru/",
			},
			want: Want{
				code: 307,
			},
		},
	}

	for _, tt := range tests {

		server := New(conf, tt.storage)
		router := server.SetupRouter()

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, tt.want.code)
		})
	}
}
