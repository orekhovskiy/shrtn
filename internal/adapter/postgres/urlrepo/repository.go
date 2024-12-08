package urlrepo

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/orekhovskiy/shrtn/config"
)

type PostgresURLRepository struct {
	db *sql.DB
}

func NewRepository(config config.Config, db *sql.DB) (*PostgresURLRepository, error) {
	createTableIfNotExists := `
	CREATE TABLE IF NOT EXISTS url_records (
		id SERIAL PRIMARY KEY,
		uuid TEXT NOT NULL,
		short_url TEXT UNIQUE NOT NULL,
		original_url TEXT UNIQUE NOT NULL,
		user_id TEXT NOT NULL,
		is_deleted BOOLEAN DEFAULT FALSE
	);`
	if _, err := db.Exec(createTableIfNotExists); err != nil {
		return nil, err
	}

	createIndexIfNotExists := "CREATE INDEX IF NOT EXISTS idx_user_id ON url_records (user_id);"
	if _, err := db.Exec(createIndexIfNotExists); err != nil {
		return nil, err
	}

	return &PostgresURLRepository{db: db}, nil
}
