package urlrepo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/orekhovskiy/shrtn/internal/entity"
)

func (r *FileURLRepository) MarkURLsAsDeleted(shortURLs []string, userID string) error {
	// Lock the repository to ensure thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if there are URLs to delete
	if len(shortURLs) == 0 {
		return nil
	}

	// Read all records from the file and filter out the ones to be marked as deleted
	updatedRecords, err := r.filterDeletedURLs(shortURLs, userID)
	if err != nil {
		return err
	}

	// Mark the URLs as deleted in memory
	for _, shortURL := range shortURLs {
		if record, exists := r.records[shortURL]; exists && record.UserID == userID {
			record.IsDeleted = true
			r.records[shortURL] = record
		}
	}

	// If file path is empty, no need to save the changes to the file
	if r.filePath == "" {
		return nil
	}

	// Overwrite the file with the updated records
	if err := r.saveUpdatedRecordsToFile(updatedRecords); err != nil {
		return err
	}

	return nil
}

func (r *FileURLRepository) filterDeletedURLs(shortURLs []string, userID string) ([]entity.URLRecord, error) {
	var updatedRecords []entity.URLRecord

	// Read the current file and filter the records
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// Process each record in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record entity.URLRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}

		// Mark the record as deleted if it's in the provided short URLs list
		if contains(shortURLs, record.ShortURL) && record.UserID == userID {
			record.IsDeleted = true
		}

		updatedRecords = append(updatedRecords, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file: %w", err)
	}

	return updatedRecords, nil
}

func (r *FileURLRepository) saveUpdatedRecordsToFile(updatedRecords []entity.URLRecord) error {
	// Open the file for writing
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write each updated record to the file
	for _, record := range updatedRecords {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("error marshaling record: %w", err)
		}

		if _, err := writer.WriteString(string(data) + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	// Ensure all data is written to the file
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func contains(shortURLs []string, shortURL string) bool {
	for _, url := range shortURLs {
		if url == shortURL {
			return true
		}
	}
	return false
}
