package urlservice

import (
	"crypto/md5"
	"encoding/hex"
)

func (s Service) Save(url string) string {
	hash := md5.Sum([]byte(url))
	id := hex.EncodeToString(hash[:])[:7]
	s.urlRepository.Save(id, url)

	return id
}
