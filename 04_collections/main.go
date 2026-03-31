// Topic 04: Collections — Arrays, Slices, Maps
// Run: go run 04_collections/main.go
//
// JS analogy:
//   Go array  ≈ fixed-size TypedArray (rare in practice)
//   Go slice  ≈ JS Array (dynamic, most common)
//   Go map    ≈ JS Map / plain object {}

package main

import (
	"fmt"
	"sort"
	"slices" // Go 1.21+
)

func main() {
	// =========================================================================
	// PART 1: ARRAYS — fixed size, value type
	// =========================================================================
	// Unlike JS arrays, Go arrays have a fixed length that is part of their TYPE.
	// [3]int and [4]int are different types and are NOT interchangeable.
	// Arrays are passed by VALUE (a copy is made) — unlike JS objects/arrays.
	// In practice, you'll use slices almost always. Arrays are mainly used
	// as the backing storage for slices.

	// Declaration
	var arr [3]int            // [0 0 0] — zero values
	arr2 := [3]int{1, 2, 3}  // literal
	arr3 := [...]int{4, 5, 6} // ... lets compiler count the elements

	fmt.Println(arr, arr2, arr3)  // [0 0 0] [1 2 3] [4 5 6]
	fmt.Println(len(arr2))         // 3
	fmt.Println(arr2[0], arr2[2])  // 1 3

	// Arrays are value types — copying makes a full independent copy
	a := [3]int{1, 2, 3}
	b := a       // full copy
	b[0] = 99
	fmt.Println(a, b) // [1 2 3] [99 2 3] — a is unchanged

	// 2D array
	matrix := [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(matrix[1][2]) // 6

	// =========================================================================
	// PART 2: SLICES — dynamic, reference type (use these in practice)
	// =========================================================================
	// A slice is a view into an underlying array: {pointer, length, capacity}
	// JS Array equivalent — append, index, range over them.

	// --- Creating slices ---

	// Slice literal (most common)
	s := []int{10, 20, 30, 40, 50}

	// make(type, length, capacity) — allocates backing array
	s2 := make([]int, 3)    // [0 0 0], len=3, cap=3
	s3 := make([]int, 3, 5) // [0 0 0], len=3, cap=5 (room to grow without realloc)

	fmt.Println(s, len(s), cap(s))    // [10 20 30 40 50] 5 5
	fmt.Println(s2, len(s2), cap(s2)) // [0 0 0] 3 3
	fmt.Println(s3, len(s3), cap(s3)) // [0 0 0] 3 5

	// nil slice (declared but not initialized) — len and cap are 0
	var nilSlice []int
	fmt.Println(nilSlice == nil, len(nilSlice)) // true 0
	// nil slices are safe to append to:
	nilSlice = append(nilSlice, 1, 2, 3)
	fmt.Println(nilSlice) // [1 2 3]

	// --- Indexing and slicing ---
	// s[low:high] — includes low, excludes high (same as JS .slice(low, high))
	s = []int{10, 20, 30, 40, 50}
	fmt.Println(s[1:3])  // [20 30]  (index 1 and 2)
	fmt.Println(s[:2])   // [10 20]  (from start to index 2)
	fmt.Println(s[3:])   // [40 50]  (from index 3 to end)
	fmt.Println(s[:])    // [10 20 30 40 50] (full slice)

	// --- append — like JS push(), but returns a new slice header
	// IMPORTANT: always reassign the result of append
	s = append(s, 60)           // append one
	s = append(s, 70, 80, 90)   // append multiple
	fmt.Println(s)               // [10 20 30 40 50 60 70 80 90]

	// Spread another slice to append all its elements (like JS spread)
	more := []int{100, 110}
	s = append(s, more...) // ... required
	fmt.Println(s)

	// --- GOTCHA: slices share backing arrays ---
	original := []int{1, 2, 3, 4, 5}
	slice1 := original[1:4] // [2 3 4] — shares memory with original!
	slice1[0] = 99
	fmt.Println(original) // [1 99 3 4 5] — original was modified!

	// Use copy() to get an independent copy
	dst := make([]int, len(original))
	copy(dst, original)
	dst[0] = 0
	fmt.Println(original, dst) // original unchanged

	// copy(dst, src) — copies min(len(dst), len(src)) elements
	small := make([]int, 2)
	copy(small, []int{1, 2, 3, 4}) // only copies 2 elements
	fmt.Println(small)              // [1 2]

	// --- Common slice operations ---

	// Delete element at index i (no built-in delete for slices)
	// Technique: replace with last, then shrink (order doesn't matter)
	sl := []int{1, 2, 3, 4, 5}
	i := 2 // delete index 2
	sl[i] = sl[len(sl)-1]
	sl = sl[:len(sl)-1]
	fmt.Println(sl) // [1 2 5 4] — order changed

	// Order-preserving delete (like JS splice)
	sl2 := []int{1, 2, 3, 4, 5}
	sl2 = append(sl2[:2], sl2[3:]...) // remove index 2
	fmt.Println(sl2)                   // [1 2 4 5]

	// Insert at index (like JS splice)
	sl3 := []int{1, 2, 4, 5}
	sl3 = append(sl3[:2], append([]int{3}, sl3[2:]...)...)
	fmt.Println(sl3) // [1 2 3 4 5]

	// sort.Slice — like JS .sort() but with explicit comparator
	unsorted := []int{5, 3, 1, 4, 2}
	sort.Slice(unsorted, func(i, j int) bool {
		return unsorted[i] < unsorted[j] // ascending
	})
	fmt.Println(unsorted) // [1 2 3 4 5]

	// Sort strings
	words := []string{"banana", "apple", "cherry"}
	sort.Strings(words)
	fmt.Println(words) // [apple banana cherry]

	// Go 1.21+ slices package (cleaner API)
	nums := []int{3, 1, 4, 1, 5, 9}
	slices.Sort(nums)
	fmt.Println(nums) // [1 1 3 4 5 9]

	idx, found := slices.BinarySearch(nums, 4)
	fmt.Println(idx, found) // 3 true

	// Contains check (Go 1.21+)
	fmt.Println(slices.Contains(nums, 5))  // true
	fmt.Println(slices.Contains(nums, 99)) // false

	// =========================================================================
	// PART 3: MAPS — key-value store
	// =========================================================================
	// JS equivalent: Map or plain {}
	// Key MUST be a comparable type (string, int, bool, struct — NOT slice or map)
	// Maps are reference types (like JS objects) — changes are visible through all refs

	// --- Creating maps ---

	// Map literal
	scores := map[string]int{
		"Alice": 95,
		"Bob":   88,
		"Carol": 92,
	}

	// make — empty map
	phoneBook := make(map[string]string)

	// --- CRUD ---

	// Create / Update
	phoneBook["Alice"] = "555-1234"
	phoneBook["Bob"] = "555-5678"
	phoneBook["Alice"] = "555-9999" // update (same as JS object)

	// Read
	fmt.Println(phoneBook["Alice"]) // 555-9999
	fmt.Println(scores["Bob"])      // 88

	// Read with existence check — CRITICAL PATTERN in Go
	// If key doesn't exist, you get the zero value (0 for int, "" for string)
	// To distinguish "missing" from "zero value", use the two-value form:
	val, ok := scores["Dave"]
	if ok {
		fmt.Println("Dave's score:", val)
	} else {
		fmt.Println("Dave not found, got zero value:", val) // 0
	}

	// Delete
	delete(scores, "Bob")
	fmt.Println(scores) // map[Alice:95 Carol:92]

	// Iterate (order NOT guaranteed — Go randomizes map iteration on purpose)
	for name, score := range scores {
		fmt.Printf("%s: %d\n", name, score)
	}

	// --- Map patterns ---

	// Count frequency (like a JS reduce to object)
	text := "hello world"
	freq := make(map[rune]int)
	for _, ch := range text {
		freq[ch]++ // zero value 0 makes this safe without initialization
	}
	fmt.Println(freq)

	// Group by (like JS reduce to Map of arrays)
	people := []string{"Alice", "Bob", "Anna", "Charlie", "Beth"}
	byInitial := make(map[string][]string)
	for _, p := range people {
		initial := string(p[0])
		byInitial[initial] = append(byInitial[initial], p)
	}
	fmt.Println(byInitial) // map[A:[Alice Anna] B:[Bob Beth] C:[Charlie]]

	// Set pattern — use map[T]struct{} (struct{} uses zero bytes)
	set := make(map[string]struct{})
	set["go"] = struct{}{}
	set["python"] = struct{}{}

	if _, exists := set["go"]; exists {
		fmt.Println("go is in the set")
	}

	// Length
	fmt.Println(len(scores))    // 2
	fmt.Println(len(phoneBook)) // 2

	// nil map — reading is safe (returns zero), writing PANICS
	var nilMap map[string]int
	fmt.Println(nilMap["key"]) // 0 (safe)
	// nilMap["key"] = 1        // ← PANIC: assignment to entry in nil map

	// =========================================================================
	// EXERCISES
	// =========================================================================

	// EXERCISE 1 (Slices):
	// Given []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	// a) Find the max value using a loop
	// b) Sort it in descending order
	// c) Print the top 3

	// EXERCISE 2 (Slices):
	// Write a function removeDuplicates(s []int) []int that returns a new
	// slice with duplicates removed (preserve order of first occurrence).
	// Input: [1 2 3 2 1 4 5 4]
	// Output: [1 2 3 4 5]

	// EXERCISE 3 (Maps):
	// Given a string sentence, count the frequency of each word.
	// sentence := "the quick brown fox jumps over the lazy dog the fox"
	// Print each word and its count, sorted alphabetically.

	// EXERCISE 4 (Maps):
	// Write a function `twoSum(nums []int, target int) (int, int)` that
	// returns the indices of two numbers that add to target.
	// Use a map for O(n) solution. (Classic LeetCode #1)

	// EXERCISE 5 (Challenge — Slices + Maps):
	// Given []string{"apple","banana","apple","cherry","banana","apple"}
	// Return a map of word → count, then find the most frequent word.

	_ = matrix // suppress unused warning
	_ = b
	_ = a
}
