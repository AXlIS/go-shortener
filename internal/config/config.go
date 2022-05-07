package config

import "os"

// Config ...
type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
