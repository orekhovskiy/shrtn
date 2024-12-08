package urlservice

import (
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (s *URLShortenerService) ProcessBatch(batch []entity.BatchRequest, userID string) ([]entity.BatchResponse, error) {
	correlationMap := s.buildCorrelationMap(batch)
	records := s.generateURLRecords(batch)

	// Only save valid records (those where short URL creation was successful)
	savedRecords, err := s.urlRepository.SaveMany(records, userID)
	if err != nil {
		return nil, err
	}

	responses, err := s.buildBatchResponses(savedRecords, correlationMap)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *URLShortenerService) buildCorrelationMap(batch []entity.BatchRequest) map[string]string {
	correlationMap := make(map[string]string, len(batch))
	for _, req := range batch {
		correlationMap[req.OriginalURL] = req.CorrelationID
	}
	return correlationMap
}

func (s *URLShortenerService) generateURLRecords(batch []entity.BatchRequest) []entity.URLRecord {
	records := []entity.URLRecord{}
	for _, req := range batch {
		// Attempt to create the short URL
		shortURL, err := s.createShortURL(req.OriginalURL)
		if err != nil {
			continue
		}

		// If successful, create a new record to be saved
		records = append(records, entity.URLRecord{
			UUID:        uuid.New().String(),
			OriginalURL: req.OriginalURL,
			ShortURL:    shortURL,
		})
	}
	return records
}

func (s *URLShortenerService) buildBatchResponses(
	records []entity.URLRecord,
	correlationMap map[string]string,
) ([]entity.BatchResponse, error) {
	responses := make([]entity.BatchResponse, 0, len(records))
	for _, record := range records {
		correlationID := correlationMap[record.OriginalURL]
		shortURL, err := s.BuildURL(record.ShortURL)
		if err != nil {
			return nil, err
		}
		responses = append(responses, entity.BatchResponse{
			ShortURL:      shortURL,
			CorrelationID: correlationID,
		})
	}
	return responses, nil
}
