package storage

import "errors"

type URLWorker interface {
	AddValue(key, value string) error
	GetValue(key string) (string, error)
}

type Storage struct {
	List map[string]string
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		List: make(map[string]string),
	}
}

func (s *Storage) AddValue(key, value string) error {
	s.List[key] = value
	return nil
}

func (s *Storage) GetValue(key string) (string, error) {
	if value, found := s.List[key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}
