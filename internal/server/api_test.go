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
				response: `{"url":"http://localhost:8080/KRJARhJf5S","message":"Short Url was created"}`,
			},
		},
	}

	conf := &config.Config{
		Port: ":8080",
	}

	for _, tt := range tests {

		server := New(conf, tt.storage)

		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)

			request := httptest.NewRequest(http.MethodPost, tt.request, reader)
			w := httptest.NewRecorder()
			h := server.APIHandlerUrl()

			h.ServeHTTP(w, request)
			resp := w.Result()

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, resp.StatusCode, tt.want.code)
			assert.Equal(t, string(body), strings.Replace(tt.want.response, "localhost:8080", request.Host, 1))
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
			request: "/KRJARhJf5S",
			storage: map[string]string{
				"KRJARhJf5S": "https://www.yandex.ru/",
			},
			want: Want{
				code: 307,
			},
		},
	}

	conf := &config.Config{
		Port: ":8080",
	}

	for _, tt := range tests {

		server := New(conf, tt.storage)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := server.APIHandlerUrl()

			h.ServeHTTP(w, request)
			resp := w.Result()

			assert.Equal(t, resp.StatusCode, tt.want.code)
		})
	}
}
