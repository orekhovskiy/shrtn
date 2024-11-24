package urlservice

import "strings"

func (s *URLShortenerService) BuildURL(uri string) string {
	var builder strings.Builder
	builder.WriteString(s.options.BaseURL)
	builder.WriteString("/")
	builder.WriteString(uri)
	return builder.String()
}
