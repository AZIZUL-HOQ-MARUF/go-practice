// Package stringutil provides string utility functions.
package stringutil

import "strings"

// Reverse returns the reverse of s, correctly handling Unicode.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome reports whether s reads the same forwards and backwards.
// Case-insensitive, ignores spaces.
func IsPalindrome(s string) bool {
	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// WordCount returns a map of word frequencies in s.
func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(s) {
		counts[strings.ToLower(word)]++
	}
	return counts
}
