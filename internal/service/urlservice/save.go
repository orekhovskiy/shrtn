package urlservice

import (
	"crypto/md5"
	"encoding/hex"
)

func (s Service) Save(url string) (string, error) {
	hash := md5.Sum([]byte(url))
	id := hex.EncodeToString(hash[:])[:7]
	err := s.urlRepository.Save(id, url)
	if err != nil {
		return "", err
	}

	return id, nil
}
