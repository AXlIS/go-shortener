package storage

import (
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/jmoiron/sqlx"
	"log"
)

type DatabaseStorage struct {
	URLWorker
	config *config.Config
	db     *sqlx.DB
}

func NewDatabaseStorage(db *sqlx.DB, config *config.Config) *DatabaseStorage {
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( 
												id         SERIAL PRIMARY KEY,
												user_id    VARCHAR(128) NOT NULL,
												short_url  VARCHAR(128) NOT NULL,
												base_url   VARCHAR(128) NOT NULL,
                                                created_at timestamp NOT NULL DEFAULT NOW(),
	                                            UNIQUE (user_id, short_url)
                                           );`, urlsTable)
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("create table error: %s", err.Error())
	}

	return &DatabaseStorage{
		db:     db,
		config: config,
	}
}

func (s *DatabaseStorage) GetValue(key string) (string, error) {

	var URL string

	getValueQuery := fmt.Sprintf(`SELECT base_url FROM %s WHERE short_url = $1 LIMIT 1`, urlsTable)
	row := s.db.QueryRow(getValueQuery, fmt.Sprintf("%s/%s", s.config.BaseURL, key))
	if err := row.Scan(&URL); err != nil {
		return "", err
	}

	return URL, nil
}

func (s *DatabaseStorage) AddValue(key, value, userId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	addValueQuery := fmt.Sprintf(`INSERT INTO %s (user_id, short_url, base_url) VALUES ($1, $2, $3)`, urlsTable)
	_, err = s.db.Exec(addValueQuery, userId, fmt.Sprintf("%s/%s", s.config.BaseURL, key), value)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *DatabaseStorage) GetAllValues(userId string) ([]u.URLItem, error) {
	var URLS []u.URLItem
	query := fmt.Sprintf(`SELECT short_url, base_url FROM %s WHERE user_id = $1`, urlsTable)

	if err := s.db.Select(&URLS, query, userId); err != nil {
		return nil, err
	}

	return URLS, nil
}

func (s *DatabaseStorage) Ping() (bool, error) {
	err := s.db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}
