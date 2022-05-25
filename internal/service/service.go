package service

import (
	"fmt"
	u "github.com/AXlIS/go-shortener"
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

func (s *Service) AddURL(url, userId string) (string, error) {
	shortURL := utils.GenerateString(url)
	if err := s.storage.AddValue(shortURL, url, userId); err != nil {
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

func (s *Service) GetAllURLS(userId string) ([]u.URLItem, error) {
	urls, err := s.storage.GetAllValues(userId)
	return urls, err
}
