package snapshot

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ktigay/short-url/internal/log"
	"github.com/ktigay/short-url/internal/storage"
)

// FileSnapshot структура для сохранения снапшота в файле.
type FileSnapshot struct {
	filePath string
}

// NewFileSnapshot конструктор.
func NewFileSnapshot(filePath string) *FileSnapshot {
	return &FileSnapshot{filePath: filePath}
}

// Read чтение снапшота из файла.
func (f *FileSnapshot) Read() ([]storage.Entity, error) {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Logger.Error().Err(err).Str("file", f.filePath).Msg("failed to close file")
		}
	}()
	var all = make([]storage.Entity, 0)

	dec := json.NewDecoder(file)

	for {
		var e storage.Entity
		if err = dec.Decode(&e); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		all = append(all, e)
	}

	return all, nil
}

// WriteAll запись данных в файл.
func (f *FileSnapshot) Write(entities []storage.Entity) error {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return err
	}

	writer, err := NewAtomicFileWriter(f.filePath)
	if err != nil {
		return err
	}

	for _, el := range entities {
		if err = writer.Write(el); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}
	// перезапись происходит только при успешном закрытии writer.
	return writer.Close()
}

func ensureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dirName, os.ModeDir)
}
