package config

// Config - конфигурация сервера.
type Config struct {
	ServerHost string `env:"SERVER_ADDRESS"`
	ServerURL  string `env:"BASE_URL"`
}
