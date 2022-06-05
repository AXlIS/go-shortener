package server

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/mocks"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/golang/mock/gomock"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockURLWorker(ctrl)
	s.EXPECT().AddValue("KRJARhJf5S", "https://www.yandex.ru/", gomock.Any()).Return(nil)

	tests := []struct {
		name    string
		request string
		body    string
		want    Want
	}{
		{
			name:    "POST 201 OK test",
			request: "/api/shorten",
			body:    "{\"url\":\"https://www.yandex.ru/\"}",
			want: Want{
				code:     201,
				response: "{\"result\":\"http://localhost:8080/KRJARhJf5S\"}",
			},
		},
	}

	conf := config.NewConfig("http://localhost:8080")
	services := service.NewService(s, conf)
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockURLWorker(ctrl)
	s.EXPECT().GetValue("VzGUU3fuyV").Return("https://www.yandex.ru/", nil)

	tests := []struct {
		name    string
		request string
		want    Want
	}{
		{
			name:    "GET 200 url test",
			request: "/VzGUU3fuyV",
			want: Want{
				code: 307,
			},
		},
	}

	for _, tt := range tests {
		conf := config.NewConfig("http://localhost:8080")
		services := service.NewService(s, conf)
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockURLWorker(ctrl)
	s.EXPECT().AddValue("KRJARhJf5S", "https://www.yandex.ru/", gomock.Any()).Return(nil)

	tests := []struct {
		name    string
		request string
		body    string
		want    Want
	}{
		{
			name:    "POST 201 OK test",
			request: "/",
			body:    "https://www.yandex.ru/",
			want: Want{
				code:     201,
				response: "http://localhost:8080/KRJARhJf5S",
			},
		},
	}

	conf := config.NewConfig("http://localhost:8080")
	services := service.NewService(s, conf)
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
