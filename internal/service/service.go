package service

import (
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
)

type Service struct {
	config  *config.Config
	storage storage.URLWorker
}

func NewService(storage storage.URLWorker, config *config.Config) *Service {
	return &Service{
		storage: storage,
		config:  config,
	}
}

func (s *Service) AddURL(url, userId string) (string, error) {
	shortURL := utils.GenerateString(url)
	if err := s.storage.AddValue(shortURL, url, userId); err != nil {
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

func (s *Service) AddBatchURL(urls []*u.ShortenBatchInput, userId string) ([]u.ShortenBatchResponse, error) {
	var shortenURLS []u.ShortenBatchResponse

	for _, item := range urls {
		shortURL := utils.GenerateString(item.OriginalURL)
		item.ShortenURL, item.UserID = shortURL, userId
		shortenURLS = append(shortenURLS, u.ShortenBatchResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.config.BaseURL, shortURL),
		})
	}

	if err := s.storage.AddBatch(urls); err != nil {
		return nil, err
	}

	return shortenURLS, nil
}

func (s *Service) Ping() (bool, error) {
	ping, err := s.storage.Ping()
	return ping, err
}
