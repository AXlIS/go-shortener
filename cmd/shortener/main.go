package main

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"log"
)

func main() {
	storage := store.NewStorage()
	services := service.NewService(storage)
	handlers := handler.NewHandler(services)

	s := new(server.Server)

	if err := s.Start(config.GetEnv("SERVER_ADDRESS", "8080"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
