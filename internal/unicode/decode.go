package unicode

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// DecodeEscapes converts \uXXXX and \UXXXXXXXX sequences to runes.
func DecodeEscapes(value string) string {
	var b strings.Builder
	b.Grow(len(value))

	for i := 0; i < len(value); {
		if value[i] != '\\' || i+1 >= len(value) {
			b.WriteByte(value[i])
			i++
			continue
		}

		switch value[i+1] {
		case 'u':
			if i+6 <= len(value) {
				r, err := parseHexRune(value[i+2 : i+6])
				if err == nil {
					b.WriteRune(r)
					i += 6
					continue
				}
			}
		case 'U':
			if i+10 <= len(value) {
				r, err := parseHexRune(value[i+2 : i+10])
				if err == nil {
					b.WriteRune(r)
					i += 10
					continue
				}
			}
		}

		b.WriteByte(value[i])
		i++
	}

	return b.String()
}

func parseHexRune(hex string) (rune, error) {
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, err
	}

	r := rune(v)
	if !utf8.ValidRune(r) {
		return 0, fmt.Errorf("invalid rune: %U", r)
	}

	return r, nil
}
