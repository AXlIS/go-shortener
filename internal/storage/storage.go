package storage

import (
	"errors"
	u "github.com/AXlIS/go-shortener"
)

type URLWorker interface {
	AddValue(key, value, userId string) error
	GetValue(key, userId string) (string, error)
	GetAllValues(userId string) ([]u.URLItem, error)
}

type Storage struct {
	URLWorker
	List map[string]map[string]string
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		List: make(map[string]map[string]string),
	}
}

func (s *Storage) AddValue(key, value, userId string) error {
	if _, found := s.List[userId]; !found {
		s.List[userId] = make(map[string]string)
	}
	s.List[userId][key] = value
	return nil
}

func (s *Storage) GetValue(key, userId string) (string, error) {
	if value, found := s.List[userId][key]; found {
		return value, nil
	}
	return "", errors.New("the map didn't contains this key")
}

func (s *Storage) GetAllValues(userId string) ([]u.URLItem, error) {
	var items []u.URLItem

	if _, found := s.List[userId]; !found {
		return items, errors.New("this user haven't got any urls")
	}

	for key, value := range s.List[userId] {
		items = append(items, u.URLItem{ShortURL: key, OriginalURL: value})
	}

	return items, nil
}
