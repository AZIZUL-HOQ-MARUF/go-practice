// Topic 02: Control Flow — for, if, switch, range
// Run: go run 02_control_flow/main.go
//
// Big difference from JS: Go has ONLY the `for` keyword.
// No while, no do-while, no forEach built-in — `for` does everything.

package main

import "fmt"

func main() {
	// -------------------------------------------------------------------------
	// 1. for — three forms
	// -------------------------------------------------------------------------

	// Form 1: C-style for (same as JS)
	for i := 0; i < 3; i++ {
		fmt.Print(i, " ") // 0 1 2
	}
	fmt.Println()

	// Form 2: condition-only = while loop in JS
	// JS: while (n < 5) { ... }
	n := 0
	for n < 5 {
		fmt.Print(n, " ") // 0 1 2 3 4
		n++
	}
	fmt.Println()

	// Form 3: infinite loop (JS: while(true) or for(;;))
	count := 0
	for {
		if count == 3 {
			break
		}
		fmt.Print(count, " ") // 0 1 2
		count++
	}
	fmt.Println()

	// continue works just like JS
	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			continue // skip evens
		}
		fmt.Print(i, " ") // 1 3 5
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// 2. range — iterate over slices, maps, strings, channels
	// JS equiv: for...of or forEach
	// -------------------------------------------------------------------------

	fruits := []string{"apple", "banana", "cherry"}

	// range gives index + value (like JS entries())
	for i, v := range fruits {
		fmt.Printf("[%d] %s\n", i, v)
	}

	// Ignore index with _
	for _, fruit := range fruits {
		fmt.Println(fruit)
	}

	// range over map — order is NOT guaranteed (unlike JS Map)
	scores := map[string]int{"Alice": 95, "Bob": 88, "Carol": 92}
	for name, score := range scores {
		fmt.Printf("%s: %d\n", name, score)
	}

	// range over string — gives rune (Unicode code point), not byte index
	// JS: for (const char of "hello") — same idea
	for i, r := range "hello" {
		fmt.Printf("index=%d char=%c\n", i, r)
	}

	// range with only index (value discarded)
	nums := []int{10, 20, 30}
	for i := range nums {
		fmt.Print(nums[i], " ") // 10 20 30
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// 3. if / else
	// -------------------------------------------------------------------------

	x := 15

	if x > 10 {
		fmt.Println("big")
	} else if x > 5 {
		fmt.Println("medium")
	} else {
		fmt.Println("small")
	}

	// Key Go feature: if with init statement
	// The variable declared here is scoped ONLY to the if/else block.
	// Very common with error returns (covered in 08_errors):
	//   if err := someFunc(); err != nil { ... }
	if val := x * 2; val > 20 {
		fmt.Println("doubled is big:", val) // val only accessible here
	}
	// fmt.Println(val) // ← would not compile — val is out of scope

	// -------------------------------------------------------------------------
	// 4. switch — cleaner multi-branch than JS
	// Key difference: NO fallthrough by default (JS falls through without break)
	// -------------------------------------------------------------------------

	day := "Monday"

	switch day {
	case "Saturday", "Sunday": // multiple values in one case
		fmt.Println("Weekend")
	case "Monday":
		fmt.Println("Start of work week")
	case "Friday":
		fmt.Println("Almost weekend!")
	default:
		fmt.Println("Midweek")
	}

	// switch with no condition = cleaner if-else chain
	temp := 72
	switch {
	case temp < 32:
		fmt.Println("Freezing")
	case temp < 60:
		fmt.Println("Cold")
	case temp < 80:
		fmt.Println("Comfortable") // ← this prints
	default:
		fmt.Println("Hot")
	}

	// explicit fallthrough (rare, unlike JS where it's the default)
	switch 2 {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
		fallthrough // explicitly continue to next case
	case 3:
		fmt.Println("three") // also prints
	case 4:
		fmt.Println("four") // does NOT print (fallthrough only goes one level)
	}

	// switch on type (type switch — covered more in 07_interfaces)
	var anything interface{} = 42
	switch v := anything.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %s\n", v)
	default:
		fmt.Printf("unknown type: %T\n", v)
	}

	// -------------------------------------------------------------------------
	// 5. Labeled break/continue (for nested loops — JS has this too)
	// -------------------------------------------------------------------------

outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer // breaks the outer loop entirely
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
	// (0,0) (0,1) (0,2) (1,0)  — stops when i=1,j=1

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Write a for loop that prints the sum of integers 1 to 100.
	// Expected: 5050

	// EXERCISE 2:
	// Given this slice, use range to build a new slice containing only even numbers.
	// input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// expected: [2 4 6 8 10]

	// EXERCISE 3:
	// Write a switch statement that converts a numeric HTTP status code (200, 404, 500)
	// to a human-readable message. Default: "Unknown status".

	// EXERCISE 4:
	// FizzBuzz in Go: for 1..30, print "Fizz" if divisible by 3,
	// "Buzz" if by 5, "FizzBuzz" if by both, else the number.
	// Use for + if/else (then try with switch {})

	// EXERCISE 5:
	// Using a labeled break, find and print the first pair (i, j) where
	// i*j > 20 from a 5x5 grid (i=1..5, j=1..5).

	// EXERCISE 6:
	// Use an `if` with an init statement to open a file and handle the error:
	//   if f, err := os.Open("nonexistent.txt"); err != nil {
	//       fmt.Println("error:", err)
	//   } else {
	//       defer f.Close()
	//       fmt.Println("opened:", f.Name())
	//   }
	// Notice `f` and `err` are scoped to the if/else block only.
	// (Import "os" for this exercise.)

	// EXERCISE 7:
	// Write a switch that maps a rune to its keyboard category.
	// Cases: 'a'-'z' and 'A'-'Z' → "letter", '0'-'9' → "digit",
	// ' ', '\t', '\n' → "whitespace", '+','-','*','/' → "operator", default → "other"
	// Hint: a single case can list multiple values: case 'a', 'e', 'i', 'o', 'u':
	// Test with at least 5 different inputs.

	// EXERCISE 8:
	// Write a type switch function describe(v any) string that returns:
	//   int    → "integer: <value>"
	//   string → "string of length <n>"
	//   bool   → "boolean: <value>"
	//   []int  → "int slice with <n> elements"
	//   nil    → "nil"
	//   default → "unknown type"
	// Call it with 5 different values and print the results.

	// EXERCISE 9:
	// Use `range` over the string "Hello, 世界" and print each character's
	// byte index and rune value using Printf("%d: %c\n", index, r).
	// Then separately loop over it as []byte and show the byte values.
	// Observe: how many iterations does range-over-string give vs range-over-bytes?
}
