// Topic 12: Testing
// Run tests: go test ./12_testing/
// Run with output: go test -v ./12_testing/
// Run specific test: go test -v -run TestAdd ./12_testing/
// Run benchmarks: go test -bench=. ./12_testing/
// Coverage: go test -cover ./12_testing/
//
// JS analogy: go test ≈ Jest/Vitest, but built into the language — no install needed.
//
// Google-level testing discipline:
// - Table-driven tests are THE standard (not individual test cases)
// - Tests live in the same package (white-box) or _test package (black-box)
// - Every exported function should have a test
// - Benchmarks live alongside tests

package testing_examples

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================
// CODE UNDER TEST — functions we will test
// ============================================================

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// Divide returns a/b or an error if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// ReverseString reverses a string (rune-safe).
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if s reads the same forwards and backwards.
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

// FizzBuzz returns the FizzBuzz string for n.
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprintf("%d", n)
	}
}

// Stack is a generic-free integer stack for testing demonstration.
type Stack struct {
	items []int
}

func NewStack() *Stack { return &Stack{} }

func (s *Stack) Push(v int) { s.items = append(s.items, v) }

func (s *Stack) Pop() (int, error) {
	if len(s.items) == 0 {
		return 0, errors.New("pop from empty stack")
	}
	n := len(s.items) - 1
	v := s.items[n]
	s.items = s.items[:n]
	return v, nil
}

func (s *Stack) Peek() (int, error) {
	if len(s.items) == 0 {
		return 0, errors.New("peek from empty stack")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) Len() int      { return len(s.items) }
func (s *Stack) IsEmpty() bool { return len(s.items) == 0 }
