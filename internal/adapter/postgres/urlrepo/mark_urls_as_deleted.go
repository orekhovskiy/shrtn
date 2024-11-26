package urlrepo

import "fmt"

func (r *PostgresURLRepository) MarkURLsAsDeleted(shortURLs []string, userID string) error {
	_, err := r.db.Exec(markURLAsDeleted, shortURLs, userID)
	if err != nil {
		return fmt.Errorf("failed to mark URLs as deleted: %w", err)
	}
	return nil
}
