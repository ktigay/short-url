package middleware

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"

	"github.com/ktigay/short-url/internal"
	"github.com/ktigay/short-url/internal/compress"
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
			Str("duration", fmt.Sprint(duration)).
			Str("headers", fmt.Sprint(r.Header)).
			Msg("request")

		l.Info().
			Int("status", rd.Status).
			Int("size", rd.Size).
			Str("headers", fmt.Sprint(lw.Header())).
			Msg("response")
	})
}

var acceptTypes = []string{"text/html", "application/json", "*/*"}

// CompressHandler обработчик сжатия данных.
func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentEncoding := r.Header.Get("Content-Encoding")
		if ceAlg := compress.TypeFromString(contentEncoding); ceAlg != "" {
			cr, err := compress.ReaderFactory(ceAlg, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		accept := r.Header.Get("Accept")
		isAccepted := func() bool {
			for _, acceptType := range acceptTypes {
				if strings.Contains(accept, acceptType) {
					return true
				}
			}
			return false
		}()

		if isAccepted {
			if aeAlg := compress.TypeFromString(acceptEncoding); string(aeAlg) != "" {
				cw, _ := compress.NewHTTPWriter(aeAlg, w)
				w = cw
				defer internal.Quite(cw.Close)
			}
		}

		next.ServeHTTP(w, r)
	})
}
