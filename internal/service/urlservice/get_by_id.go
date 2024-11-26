package urlservice

import (
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (s URLShortenerService) GetByID(id string) (*entity.Result, error) {
	record, err := s.urlRepository.GetByID(id)

	if err != nil {
		return nil, err
	}

	if record.IsDeleted {
		return &entity.Result{
			Success:     false,
			OriginalURL: "",
		}, nil
	}

	return &entity.Result{
		Success:     true,
		OriginalURL: record.OriginalURL,
	}, nil
}
