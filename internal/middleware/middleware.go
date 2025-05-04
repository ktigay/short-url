package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"

	ihttp "github.com/ktigay/short-url/internal/http"
)

// WithContentType устанавливает в ResponseWriter Content-Type.
func WithContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// WithLogging логирует запрос.
func WithLogging(l *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &ihttp.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := ihttp.NewWriter(w, rd)
		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		l.Info().
			Str("requestURI", r.RequestURI).
			Str("method", r.Method).
			Int64("duration", int64(duration)).
			Msg("request")

		l.Info().
			Int("status", rd.Status).
			Int("size", rd.Size).
			Msg("response")
	})
}
