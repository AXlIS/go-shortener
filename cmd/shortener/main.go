package main

import (
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/server"
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

	s := server.New(conf, storage)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
