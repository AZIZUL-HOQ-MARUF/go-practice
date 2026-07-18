// Topic 13: Standard Library Toolkit
// Run: go run 13_stdlib_toolkit/main.go
//
// The 20% of stdlib you'll use 80% of the time in production Go.
// Focus: sort, math, container/heap, time, os/bufio, log/slog
//
// For a migration project at Google: these packages replace
// utility functions you'd normally reach for lodash/moment/etc.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"log/slog"
	"math"
	"math/rand/v2"
	"os"
	"sort"
	"strings"
	"time"
)

// =========================================================================
// PART 1: sort — beyond the basics
// =========================================================================

type Employee struct {
	Name   string
	Salary int
	Dept   string
}

func sortDemo() {
	fmt.Println("=== SORT ===")

	// sort.Slice — comparator-based, unstable
	employees := []Employee{
		{"Alice", 90000, "Engineering"},
		{"Bob", 85000, "Marketing"},
		{"Carol", 90000, "Engineering"},
		{"Dave", 75000, "HR"},
	}

	// Sort by salary descending, then name ascending (multi-key sort)
	sort.Slice(employees, func(i, j int) bool {
		if employees[i].Salary != employees[j].Salary {
			return employees[i].Salary > employees[j].Salary // desc
		}
		return employees[i].Name < employees[j].Name // asc
	})
	for _, e := range employees {
		fmt.Printf("  %s: %d\n", e.Name, e.Salary)
	}

	// sort.SliceStable — stable (preserves relative order of equal elements)
	sort.SliceStable(employees, func(i, j int) bool {
		return employees[i].Dept < employees[j].Dept
	})

	// sort.Search — binary search on ANY sorted structure (returns insertion point)
	// JS: Array.prototype.findIndex with a sorted array
	sorted := []int{1, 3, 5, 7, 9, 11, 13}
	target := 7
	idx := sort.SearchInts(sorted, target) // specialized for []int
	fmt.Printf("  7 found at index %d\n", idx) // 3

	// sort.Search (generic) — find first index where f(i) is true
	idx2 := sort.Search(len(sorted), func(i int) bool {
		return sorted[i] >= 7
	})
	fmt.Printf("  first >= 7 at index %d, value %d\n", idx2, sorted[idx2])

	// Check if sorted
	fmt.Println("  is sorted:", sort.IntsAreSorted(sorted)) // true

	// Sort strings, floats
	words := []string{"banana", "apple", "cherry", "apricot"}
	sort.Strings(words)
	fmt.Println("  sorted words:", words)

	floats := []float64{3.1, 1.4, 1.5, 9.2, 6.5}
	sort.Float64s(floats)
	fmt.Println("  sorted floats:", floats)
}

// =========================================================================
// PART 2: container/heap — priority queue
// Essential for LeetCode and production scheduling problems.
// Go's heap is a min-heap by default. For max-heap, negate values.
// =========================================================================

