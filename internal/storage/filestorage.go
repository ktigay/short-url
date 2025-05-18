package storage

type Snapshot interface {
	Read() ([]Entity, error)
	Write([]Entity) error
}

type FileStorage struct {
	mem      *MemStorage
	snapshot Snapshot
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

	if err = f.snapshot.Write([]Entity{*om}); err != nil {
		_ = f.mem.Unlink(key)
		return nil, err
	}

	return om, nil
}

func (f *FileStorage) ShortLink(v string) string {
	return f.mem.ShortLink(v)
}

func (f *FileStorage) restore() error {
	if f.mem == nil {
		return nil
	}

	data, err := f.snapshot.Read()
	if err != nil {
		return err
	}
	for _, m := range data {
		f.mem.links[m.ShortURL] = m
	}

	return nil
}

func NewFileStorage(snapshot Snapshot, mem *MemStorage, restore bool) (*FileStorage, error) {
	st := FileStorage{
		mem:      mem,
		snapshot: snapshot,
	}
	if restore {
		if err := st.restore(); err != nil {
			return nil, err
		}
	}

	return &st, nil
}
