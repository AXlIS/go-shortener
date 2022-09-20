package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

// Config ...
type Config struct {
	BaseURL       string
	TrustedSubnet net.IPNet
}

type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

func NewJSONConfig() *JSONConfig {
	var config JSONConfig

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(filepath.Join(pwd, "config.json"), os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func NewConfig(URL, trustedSubnet string) *Config {
	_, network, err := net.ParseCIDR(trustedSubnet)

	if err != nil {
		log.Printf("WARN can`t parse TRUSTED_SUBNET")
	}

	return &Config{
		BaseURL:       URL,
		TrustedSubnet: *network,
	}
}

type Env interface {
	string | bool
}

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetBoolEnv(key string, defaultValue bool) bool {
	_, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return exists

}
