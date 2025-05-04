package main

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"github.com/rs/zerolog"
	"net/http"
	"os"

	"github.com/ktigay/short-url/internal/generator"
	"github.com/ktigay/short-url/internal/logger"
	"github.com/ktigay/short-url/internal/middleware"
	"github.com/ktigay/short-url/internal/shorturl"
	"github.com/ktigay/short-url/internal/storage"
)

func main() {
	var l = logger.Initialize()

	router := mux.NewRouter()

	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		l.Fatal().Err(err).Msg("failed to parse flags")
	}

	registerMiddlewares(l, router)
	registerHandlers(cfg, l, router)

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

func registerHandlers(config *config.Config, l *zerolog.Logger, router *mux.Router) {
	u := shorturl.NewShortURL(config, storage.NewMemStorage(), generator.NewRandStringGenerator(), l)

	router.HandleFunc("/", u.PutHandler).Methods(http.MethodPost)
	router.HandleFunc("/{path:.*}", u.GetHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/shorten", u.PutJSONHandler).Methods(http.MethodPost)
}
