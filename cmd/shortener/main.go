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
	"os"
)

var (
	fileStoragePath, serverAddress, baseURL string
	envPath = "./.env"
)

func init() {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		_, _ = os.Create(envPath)
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	flag.StringVar(&fileStoragePath, "f", config.GetEnv("FILE_STORAGE_PATH", "./storage.json"), "path to file")
	flag.StringVar(&serverAddress, "a", config.GetEnv("SERVER_ADDRESS", ":8080"), "port")
	flag.StringVar(&baseURL, "b", config.GetEnv("BASE_URL", "http://localhost:8080"), "base url")
	flag.Parse()

	env := map[string]string{
		"BASE_URL":          baseURL,
		"SERVER_ADDRESS":    serverAddress,
		"FILE_STORAGE_PATH": fileStoragePath,
	}

	if err := godotenv.Write(env, envPath); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

func main() {
	var (
		storage store.URLWorker
		err     error
	)

	if fileStoragePath != "" {
		storage, err = store.NewFileStorage(fileStoragePath)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	} else {
		storage = store.NewStorage()
	}

	services := service.NewService(storage)
	handlers := handler.NewHandler(services)

	s := new(server.Server)

	if err := s.Start(serverAddress, handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