// MinHeap of ints (implements heap.Interface)
type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] } // min at top
func (h MinHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)          { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MaxHeap — negate values to reverse ordering
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] } // max at top
func (h MaxHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x any)          { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Priority queue with custom items (common in production)
type Task struct {
	Name     string
	Priority int
}

type TaskQueue []Task

func (pq TaskQueue) Len() int           { return len(pq) }
func (pq TaskQueue) Less(i, j int) bool { return pq[i].Priority > pq[j].Priority } // higher = first
func (pq TaskQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *TaskQueue) Push(x any)         { *pq = append(*pq, x.(Task)) }
func (pq *TaskQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func heapDemo() {
	fmt.Println("=== HEAP ===")

	// Min-heap
	h := &MinHeap{5, 3, 1, 4, 2}
	heap.Init(h)

	heap.Push(h, 0)
	fmt.Print("  min-heap pops: ")
	for h.Len() > 0 {
		fmt.Print(heap.Pop(h).(int), " ") // 0 1 2 3 4 5
	}
	fmt.Println()

	// Max-heap
	mh := &MaxHeap{5, 3, 1, 4, 2}
	heap.Init(mh)
	fmt.Print("  max-heap pops: ")
	for mh.Len() > 0 {
		fmt.Print(heap.Pop(mh).(int), " ") // 5 4 3 2 1
	}
	fmt.Println()

	// Task priority queue
	pq := &TaskQueue{
		{"deploy", 3},
		{"fix-bug", 5},
		{"write-docs", 1},
		{"code-review", 4},
	}
	heap.Init(pq)

	fmt.Println("  tasks by priority:")
	for pq.Len() > 0 {
		task := heap.Pop(pq).(Task)
		fmt.Printf("    [%d] %s\n", task.Priority, task.Name)
	}

	// K largest elements (classic interview problem)
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	k := 3
	kh := &MinHeap{}
	heap.Init(kh)
	for _, n := range nums {
		heap.Push(kh, n)
		if kh.Len() > k {
			heap.Pop(kh) // remove smallest, keeping k largest
		}
	}
	fmt.Printf("  %d largest: %v\n", k, []int(*kh))
}

// =========================================================================
// PART 3: math — constants and functions
// =========================================================================

func mathDemo() {
	fmt.Println("=== MATH ===")

	// Constants
	fmt.Printf("  Pi=%.10f\n", math.Pi)
	fmt.Printf("  E=%.10f\n", math.E)
	fmt.Printf("  MaxInt=%d\n", math.MaxInt)
	fmt.Printf("  MinInt=%d\n", math.MinInt)
	fmt.Printf("  MaxFloat64=%.3e\n", math.MaxFloat64)
	fmt.Printf("  SmallestFloat64=%.3e\n", math.SmallestNonzeroFloat64)

	// Common functions
	fmt.Printf("  Abs(-5)=%v\n", math.Abs(-5))
	fmt.Printf("  Ceil(1.1)=%v\n", math.Ceil(1.1))   // 2
	fmt.Printf("  Floor(1.9)=%v\n", math.Floor(1.9)) // 1
	fmt.Printf("  Round(1.5)=%v\n", math.Round(1.5)) // 2
	fmt.Printf("  Sqrt(16)=%v\n", math.Sqrt(16))     // 4
	fmt.Printf("  Pow(2,10)=%v\n", math.Pow(2, 10))  // 1024
	fmt.Printf("  Log2(1024)=%v\n", math.Log2(1024)) // 10
	fmt.Printf("  Log10(1000)=%v\n", math.Log10(1000)) // 3
	fmt.Printf("  Log(math.E)=%v\n", math.Log(math.E)) // 1 (natural log)
	fmt.Printf("  Max(3,5)=%v\n", math.Max(3, 5))
	fmt.Printf("  Min(3,5)=%v\n", math.Min(3, 5))
	fmt.Printf("  Mod(10,3)=%v\n", math.Mod(10, 3)) // 1 (float modulo)

	// Infinity and NaN — useful as initial values in algorithms
	inf := math.Inf(1)   // positive infinity
	ninf := math.Inf(-1) // negative infinity
	nan := math.NaN()

	fmt.Printf("  Inf=%v, -Inf=%v, NaN=%v\n", inf, ninf, nan)
	fmt.Printf("  IsInf=%v, IsNaN=%v\n", math.IsInf(inf, 1), math.IsNaN(nan))

	// Min/Max trick for algorithm initialization
	minVal := math.MaxInt
	maxVal := math.MinInt
	data := []int{5, 2, 8, 1, 9, 3}
	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	fmt.Printf("  min=%d max=%d\n", minVal, maxVal)

	// Random (Go 1.22+ math/rand/v2)
	fmt.Printf("  random [0,10): %d\n", rand.IntN(10))
	fmt.Printf("  random float [0,1): %.4f\n", rand.Float64())
}

// =========================================================================
// PART 4: time — durations, formatting, comparison
// =========================================================================

func timeDemo() {
	fmt.Println("=== TIME ===")

	now := time.Now()
	fmt.Println("  now:", now.Format(time.RFC3339))

	// Duration arithmetic
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	fmt.Printf("  duration: %v\n", d)  // 2h30m15s
	fmt.Printf("  in seconds: %.0f\n", d.Seconds())
	fmt.Printf("  in minutes: %.0f\n", d.Minutes())

	// Time arithmetic
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)
	fmt.Println("  tomorrow:", tomorrow.Format("2006-01-02"))
	fmt.Println("  yesterday:", yesterday.Format("2006-01-02"))
	// Go's reference time: Mon Jan 2 15:04:05 MST 2006 — memorize this!
	// 2006=year, 01=month, 02=day, 15=hour, 04=minute, 05=second

	// Duration since / until
	start := time.Now()
	time.Sleep(1 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("  elapsed: %v\n", elapsed)

	// Parsing
	t, err := time.Parse("2006-01-02", "2025-01-15")
	if err == nil {
		fmt.Println("  parsed:", t.Format(time.RFC3339))
	}

	// Comparison
	fmt.Println("  tomorrow.After(now):", tomorrow.After(now)) // true
	fmt.Println("  now.Before(tomorrow):", now.Before(tomorrow)) // true

	// Ticker (repeating) and Timer (one-shot) — used in production for polling
	ticker := time.NewTicker(10 * time.Millisecond)
	done := make(chan bool)
	count := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				count++
				if count >= 3 {
					ticker.Stop()
					done <- true
					return
				}
			}
		}
	}()
	<-done
	fmt.Printf("  ticker fired %d times\n", count)
}

// =========================================================================
// PART 5: os and bufio — file I/O and stdin reading
// =========================================================================

