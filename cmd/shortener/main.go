package main

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/generator"
	"github.com/ktigay/short-url/internal/shorturl"
	"github.com/ktigay/short-url/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	config, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	})

	u := shorturl.NewShortURL(config, storage.NewMemStorage(), generator.NewRandStringGenerator())

	router.HandleFunc("/", u.PutHandler).Methods(http.MethodPost)
	router.HandleFunc("/{path:.*}", u.GetHandler).Methods(http.MethodGet)

	if err := http.ListenAndServe(config.ServerHost, router); err != nil {
		log.Fatal(err)
	}
}
