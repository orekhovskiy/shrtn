package urlservice

func (s *URLShortenerService) MarkURLsAsDeleted(shortURLs []string, userID string) error {
	updateChan := make(chan []string)
	errors := make(chan error)

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
		go func() {
			for batch := range updateChan {
				if err := s.urlRepository.MarkURLsAsDeleted(batch, userID); err != nil {
					errors <- err
				}
			}
		}()
	}

	close(errors)
	return nil
}
