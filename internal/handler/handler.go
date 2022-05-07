package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Handler struct {
	service *service.Service
}

type ShortenInput struct {
	URL string `json:"url"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/", h.CreateShorten)
	router.GET("/:id", h.GetShorten)

	api := router.Group("/api")
	{
		api.POST("/shorten", h.CreateJSONShorten)
	}

	return router
}

func (h *Handler) CreateJSONShorten(c *gin.Context) {
	var input ShortenInput

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortURL, err := h.service.AddURL(input.URL)
	if err != nil {
		fmt.Println(3)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, map[string]string{
		"result": fmt.Sprintf("%s/%s", config.GetEnv("BASE_URL", "http://localhost:8080"), shortURL),
	})
}

func (h *Handler) GetShorten(c *gin.Context) {
	key := c.Params.ByName("id")
	if key == "" {
		errorResponse(c, http.StatusBadRequest, "The query parameter is missing")
		return
	}

	url, err := h.service.GetURL(key)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.Header("Location", url)
	c.Status(http.StatusTemporaryRedirect)
}

func (h *Handler) CreateShorten(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortURL, err := h.service.AddURL(string(body))
	if err != nil {
		fmt.Println(3)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.String(http.StatusCreated, fmt.Sprintf("%s/%s", config.GetEnv("BASE_URL", "http://localhost:8080"), shortURL))
}
