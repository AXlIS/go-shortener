package storage

import "errors"

type Storage map[string]string

// NewStorage ...
func NewStorage() Storage {
	return make(Storage)
}


func (s *Storage) AddValue(key, value string) {
	(*s)[key] = value
}

func (s *Storage) GetValue(key string) (string, error) {
	if value, found := (*s)[key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}
