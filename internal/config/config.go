package config

// Config - конфигурация сервера.
type Config struct {
	ServerHost      string `env:"SERVER_ADDRESS"`
	ServerURL       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}
