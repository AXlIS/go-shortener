package config

import "os"

// Config ...
type Config struct {
	BaseURL string
}

func NewConfig(URL string) *Config {
	return &Config{
		BaseURL: URL,
	}
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
