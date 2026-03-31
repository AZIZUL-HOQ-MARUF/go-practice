// Topic 11: Generics (Go 1.18+)
// Run: go run 11_generics/main.go
//
// JS/TS: you know generics from TypeScript: function identity<T>(x: T): T
// Go generics are very similar in syntax. They were added in Go 1.18 (2022).
//
// Why generics matter for a migration project:
// Without generics, you'd write separate Stack[int], Stack[string], etc.
// With generics, one implementation works for all types safely.

package main

import (
	"cmp"    // Go 1.21 — Ordered constraint and comparison helpers
	"fmt"
	"slices"
)

// =========================================================================
// PART 1: Generic functions
// =========================================================================

// Without generics: one function per type (pre-1.18 pain)
// func maxInt(a, b int) int { if a > b { return a }; return b }
// func maxFloat(a, b float64) float64 { ... }

// With generics: one function, type-safe for any ordered type
// [T cmp.Ordered] is a constraint — T must support <, >, ==
// cmp.Ordered covers: int, float64, string, and all numeric types
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Multiple type parameters
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](slice []T, pred func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = fn(acc, v)
	}
	return acc
}

// Contains — requires comparable constraint (supports ==)
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// Keys — extract map keys as a slice
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values — extract map values as a slice
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// =========================================================================
// PART 2: Generic types — data structures
// =========================================================================

// Generic Stack — works for any type T
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T // zero value of T
	if len(s.items) == 0 {
		return zero, false
	}
	n := len(s.items) - 1
	v := s.items[n]
	s.items = s.items[:n]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int    { return len(s.items) }
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }

// Generic Queue (FIFO)
type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(v T) {
	q.items = append(q.items, v)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

func (q *Queue[T]) Len() int { return len(q.items) }

// Generic Pair
type Pair[A, B any] struct {
	First  A
	Second B
}

func NewPair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}

// Generic Optional (like Rust's Option<T> or TypeScript's T | null)
type Optional[T any] struct {
	value *T
}

func Some[T any](v T) Optional[T] { return Optional[T]{value: &v} }
func None[T any]() Optional[T]    { return Optional[T]{} }

func (o Optional[T]) IsPresent() bool  { return o.value != nil }
func (o Optional[T]) Get() (T, bool) {
	if o.value == nil {
		var zero T
		return zero, false
	}
	return *o.value, true
}
func (o Optional[T]) OrElse(def T) T {
	if o.value == nil {
		return def
	}
	return *o.value
}

// =========================================================================
// PART 3: Constraints
// =========================================================================

