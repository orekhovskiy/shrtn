package urlrepo

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/orekhovskiy/shrtn/config"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(config config.Config) (*Repository, error) {
	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableIfNotExists); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}
