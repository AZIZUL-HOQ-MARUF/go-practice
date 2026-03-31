// Topic 08: Error Handling
// Run: go run 08_errors/main.go
//
// JS: try/catch/throw with exceptions
// Go: errors are VALUES returned alongside results — no exceptions.
//
// The pattern is: func doThing() (Result, error)
//   caller checks: if err != nil { handle it }
//
// This is intentional — Go forces you to think about failures
// at every step, making error handling explicit and visible.

package main

import (
	"errors"
	"fmt"
	"strconv"
)

// -------------------------------------------------------------------------
// 1. The error interface
// error is a built-in interface: type error interface { Error() string }
// nil means "no error"
// -------------------------------------------------------------------------

// Basic error creation
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero") // simple error
	}
	return a / b, nil
}

// fmt.Errorf creates formatted errors (like errors.New but with formatting)
func parseAge(s string) (int, error) {
	age, err := strconv.Atoi(s) // standard library returns errors this way
	if err != nil {
		return 0, fmt.Errorf("parseAge: invalid input %q: %w", s, err) // %w wraps
	}
	if age < 0 || age > 150 {
		return 0, fmt.Errorf("parseAge: age %d is out of range [0, 150]", age)
	}
	return age, nil
}

// -------------------------------------------------------------------------
// 2. Custom error types — carry structured data
// -------------------------------------------------------------------------

type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %q (value=%v): %s",
		e.Field, e.Value, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

func validateUsername(name string) error {
	if len(name) < 3 {
		return &ValidationError{
			Field:   "username",
			Value:   name,
			Message: "must be at least 3 characters",
		}
	}
	if len(name) > 20 {
		return &ValidationError{
			Field:   "username",
			Value:   name,
			Message: "must be at most 20 characters",
		}
	}
	return nil
}

func getUser(id int) (string, error) {
	db := map[int]string{1: "Alice", 2: "Bob"}
	user, ok := db[id]
	if !ok {
		return "", &NotFoundError{Resource: "user", ID: id}
	}
	return user, nil
}

// -------------------------------------------------------------------------
// 3. Error wrapping with %w — chain of context
// Use fmt.Errorf("context: %w", err) to wrap an error.
// The wrapped error can be unwrapped with errors.Is / errors.As
// -------------------------------------------------------------------------

var ErrPermissionDenied = errors.New("permission denied") // sentinel error

func readFile(path string, userRole string) (string, error) {
	if userRole != "admin" {
		return "", fmt.Errorf("readFile(%s): %w", path, ErrPermissionDenied)
	}
	return "file contents", nil
}

// -------------------------------------------------------------------------
// 4. errors.Is — checks if ANY error in the chain matches a target
// Like JS: err instanceof SomeErrorClass, but for wrapped chains
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// 5. errors.As — extracts a specific error type from the chain
// Like JS: if (err instanceof ValidationError) { err.field }
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// 6. Sentinel errors — pre-declared errors for comparison
// -------------------------------------------------------------------------

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrTimeout   = errors.New("timeout")
)

// -------------------------------------------------------------------------
// 7. panic and recover — Go's last resort (not for normal error handling)
// panic: like throwing an unrecoverable exception
// recover: catch a panic inside a deferred function (like a last-resort catch)
// Use cases: truly unrecoverable states, library code protecting callers
// -------------------------------------------------------------------------

func mustDivide(a, b float64) float64 {
	if b == 0 {
		panic("mustDivide: denominator is zero") // panics stop execution
	}
	return a / b
}

// safeDiv wraps mustDivide and recovers from any panic
func safeDiv(a, b float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	result = mustDivide(a, b)
	return
}

// -------------------------------------------------------------------------
// 8. Multi-error / error aggregation (Go 1.20+)
// -------------------------------------------------------------------------

func validateUser(name string, age int) error {
	var errs []error

	if len(name) < 2 {
		errs = append(errs, fmt.Errorf("name too short"))
	}
	if age < 0 {
		errs = append(errs, fmt.Errorf("age cannot be negative"))
	}
	if age > 150 {
		errs = append(errs, fmt.Errorf("age too large"))
	}

	return errors.Join(errs...) // Go 1.20+: join multiple errors
}

