// Topic 15: context.Context — Cancellation, Timeouts, Request Scoping
// Run: go run 15_context/main.go
//
// context is one of the most important packages for production Go.
// Google invented it and it's used in EVERY Go service.
//
// JS analogy: AbortController + AbortSignal, but much more pervasive.
// In Go, context flows through the entire call stack of a request.
//
// The rule: if your function does I/O, calls another service, or
// runs for more than a millisecond — it should accept ctx context.Context
// as its FIRST parameter.
//
// func DoSomething(ctx context.Context, other args...) (Result, error)
//
// This is enforced by convention at Google and most serious Go shops.

package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// =========================================================================
// PART 1: Context basics — the four constructors
// =========================================================================

// context.Background() — root context, used at the top of call chains
//   (main, server handlers, top-level goroutines)
// context.TODO()       — placeholder when you're not sure yet
//   (use during refactoring; linters can flag it)
// context.WithCancel  — manual cancellation
// context.WithTimeout — deadline relative to now
// context.WithDeadline — absolute deadline
// context.WithValue   — attach request-scoped values

// =========================================================================
// PART 2: WithCancel — manual cancellation
// =========================================================================

// simulateWork does work and respects cancellation.
// This is the pattern for any long-running or looping function.
func simulateWork(ctx context.Context, name string) error {
	for i := 0; ; i++ {
		// ALWAYS check ctx.Done() in loops and between I/O operations
		select {
		case <-ctx.Done():
			fmt.Printf("  %s: cancelled after %d iterations: %v\n", name, i, ctx.Err())
			return ctx.Err() // return the reason: Canceled or DeadlineExceeded
		default:
			// do actual work
			fmt.Printf("  %s: working... iteration %d\n", name, i)
			time.Sleep(10 * time.Millisecond)
			if i >= 2 {
				// Simulate completing the work
				fmt.Printf("  %s: done!\n", name)
				return nil
			}
		}
	}
}

// =========================================================================
// PART 3: WithTimeout — most common in production
// Used for: DB queries, HTTP calls, RPC calls
// =========================================================================

// fetchData simulates an external service call with timeout.
func fetchData(ctx context.Context, url string) (string, error) {
	// Simulate variable latency
	done := make(chan string, 1)
	go func() {
		time.Sleep(50 * time.Millisecond) // simulated network call
		done <- "data from " + url
	}()

	select {
	case result := <-done:
		return result, nil
	case <-ctx.Done():
		return "", fmt.Errorf("fetchData(%s): %w", url, ctx.Err())
	}
}

// queryDB simulates a DB query that can be cancelled.
func queryDB(ctx context.Context, query string) ([]string, error) {
	// In real code: rows, err := db.QueryContext(ctx, query)
	// The *Context variants of stdlib functions all accept context.

	select {
	case <-time.After(20 * time.Millisecond): // simulate query time
		return []string{"row1", "row2", "row3"}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("queryDB: %w", ctx.Err())
	}
}

// =========================================================================
// PART 4: WithValue — request-scoped values
// Use sparingly. Only for cross-cutting concerns like:
// - Request ID
// - User ID / auth token
// - Trace ID
// DO NOT use for passing business logic parameters — use function args.
// =========================================================================

// Key types for context values — use unexported types to avoid collisions.
// NEVER use string or int as a key — any package could collide.
type contextKey string

const (
	RequestIDKey contextKey = "requestID"
	UserIDKey    contextKey = "userID"
)

// WithRequestID returns a context with the request ID attached.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// RequestID extracts the request ID from context. Returns "" if not set.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string) // type assertion with ok
	return id
}

func WithUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, UserIDKey, id)
}

func UserID(ctx context.Context) int {
	id, _ := ctx.Value(UserIDKey).(int)
	return id
}

// =========================================================================
// PART 5: Propagating context through call chains
// This is the key pattern — every function passes ctx down.
// =========================================================================

// Handler simulates an HTTP request handler (top of call chain).
func Handler(ctx context.Context, reqID string) error {
	// Attach request-scoped data
	ctx = WithRequestID(ctx, reqID)
	ctx = WithUserID(ctx, 42)

	// Add a timeout for the entire request
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel() // ALWAYS defer cancel to release resources

	return processRequest(ctx)
}

func processRequest(ctx context.Context) error {
	reqID := RequestID(ctx)
	userID := UserID(ctx)
	fmt.Printf("  processing request %s for user %d\n", reqID, userID)

	// Query DB — passes ctx so it can be cancelled
	rows, err := queryDB(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("processRequest: %w", err)
	}

	fmt.Printf("  got %d rows\n", len(rows))

	// Call external service — also passes ctx
	data, err := fetchData(ctx, "https://api.example.com/data")
	if err != nil {
		return fmt.Errorf("processRequest: %w", err)
	}

	fmt.Printf("  got data: %s\n", data)
	return nil
}

// =========================================================================
// PART 6: Goroutine lifecycle with context
// =========================================================================

