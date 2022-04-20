package config

// Config ...
type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}
