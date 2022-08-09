package storage

import (
	"errors"
	"fmt"
	"log"

	urls "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
)

type URLWorker interface {
	AddValue(key, value, userID string) error
	AddBatch(input []*urls.ShortenBatchInput) error
	GetValue(key string) (string, error)
	GetAllValues(userID string) ([]urls.Item, error)
	DeleteValues(urls []string, userID string)
	Ping() (bool, error)
}

type Storage struct {
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

func (s *Storage) AddValue(key, value, userID string) error {
	if _, found := s.List[userID]; !found {
		s.List[userID] = make(map[string]string)
	}
	s.List[userID][key] = value
	return nil
}

func (s *Storage) AddBatch(input []*urls.ShortenBatchInput) error {

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

func (s *Storage) GetAllValues(userID string) ([]urls.Item, error) {
	var items []urls.Item

	if _, found := s.List[userID]; !found {
		return items, errors.New("this user haven't got any urls")
	}

	for key, value := range s.List[userID] {
		items = append(items, urls.Item{ShortURL: fmt.Sprintf("%s/%s", s.Config.BaseURL, key), OriginalURL: value})
	}

	return items, nil
}

func (s *Storage) Ping() (bool, error) {
	return false, errors.New("storage in memory is active")
}

func (s *Storage) DeleteValues(urls []string, userID string) {
	log.Println("Delete from storage")
}
