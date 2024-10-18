package errors

import "fmt"

type URLConflictError struct {
	OriginalURL string
	ShortURL    string
}

func (e *URLConflictError) Error() string {
	return fmt.Sprintf("URL conflict: original URL %s already exists with short URL %s", e.OriginalURL, e.ShortURL)
}
