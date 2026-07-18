// Topic 05: Structs, Methods, Embedding
// Run: go run 05_structs/main.go
//
// JS analogy: struct ≈ class / object literal in TypeScript
// Go has NO classes, NO inheritance. Instead:
//   - Structs hold data
//   - Methods are functions attached to a struct (via receiver)
//   - Composition via embedding (prefer over inheritance)

package main

import (
	"fmt"
	"math"
)

// -------------------------------------------------------------------------
// 1. Struct definition
// -------------------------------------------------------------------------

// Capital letter = exported (public), lowercase = unexported (private)
type Point struct {
	X float64 // exported field
	Y float64 // exported field
}

type Person struct {
	Name string
	Age  int
	Email string
}

// -------------------------------------------------------------------------
// 2. Methods — functions with a receiver
// JS: class Person { greet() { return `Hello, ${this.name}` } }
// Go: func (p Person) Greet() string { return "Hello, " + p.Name }
// -------------------------------------------------------------------------

// Value receiver — receives a COPY of the struct (reads only)
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s and I'm %d years old.", p.Name, p.Age)
}

// Pointer receiver — receives a pointer (for mutation or large structs)
// Use pointer receiver when you need to modify the struct.
// Rule of thumb: if any method uses a pointer receiver, make them ALL pointer receivers.
func (p *Person) HaveBirthday() {
	p.Age++ // modifies the original
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name:%s, Age:%d}", p.Name, p.Age)
}

// -------------------------------------------------------------------------
// 3. Methods on Point
// -------------------------------------------------------------------------

func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

// -------------------------------------------------------------------------
// 4. Embedding — composition over inheritance
// JS: class Manager extends Employee {...}
// Go: type Manager struct { Employee; ... } (embed, not inherit)
// -------------------------------------------------------------------------

type Employee struct {
	Name   string
	Salary float64
}

func (e Employee) Details() string {
	return fmt.Sprintf("%s earns %.2f", e.Name, e.Salary)
}

func (e *Employee) Raise(amount float64) {
	e.Salary += amount
}

type Manager struct {
	Employee        // embedded — Manager "has" an Employee (not "is")
	Department string
	Reports    []string
}

// Manager can override Employee's method
func (m Manager) Details() string {
	return fmt.Sprintf("%s (Manager of %s) earns %.2f",
		m.Name, m.Department, m.Salary)
}

// -------------------------------------------------------------------------
// 5. Struct tags — metadata on fields (used by JSON, DB libs, etc.)
// -------------------------------------------------------------------------

// import "encoding/json" would use these tags
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // "-" means omit from JSON
	Email    string `json:"email,omitempty"` // omit if empty
}

// -------------------------------------------------------------------------
// 6. Anonymous structs — one-off structures without naming a type
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// 7. Constructor functions (Go convention — no `new` keyword in classes)
// -------------------------------------------------------------------------

func NewPerson(name string, age int) *Person {
	if age < 0 {
		age = 0
	}
	return &Person{Name: name, Age: age}
}

func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// -------------------------------------------------------------------------
// 8. Struct comparison
// -------------------------------------------------------------------------

