package storage

import "github.com/ktigay/short-url/internal/snapshot"

type FileStorage struct {
	mem  *MemStorage
	path string
}

func (f *FileStorage) Link(key string) (*Entity, error) {
	return f.mem.Link(key)
}

func (f *FileStorage) Unlink(key string) error {
	return f.mem.Unlink(key)
}

func (f *FileStorage) PutLink(key string, value string) (*Entity, error) {
	om, err := f.mem.Link(key)
	if err != nil {
		return nil, err
	}

	if om != nil {
		return om, nil
	}

	om, err = f.mem.PutLink(key, value)

	if err != nil {
		return nil, err
	}

	if err = snapshot.FileWrite[Entity](f.path, om); err != nil {
		_ = f.mem.Unlink(key)
		return nil, err
	}

	return om, nil
}

func (f *FileStorage) ShortLink(v string) string {
	return f.mem.ShortLink(v)
}

func NewFileStorage(path string, mem *MemStorage) *FileStorage {
	return &FileStorage{
		mem:  mem,
		path: path,
	}
}
