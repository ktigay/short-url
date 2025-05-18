package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"github.com/ktigay/short-url/internal/generator"
	"github.com/ktigay/short-url/internal/log"
	"github.com/ktigay/short-url/internal/middleware"
	"github.com/ktigay/short-url/internal/shorturl"
	"github.com/ktigay/short-url/internal/snapshot"
	"github.com/ktigay/short-url/internal/storage"
)

func main() {
	log.Initialize()

	var (
		cfg    *config.Config
		err    error
		st     shorturl.StorageInterface
		router *mux.Router
	)

	router = mux.NewRouter()

	if cfg, err = parseFlags(os.Args[1:]); err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to parse flags")
	}

	if st, err = initStorage(cfg.Restore, cfg.FileStoragePath); err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to initialize storage")
	}

	registerMiddlewares(router)
	registerHandlers(st, cfg, router)

	if err = http.ListenAndServe(cfg.ServerHost, router); err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to start server")
	}
}

func registerMiddlewares(router *mux.Router) {
	router.Use(
		middleware.WithContentType,
		func(next http.Handler) http.Handler {
			return middleware.WithLogging(next)
		},
		middleware.CompressHandler,
	)
}

func registerHandlers(st shorturl.StorageInterface, config *config.Config, router *mux.Router) {
	u := shorturl.NewShortURL(config, st, generator.NewRandStringGenerator())

	router.HandleFunc("/", u.PutHandler).Methods(http.MethodPost)
	router.HandleFunc("/{path:.*}", u.GetHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/shorten", u.PutJSONHandler).Methods(http.MethodPost)
}

func initStorage(restore bool, filePath string) (shorturl.StorageInterface, error) {
	var m []storage.Entity

	return storage.NewFileStorage(
		snapshot.NewFileSnapshot(filePath),
		storage.NewMemStorage(m),
		restore,
	)
}
