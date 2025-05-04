package shorturl

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"github.com/ktigay/short-url/internal/generator"
	"github.com/ktigay/short-url/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortUrl_GetHandler(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "Bad_request",
			args: args{
				url: "/",
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        false,
		},
		{
			name: "Not_Found",
			args: args{
				url: "/AdsKfd",
			},
			wantStatusCode: http.StatusNotFound,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			s := storage.NewMemStorage()

			h := NewShortURL(cfg, s, generator.NewRandStringGenerator())

			router := mux.NewRouter()
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("content-type", "text/plain; charset=utf-8")
					next.ServeHTTP(w, r)
				})
			})
			router.HandleFunc("/{path:.*}", h.GetHandler).Methods(http.MethodGet)

			svr := httptest.NewServer(router)
			defer svr.Close()

			resp, err := http.Get(svr.URL + tt.args.url)
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantErr, err != nil)

			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()
		})
	}
}

func TestShortUrl_PutHandler(t *testing.T) {
	type args struct {
		url       string
		body      string
		shortURL  string
		serverURL string
	}
	tests := []struct {
		name           string
		args           args
		want           string
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "Positive_test",
			args: args{
				serverURL: "http://localhost:8080",
				url:       "/",
				body:      "http://asssddsd.dd/asdd/asddd/ddd.html",
				shortURL:  "DfhGfd",
			},
			want:           "http://localhost:8080/DfhGfd",
			wantStatusCode: http.StatusCreated,
			wantErr:        false,
		},
		{
			name: "Bad_request",
			args: args{
				serverURL: "http://localhost:8080",
				url:       "/",
				body:      "",
			},
			want:           "",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ServerURL: tt.args.serverURL,
			}
			s := storage.NewMemStorage()

			g := &MockGenerator{str: tt.args.shortURL}
			h := NewShortURL(cfg, s, g)

			router := mux.NewRouter()
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("content-type", "text/plain; charset=utf-8")
					next.ServeHTTP(w, r)
				})
			})
			router.HandleFunc("/", h.PutHandler).Methods(http.MethodPost)

			svr := httptest.NewServer(router)
			defer svr.Close()

			resp, err := http.Post(svr.URL+tt.args.url, "text/plain", strings.NewReader(tt.args.body))

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantErr, err != nil)

			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, string(body))

			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()
		})
	}
}

type MockGenerator struct {
	str string
}

func (m *MockGenerator) Generate(_ int, _ int) string {
	return m.str
}
