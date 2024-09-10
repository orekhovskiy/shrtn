package urlservice

import (
	"github.com/orekhovskiy/shrtn/pkg/util"
)

func (s Service) Save(url string) string {
	id := util.ShortenURL(url)
	s.urlRepository.Save(id, url)

	return id
}
