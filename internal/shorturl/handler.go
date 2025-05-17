package shorturl

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"github.com/ktigay/short-url/internal/storage"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

const (
	minShortLinkLength = 6
	maxShortLinkLength = 8
)

// StorageInterface - интерфейс хранилища.
type StorageInterface interface {
	Link(key string) (*storage.Entity, error)
	Unlink(key string) error
	PutLink(key string, value string) (*storage.Entity, error)
	ShortLink(v string) string
}

type StringGeneratorInterface interface {
	Generate(min int, mix int) string
}

// ShortURL - структура обработчиков коротких ссылок.
type ShortURL struct {
	config    *config.Config
	storage   StorageInterface
	generator StringGeneratorInterface
	logger    *zerolog.Logger
}

// NewShortURL - конструктор.
func NewShortURL(config *config.Config, storage StorageInterface, gen StringGeneratorInterface, l *zerolog.Logger) *ShortURL {
	return &ShortURL{
		config:    config,
		storage:   storage,
		generator: gen,
		logger:    l,
	}
}

// PutHandler - сохранение ссылки
func (s *ShortURL) PutHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link := string(body)
	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortLink string
	if sh := s.storage.ShortLink(link); sh != "" {
		shortLink = sh
	} else {
		shortLink = s.generator.Generate(minShortLinkLength, maxShortLinkLength)
		_, _ = s.storage.PutLink(shortLink, link)
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(s.config.ServerURL + "/" + shortLink))
}

// GetHandler - получение ссылки
func (s *ShortURL) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["path"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := s.storage.Link(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if link == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
}

func (s *ShortURL) PutJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqJSON = new(struct {
		URL string `json:"url"`
	})

	if err = json.Unmarshal(body, &reqJSON); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortLink string
	if sh := s.storage.ShortLink(reqJSON.URL); sh != "" {
		shortLink = sh
	} else {
		shortLink = s.generator.Generate(minShortLinkLength, maxShortLinkLength)
		_, _ = s.storage.PutLink(shortLink, reqJSON.URL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		Result string `json:"result"`
	}{
		Result: s.config.ServerURL + "/" + shortLink,
	}); err != nil {
		s.logger.Error().Err(err).Msg("Failed to write response")
	}
}
