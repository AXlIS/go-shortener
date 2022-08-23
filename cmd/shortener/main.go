package main

import (
	"flag"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	_ "net/http/pprof"

	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
)

// @title Go Shortener App API
// @version 1.0
// description Service for shorting URLS

// host localhost:8080
// BasePath /

var (
	fileStoragePath, serverAddress, baseURL, databaseDsn string
	buildVersion                                         string = "N/A"
	buildDate                                            string = "N/A"
	buildCommit                                          string = "N/A"
)

const (
	addr = ":8000" // адрес сервера
)

func init() {

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	flag.StringVar(&fileStoragePath, "f", "storage.json", "path to file")
	flag.StringVar(&serverAddress, "a", ":8080", "port")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&databaseDsn, "d", "", "database address")
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

	conf := config.NewConfig(baseURL)

	if databasePath := config.GetEnv("DATABASE_DSN", databaseDsn); databasePath != "" {
		db, err := store.NewPostgresDB(databasePath)
		if err != nil {
			log.Fatalf("faild to initialize db: %s", err.Error())
		}

		storage = store.NewDatabaseStorage(db, conf)

	} else if filePath := config.GetEnv("FILE_STORAGE_PATH", fileStoragePath); filePath != "" {
		storage, err = store.NewFileStorage(filePath, conf)

		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	} else {
		storage = store.NewStorage(conf)
	}

	services := service.NewService(storage, conf)
	handlers := handler.NewHandler(services, conf)

	s := new(server.Server)

	if err := s.Start(config.GetEnv("SERVER_ADDRESS", serverAddress), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
