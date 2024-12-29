package utils

import (
	"strings"
	"unicode"
)

func NormalizeString(s string) string {
	// Önce tüm stringi uppercase yapalım
	s = strings.ToUpper(s)

	// Türkçe karakterleri değiştirelim
	replacements := map[string]string{
		"İ": "I", // Türkçe büyük İ -> I
		"Ş": "S",
		"Ğ": "G",
		"Ü": "U",
		"Ö": "O",
		"Ç": "C",
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	// Unicode normalizasyonu
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r == 'İ' {
			result = append(result, 'I')
		} else if unicode.IsLetter(r) {
			result = append(result, unicode.ToUpper(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
} 