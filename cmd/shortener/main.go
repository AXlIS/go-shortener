package main

import (
	"fmt"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}
}

func main() {
	var (
		storage store.URLWorker
		err     error
	)

	fmt.Println(2, config.GetEnv("FILE_STORAGE_PATH", ""))

	if filePath := config.GetEnv("FILE_STORAGE_PATH", ""); filePath != "" {
		storage, err = store.NewFileStorage(filePath)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	} else {
		storage = store.NewStorage()
	}

	services := service.NewService(storage)
	handlers := handler.NewHandler(services)

	s := new(server.Server)

	if err := s.Start(config.GetEnv("SERVER_ADDRESS", "8080"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
