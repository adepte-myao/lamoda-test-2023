package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB(connURL string) (*sql.DB, func(db *sql.DB), error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	return db, closePostgresDb, nil
}

func closePostgresDb(db *sql.DB) {
	_ = db.Close()
}
