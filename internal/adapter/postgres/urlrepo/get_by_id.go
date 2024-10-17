package urlrepo

import "errors"

func (r *Repository) GetByID(_ string) (string, error) {
	return "", errors.New("not implemented")
}
