package urlrepo

import (
	"errors"
)

func (r *Repository) Ping() error {
	return errors.New("no database connection")
}
