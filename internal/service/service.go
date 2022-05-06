package service

import (
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
)

type Service struct {
	storage *storage.Storage
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddURL(url string) string {
	shortURL := utils.GenerateShortURL(url)
	s.storage.AddValue(shortURL, url)
	return shortURL
}

func (s *Service) GetURL(key string) (string, error){
	url, err := s.storage.GetValue(key)
	if err != nil {
		return "", err
	}

	return url, nil
}
