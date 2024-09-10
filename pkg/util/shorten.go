package util

import (
	"crypto/md5"
	"encoding/hex"
)

func ShortenURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hasher.Sum(nil)
	id := hex.EncodeToString(hash)[:7]
	return id
}
