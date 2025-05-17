package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"github.com/ktigay/short-url/internal/snapshot"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"

	"github.com/ktigay/short-url/internal/generator"
	"github.com/ktigay/short-url/internal/logger"
	"github.com/ktigay/short-url/internal/middleware"
	"github.com/ktigay/short-url/internal/shorturl"
	"github.com/ktigay/short-url/internal/storage"
)

func init() {
	path := "./cache"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	var l = logger.Initialize()

	router := mux.NewRouter()

	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		l.Fatal().Err(err).Msg("failed to parse flags")
	}

	st, err := initStorage(cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to initialize storage")
	}

	registerMiddlewares(l, router)
	registerHandlers(st, cfg, l, router)

	if err := http.ListenAndServe(cfg.ServerHost, router); err != nil {
		l.Fatal().Err(err).Msg("failed to start server")
	}
}

func registerMiddlewares(l *zerolog.Logger, router *mux.Router) {
	router.Use(
		middleware.WithContentType,
		func(next http.Handler) http.Handler {
			return middleware.WithLogging(l, next)
		},
		middleware.CompressHandler,
	)
}

func registerHandlers(st shorturl.StorageInterface, config *config.Config, l *zerolog.Logger, router *mux.Router) {
	u := shorturl.NewShortURL(config, st, generator.NewRandStringGenerator(), l)

	router.HandleFunc("/", u.PutHandler).Methods(http.MethodPost)
	router.HandleFunc("/{path:.*}", u.GetHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/shorten", u.PutJSONHandler).Methods(http.MethodPost)
}

func initStorage(config *config.Config) (shorturl.StorageInterface, error) {
	var m []storage.Entity

	if config.Restore {
		s, err := snapshot.FileReadAll[storage.Entity](config.FileStoragePath)
		if err != nil {
			return nil, err
		}
		m = s
	}

	return storage.NewFileStorage(config.FileStoragePath, storage.NewMemStorage(m)), nil
}
