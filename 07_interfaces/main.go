// Topic 07: Interfaces
// Run: go run 07_interfaces/main.go
//
// JS analogy: like TypeScript interfaces, but Go's are IMPLICIT.
// You don't say "class Dog implements Animal" — if a type has the
// required methods, it automatically satisfies the interface.
// This is called "structural typing" or duck typing with compile-time safety.

package main

import (
	"fmt"
	"math"
	"sort"
)

// -------------------------------------------------------------------------
// 1. Defining and implementing interfaces
// -------------------------------------------------------------------------

// Interface: just a set of method signatures
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width, Height float64
}

type Triangle struct {
	A, B, C float64 // side lengths
}

// Circle implements Shape — no "implements" keyword needed
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle implements Shape
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Triangle implements Shape (Heron's formula for area)
func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// Function that accepts ANY Shape — polymorphism via interface
func printShapeInfo(s Shape) {
	fmt.Printf("%T — area=%.2f perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// -------------------------------------------------------------------------
// 2. Interface composition — embed interfaces into larger interfaces
// -------------------------------------------------------------------------

type Stringer interface {
	String() string
}

type Saver interface {
	Save() error
}

// Composed interface
type StringSaver interface {
	Stringer
	Saver
}

// -------------------------------------------------------------------------
// 3. The Stringer interface (fmt.Stringer) — built-in Go convention
// Implement String() string and fmt.Println will use it automatically
// -------------------------------------------------------------------------

type Temperature struct {
	Celsius float64
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius, t.Celsius*9/5+32)
}

// -------------------------------------------------------------------------
// 4. The error interface — just one method
// type error interface { Error() string }
// -------------------------------------------------------------------------

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

// -------------------------------------------------------------------------
// 5. Empty interface — any{} / interface{}
// Equivalent to TypeScript's `any`. Holds a value of any type.
// Use sparingly — prefer concrete types or typed interfaces.
// -------------------------------------------------------------------------

func printAnything(v any) {
	fmt.Printf("value=%v  type=%T\n", v, v)
}

// -------------------------------------------------------------------------
// 6. Type assertion — extract concrete type from interface
// JS: instanceof check / TypeScript type narrowing
// -------------------------------------------------------------------------

func describeShape(s Shape) string {
	// Type assertion: s.(Circle) panics if s is not a Circle
	// Safe form: value, ok := s.(Circle)
	if c, ok := s.(Circle); ok {
		return fmt.Sprintf("Circle with radius %.2f", c.Radius)
	}
	if r, ok := s.(Rectangle); ok {
		return fmt.Sprintf("Rectangle %gx%g", r.Width, r.Height)
	}
	return "Unknown shape"
}

// -------------------------------------------------------------------------
// 7. Type switch — cleaner than chained type assertions
// -------------------------------------------------------------------------

func classify(i any) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case string:
		return fmt.Sprintf("string: %q (len=%d)", v, len(v))
	case bool:
		return fmt.Sprintf("bool: %v", v)
	case []int:
		return fmt.Sprintf("[]int with %d elements", len(v))
	case Shape:
		return fmt.Sprintf("Shape with area=%.2f", v.Area())
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

// -------------------------------------------------------------------------
// 8. Interface with sort.Interface — implementing stdlib interfaces
// sort.Interface requires: Len() int, Less(i,j int) bool, Swap(i,j int)
// -------------------------------------------------------------------------

type Person struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface for []Person sorted by Age
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByName implements sort.Interface for []Person sorted by Name
type ByName []Person

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// -------------------------------------------------------------------------
// 9. Interface nil gotcha
// -------------------------------------------------------------------------

type MyError struct{ msg string }
func (e *MyError) Error() string { return e.msg }

func riskyOp(fail bool) error {
	var err *MyError // typed nil pointer
	if fail {
		err = &MyError{"something went wrong"}
	}
	return err
	// BUG: even when fail=false, this returns a non-nil error interface!
	// The interface holds (type=*MyError, value=nil) which is != nil
	// Fix: return nil directly, not a typed nil pointer
}

func riskyOpFixed(fail bool) error {
	if fail {
		return &MyError{"something went wrong"}
	}
	return nil // return untyped nil — interface is truly nil
}

func main() {
	// --- Polymorphism via interface ---
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{A: 3, B: 4, C: 5},
	}

	for _, s := range shapes {
		printShapeInfo(s)
	}

	fmt.Printf("Total area: %.2f\n", totalArea(shapes))

	// --- fmt.Stringer ---
	t := Temperature{Celsius: 100}
	fmt.Println(t) // uses String() automatically: 100.0°C (212.0°F)

	// --- Type assertion ---
	for _, s := range shapes {
		fmt.Println(describeShape(s))
	}

	// --- Type switch ---
	values := []any{42, "hello", true, []int{1, 2, 3}, Circle{Radius: 3}, nil}
	for _, v := range values {
		fmt.Println(classify(v))
	}

	// --- Empty interface ---
	printAnything(42)
	printAnything("hello")
	printAnything([]int{1, 2, 3})

	// --- sort.Interface ---
	people := []Person{
		{"Charlie", 30},
		{"Alice", 25},
		{"Bob", 35},
		{"Dave", 28},
	}

	sort.Sort(ByAge(people))
	fmt.Println("By age:", people)

	sort.Sort(ByName(people))
	fmt.Println("By name:", people)

	// Go 1.8+: sort.Slice is more convenient (no need to implement interface)
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println("By age (Slice):", people)

	// --- Interface nil gotcha ---
	err := riskyOp(false)
	if err != nil {
		fmt.Println("BUG — got non-nil error even though nothing failed:", err)
		// This prints! The interface wraps a typed nil *MyError
	}

	err2 := riskyOpFixed(false)
	if err2 != nil {
		fmt.Println("Still a bug")
	} else {
		fmt.Println("Correct: no error") // This prints
	}

	// --- Checking what's inside an interface ---
	var s Shape = Circle{Radius: 3}
	c, ok := s.(Circle)
	fmt.Println(c, ok) // {3} true

	_, ok2 := s.(Rectangle)
	fmt.Println(ok2) // false (no panic with two-value form)

	// Panics! Only use one-value assertion when you're 100% certain:
	// c2 := s.(Rectangle) // panic: interface conversion

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Define an Animal interface with Speak() string and Name() string.
	// Implement it for Dog, Cat, and Duck structs.
	// Write a function MakeNoise(animals []Animal) that prints each animal's
	// name and what it says.

	// EXERCISE 2:
	// Implement the io.Writer interface (Write(p []byte) (n int, err error))
	// for a `CountingWriter` struct that counts how many bytes have been written.
	// fmt.Fprintf(cw, "hello %s", "world") should work.

	// EXERCISE 3:
	// Define a `Calculator` interface with Add, Sub, Mul, Div methods.
	// Implement it with a `BasicCalc` struct.
	// Write a RunCalc(c Calculator, a, b float64) function that prints all 4.

	// EXERCISE 4:
	// Implement sort.Interface for a custom type `[]Student` where
	// Student has Name string and GPA float64.
	// Sort descending by GPA, then ascending by Name as tiebreaker.

	// EXERCISE 5 (Challenge):
	// Build a simple plugin system:
	//   type Transformer interface { Transform(s string) string; Name() string }
	// Implement: UpperCase, Reverse, PigLatin transformers.
	// Write Pipeline(input string, transformers ...Transformer) string
	// that applies them in sequence.
}
