// Test file — MUST end in _test.go
// Package name: same as source (white-box) or "package_test" (black-box)
// go test discovers all *_test.go files automatically.
//
// Rules:
// - Test functions: func TestXxx(t *testing.T)
// - Benchmark functions: func BenchmarkXxx(b *testing.B)
// - Example functions: func ExampleXxx() (verified by go test)
// - Helper functions: func helperXxx(t *testing.T) (not a test itself)

package testing_examples

import (
	"fmt"
	"testing"
)

// ============================================================
// 1. SIMPLE TEST — just to understand the structure
// ============================================================

func TestAdd_Simple(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
		// t.Errorf — marks test failed, continues execution
		// t.Fatalf — marks test failed, stops this test function immediately
	}
}

// ============================================================
// 2. TABLE-DRIVEN TESTS — THE Go standard
// This is how tests are written at Google and everywhere serious.
// JS equivalent: test.each() in Jest
// ============================================================

func TestAdd(t *testing.T) {
	// Each row is a test case
	tests := []struct {
		name string // describes what's being tested
		a, b int
		want int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		// t.Run creates a SUBTEST — shows as "TestAdd/positive_numbers" in output
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// ============================================================
// 3. TESTING ERRORS
// ============================================================

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool // whether we expect an error
	}{
		{"normal division", 10, 2, 5, false},
		{"fractional result", 1, 3, 1.0 / 3.0, false},
		{"divide by zero", 5, 0, 0, true},
		{"negative divisor", -10, 2, -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Fatalf("Divide(%v, %v) error = %v; wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}

			// Skip value check if we expected an error
			if tt.wantErr {
				return
			}

			// Float comparison — never use == for floats
			const epsilon = 1e-9
			diff := got - tt.want
			if diff < -epsilon || diff > epsilon {
				t.Errorf("Divide(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// ============================================================
// 4. TABLE-DRIVEN TEST FOR STRING FUNCTIONS
// ============================================================

func TestReverseString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"", ""},            // edge case: empty string
		{"a", "a"},          // edge case: single char
		{"racecar", "racecar"}, // palindrome
		{"hello 世界", "界世 olleh"}, // Unicode — critical to test
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ReverseString(tt.input)
			if got != tt.want {
				t.Errorf("ReverseString(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"racecar", true},
		{"hello", false},
		{"A man a plan a canal Panama", true},
		{"", true},  // empty string is palindrome
		{"a", true}, // single char
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsPalindrome(tt.input)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "1"},
		{3, "Fizz"},
		{5, "Buzz"},
		{15, "FizzBuzz"},
		{30, "FizzBuzz"},
		{9, "Fizz"},
		{10, "Buzz"},
		{7, "7"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FizzBuzz(tt.n)
			if got != tt.want {
				t.Errorf("FizzBuzz(%d) = %q; want %q", tt.n, got, tt.want)
			}
		})
	}
}

// ============================================================
// 5. TESTING A TYPE (Stack)
// ============================================================

func TestStack(t *testing.T) {
	t.Run("new stack is empty", func(t *testing.T) {
		s := NewStack()
		if !s.IsEmpty() {
			t.Error("new stack should be empty")
		}
		if s.Len() != 0 {
			t.Errorf("new stack len = %d; want 0", s.Len())
		}
	})

	t.Run("push and pop", func(t *testing.T) {
		s := NewStack()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		if s.Len() != 3 {
			t.Errorf("len after 3 pushes = %d; want 3", s.Len())
		}

		got, err := s.Pop()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 3 {
			t.Errorf("pop = %d; want 3", got)
		}
	})

	t.Run("pop from empty stack returns error", func(t *testing.T) {
		s := NewStack()
		_, err := s.Pop()
		if err == nil {
			t.Error("expected error when popping empty stack, got nil")
		}
	})

	t.Run("peek does not remove element", func(t *testing.T) {
		s := NewStack()
		s.Push(42)

		got, err := s.Peek()
		if err != nil || got != 42 {
			t.Errorf("Peek() = %d, %v; want 42, nil", got, err)
		}
		if s.Len() != 1 {
			t.Error("Peek should not change the stack length")
		}
	})
}

// ============================================================
// 6. TEST HELPERS — t.Helper() marks the function as a helper
// so failures report the CALLER's line number, not the helper's
// ============================================================

func assertEq(t *testing.T, got, want int) {
	t.Helper() // critical — without this, error points to this line, not the caller
	if got != want {
		t.Errorf("got %d; want %d", got, want)
	}
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddWithHelper(t *testing.T) {
	assertEq(t, Add(2, 3), 5)
	assertEq(t, Add(-1, 1), 0)
}

// ============================================================
// 7. SETUP AND TEARDOWN — TestMain for package-level setup
// (like beforeAll/afterAll in Jest)
// ============================================================

// func TestMain(m *testing.M) {
//     // setup code here (e.g., start a test database)
//     fmt.Println("setting up test suite")
//
//     code := m.Run() // run all tests
//
//     // teardown code here
//     fmt.Println("tearing down test suite")
//     os.Exit(code)
// }

// ============================================================
// 8. EXAMPLE FUNCTIONS — verified as tests, appear in docs
// ============================================================

func ExampleAdd() {
	fmt.Println(Add(2, 3))
	// Output:
	// 5
}

func ExampleFizzBuzz() {
	fmt.Println(FizzBuzz(3))
	fmt.Println(FizzBuzz(5))
	fmt.Println(FizzBuzz(15))
	// Output:
	// Fizz
	// Buzz
	// FizzBuzz
}

// ============================================================
// 9. BENCHMARKS — go test -bench=. -benchmem
// ============================================================

func BenchmarkAdd(b *testing.B) {
	// b.N is determined automatically by the test runner
	for i := 0; i < b.N; i++ {
		Add(2, 3)
	}
}

func BenchmarkReverseString(b *testing.B) {
	s := "hello world, this is a benchmark test string"
	for i := 0; i < b.N; i++ {
		ReverseString(s)
	}
}

// Benchmark comparison: two implementations
func BenchmarkReverseStringShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReverseString("hi")
	}
}

func BenchmarkReverseStringLong(b *testing.B) {
	long := string(make([]byte, 10000))
	for i := 0; i < b.N; i++ {
		ReverseString(long)
	}
}

// ============================================================
// 10. PARALLEL TESTS — t.Parallel() runs tests concurrently
// Useful for tests involving I/O or sleeps
// ============================================================

func TestParallelExample(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"fizz", 3, "Fizz"},
		{"buzz", 5, "Buzz"},
		{"fizzbuzz", 15, "FizzBuzz"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable (pre-Go 1.22 requirement)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // run this subtest concurrently
			got := FizzBuzz(tt.n)
			if got != tt.want {
				t.Errorf("FizzBuzz(%d) = %q; want %q", tt.n, got, tt.want)
			}
		})
	}
}
