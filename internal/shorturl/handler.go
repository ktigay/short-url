package shorturl

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/short-url/internal/config"
	"io"
	"net/http"
)

const (
	minShortLinkLength = 6
	maxShortLinkLength = 8
)

// StorageInterface - интерфейс хранилища.
type StorageInterface interface {
	Link(key string) (string, error)
	Unlink(key string) error
	PutLink(key string, value string)
}

type StringGeneratorInterface interface {
	Generate(min int, mix int) string
}

// ShortURL - структура обработчиков коротких ссылок.
type ShortURL struct {
	config    *config.Config
	storage   StorageInterface
	generator StringGeneratorInterface
}

// NewShortURL - конструктор.
func NewShortURL(config *config.Config, storage StorageInterface, gen StringGeneratorInterface) *ShortURL {
	return &ShortURL{
		config:    config,
		storage:   storage,
		generator: gen,
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

	shortLink := s.generator.Generate(minShortLinkLength, maxShortLinkLength)
	s.storage.PutLink(shortLink, link)

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
	if link == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}
