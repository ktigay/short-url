package config

// Config - конфигурация сервера.
type Config struct {
	ServerHost string `env:"ADDRESS"`
	ServerURL  string
}
