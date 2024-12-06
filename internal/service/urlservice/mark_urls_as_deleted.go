package urlservice

import (
	"sync"
)

func (s *URLShortenerService) MarkURLsAsDeleted(shortURLs []string, userID string) []error {
	updateChan := make(chan []string)
	errorChan := make(chan error)
	var wg sync.WaitGroup

	// Fan-In
	go func() {
		batchSize := 10
		for i := 0; i < len(shortURLs); i += batchSize {
			end := i + batchSize
			if end > len(shortURLs) {
				end = len(shortURLs)
			}
			updateChan <- shortURLs[i:end]
		}
		close(updateChan)
	}()

	// Workers
	const workerCount = 3
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range updateChan {
				if err := s.urlRepository.MarkURLsAsDeleted(batch, userID); err != nil {
					errorChan <- err
				}
			}
		}()
	}

	// Close errorChan after all workers complete
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	return errors
}
