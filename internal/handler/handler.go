package handler

import (
	"encoding/json"
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
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
		api.POST("/shorten", h.CreateJSONShorten)

		shorten := api.Group("/shorten")
		{
			shorten.POST("/batch", h.CreateJSONShortenBatch)
		}

		user := api.Group("user")
		{
			user.GET("urls", h.GetAllShortens)
			user.DELETE("urls", h.DeleteShortens)
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

	userID := GetUserID(c)

	if err := json.Unmarshal(body, &input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortURL, err := h.service.AddURL(input.URL, userID)

	c.Header("content-type", "application/json")
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			c.JSON(http.StatusConflict, map[string]string{
				"result": fmt.Sprintf("%s/%s", h.config.BaseURL, shortURL),
			})
			return
		}
	}

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	userID := GetUserID(c)

	urls, err := h.service.AddBatchURL(input, userID)
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

	if url == "" {
		log.Println(key)
		c.Status(http.StatusGone)
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

	userID := GetUserID(c)

	shortURL, err := h.service.AddURL(string(body), userID)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			c.String(http.StatusConflict, fmt.Sprintf("%s/%s", h.config.BaseURL, shortURL))
			return
		}
	}

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusCreated, fmt.Sprintf("%s/%s", h.config.BaseURL, shortURL))
}

func (h *Handler) GetAllShortens(c *gin.Context) {
	userID := GetUserID(c)

	items, err := h.service.GetAllURLS(userID)
	if err != nil {
		errorResponse(c, http.StatusNoContent, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, items)
}

func (h *Handler) DeleteShortens(c *gin.Context) {
	userID := GetUserID(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input []string
	if err := json.Unmarshal(body, &input); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("start delete process")

	h.service.DeleteURLS(input, userID)

	log.Println("finish delete process")

	c.Header("content-type", "application/json")
	c.Status(http.StatusAccepted)
}

func (h *Handler) GetPing(c *gin.Context) {
	ping, err := h.service.Ping()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"active": ping})
}
