package service

import (
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/storage"
	"github.com/AXlIS/go-shortener/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
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

func (s *Service) AddURL(url, userID string) (string, error) {
	shortURL := utils.GenerateString(url)
	err := s.storage.AddValue(shortURL, url, userID)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return shortURL, err
		}
	}

	if err != nil {
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

func (s *Service) GetAllURLS(userID string) ([]u.Item, error) {
	urls, err := s.storage.GetAllValues(userID)
	return urls, err
}

func (s *Service) AddBatchURL(urls []*u.ShortenBatchInput, userID string) ([]u.ShortenBatchResponse, error) {
	shortenURLS := make([]u.ShortenBatchResponse, 0, len(urls))

	for _, item := range urls {
		shortURL := utils.GenerateString(item.OriginalURL)
		item.ShortenURL, item.UserID = shortURL, userID
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
