package storage

import "github.com/jmoiron/sqlx"

const (
	urlsTable = "urls"
)

func NewPostgresDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
