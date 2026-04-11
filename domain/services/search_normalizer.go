package services

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeQuery applies the normalization pipeline to a raw search query:
// 1. Trim whitespace
// 2. Lowercase
// 3. Remove accents (NFD decomposition + strip combining marks)
// 4. Strip special characters (keep only letters, digits, spaces)
// 5. Collapse multiple spaces into one
func NormalizeQuery(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.ToLower(s)
	s = removeAccents(s)
	s = specialCharsRegex.ReplaceAllString(s, " ")
	s = spacesRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

// removeAccents strips diacritical marks from a string by decomposing into NFD
// and removing Unicode category Mn (Mark, Nonspacing) characters.
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

var (
	specialCharsRegex = regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	spacesRegex       = regexp.MustCompile(`\s+`)
)
