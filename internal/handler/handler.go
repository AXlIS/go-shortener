package handler

import (
	"encoding/json"
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Handler struct {
	service *service.Service
	config  *config.Config
}

type ShortenInput struct {
	URL string `json:"url"`
}

func NewHandler(service *service.Service, conf *config.Config) *Handler {
	return &Handler{service: service, config: conf}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(CookieHandler())
	router.Use(DecompressBody())
	router.Use(gzip.Gzip(gzip.BestCompression))

	router.POST("/", h.CreateShorten)
	router.GET("/:id", h.GetShorten)
	router.GET("/ping", h.GetPing)

	api := router.Group("/api")
	{
		shorten := api.Group("/shorten")
		{
			shorten.POST("/", h.CreateJSONShorten)
			shorten.POST("/batch", h.CreateJSONShortenBatch)
		}

		user := api.Group("user")
		{
			user.GET("urls", h.GetAllShortens)
		}
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

	userId := GetUserId(c)

	if err := json.Unmarshal(body, &input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortURL, err := h.service.AddURL(input.URL, userId)
	if err != nil {
		fmt.Println(3)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, map[string]string{
		"result": fmt.Sprintf("%s/%s", h.config.BaseURL, shortURL),
	})
}

func (h *Handler) CreateJSONShortenBatch(c *gin.Context) {
	var input []*u.ShortenBatchInput

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(input) == 0 {
		errorResponse(c, http.StatusBadRequest, "the list is empty")
		return
	}

	userId := GetUserId(c)

	urls, err := h.service.AddBatchURL(input, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, urls)
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

	userId := GetUserId(c)

	shortURL, err := h.service.AddURL(string(body), userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.String(http.StatusCreated, fmt.Sprintf("%s/%s", h.config.BaseURL, shortURL))
}

func (h *Handler) GetAllShortens(c *gin.Context) {
	userId := GetUserId(c)

	items, err := h.service.GetAllURLS(userId)
	if err != nil {
		errorResponse(c, http.StatusNoContent, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetPing(c *gin.Context) {
	ping, err := h.service.Ping()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"active": ping})
}
