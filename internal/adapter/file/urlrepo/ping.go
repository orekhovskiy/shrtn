package urlrepo

import (
	"errors"
)

func (r *FileURLRepository) Ping() error {
	return errors.New("no database connection")
}
