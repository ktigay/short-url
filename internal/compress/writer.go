package compress

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

var encodings = map[Type]string{
	Gzip:    "gzip",
	Deflate: "deflate",
	Br:      "br",
}

// HTTPWriter структура для обработки сжатия ответа.
type HTTPWriter struct {
	writer     http.ResponseWriter
	compressor io.WriteCloser
}

// NewHTTPWriter конструктор.
func NewHTTPWriter(t Type, w http.ResponseWriter) (*HTTPWriter, error) {
	cmp, _ := compressor(t, w)
	if cmp != nil {
		w.Header().Set("Content-Encoding", encodings[t])
	}

	httpWr := HTTPWriter{
		writer:     w,
		compressor: cmp,
	}

	return &httpWr, nil
}

// Header возвращает заголовоки.
func (c *HTTPWriter) Header() http.Header {
	return c.writer.Header()
}

// Write записывает данные.
func (c *HTTPWriter) Write(p []byte) (int, error) {
	if c.compressor == nil {
		return c.writer.Write(p)
	}
	return c.compressor.Write(p)
}

// WriteHeader устанавливает заголовок ответа.
func (c *HTTPWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

// Close закрывает HTTPWriter и досылает все данные из буфера.
func (c *HTTPWriter) Close() error {
	if c.compressor == nil {
		return nil
	}
	return c.compressor.Close()
}

func compressor(t Type, w io.Writer) (io.WriteCloser, error) {
	switch t {
	case Gzip:
		return gzip.NewWriter(w), nil
	case Deflate:
		return zlib.NewWriter(w), nil
	case Br:
		return brotli.NewWriter(w), nil
	}
	return nil, fmt.Errorf("unsupported compress type: %v", t)
}
