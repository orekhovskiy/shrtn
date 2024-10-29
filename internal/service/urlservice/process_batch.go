package urlservice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (s *URLShortenerService) ProcessBatch(batch []entity.BatchRequest) ([]entity.BatchResponse, error) {
	records := make([]entity.URLRecord, len(batch))
	responses := make([]entity.BatchResponse, 0)

	for i, req := range batch {
		hash := sha256.Sum256([]byte(req.OriginalURL))
		shortURL := hex.EncodeToString(hash[:])[:7]
		records[i] = entity.URLRecord{
			UUID:        uuid.New().String(),
			OriginalURL: req.OriginalURL,
			ShortURL:    shortURL,
		}
	}
	savedRecords, err := s.urlRepository.SaveMany(records)
	if err != nil {
		return nil, err
	}

	for _, record := range savedRecords {
		var correlationID string
		for _, req := range batch {
			if req.OriginalURL == record.OriginalURL {
				correlationID = req.CorrelationID
				break
			}
		}
		responses = append(responses, entity.BatchResponse{
			ShortURL:      fmt.Sprintf("%s/%s", s.options.BaseURL, record.ShortURL),
			CorrelationID: correlationID,
		})
	}

	return responses, nil
}
