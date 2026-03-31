// Topic 01: Basics — Variables, Types, Zero Values, Constants
// Run: go run 01_basics/main.go
//
// JS analogy: think of this as replacing let/const/var with a
// typed system where every variable has a default "zero" value
// instead of undefined.

package main

import "fmt"

func main() {
	// -------------------------------------------------------------------------
	// 1. Variable declaration — three styles
	// -------------------------------------------------------------------------

	// Style 1: var with explicit type (like TS: let name: string = "Alice")
	var name string = "Alice"

	// Style 2: var with type inference (Go infers string from the value)
	var age = 30

	// Style 3: short declaration with := (most common inside functions)
	// JS equiv: let city = "Dhaka"
	city := "Dhaka"

	fmt.Println(name, age, city) // Alice 30 Dhaka

	// -------------------------------------------------------------------------
	// 2. Zero values — Go has NO undefined or null by default
	// Every type starts at its zero value when declared without assignment.
	// -------------------------------------------------------------------------

	var i int     // 0
	var f float64 // 0.0
	var b bool    // false
	var s string  // "" (empty string)

	fmt.Printf("int=%d, float64=%f, bool=%t, string=%q\n", i, f, b, s)
	// int=0, float64=0.000000, bool=false, string=""

	// -------------------------------------------------------------------------
	// 3. Basic types cheat sheet
	// -------------------------------------------------------------------------
	// int, int8, int16, int32, int64
	// uint, uint8 (byte), uint16, uint32, uint64
	// float32, float64
	// complex64, complex128
	// bool
	// string
	// rune  (alias for int32, represents a Unicode code point — like JS's char)
	// byte  (alias for uint8)

	var score int64 = 9_000_000 // underscores for readability (Go 1.13+)
	var pi float64 = 3.14159
	var letter rune = 'A' // single quotes = rune (not string)
	var initial byte = 'Z'

	fmt.Println(score, pi, letter, initial) // 9000000 3.14159 65 90

	// -------------------------------------------------------------------------
	// 4. Constants — evaluated at compile time
	// JS: const MAX = 100 (runtime constant, not compile-time)
	// Go: const is truly compile-time
	// -------------------------------------------------------------------------

	const MaxRetries = 3
	const AppName = "GoLearner"
	const Pi = 3.14159265358979

	fmt.Println(MaxRetries, AppName, Pi)

	// iota: auto-incrementing constant generator inside const blocks
	// (like an enum)
	const (
		Sunday    = iota // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
	)
	fmt.Println(Sunday, Monday, Tuesday, Wednesday) // 0 1 2 3

	// Bit-flag pattern with iota
	const (
		Read    = 1 << iota // 1 (001)
		Write               // 2 (010)
		Execute             // 4 (100)
	)
	fmt.Printf("R=%d W=%d X=%d, RW=%d\n", Read, Write, Execute, Read|Write)

	// -------------------------------------------------------------------------
	// 5. Type conversion — explicit only, NO implicit coercion (unlike JS!)
	// JS: 1 + "2" = "12"  ← implicit, dangerous
	// Go: you must convert manually
	// -------------------------------------------------------------------------

	var x int = 42
	var y float64 = float64(x) // explicit cast
	var z int = int(y * 1.5)   // truncates, does NOT round

	fmt.Println(x, y, z) // 42 42 63

	// string <-> int requires strconv package (covered in later topics)
	// This WON'T work: string(65) gives "A" (rune→string), not "65"
	fmt.Println(string(rune(65))) // "A"

	// -------------------------------------------------------------------------
	// 6. fmt verbs — Go's printf
	// -------------------------------------------------------------------------
	// %v  = default format (like JS template literals)
	// %T  = type of the variable
	// %d  = integer decimal
	// %f  = float, %.2f = 2 decimal places
	// %s  = string, %q = quoted string
	// %t  = bool
	// %b  = binary, %x = hex
	// %p  = pointer address

	num := 255
	fmt.Printf("decimal=%d  binary=%b  hex=%x  type=%T\n", num, num, num, num)

	// -------------------------------------------------------------------------
	// 7. Multiple assignment (JS destructuring equivalent)
	// -------------------------------------------------------------------------

	a, b2 := 10, 20
	fmt.Println(a, b2) // 10 20

	// Swap without temp variable (idiomatic Go)
	a, b2 = b2, a
	fmt.Println(a, b2) // 20 10

	// Blank identifier _ discards a value (JS has no direct equivalent)
	_, second := "first", "second"
	fmt.Println(second) // "second"

	// -------------------------------------------------------------------------
	// EXERCISES — complete these:
	// -------------------------------------------------------------------------

	fmt.Println("\n--- EXERCISES ---")

	// EXERCISE 1: (done)
	// Declare a variable `temperature` of type float64 using := with value 36.6
	// Then print it with 1 decimal place using fmt.Printf

	temperature := 36.6
	fmt.Printf("Temperature: %.1f\n", temperature)

	// EXERCISE 2:
	// Create a const block with iota representing card suits:
	// Hearts=0, Diamonds=1, Clubs=2, Spades=3
	// Print all four.

	// EXERCISE 3:
	// Declare an int variable `meters` = 100
	// Convert it to float64, multiply by 3.28084 to get feet, print result
	// Expected: 328.084000 feet (or similar)

	// EXERCISE 4:
	// Using fmt.Printf, print your name, age, and city on one line using %s and %d
	// Output should look like: "Name: Alice, Age: 30, City: Dhaka"
}
