package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"io"
	"os"
)

type FileStorage struct {
	URLWorker
	FilePath string
	List     map[string]map[string]string
	Config   *config.Config
}

func NewFileStorage(filePath string, config *config.Config) (*FileStorage, error) {

	var storage = &FileStorage{FilePath: filePath, Config: config}

	_ = os.Mkdir("/tmp", 0750)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		storage.List = make(map[string]map[string]string)
		return storage, nil
	}

	if err := json.Unmarshal(data, &storage.List); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *FileStorage) AddValue(key, value, userId string) error {

	if _, found := s.List[userId]; !found {
		s.List[userId] = make(map[string]string)
	}
	s.List[userId][key] = value

	file, err := os.OpenFile(s.FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.List, "", "	")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return file.Close()
}

func (s *FileStorage) GetValue(key string) (string, error) {
	for _, dict := range s.List {
		if value, found := dict[key]; found {
			return value, nil
		}
	}

	return "", errors.New("storage didn't contains this key")
}

func (s *FileStorage) GetAllValues(userId string) ([]u.URLItem, error) {
	var items []u.URLItem

	if _, found := s.List[userId]; !found {
		return items, errors.New("this user haven't got any urls")
	}

	for key, value := range s.List[userId] {
		items = append(items, u.URLItem{ShortURL: fmt.Sprintf("%s/%s", s.Config.BaseURL, key), OriginalURL: value})
	}

	return items, nil
}