func osDemo() {
	fmt.Println("=== OS & BUFIO ===")

	// Command line args (like process.argv in Node.js)
	fmt.Println("  program:", os.Args[0])
	// os.Args[1:] would be user-provided args

	// Environment variables
	path := os.Getenv("PATH")
	fmt.Printf("  PATH starts with: %s...\n", path[:min(40, len(path))])

	// Write to stdout (os.Stdout implements io.Writer)
	fmt.Fprintln(os.Stdout, "  hello from os.Stdout")

	// Write/read files
	tmpFile, err := os.CreateTemp("", "golearn-*.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // cleanup

	// Write
	lines := []string{"line 1", "line 2", "line 3"}
	for _, l := range lines {
		fmt.Fprintln(tmpFile, l)
	}
	tmpFile.Close()

	// Read with bufio.Scanner (line by line — efficient for large files)
	f, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	fmt.Print("  file contents: ")
	for scanner.Scan() {
		fmt.Print(scanner.Text(), " | ")
	}
	fmt.Println()

	// bufio.NewReader — read stdin efficiently (common in competitive programming)
	// reader := bufio.NewReader(os.Stdin)
	// line, _ := reader.ReadString('\n')

	// bufio.Writer — buffer writes for performance
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(writer, "  buffered output\n")
	writer.Flush() // must flush to actually write
}

// =========================================================================
// PART 6: log/slog — structured logging (Go 1.21+)
// Critical for production services. Replace fmt.Println with slog.
// =========================================================================

func loggingDemo() {
	fmt.Println("=== STRUCTURED LOGGING (slog) ===")

	// Default logger (text format)
	slog.Info("server starting", "port", 8080, "env", "production")
	slog.Warn("high memory usage", "used_mb", 512, "limit_mb", 1024)
	slog.Error("database connection failed", "host", "db.example.com", "err", "timeout")

	// JSON logger — standard in production (log aggregators like Stackdriver parse JSON)
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(jsonHandler)

	logger.Info("request received",
		"method", "GET",
		"path", "/api/users",
		"latency_ms", 45,
		"status", 200,
	)

	// With — create child logger with persistent fields
	// Common: add request ID, user ID to all logs in a request handler
	requestLogger := logger.With(
		"request_id", "abc-123",
		"user_id", 42,
	)
	requestLogger.Info("user fetched")
	requestLogger.Info("response sent", "status", 200)

	// Log levels
	// slog.Debug only shows if level is set to Debug
	logger.Debug("detailed debug info", "query", "SELECT *")
	logger.Info("informational")
	logger.Warn("something to watch")
	logger.Error("something went wrong", "err", "connection refused")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	sortDemo()
	fmt.Println()
	heapDemo()
	fmt.Println()
	mathDemo()
	fmt.Println()
	timeDemo()
	fmt.Println()
	osDemo()
	fmt.Println()
	loggingDemo()

	// =========================================================================
	// EXERCISES
	// =========================================================================

	// EXERCISE 1 (sort):
	// Given []string{"banana","Apple","cherry","apricot","BERRY"},
	// sort case-insensitively using sort.Slice + strings.ToLower.

	// EXERCISE 2 (heap — LeetCode style):
	// Find the median from a data stream using two heaps:
	// a max-heap for the lower half, min-heap for the upper half.
	// addNum(num int), findMedian() float64

	// EXERCISE 3 (heap — LeetCode style):
	// Merge K sorted lists:
	// Input: [][]int{{1,4,7},{2,5,8},{3,6,9}}
	// Output: [1 2 3 4 5 6 7 8 9]
	// Use a min-heap.

	// EXERCISE 4 (math):
	// Implement isPrime(n int) bool and
	// sieve(n int) []int (Sieve of Eratosthenes) to find all primes up to n.

	// EXERCISE 5 (time):
	// Write a Stopwatch struct with Start(), Stop(), Elapsed() time.Duration methods.
	// Bonus: Lap() that records split times.

	// EXERCISE 6 (os/bufio):
	// Read a file line by line, count lines/words/bytes (like wc command).
	// func wcFile(path string) (lines, words, bytes int, err error)

	// EXERCISE 7 (time.NewTicker):
	// Write a function runWithTicker(ctx context.Context, interval time.Duration, fn func())
	// that calls fn every interval until ctx is cancelled.
	// Test it: create a context with 350ms timeout, tick every 100ms.
	// Verify fn is called ~3 times before the context cancels.
	// (Import "context" for this exercise.)

	// EXERCISE 8 (slog — structured logging):
	// Create a JSON slog logger. Write a function processOrder(log *slog.Logger, orderID int, userID string)
	// that logs:
	//   - "processing order" at Info level with orderID and userID as fields
	//   - "order complete" at Info level with a duration field (use time.Since)
	//   - if orderID < 0: "invalid order" at Warn level with an error field
	// Call it with valid and invalid inputs and observe the JSON output.
}

// Compile check: ensure imports are used
var _ = strings.Join
