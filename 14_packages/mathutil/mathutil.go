// Package mathutil provides mathematical utility functions.
// This demonstrates a well-structured Go package.
package mathutil

import "math"

// Pi is the mathematical constant π.
const Pi = math.Pi

// MaxPrime is the upper bound for prime checks in this package.
const MaxPrime = 1000

// unexported — not accessible from outside this package
const defaultPrecision = 2

// Add returns the sum of a and b.
func Add(a, b int) int {
	return a + b
}

// Multiply returns the product of a and b.
func Multiply(a, b int) int {
	return a * b
}

// IsPrime reports whether n is a prime number.
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Fibonacci returns the nth Fibonacci number (0-indexed).
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Calculator is a stateful calculator. Its internal value is unexported.
type Calculator struct {
	value float64 // unexported — encapsulated
}

// NewCalculator creates a new Calculator starting at 0.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add adds n to the calculator's current value.
func (c *Calculator) Add(n float64) *Calculator {
	c.value += n
	return c // return self for method chaining
}

// Multiply multiplies the calculator's current value by n.
func (c *Calculator) Multiply(n float64) *Calculator {
	c.value *= n
	return c
}

// Result returns the current value.
func (c *Calculator) Result() float64 {
	return c.value
}

// Reset sets the value back to zero.
func (c *Calculator) Reset() *Calculator {
	c.value = 0
	return c
}
