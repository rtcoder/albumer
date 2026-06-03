package domain

import (
	"strings"
	"unicode"
)

func SearchKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var out strings.Builder
	for _, r := range value {
		for _, mapped := range foldRune(r) {
			if unicode.IsLetter(mapped) || unicode.IsDigit(mapped) {
				out.WriteRune(mapped)
			}
		}
	}
	return out.String()
}

func foldRune(r rune) string {
	switch r {
	case 'ą':
		return "a"
	case 'ć':
		return "c"
	case 'ę':
		return "e"
	case 'ł':
		return "l"
	case 'ń':
		return "n"
	case 'ó':
		return "o"
	case 'ś':
		return "s"
	case 'ż', 'ź':
		return "z"
	default:
		return string(r)
	}
}
