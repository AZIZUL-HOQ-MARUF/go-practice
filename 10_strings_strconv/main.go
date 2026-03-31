// Topic 10: Strings & strconv
// Run: go run 10_strings_strconv/main.go
//
// JS: strings are UTF-16, have .length, .slice(), .split() etc.
// Go: strings are immutable byte slices (UTF-8 encoded).
//     len(s) returns BYTE count, not character count.
//     Iterating with range gives runes (Unicode code points), not bytes.
//
// This distinction matters in production — emojis, CJK chars, Arabic etc.
// are multi-byte in UTF-8 and will burn you if you index naively.

package main

import (
	"fmt"
	"strings"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func main() {
	// =========================================================================
	// PART 1: STRING FUNDAMENTALS
	// =========================================================================

	// Strings are immutable byte slices
	s := "Hello, 世界" // "World" in Chinese — multi-byte characters

	fmt.Println(len(s))             // 13 — BYTES, not characters!
	fmt.Println(utf8.RuneCountInString(s)) // 9 — actual character count

	// Byte indexing (raw bytes — dangerous with multi-byte chars)
	fmt.Printf("byte[0]=%d (%c)\n", s[0], s[0]) // 72 (H)

	// Safe character iteration — range decodes UTF-8 into runes
	for i, r := range s {
		fmt.Printf("byte_index=%d rune=%c code=%d\n", i, r, r)
	}
	// Note: byte index jumps by 3 for each Chinese character (3 bytes each)

	// String vs []byte vs []rune
	str := "hello"
	bytes := []byte(str)    // mutable byte slice — use for binary/ASCII ops
	runes := []rune(str)    // mutable rune slice — use for Unicode ops

	bytes[0] = 'H'
	fmt.Println(string(bytes)) // Hello

	runes[0] = 'H'
	fmt.Println(string(runes)) // Hello

	// Raw string literals (backtick) — no escape sequences (like JS template literals)
	raw := `first line
second line
path: C:\Users\Go`
	fmt.Println(raw)

	// =========================================================================
	// PART 2: strings PACKAGE — your daily driver
	// =========================================================================

	// --- Searching ---
	fmt.Println(strings.Contains("seafood", "foo"))     // true
	fmt.Println(strings.HasPrefix("seafood", "sea"))    // true
	fmt.Println(strings.HasSuffix("seafood", "food"))   // true
	fmt.Println(strings.Count("cheese", "e"))           // 3
	fmt.Println(strings.Index("chicken", "ken"))        // 4 (-1 if not found)
	fmt.Println(strings.LastIndex("go gopher", "go"))   // 3

	// --- Transformation ---
	fmt.Println(strings.ToUpper("gopher"))              // GOPHER
	fmt.Println(strings.ToLower("GOPHER"))              // gopher
	fmt.Println(strings.Title("hello world"))           // Hello World (deprecated → use cases package)
	fmt.Println(strings.TrimSpace("  hello  "))         // "hello"
	fmt.Println(strings.Trim("***hello***", "*"))       // hello
	fmt.Println(strings.TrimLeft("***hello***", "*"))   // hello***
	fmt.Println(strings.TrimRight("***hello***", "*"))  // ***hello
	fmt.Println(strings.TrimPrefix("foobar", "foo"))    // bar
	fmt.Println(strings.TrimSuffix("foobar", "bar"))    // foo
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", 2))  // moo moo oink
	fmt.Println(strings.ReplaceAll("oink oink", "oink", "moo"))        // moo moo

	// --- Splitting and joining ---
	// JS: "a,b,c".split(",")
	parts := strings.Split("a,b,c", ",")
	fmt.Println(parts) // [a b c]

	// Split with limit (JS: "a,b,c".split(",", 2))
	fmt.Println(strings.SplitN("a,b,c", ",", 2)) // [a b,c]

	// SplitAfter — include the delimiter
	fmt.Println(strings.SplitAfter("a,b,c", ",")) // [a, b, c]

	// Fields — split on any whitespace, filter empty strings
	fmt.Println(strings.Fields("  foo bar  baz   ")) // [foo bar baz]

	// Join — JS: array.join(", ")
	fmt.Println(strings.Join([]string{"a", "b", "c"}, ", ")) // a, b, c

	// Repeat — JS: "ha".repeat(3)
	fmt.Println(strings.Repeat("ha", 3)) // hahaha

	// Cut — split on first occurrence (Go 1.18+) — very handy
	before, after, found := strings.Cut("user@example.com", "@")
	fmt.Println(before, after, found) // user example.com true

	// --- Comparison ---
	fmt.Println(strings.EqualFold("Go", "go"))    // true (case-insensitive)
	fmt.Println(strings.Compare("a", "b"))        // -1

	// --- Rune-based operations (unicode-safe) ---
	fmt.Println(strings.Map(unicode.ToUpper, "hello 世界")) // HELLO 世界
	fmt.Println(strings.IndexRune("hello", 'l'))            // 2
	fmt.Println(strings.ContainsRune("hello", 'z'))         // false
	fmt.Println(strings.ContainsAny("hello", "aeiou"))      // true

	// =========================================================================
	// PART 3: strings.Builder — efficient string concatenation
	// =========================================================================
	// JS: array.join("") is efficient. In Go, naive `s += x` in a loop
	// allocates a new string each time. Use strings.Builder instead.

	// WRONG — allocates on every iteration (O(n²)):
	// result := ""
	// for i := 0; i < 1000; i++ { result += strconv.Itoa(i) }

	// RIGHT — single allocation amortized:
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString(strconv.Itoa(i))
		if i < 4 {
			sb.WriteByte(',')
		}
	}
	fmt.Println(sb.String()) // 0,1,2,3,4
	sb.Reset()               // reuse the builder

	// Builder also accepts runes and bytes
	sb.WriteRune('A')
	sb.WriteByte('B')
	sb.WriteString("CD")
	fmt.Println(sb.String()) // ABCD

	// fmt.Fprintf works with Builder (implements io.Writer)
	fmt.Fprintf(&sb, "-%d", 42)
	fmt.Println(sb.String()) // ABCD-42

	// =========================================================================
	// PART 4: strings.Reader — treat a string as an io.Reader
	// =========================================================================

	r := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	fmt.Printf("read %d bytes: %s\n", n, buf[:n]) // read 5 bytes: hello

	// =========================================================================
	// PART 5: strconv — convert between strings and primitive types
	// =========================================================================
	// JS: parseInt, parseFloat, String(), Number()

	// int → string
	fmt.Println(strconv.Itoa(42))            // "42"
	fmt.Println(strconv.FormatInt(255, 16))  // "ff" (base 16)
	fmt.Println(strconv.FormatInt(8, 2))     // "1000" (base 2)
	fmt.Println(strconv.FormatFloat(3.14159, 'f', 2, 64)) // "3.14"

	// string → int (returns error — always check it)
	n2, err := strconv.Atoi("42")
	if err == nil {
		fmt.Println(n2 + 1) // 43
	}

	_, err = strconv.Atoi("abc")
	fmt.Println(err) // strconv.Atoi: parsing "abc": invalid syntax

	// string → float
	f, err := strconv.ParseFloat("3.14", 64)
	fmt.Println(f, err) // 3.14 <nil>

	// string → bool
	b, _ := strconv.ParseBool("true")
	fmt.Println(b) // true
	b, _ = strconv.ParseBool("1") // "1", "t", "T", "TRUE" all work
	fmt.Println(b) // true

	// bool/float → string
	fmt.Println(strconv.FormatBool(true))       // "true"
	fmt.Println(strconv.FormatFloat(3.14, 'g', -1, 64)) // "3.14"

	// Quote / Unquote — useful for escaping
	fmt.Println(strconv.Quote("Hello, 世界")) // "\"Hello, 世界\""

	// AppendInt — append formatted int to byte slice (avoids allocation)
	buf2 := []byte("val=")
	buf2 = strconv.AppendInt(buf2, 42, 10)
	fmt.Println(string(buf2)) // val=42

	// =========================================================================
	// PART 6: unicode package — character classification
	// =========================================================================

	fmt.Println(unicode.IsLetter('A'))   // true
	fmt.Println(unicode.IsDigit('5'))    // true
	fmt.Println(unicode.IsSpace(' '))    // true
	fmt.Println(unicode.IsUpper('A'))    // true
	fmt.Println(unicode.IsLower('a'))    // true
	fmt.Println(unicode.ToUpper('a'))    // 65 (rune for 'A')
	fmt.Println(string(unicode.ToUpper('a'))) // A

	// =========================================================================
	// EXERCISES
	// =========================================================================

	// EXERCISE 1:
	// Write isPalindrome(s string) bool that works correctly with Unicode.
	// "racecar" → true, "A man a plan a canal Panama" → false (spaces matter)
	// Bonus: handle case-insensitive, ignore spaces: "A man..." → true

	// EXERCISE 2:
	// Write wordCount(s string) map[string]int that counts word frequency.
	// Use strings.Fields and strings.ToLower.
	// "the fox and the dog" → map[and:1 dog:1 fox:1 the:2]

	// EXERCISE 3:
	// Write reverseWords(s string) string — reverse word ORDER (not characters).
	// "hello world foo" → "foo world hello"

	// EXERCISE 4:
	// Write a function that converts a camelCase string to snake_case.
	// "helloWorldFoo" → "hello_world_foo"
	// Hint: use strings.Builder and unicode.IsUpper

	// EXERCISE 5 (LeetCode style):
	// Implement strStr(haystack, needle string) int
	// Return the index of first occurrence of needle in haystack, or -1.
	// Bonus: implement it WITHOUT using strings.Index (KMP or brute force)

	// EXERCISE 6 (Production):
	// Parse a raw HTTP-like header string:
	// "Content-Type: application/json\r\nAuthorization: Bearer abc123\r\n"
	// Return a map[string]string of header name → value.
	// Use strings.Split, strings.Cut, strings.TrimSpace.
}
