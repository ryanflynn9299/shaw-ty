package encoder

import (
	"slices"
	"strings"
)

const BASE63_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+"

// Base63Encode radix-63 encodes the provided integer using the characters [A-Z][a-z][0-9](+)
func Base63Encode(id int64) string {
	if id == int64(0) {
		return "A"
	}

	var encoded []rune
	for id > 0 {
		rem := id % 63
		id = id / 63
		encoded = append(encoded, rune(BASE63_CHARS[rem]))
	}

	// Reverse the slice to correct the order
	slices.Reverse(encoded)
	return string(encoded)
}

// Base63Decode radix-63 decodes the provided encoded string using the characters [A-Z][a-z][0-9](+) to the integer it represents
func Base63Decode(s string) int64 {
	var id int64
	for _, c := range s {
		index := strings.Index(BASE63_CHARS, string(c))
		if index == -1 {
			panic("Character not found: " + string(c))
		}
		id = id*63 + int64(index)
	}
	return id
}
