package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/ktigay/short-url/internal/config"
)

const (
	defaultServerPort      = "8080"
	defaultServerHost      = ":" + defaultServerPort
	defaultServerURL       = "http://localhost:" + defaultServerPort
	defaultFileStoragePath = "./cache/storage.txt"
	defaultRestore         = true
)

func parseFlags(args []string) (*config.Config, error) {
	cfg := &config.Config{}

	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)
	flags.StringVar(&cfg.ServerHost, "a", defaultServerHost, "address and port to run server")
	flags.StringVar(&cfg.ServerURL, "b", defaultServerURL, "base server URL for short link")
	flags.StringVar(&cfg.FileStoragePath, "f", defaultFileStoragePath, "file storage path")
	flags.BoolVar(&cfg.Restore, "r", defaultRestore, "restore data from storage")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.ServerHost == "" {
		return nil, fmt.Errorf("host flag is required")
	}

	return cfg, nil
}
