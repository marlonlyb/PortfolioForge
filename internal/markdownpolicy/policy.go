package markdownpolicy

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidSourceURL   = errors.New("markdown source url invalid")
	ErrSourceHostRejected = errors.New("markdown source host rejected")
)

var allowedHosts = map[string]struct{}{
	"mlbautomation.com":     {},
	"www.mlbautomation.com": {},
}

func AllowedHosts() []string {
	return []string{"mlbautomation.com", "www.mlbautomation.com"}
}

func ValidateSourceURL(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme != "https" || strings.TrimSpace(parsed.Host) == "" {
		return ErrInvalidSourceURL
	}

	if !IsAllowedHost(parsed.Hostname()) {
		return ErrSourceHostRejected
	}

	return nil
}

func IsAllowedSourceURL(raw string) bool {
	return ValidateSourceURL(raw) == nil && strings.TrimSpace(raw) != ""
}

func IsAllowedHost(host string) bool {
	_, ok := allowedHosts[strings.ToLower(strings.TrimSpace(host))]
	return ok
}

func SanitizeSourceURL(raw string) string {
	if !IsAllowedSourceURL(raw) {
		return ""
	}

	return strings.TrimSpace(raw)
}
