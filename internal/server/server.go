package server

import (
	"fmt"
	"github.com/AXlIS/go-shortener/internal/config"
	s "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Response struct {
	ShortURL string `json:"url"`
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

// SetupRouter ...
func (s *APIServer) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/", s.PostURLHandler)
	router.GET("/:id", s.GetURLHandler)

	return router
}

// Start ...
func (s *APIServer) Start() error {
	r := s.SetupRouter()

	return r.Run(s.config.Port)
}

func (s *APIServer) PostURLHandler(c *gin.Context) {

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	url := string(body)
	shortURL := utils.GenerateShortURL(url)
	s.storage.AddValue(shortURL, url)

	shortURL = fmt.Sprintf("http://%s/%s", c.Request.Host, shortURL)

	c.Header("content-type", "application/json")
	c.String(http.StatusCreated, shortURL)
}

func (s *APIServer) GetURLHandler(c *gin.Context) {
	key := c.Params.ByName("id")
	if key == "" {
		c.String(http.StatusBadRequest, "The query parameter is missing")
		return
	}

	url, err := s.storage.GetValue(key)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Header("Location", url)
	c.Status(http.StatusTemporaryRedirect)
}
