package main

import (
	"context"
	"flag"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

// @title Go Shortener App API
// @version 1.0
// description Service for shorting URLS

// host localhost:8080
// BasePath /

var (
	fileStoragePath, serverAddress, baseURL, databaseDsn, trustedSubnet string
	buildVersion                                                        string = "N/A"
	buildDate                                                           string = "N/A"
	buildCommit                                                         string = "N/A"
	tls                                                                 bool
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

	JSONConfig := config.NewJSONConfig()

	flag.StringVar(&fileStoragePath, "f", JSONConfig.FileStoragePath, "path to file")
	flag.StringVar(&serverAddress, "a", ":8080", "port")
	flag.StringVar(&baseURL, "b", JSONConfig.FileStoragePath, "base url")
	flag.StringVar(&databaseDsn, "d", JSONConfig.DatabaseDSN, "database address")
	flag.StringVar(&trustedSubnet, "t", JSONConfig.TrustedSubnet, "trusted subnet")
	flag.BoolVar(&tls, "s", JSONConfig.EnableHTTPS, "enable https")
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

	conf := config.NewConfig(baseURL, trustedSubnet)

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

	go func() {
		if err := s.Start(config.GetEnv("SERVER_ADDRESS", serverAddress), handlers.InitRoutes(), config.GetBoolEnv("ENABLE_HTTPS", tls)); err != nil {
			log.Fatalf("Error occured while running https server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := s.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

}
