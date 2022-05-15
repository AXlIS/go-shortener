package main

import (
	"flag"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/joho/godotenv"
	"log"
)

var (
	fileStoragePath, serverAddress, baseURL string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	flag.StringVar(&fileStoragePath, "f", "./storage.json", "path to file")
	flag.StringVar(&serverAddress, "a", ":8080", "port")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base url")
	flag.Parse()

	if path := config.GetEnv("BASE_URL", ""); path != "" {
		baseURL = path
	}
}

func main() {
	var (
		storage store.URLWorker
		err     error
	)

	if filePath := config.GetEnv("FILE_STORAGE_PATH", fileStoragePath); filePath != "" {
		storage, err = store.NewFileStorage(filePath)

		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	} else {
		storage = store.NewStorage()
	}

	services := service.NewService(storage)
	conf := config.NewConfig(baseURL)
	handlers := handler.NewHandler(services, conf)

	s := new(server.Server)

	if err := s.Start(config.GetEnv("SERVER_ADDRESS", serverAddress), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