// Custom constraint — type set using interface with ~
// ~ means "any type whose underlying type is int"
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func Abs[T Number](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Custom type satisfies Integer constraint
type MyInt int

func (m MyInt) Double() MyInt { return m * 2 }

// =========================================================================
// PART 4: Type inference — Go infers type params from arguments
// =========================================================================

func Zip[A, B any](as []A, bs []B) []Pair[A, B] {
	n := len(as)
	if len(bs) < n {
		n = len(bs)
	}
	result := make([]Pair[A, B], n)
	for i := 0; i < n; i++ {
		result[i] = NewPair(as[i], bs[i])
	}
	return result
}

func main() {
	// -------------------------------------------------------------------------
	// Generic functions
	// -------------------------------------------------------------------------

	fmt.Println(Max(3, 7))         // 7 — type inferred as int
	fmt.Println(Max(3.14, 2.72))   // 3.14 — float64
	fmt.Println(Max("apple", "banana")) // banana — string

	// Explicit type parameter (needed when compiler can't infer)
	fmt.Println(Min[int](10, 20))  // 10

	// Map, Filter, Reduce — the functional trio
	nums := []int{1, 2, 3, 4, 5}

	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println(doubled) // [2 4 6 8 10]

	// Map to different type
	strs := Map(nums, func(n int) string { return fmt.Sprintf("item%d", n) })
	fmt.Println(strs) // [item1 item2 item3 item4 item5]

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens) // [2 4]

	total := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println(total) // 15

	// Contains
	fmt.Println(Contains([]string{"go", "rust", "ts"}, "go"))  // true
	fmt.Println(Contains([]int{1, 2, 3}, 4))                    // false

	// Keys / Values
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := Keys(m)
	slices.Sort(keys)
	fmt.Println(keys) // [a b c]

	// -------------------------------------------------------------------------
	// Generic Stack
	// -------------------------------------------------------------------------

	var intStack Stack[int]
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	for !intStack.IsEmpty() {
		v, _ := intStack.Pop()
		fmt.Print(v, " ") // 3 2 1 (LIFO)
	}
	fmt.Println()

	// String stack
	var strStack Stack[string]
	strStack.Push("go")
	strStack.Push("is")
	strStack.Push("great")
	top, _ := strStack.Peek()
	fmt.Println("top:", top) // great
	fmt.Println("len:", strStack.Len()) // 3

	// -------------------------------------------------------------------------
	// Generic Queue
	// -------------------------------------------------------------------------

	var q Queue[string]
	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	for q.Len() > 0 {
		v, _ := q.Dequeue()
		fmt.Print(v, " ") // first second third (FIFO)
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// Pair and Optional
	// -------------------------------------------------------------------------

	p := NewPair("Alice", 95)
	fmt.Printf("Name=%s Score=%d\n", p.First, p.Second)

	zipped := Zip([]string{"a", "b", "c"}, []int{1, 2, 3})
	for _, pair := range zipped {
		fmt.Printf("(%s, %d) ", pair.First, pair.Second)
	}
	fmt.Println()

	opt := Some(42)
	fmt.Println(opt.IsPresent())  // true
	fmt.Println(opt.OrElse(0))    // 42

	empty := None[int]()
	fmt.Println(empty.IsPresent()) // false
	fmt.Println(empty.OrElse(-1))  // -1

	// -------------------------------------------------------------------------
	// Number constraint
	// -------------------------------------------------------------------------

	fmt.Println(Sum([]int{1, 2, 3, 4, 5}))          // 15
	fmt.Println(Sum([]float64{1.1, 2.2, 3.3}))       // 6.6
	fmt.Println(Abs(-42))                             // 42
	fmt.Println(Abs(-3.14))                           // 3.14

	// Custom type satisfying constraint
	var x MyInt = 5
	fmt.Println(Sum([]MyInt{x, 10, 15}))             // 30
	fmt.Println(x.Double())                           // 10

	// -------------------------------------------------------------------------
	// Comparison with cmp package (Go 1.21+)
	// -------------------------------------------------------------------------

	fmt.Println(cmp.Compare(3, 5))    // -1
	fmt.Println(cmp.Compare(5, 5))    //  0
	fmt.Println(cmp.Compare(7, 5))    //  1
	fmt.Println(cmp.Compare("a", "b")) // -1

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Write generic Reverse[T any](s []T) []T that returns a reversed copy.
	// Test with []int and []string.

	// EXERCISE 2:
	// Write generic Unique[T comparable](s []T) []T that removes duplicates,
	// preserving order of first occurrence.

	// EXERCISE 3:
	// Implement a generic Set[T comparable] type with:
	// Add(v T), Remove(v T), Contains(v T) bool, Len() int, ToSlice() []T
	// Backed by map[T]struct{}.

	// EXERCISE 4:
	// Write generic ChunkBy[T any](slice []T, size int) [][]T that splits
	// a slice into chunks of given size.
	// ChunkBy([]int{1,2,3,4,5}, 2) → [[1 2] [3 4] [5]]

	// EXERCISE 5 (Production):
	// Write generic Must[T any](v T, err error) T that panics if err != nil,
	// otherwise returns v. (Pattern used in initialization code)
	// val := Must(strconv.Atoi("42")) // returns 42
	// val2 := Must(strconv.Atoi("x")) // panics
}
