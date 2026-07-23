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
	"slices" // Go 1.21+
	"sort"
	"strings"
	"maps"
)

type StringSet map[string]struct{}

func Add(set StringSet, key string) {
	set[key] = struct{}{}
}

func Remove(set StringSet, key string) {
	delete(set, key)
}

func Contains(set StringSet, key string) bool {
	_, exists := set[key]
	return exists
}

// Union merges two sets into a new set (A ∪ B)
func Union(a, b StringSet) StringSet {
	result := make(StringSet)

	for k := range a {
		result[k] = struct{}{}
	}

	for k := range b {
		result[k] = struct{}{}
	}
	return result
}

func Keys(set StringSet) []string {
	keys := make([]string, 0 , len(set))
	for key := range set {
		keys = append(keys, key)
	}
	return keys
}


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
	arr2 := [3]int{1, 2, 3}   // literal
	arr3 := [...]int{4, 5, 6} // ... lets compiler count the elements

	fmt.Println(arr, arr2, arr3)  // [0 0 0] [1 2 3] [4 5 6]
	fmt.Println(len(arr2))        // 3
	fmt.Println(arr2[0], arr2[2]) // 1 3

	// Arrays are value types — copying makes a full independent copy
	a := [3]int{1, 2, 3}
	b := a // full copy
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
	fmt.Println(s[1:3]) // [20 30]  (index 1 and 2)
	fmt.Println(s[:2])  // [10 20]  (from start to index 2)
	fmt.Println(s[3:])  // [40 50]  (from index 3 to end)
	fmt.Println(s[:])   // [10 20 30 40 50] (full slice)

	// --- append — like JS push(), but returns a new slice header
	// IMPORTANT: always reassign the result of append
	s = append(s, 60)         // append one
	s = append(s, 70, 80, 90) // append multiple
	fmt.Println(s)            // [10 20 30 40 50 60 70 80 90]

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
	fmt.Println(small)             // [1 2]

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
	fmt.Println(sl2)                  // [1 2 4 5]

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

	fmt.Println("Exercise ----- 1")
	// EXERCISE 1 (Slices):
	// Given []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	// a) Find the max value using a loop
	// b) Sort it in descending order
	// c) Print the top 3
	sl1 := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}

	var slMax int // a
	for _, v := range sl1 {
		if v > slMax {
			slMax = v
		}
	}
	fmt.Println(slMax)

	// b
	sort.Slice(sl1, func(a, b int) bool { return sl1[a] > sl1[b] })
	fmt.Println(sl1)

	//c
	fmt.Println(sl1[:3])

	fmt.Println("Exercise ----- 2")
	// EXERCISE 2 (Slices):
	// Write a function removeDuplicates(s []int) []int that returns a new
	// slice with duplicates removed (preserve order of first occurrence).
	// Input: [1 2 3 2 1 4 5 4]
	// Output: [1 2 3 4 5]
	removeDuplicates := func (s []int) []int {
		set := make(map[int]struct{})
		result := make([]int, 0, len(s)) // pre-allocate capacity to make performant

		for _, v := range s {
			if _, exist := set[v]; !exist {
				set[v] = struct{}{}
				result = append(result, v)
			} 
		}

		return result
	}

	fmt.Println(removeDuplicates([]int{1, 2, 3, 2, 1, 4, 5, 4}))

	fmt.Println("Exercise ----- 3")

	// EXERCISE 3 (Maps):
	// Given a string sentence, count the frequency of each word.
	// sentence := "the quick brown fox jumps over the lazy dog the fox"
	// Print each word and its count, sorted alphabetically.
	sentence := "the quick brown fox jumps over the lazy dog the fox"
	wordSl := strings.Fields(sentence)
	wordMap := make(map[string]int)

	for _, word := range wordSl {
		wordMap[word]++
	}

	mapKeys := slices.Collect(maps.Keys(wordMap))
	slices.Sort(mapKeys)

	for _, word := range mapKeys {
		fmt.Printf("%s: %d\n",word, wordMap[word])
	}

	fmt.Println("Exercise ----- 4")
	// EXERCISE 4 (Maps):
	// Write a function `twoSum(nums []int, target int) (int, int)` that
	// returns the indices of two numbers that add to target.
	// Use a map for O(n) solution. (Classic LeetCode #1)
	twoSum := func (nums []int, target int) (int, int) {
		m := make(map[int]int)
		for i, v := range nums {
			if j, ok := m[target - v]; ok {
				return j , i
			}
			m[v] = i
		}
		return -1, -1
	}

	fmt.Println(twoSum([]int {1, 4, 5, 6, 8, 3, 9}, 10))

	fmt.Println("Exercise ----- 5")
	// EXERCISE 5 (Challenge — Slices + Maps):
	// Given []string{"apple","banana","apple","cherry","banana","apple"}
	// Return a map of word → count, then find the most frequent word.
	wordFrequency := func(words []string) (map[string]int, string, int) {
		m := make(map[string]int)
		freqWord := ""
		freqWordCount := 0

		for _, v := range words {
			m[v]++
			if m[v] > freqWordCount {
				freqWord = v
				freqWordCount = m[v]
			}
		}
		return m, freqWord, freqWordCount
	}


	wMap, word, count := wordFrequency([]string{"apple","banana","apple","cherry","banana","apple"})
	fmt.Println(wMap, word, count)


	fmt.Println("Exercise ----- 6")
	// EXERCISE 6 (Slice internals):
	// Demonstrate the shared backing array gotcha:
	//   a := []int{1, 2, 3, 4, 5}
	//   b := a[1:3]   // b shares memory with a
	//   b[0] = 99
	// Print a and b after the modification. Does a change?
	// Then fix it by using copy() to make b independent, and repeat.
	 arr1 := []int{1, 2, 3, 4, 5}
	 ar2 := arr1[1:3]   // b shares memory with a
	 ar2[0] = 99
	 fmt.Println(arr1, ar2)

	 ar1 := []int{1, 2, 3, 4, 5}
	 copied := make([]int, 2, len(ar1))
	 copy(copied, ar1)
	 copied[0] = 55
	 fmt.Println(ar1, copied)


	fmt.Println("Exercise ----- 7")
	// EXERCISE 7 (Set pattern):
	// Implement a string set using map[string]struct{}.
	// Functions: Add(set, val), Remove(set, val), Contains(set, val) bool, Union(a, b) set
	// Use it to find the unique words in:
	//   s1 := "the cat sat on the mat"
	//   s2 := "the dog sat on the log"
	// Print the union of both word sets.

	str1 := "the cat sat on the mat"
	str2 := "the dog sat on the log"

	set1 := make(StringSet)
	set2 := make(StringSet)

	for _, word := range strings.Fields(str1) {
		Add(set1, word)
	}

	fmt.Println("Set1 Keys", Keys(set1))
	
	for _, word := range strings.Fields(str2) {
		Add(set2, word)
	}
	
	fmt.Println("Set2 Keys", Keys(set2))

	unionSet := Union(set1, set2)
	fmt.Println("Union keys", Keys(unionSet))

	fmt.Println("Exercise ----- 8")

	// EXERCISE 8 (make + capacity):
	// Write a function buildSlice(n int) []int that uses make([]int, 0, n) to
	// pre-allocate capacity, then appends n*n, (n-1)*(n-1), ..., 1*1 to it.
	// Print the resulting slice, its length, and its capacity.
	// Explain in a comment why pre-allocating capacity is faster for large n.

	buildSlice := func(n int) []int {
		// Pre-allocating capacity guarantees the underlying backing array 
		// is large enough upfront.
		sl := make([]int, 0, n)
		for ; n >= 1; n-- {
			sl = append(sl, n*n)
		}
		return sl
	}
	result := buildSlice(10)

	// Print slice contents, length, and capacity
	fmt.Println("Slice:", result)
	fmt.Printf("Length: %d, Capacity: %d\n", len(result), cap(result))


	_ = matrix // suppress unused warning
	_ = b
	_ = a
}
