package handler

import (
	"fmt"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *service.Service
}

type ShortenInput struct {
	Url string `json:"url"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/shorten", h.CreateShorten)
		api.GET("/:id", h.GetShorten)
	}

	return router
}

func (h *Handler) CreateShorten(c *gin.Context) {
	var input ShortenInput
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	shortURL := h.service.AddURL(input.Url)

	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, map[string]string{
		"result": fmt.Sprintf("http://%s/%s", c.Request.Host, shortURL),
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