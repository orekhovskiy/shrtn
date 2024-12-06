package urlrepo

import "fmt"

func (r *PostgresURLRepository) MarkURLsAsDeleted(shortURLs []string, userID string) error {
	markURLAsDeleted := "UPDATE url_records SET is_deleted = TRUE WHERE short_url = ANY($1) AND user_id = $2;"
	_, err := r.db.Exec(markURLAsDeleted, shortURLs, userID)
	if err != nil {
		return fmt.Errorf("failed to mark URLs as deleted: %w", err)
	}
	return nil
}
