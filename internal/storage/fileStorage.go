package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type FileStorage struct {
	FilePath string
	List     map[string]string
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fmt.Println(1, filePath)
	var storage = &FileStorage{FilePath: filePath}

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		storage.List = make(map[string]string)
		return storage, nil
	}

	if err := json.Unmarshal(data, &storage.List); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *FileStorage) AddValue(key, value string) error {
	s.List[key] = value

	file, err := os.OpenFile(s.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s.List)
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
	if value, found := s.List[key]; found {
		return value, nil
	}
	return "", errors.New("storage didn't contains this key")
}
