package main

import (
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	storage := store.NewStorage()
	services := service.NewService(storage)
	handlers := handler.NewHandler(services)

	s := new(server.Server)

	if err := s.Start(os.Getenv("SERVER_ADDRESS"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
