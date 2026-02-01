package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Normalize unicode characters
	s = norm.NFD.String(s)

	// Remove diacritics
	var builder strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		builder.WriteRune(r)
	}
	s = builder.String()

	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace Vietnamese special characters
	replacements := map[string]string{
		"đ": "d",
		"Đ": "d",
	}
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")

	return s
}

// CountWords counts the number of words in a string
func CountWords(s string) int {
	fields := strings.Fields(s)
	return len(fields)
}
