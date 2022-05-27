package storage

import (
	"errors"
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
)

type URLWorker interface {
	AddValue(key, value, userId string) error
	AddBatch(input []*u.ShortenBatchInput) error
	GetValue(key string) (string, error)
	GetAllValues(userId string) ([]u.URLItem, error)
	Ping() (bool, error)
}

type Storage struct {
	URLWorker
	List   map[string]map[string]string
	Config *config.Config
}

// NewStorage ...
func NewStorage(config *config.Config) *Storage {
	return &Storage{
		List:   make(map[string]map[string]string),
		Config: config,
	}
}

func (s *Storage) AddValue(key, value, userId string) error {
	if _, found := s.List[userId]; !found {
		s.List[userId] = make(map[string]string)
	}
	s.List[userId][key] = value
	return nil
}

func (s *Storage) AddBatch(input []*u.ShortenBatchInput) error {

	if _, found := s.List[input[0].UserID]; !found {
		s.List[input[0].UserID] = make(map[string]string)
	}

	for _, item := range input {
		s.List[item.UserID][item.ShortenURL] = item.OriginalURL
	}

	return nil
}

func (s *Storage) GetValue(key string) (string, error) {

	for _, dict := range s.List {
		if value, found := dict[key]; found {
			return value, nil
		}
	}

	return "", errors.New("the map didn't contains this key")
}

func (s *Storage) GetAllValues(userId string) ([]u.URLItem, error) {
	var items []u.URLItem

	if _, found := s.List[userId]; !found {
		return items, errors.New("this user haven't got any urls")
	}

	for key, value := range s.List[userId] {
		items = append(items, u.URLItem{ShortURL: fmt.Sprintf("%s/%s", s.Config.BaseURL, key), OriginalURL: value})
	}

	return items, nil
}

func (s *Storage) Ping() (bool, error) {
	return false, errors.New("storage in memory is active")
}
