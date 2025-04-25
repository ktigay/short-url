package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/ktigay/short-url/internal/config"
)

const (
	defaultServerPort = "8080"
	defaultServerHost = ":" + defaultServerPort
	defaultAServerURL = "http://localhost:" + defaultServerPort
)

func parseFlags(args []string) (*config.Config, error) {
	cfg := &config.Config{
		ServerURL: defaultAServerURL,
	}
	flags := flag.NewFlagSet("server flags", flag.ContinueOnError)
	flags.StringVar(&cfg.ServerHost, "a", defaultServerHost, "address and port to run server")

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