// Worker runs until ctx is done.
func Worker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Channel closed — no more jobs
				fmt.Printf("  worker %d: jobs channel closed\n", id)
				return
			}
			results <- job * job
		case <-ctx.Done():
			fmt.Printf("  worker %d: context cancelled: %v\n", id, ctx.Err())
			return
		}
	}
}

// =========================================================================
// PART 7: context.Cause (Go 1.21+) — richer cancellation reason
// =========================================================================

func main() {
	// -------------------------------------------------------------------------
	// WithCancel
	// -------------------------------------------------------------------------

	fmt.Println("=== WithCancel ===")
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 35ms
	go func() {
		time.Sleep(35 * time.Millisecond)
		cancel()
	}()

	err := simulateWork(ctx, "worker1")
	fmt.Println("  error:", err) // context.Canceled

	// -------------------------------------------------------------------------
	// WithTimeout — fast enough
	// -------------------------------------------------------------------------

	fmt.Println("\n=== WithTimeout (fast) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel2()

	data, err := fetchData(ctx2, "api.example.com")
	if err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Println("  got:", data)
	}

	// -------------------------------------------------------------------------
	// WithTimeout — too slow (deadline exceeded)
	// -------------------------------------------------------------------------

	fmt.Println("\n=== WithTimeout (deadline exceeded) ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel3()

	_, err = fetchData(ctx3, "slow.api.com")
	if err != nil {
		fmt.Println("  error:", err) // context deadline exceeded
		fmt.Println("  is DeadlineExceeded:", errors.Is(err, context.DeadlineExceeded))
	}

	// -------------------------------------------------------------------------
	// WithDeadline — absolute time
	// -------------------------------------------------------------------------

	fmt.Println("\n=== WithDeadline ===")
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx4, cancel4 := context.WithDeadline(context.Background(), deadline)
	defer cancel4()

	dl, hasDeadline := ctx4.Deadline()
	fmt.Printf("  deadline: %v ok=%v\n", dl.Format(time.RFC3339), hasDeadline)
	rows, err := queryDB(ctx4, "SELECT 1")
	fmt.Println("  rows:", rows, "err:", err)

	// -------------------------------------------------------------------------
	// WithValue and propagation
	// -------------------------------------------------------------------------

	fmt.Println("\n=== WithValue and request chain ===")
	rootCtx := context.Background()
	err = Handler(rootCtx, "req-abc-123")
	if err != nil {
		fmt.Println("  handler error:", err)
	}

	// -------------------------------------------------------------------------
	// Goroutine lifecycle with context
	// -------------------------------------------------------------------------

	fmt.Println("\n=== Goroutine lifecycle ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel5()

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Start workers
	for i := 0; i < 3; i++ {
		go Worker(ctx5, i, jobs, results)
	}

	// Send jobs
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	// Collect results
	collected := 0
	for collected < 5 {
		select {
		case r := <-results:
			fmt.Printf("  result: %d\n", r)
			collected++
		case <-ctx5.Done():
			fmt.Println("  timeout collecting results")
			goto done
		}
	}
done:

	// -------------------------------------------------------------------------
	// Context inspection
	// -------------------------------------------------------------------------

	fmt.Println("\n=== Context inspection ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel6()

	deadline2, hasDeadline := ctx6.Deadline()
	fmt.Printf("  has deadline: %v, at: %v\n", hasDeadline, deadline2.Format(time.RFC3339))
	fmt.Printf("  done channel nil: %v\n", ctx6.Done() == nil) // false (has deadline)
	fmt.Printf("  err before cancel: %v\n", ctx6.Err()) // nil

	cancel6()
	fmt.Printf("  err after cancel: %v\n", ctx6.Err()) // context canceled

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Write a function `retry(ctx context.Context, attempts int, fn func(context.Context) error) error`
	// that retries fn up to `attempts` times, stopping immediately if ctx is cancelled.

	// EXERCISE 2:
	// Write a concurrent downloader:
	// func downloadAll(ctx context.Context, urls []string) ([]string, error)
	// Launch one goroutine per URL. If any URL fails OR ctx expires,
	// cancel all remaining downloads and return the error.

	// EXERCISE 3:
	// Implement a simple in-memory cache with TTL using context:
	// type Cache struct { ... }
	// Set(ctx context.Context, key string, value any, ttl time.Duration)
	// Get(ctx context.Context, key string) (any, bool)
	// Each entry should expire after its TTL.

	// EXERCISE 4 (Production pattern):
	// Write an HTTP middleware that:
	// - Generates a request ID (use fmt.Sprintf("req-%d", time.Now().UnixNano()))
	// - Attaches it to the context via WithValue
	// - Adds a 30s timeout to the context
	// - Logs "request started" and "request finished" with the request ID

	// EXERCISE 5 (Challenge):
	// Implement a context-aware pipeline:
	// func Pipeline(ctx context.Context, input []int, stages ...func(context.Context, int) (int, error)) ([]int, error)
	// Each stage processes every item. If ctx is cancelled, stop processing.
	// If any stage returns an error, propagate it.
}
