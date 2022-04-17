package server

// Config ...
type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: ":8000",
	}
}
