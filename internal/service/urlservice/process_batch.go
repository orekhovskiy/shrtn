package urlservice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (s *URLShortenerService) ProcessBatch(batch []entity.BatchRequest, userID string) ([]entity.BatchResponse, error) {
	correlationMap := buildCorrelationMap(batch)
	records := s.generateURLRecords(batch)

	fmt.Println(correlationMap)
	fmt.Println(records)
	fmt.Println(batch)

	savedRecords, err := s.urlRepository.SaveMany(records, userID)
	if err != nil {
		return nil, err
	}

	responses := buildBatchResponses(savedRecords, correlationMap, s.BuildURL)
	return responses, nil
}

func buildCorrelationMap(batch []entity.BatchRequest) map[string]string {
	correlationMap := make(map[string]string, len(batch))
	for _, req := range batch {
		correlationMap[req.OriginalURL] = req.CorrelationID
	}
	return correlationMap
}

func (s *URLShortenerService) generateURLRecords(batch []entity.BatchRequest) []entity.URLRecord {
	records := make([]entity.URLRecord, len(batch))
	for i, req := range batch {

		records[i] = entity.URLRecord{
			UUID:        uuid.New().String(),
			OriginalURL: req.OriginalURL,
			ShortURL:    s.createShortURL(req.OriginalURL),
		}
	}
	return records
}

func buildBatchResponses(
	records []entity.URLRecord,
	correlationMap map[string]string,
	buildURL func(string) string,
) []entity.BatchResponse {
	responses := make([]entity.BatchResponse, 0, len(records))
	for _, record := range records {
		correlationID := correlationMap[record.OriginalURL]
		responses = append(responses, entity.BatchResponse{
			ShortURL:      buildURL(record.ShortURL),
			CorrelationID: correlationID,
		})
	}
	return responses
}
