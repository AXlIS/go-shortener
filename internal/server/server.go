package server

import (
	"github.com/AXlIS/go-shortener/internal/config"
	s "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

	router.POST("/", s.PostUrlHandler)
	router.GET("/:id", s.GetUrlHandler)

	return router
}

// Start ...
func (s *APIServer) Start() error {
	r := s.SetupRouter()

	return r.Run(s.config.Port)
}

func (s *APIServer) PostUrlHandler(c *gin.Context) {

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	url := string(body)
	shortUrl := utils.GenerateShortUrl(url)
	s.storage.AddValue(shortUrl, url)

	resBody := Response{
		ShortUrl: "http://" + c.Request.Host + "/" + shortUrl,
		Message:  "Short Url was created",
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, resBody)
	return
}

func (s *APIServer) GetUrlHandler(c *gin.Context) {
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
	return
}