func main() {
	// --- Creating structs ---

	// Zero value struct
	var p1 Person
	fmt.Println(p1) // {  0}

	// Struct literal (positional — fragile, avoid)
	p2 := Person{"Alice", 30, "alice@example.com"}

	// Struct literal (named fields — preferred)
	p3 := Person{
		Name:  "Bob",
		Age:   25,
		Email: "bob@example.com",
	}

	fmt.Println(p2.Name, p3.Age) // Alice 25

	// Pointer to struct
	p4 := &Person{Name: "Carol", Age: 28}
	// Go auto-dereferences for field access (no need for p4->Name)
	fmt.Println(p4.Name) // Carol (same as (*p4).Name)

	// Constructor
	p5 := NewPerson("Dave", 35)
	fmt.Println(p5.Greet())

	// --- Methods ---
	fmt.Println(p2.Greet())
	p2.HaveBirthday() // pointer receiver — modifies p2
	fmt.Println(p2.Age) // 31

	// Note: Go automatically takes the address when calling pointer receiver methods
	// So `p2.HaveBirthday()` works even though p2 is not a pointer —
	// Go rewrites it as (&p2).HaveBirthday()

	// Methods on Point
	a := NewPoint(0, 0)
	b := NewPoint(3, 4)
	fmt.Printf("distance: %.2f\n", a.Distance(b)) // 5.00

	a.Scale(2)
	fmt.Println(a) // {6 8}... wait, a was (0,0), scaled = (0,0)?
	// Let's use a non-zero point
	c := Point{X: 1, Y: 2}
	c.Scale(3)
	fmt.Println(c) // {3 6}

	// --- Embedding ---
	mgr := Manager{
		Employee: Employee{
			Name:   "Eve",
			Salary: 90000,
		},
		Department: "Engineering",
		Reports:    []string{"Frank", "Grace"},
	}

	// Promoted fields — access Employee fields directly on Manager
	fmt.Println(mgr.Name)       // Eve (from embedded Employee)
	fmt.Println(mgr.Salary)     // 90000

	// Promoted methods
	mgr.Raise(10000)            // calls Employee.Raise via promotion
	fmt.Println(mgr.Salary)     // 100000

	// Override: Manager.Details() is called, not Employee.Details()
	fmt.Println(mgr.Details())  // Eve (Manager of Engineering) earns 100000.00

	// Explicitly call the embedded method
	fmt.Println(mgr.Employee.Details()) // Eve earns 100000.00

	// --- Anonymous struct (useful for JSON, test fixtures, one-off grouping) ---
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}
	fmt.Printf("%s:%d\n", config.Host, config.Port)

	// Slice of anonymous structs — common in tests
	tests := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 4},
		{3, 9},
	}
	for _, tt := range tests {
		result := tt.input * tt.input
		fmt.Printf("square(%d) = %d, pass: %v\n", tt.input, result, result == tt.expected)
	}

	// --- Struct comparison ---
	// Structs are comparable if ALL fields are comparable (no slices/maps)
	pt1 := Point{1.0, 2.0}
	pt2 := Point{1.0, 2.0}
	pt3 := Point{3.0, 4.0}
	fmt.Println(pt1 == pt2) // true
	fmt.Println(pt1 == pt3) // false

	// --- Struct is a VALUE type (important!) ---
	original := Person{Name: "Alice", Age: 30}
	copy2 := original    // full copy
	copy2.Age = 99
	fmt.Println(original.Age, copy2.Age) // 30 99 — independent

	// To share, use pointers:
	ref := &original
	ref.Age = 99
	fmt.Println(original.Age) // 99 — modified through pointer

	_ = p1
	_ = p4
	_ = p5
	_ = p3

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Define a `Rectangle` struct with Width and Height float64.
	// Add methods: Area() float64, Perimeter() float64, IsSquare() bool.
	// Create two rectangles and compare their areas.

	// EXERCISE 2:
	// Define an `Animal` struct (Name, Sound string).
	// Define a `Dog` struct that embeds Animal and adds Breed string.
	// Add a Speak() method on Animal.
	// Create a dog, call Speak() via promotion, then override with a Dog-specific Speak().

	// EXERCISE 3:
	// Build a simple stack data structure using a struct:
	// type Stack struct { items []int }
	// Methods: Push(v int), Pop() (int, bool), Peek() (int, bool), Len() int, IsEmpty() bool
	// Write a main test that pushes 1-5, pops twice, and verifies.

	// EXERCISE 4:
	// Create a `BankAccount` struct with owner string and balance float64 (unexported).
	// Export: Deposit(amount float64), Withdraw(amount float64) error,
	//         Balance() float64, String() string.
	// Enforce: balance cannot go negative, deposits must be positive.

	// EXERCISE 5 (Challenge):
	// Implement a linked list node:
	// type Node struct { Val int; Next *Node }
	// Write functions: NewList(vals ...int) *Node, String(*Node) string,
	//                  Reverse(*Node) *Node

	// EXERCISE 6 (Stringer):
	// Define a `Color` struct with R, G, B uint8 fields.
	// Implement String() string so fmt.Println(c) prints "rgb(R, G, B)".
	// Also implement a Hex() string method that returns "#RRGGBB" (hex format).
	// Test: Color{255, 128, 0} → "rgb(255, 128, 0)" and "#FF8000"

	// EXERCISE 7 (Value copy semantics):
	// Define a Coord struct {X, Y int}.
	// a) Assign one Coord to another and modify the copy — show the original is unchanged.
	// b) Now do the same with *Coord — show that modifying through a pointer DOES change the original.
	// c) Write a function scaleCoord(p Coord, factor int) Coord (value receiver version).
	//    Write scaleCoordInPlace(p *Coord, factor int) (pointer receiver version).
	//    Call both and compare results.

	// EXERCISE 8 (Struct tags):
	// Define a Product struct with fields:
	//   ID    int    (json:"id")
	//   Name  string (json:"name")
	//   Price float64 (json:"price,omitempty")
	//   secret string  (unexported — never appears in JSON)
	// Marshal a Product to JSON (import "encoding/json").
	// Then unmarshal the JSON back and verify secret is empty string.
	// Print the JSON output.
}
