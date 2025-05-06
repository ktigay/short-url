package storage

type Entity struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type MemStorage struct {
	links map[string]Entity
}

func NewMemStorage(m []Entity) *MemStorage {
	var l = make(map[string]Entity)
	for _, e := range m {
		l[e.ShortURL] = e
	}
	return &MemStorage{
		links: l,
	}
}

func (s *MemStorage) Link(key string) (*Entity, error) {
	if v, ok := s.links[key]; ok {
		return &v, nil
	}
	return nil, nil
}

func (s *MemStorage) Unlink(key string) error {
	delete(s.links, key)
	return nil
}

func (s *MemStorage) PutLink(key string, value string) (*Entity, error) {
	l := len(s.links)
	s.links[key] = Entity{
		UUID:        l + 1,
		ShortURL:    key,
		OriginalURL: value,
	}
	e := s.links[key]

	return &e, nil
}

func (s *MemStorage) ShortLink(v string) string {
	for k, x := range s.links {
		if x.OriginalURL == v {
			return k
		}
	}
	return ""
}