func main() {
	// -------------------------------------------------------------------------
	// Basic error handling
	// -------------------------------------------------------------------------

	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10/3 = %.4f\n", result)
	}

	_, err = divide(5, 0)
	if err != nil {
		fmt.Println("Error:", err) // Error: division by zero
	}

	// -------------------------------------------------------------------------
	// Error from standard library
	// -------------------------------------------------------------------------

	age, err := parseAge("25")
	fmt.Println(age, err) // 25 <nil>

	_, err = parseAge("abc")
	fmt.Println(err) // parseAge: invalid input "abc": strconv.Atoi: ...

	_, err = parseAge("200")
	fmt.Println(err) // parseAge: age 200 is out of range [0, 150]

	// -------------------------------------------------------------------------
	// Custom error types
	// -------------------------------------------------------------------------

	err = validateUsername("ab") // too short
	if err != nil {
		fmt.Println(err) // validation failed for field "username" ...
	}

	user, err := getUser(3) // not found
	if err != nil {
		fmt.Println(err) // user with ID 3 not found
	}
	_ = user

	// -------------------------------------------------------------------------
	// errors.Is — works through wrapped error chains
	// -------------------------------------------------------------------------

	_, err = readFile("/etc/passwd", "guest")
	if errors.Is(err, ErrPermissionDenied) {
		fmt.Println("Access denied!") // prints
	}
	fmt.Println("Full error:", err) // readFile(/etc/passwd): permission denied

	// Unwrapping chain manually
	fmt.Println("Unwrapped:", errors.Unwrap(err)) // permission denied

	// -------------------------------------------------------------------------
	// errors.As — extract the concrete type from anywhere in the chain
	// -------------------------------------------------------------------------

	err = validateUsername("x")
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		// Now we have structured access to the error
		fmt.Printf("Field: %s, Value: %v\n", valErr.Field, valErr.Value)
	}

	_, err = getUser(99)
	var nfErr *NotFoundError
	if errors.As(err, &nfErr) {
		fmt.Printf("Resource: %s, ID: %d\n", nfErr.Resource, nfErr.ID)
	}

	// -------------------------------------------------------------------------
	// Sentinel errors
	// -------------------------------------------------------------------------

	// Wrap a sentinel to add context, still detectable with errors.Is
	wrapped := fmt.Errorf("operation failed: %w", ErrNotFound)
	fmt.Println(errors.Is(wrapped, ErrNotFound)) // true

	// -------------------------------------------------------------------------
	// panic and recover
	// -------------------------------------------------------------------------

	r, err := safeDiv(10, 2)
	fmt.Println(r, err) // 5 <nil>

	r, err = safeDiv(10, 0)
	fmt.Println(r, err) // 0 recovered from panic: mustDivide: denominator is zero

	// -------------------------------------------------------------------------
	// Multi-error (Go 1.20+)
	// -------------------------------------------------------------------------

	err = validateUser("", -5)
	if err != nil {
		fmt.Println(err)
		// name too short
		// age cannot be negative
	}

	// -------------------------------------------------------------------------
	// Idiomatic patterns
	// -------------------------------------------------------------------------

	// Pattern 1: early return on error (guard clauses)
	// Instead of deeply nested if-else, return early
	data, processErr := processData("input")
	if processErr != nil {
		fmt.Println("processing failed:", processErr)
		return
	}
	fmt.Println("processed:", data)

	// Pattern 2: error wrapping with context at each layer
	_, err = fetchUserProfile(42)
	fmt.Println(err) // fetchUserProfile: getUser: user with ID 42 not found

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Write a function `sqrt(x float64) (float64, error)` that returns an error
	// if x is negative. Use errors.New.

	// EXERCISE 2:
	// Define a custom `HTTPError` struct with StatusCode int and Message string.
	// Implement the error interface. Write a function that returns different
	// HTTP errors for different inputs. Use errors.As to extract the status code.

	// EXERCISE 3:
	// Write a function `parseConfig(data map[string]string) (Config, error)`
	// where Config has Name string, Port int, Debug bool.
	// Return a ValidationError for any missing or invalid field.

	// EXERCISE 4:
	// Implement retry logic:
	// func Retry(attempts int, fn func() error) error
	// If fn returns nil, stop and return nil.
	// After all attempts, return the last error.
	// Test with a function that fails 3 times then succeeds.

	// EXERCISE 5 (Challenge):
	// Build a mini result type (like Rust's Result<T, E>):
	// type Result[T any] struct { ... }
	// Methods: Ok(val T) Result[T], Err(err error) Result[T],
	//          IsOk() bool, Unwrap() T, UnwrapOr(default T) T
}

// Helper functions for examples
func processData(input string) (string, error) {
	if input == "" {
		return "", errors.New("empty input")
	}
	return "processed: " + input, nil
}

func getUser2(id int) (string, error) {
	db := map[int]string{1: "Alice", 2: "Bob"}
	user, ok := db[id]
	if !ok {
		return "", &NotFoundError{Resource: "user", ID: id}
	}
	return user, nil
}

func fetchUserProfile(userID int) (string, error) {
	user, err := getUser2(userID)
	if err != nil {
		return "", fmt.Errorf("fetchUserProfile: getUser: %w", err) // wrap with context
	}
	return "profile of " + user, nil
}
