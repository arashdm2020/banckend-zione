package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// Compiled RegEx for replacing special characters
	regExp = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	// Compiled RegEx for replacing multiple hyphens
	multipleHyphenRegExp = regexp.MustCompile(`-+`)
)

// GenerateSlug generates a slug from a string
func GenerateSlug(str string) string {
	// Convert to lowercase
	slug := strings.ToLower(str)

	// Replace special characters with hyphen
	slug = regExp.ReplaceAllString(slug, "-")

	// Replace multiple hyphens with single hyphen
	slug = multipleHyphenRegExp.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateUniqueSlug generates a unique slug by appending a timestamp
func GenerateUniqueSlug(str string) string {
	return fmt.Sprintf("%s-%d", GenerateSlug(str), time.Now().Unix())
}

// RemoveAccents removes accents from string
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		return s
	}
	return output
}

// SanitizeSlug sanitizes a slug by handling unicode characters and spaces
func SanitizeSlug(s string) string {
	// Remove accents
	s = RemoveAccents(s)
	
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	
	// Generate slug
	return GenerateSlug(s)
} 