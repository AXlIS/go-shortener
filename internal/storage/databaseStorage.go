package storage

import (
	"fmt"
	urls "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

type DatabaseStorage struct {
	config *config.Config
	db     *sqlx.DB
}

func NewDatabaseStorage(db *sqlx.DB, config *config.Config) *DatabaseStorage {
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( 
												id         SERIAL PRIMARY KEY,
												user_id    VARCHAR(32) NOT NULL,
												short_url  VARCHAR(32) NOT NULL,
												base_url   VARCHAR(255) NOT NULL,
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

func (s *DatabaseStorage) AddValue(key, value, userID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	addValueQuery := fmt.Sprintf(`INSERT INTO %s (user_id, short_url, base_url) VALUES ($1, $2, $3)`, urlsTable)
	_, err = tx.Exec(addValueQuery, userID, fmt.Sprintf("%s/%s", s.config.BaseURL, key), value)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return err
		}
	}

	if err != nil {
		err = tx.Rollback()
		log.Printf("Error: %s", err.Error())
		return err
	}

	return tx.Commit()
}

func (s *DatabaseStorage) AddBatch(input []*urls.ShortenBatchInput) error {
	for _, item := range input {
		item.ShortenURL = fmt.Sprintf("%s/%s", s.config.BaseURL, item.ShortenURL)
	}

	insertQuery := fmt.Sprintf(`INSERT INTO %s (user_id, short_url, base_url) 
                                      VALUES (:user_id, :short_url, :base_url)`, urlsTable)
	if _, err := s.db.NamedExec(insertQuery, input); err != nil {
		return err
	}

	return nil
}

func (s *DatabaseStorage) GetAllValues(userID string) ([]urls.Item, error) {
	var URLS []urls.Item
	query := fmt.Sprintf(`SELECT short_url, base_url FROM %s WHERE user_id = $1`, urlsTable)

	if err := s.db.Select(&URLS, query, userID); err != nil {
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
