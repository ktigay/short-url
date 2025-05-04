package storage

type MemStorage struct {
	links map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		links: make(map[string]string),
	}
}

func (s *MemStorage) Link(key string) (string, error) {
	return s.links[key], nil
}

func (s *MemStorage) Unlink(key string) error {
	delete(s.links, key)
	return nil
}

func (s *MemStorage) PutLink(key string, value string) {
	s.links[key] = value
}
