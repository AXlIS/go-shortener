package service

import (
	"fmt"
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
)

type Service struct {
	storage storage.URLWorker
}

func NewService(storage storage.URLWorker) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddURL(url string) (string, error) {
	shortURL := utils.GenerateShortURL(url)
	if err := s.storage.AddValue(shortURL, url); err != nil {
		fmt.Println(5)
		return "", err
	}
	return shortURL, nil
}

func (s *Service) GetURL(key string) (string, error) {
	url, err := s.storage.GetValue(key)
	if err != nil {
		return "", err
	}

	return url, nil
}
