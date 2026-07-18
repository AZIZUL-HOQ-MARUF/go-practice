// Topic 03: Functions — Multi-return, Variadic, Closures, Defer
// Run: go run 03_functions/main.go
//
// JS analogy: functions look similar but Go adds power features —
// multiple return values eliminate the need for exceptions or
// result objects, and defer gives you guaranteed cleanup.

package main

import (
	"fmt"
	"math"
)

// -------------------------------------------------------------------------
// 1. Basic function — same as JS but types are required
// -------------------------------------------------------------------------

// JS: function add(a, b) { return a + b }
func add(a int, b int) int {
	return a + b
}

// When consecutive params share the same type, shorthand:
func multiply(a, b int) int {
	return a * b
}

// -------------------------------------------------------------------------
// 2. Multiple return values — Go's killer feature
// JS has NO equivalent (you'd return an object or use callbacks)
// This is how Go handles errors: return (result, error)
// -------------------------------------------------------------------------

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil // nil = no error (like null but typed)
}

// Multiple non-error returns
func minMax(nums []int) (int, int) {
	min, max := nums[0], nums[0]
	for _, v := range nums[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// -------------------------------------------------------------------------
// 3. Named return values — document intent, enable bare return
// Use sparingly; great for short functions or complex logic
// -------------------------------------------------------------------------

func circleStats(r float64) (area, circumference float64) {
	area = math.Pi * r * r
	circumference = 2 * math.Pi * r
	return // bare return — returns named values
}

// -------------------------------------------------------------------------
// 4. Variadic functions — like JS rest parameters ...args
// JS: function sum(...nums) { return nums.reduce((a,b) => a+b, 0) }
// -------------------------------------------------------------------------

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// -------------------------------------------------------------------------
// 5. Functions as first-class values — same as JS
// -------------------------------------------------------------------------

// Function type (like a TypeScript function type signature)
type MathOp func(int, int) int

func applyOp(a, b int, op MathOp) int {
	return op(a, b)
}

// Returning a function (factory pattern)
func makeMultiplier(factor int) func(int) int {
	// factor is captured in the closure — same as JS closures
	return func(n int) int {
		return n * factor
	}
}

// -------------------------------------------------------------------------
// 6. Closures — capture surrounding variables (same concept as JS)
// -------------------------------------------------------------------------

func makeCounter() func() int {
	count := 0
	return func() int {
		count++ // closes over `count`
		return count
	}
}

// -------------------------------------------------------------------------
// 7. defer — runs when the surrounding function returns
// Think of it like a "finally" that always fires, even on panic.
// Multiple defers stack: LIFO order (last deferred, first executed)
// Most common use: close files, release locks, cleanup resources
// -------------------------------------------------------------------------

func demonstrateDefer() {
	fmt.Println("start")
	defer fmt.Println("deferred 1") // runs 3rd
	defer fmt.Println("deferred 2") // runs 2nd
	defer fmt.Println("deferred 3") // runs 1st (LIFO)
	fmt.Println("end")
	// Output:
	// start
	// end
	// deferred 3
	// deferred 2
	// deferred 1
}

// Practical defer: resource cleanup pattern
// (In real code, this would be a file or DB connection)
func processResource() {
	fmt.Println("opening resource")
	defer fmt.Println("closing resource") // guaranteed to run

	fmt.Println("doing work...")
	// even if panic happens here, defer still runs
	fmt.Println("work done")
}

// -------------------------------------------------------------------------
// 8. init function — runs before main, one per file
// No JS equivalent (closest: module-level code at import time)
// -------------------------------------------------------------------------

// func init() {  ← uncommented this would run before main()
//     fmt.Println("package initialized")
// }

func main() {
	// Basic functions
	fmt.Println(add(3, 4))      // 7
	fmt.Println(multiply(3, 4)) // 12

	// Multiple return values — MUST handle both
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("%.4f\n", result) // 3.3333
	}

	// Error case
	_, err = divide(5, 0)
	if err != nil {
		fmt.Println("Error:", err) // Error: cannot divide by zero
	}

	// Ignoring one return value with _
	min, max := minMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
	fmt.Println("min:", min, "max:", max) // min: 1 max: 9

	// Named returns
	area, circ := circleStats(5)
	fmt.Printf("area=%.2f circumference=%.2f\n", area, circ)

	// Variadic
	fmt.Println(sum(1, 2, 3))       // 6
	fmt.Println(sum(1, 2, 3, 4, 5)) // 15
	fmt.Println(sum())              // 0

	// Spread a slice into variadic — use ... suffix (like JS spread)
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...)) // 60

	// First-class functions
	result2 := applyOp(5, 3, func(a, b int) int { return a - b })
	fmt.Println(result2) // 2

	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Println(double(7), triple(7)) // 14 21

	// Closures
	counter := makeCounter()
	fmt.Println(counter()) // 1
	fmt.Println(counter()) // 2
	fmt.Println(counter()) // 3

	counter2 := makeCounter() // independent counter
	fmt.Println(counter2())   // 1 (fresh count)

	// Defer
	demonstrateDefer()
	processResource()

	// Defer with loop — captures loop variable by value at defer time
	for i := 0; i < 3; i++ {
		i := i // shadow the loop variable! (common Go gotcha)
		defer fmt.Printf("deferred loop i=%d\n", i)
	}
	// Without shadowing, all defers would print the same final value of i
	// This is a classic Go gotcha when deferring in loops

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	fmt.Println("Exercises ----")

	// EXERCISE 1:
	// Write a function `swap(a, b int) (int, int)` that returns the two values swapped.
	// Call it and print both results.
	swap := func(a, b int) (int, int) {
		return b, a
	}
	fmt.Println(swap(5, 10))

	// EXERCISE 2:
	// Write a variadic function `joinStrings(sep string, parts ...string) string`
	// that joins strings with a separator.
	// joinStrings(", ", "Go", "is", "fun") → "Go, is, fun"
	// Hint: use a for loop with range, build result string manually (or use strings.Join)
	joinStrings := func(sep string, parts ...string) string {
		result := ""
		for i, v := range parts {
			result += v
			if i != len(parts)-1 {
				result += sep
			}
		}
		return result
	}
	fmt.Println(joinStrings(", ", "Go", "is", "fun"))

	// EXERCISE 3:
	// Write a function `makeAdder(n int) func(int) int` that returns a closure.
	// add5 := makeAdder(5)
	// fmt.Println(add5(3))  // 8
	// fmt.Println(add5(10)) // 15

	// EXERCISE 4:
	// Write a function that uses defer to print "Done!" at the end regardless
	// of where in the function you return.
	// Simulate an early return condition and prove defer still fires.

	// EXERCISE 5 (Challenge):
	// Write a function `memoize(f func(int) int) func(int) int`
	// that wraps f and caches results. Prove it works with a slow fibonacci.

	// EXERCISE 6:
	// Write a function `incrementOnReturn(x int) (result int)` that uses a deferred
	// anonymous function to increment the named return value `result` by 1 after it
	// has been set.
	// If you set `result = x * 2` in the main body, check what the function returns
	// when called with 5 (should return 11). This tests your understanding of how
	// defer interacts with named return values.

	// EXERCISE 7:
	// Write a function `compose(f, g func(int) int) func(int) int` that returns a
	// closure representing the composition f(g(x)).
	// Test it with two functions: double (x * 2) and addThree (x + 3), and print
	// the result of compose(double, addThree)(5). Expected output: 16.

	// EXERCISE 8:
	// Write a function `filter(nums []int, predicate func(int) bool) []int` that
	// returns a new slice containing only the elements that satisfy the predicate.
	// Test it by filtering a slice of integers to only keep even numbers.

	// EXERCISE 9:
	// Write a recursive function `power(base, exp int) int` that computes base^exp.
	// Then write `flatten(nested [][]int) []int` that flattens a 2D slice into 1D.
	// Test: power(2, 10) → 1024, flatten([][]int{{1,2},{3,4},{5}}) → [1 2 3 4 5]

}
