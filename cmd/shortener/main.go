package main

import (
	"github.com/AXlIS/go-shortener/internal/app/server"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	conf := server.NewConfig()
	if err := viper.Unmarshal(conf); err != nil {
		log.Fatal(err)
	}

	s := server.New(conf)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
