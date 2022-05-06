package storage

import "errors"

type Storage struct {
	List map[string]string
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		List: make(map[string]string),
	}
}

func (s *Storage) AddValue(key, value string) {
	s.List[key] = value
}

func (s *Storage) GetValue(key string) (string, error) {
	if value, found := s.List[key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}
