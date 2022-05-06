package main

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/handler"
	"github.com/AXlIS/go-shortener/internal/server"
	"github.com/AXlIS/go-shortener/internal/service"
	store "github.com/AXlIS/go-shortener/internal/storage"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatal(err)
	//}

	conf := config.NewConfig()
	if err := viper.Unmarshal(conf); err != nil {
		log.Fatal(err)
	}

	storage := store.NewStorage()
	services := service.NewService(storage)
	handlers := handler.NewHandler(services)

	s := new(server.Server)

	if err := s.Start("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
