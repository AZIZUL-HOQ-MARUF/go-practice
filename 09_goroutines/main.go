// Topic 09: Goroutines & Channels — Concurrency
// Run: go run 09_goroutines/main.go
//
// JS analogy: goroutines ≈ async functions, but much cheaper (not OS threads)
// channels ≈ a typed message queue between async tasks
//
// Go's motto: "Don't communicate by sharing memory;
//              share memory by communicating."
//
// Go can run millions of goroutines. JS has one thread + event loop.
// Go's concurrency is ACTUAL parallelism (multiple CPU cores).
//
// =========================================================================
// THEORY: Concurrency concepts every Go dev must know
// =========================================================================
//
// Concurrency vs Parallelism
//   Concurrency  — structuring a program to handle multiple tasks at once
//                  (doesn't require multiple CPUs; about DESIGN)
//   Parallelism  — executing multiple tasks simultaneously on multiple CPUs
//                  (about EXECUTION)
//   Rob Pike: "Concurrency is about DEALING with lots of things at once.
//              Parallelism is about DOING lots of things at once."
//
// CSP — Communicating Sequential Processes (Tony Hoare, 1978)
//   Go's concurrency model. Instead of sharing memory, independent processes
//   communicate by passing messages through channels.
//   This is why channels exist — they ARE the communication mechanism.
//
// Data Race
//   Two goroutines access the same memory concurrently and at least one writes,
//   with no synchronization. Result is undefined/unpredictable.
//   Detect with: go run -race 09_goroutines/main.go
//
// Race Condition
//   Program correctness depends on timing/ordering of goroutines.
//   Example: read-check-write without a lock — another goroutine can
//   modify the value between your read and write.
//
// Deadlock
//   All goroutines are blocked waiting on each other — program hangs forever.
//   Requires all 4 Coffman conditions simultaneously:
//     1. Mutual exclusion  — only one goroutine holds a resource at a time
//     2. Hold and wait     — goroutine holds resource while waiting for another
//     3. No preemption     — resources can't be forcibly taken away
//     4. Circular wait     — G1 waits on G2, G2 waits on G1
//   Go detects this at runtime: "all goroutines are asleep - deadlock!"
//
// Livelock
//   Goroutines are NOT blocked — they keep running and responding to each
//   other — but no progress is made. Like two people stepping aside for
//   each other in a hallway, indefinitely.
//
// Starvation
//   A goroutine can never acquire the resources it needs because other
//   goroutines (or a poorly tuned scheduler) always win first.
// =========================================================================

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// -------------------------------------------------------------------------
// 1. Goroutines — `go func()` launches a concurrent function
// Goroutines are lightweight (few KB of stack, grows as needed)
// -------------------------------------------------------------------------

func sayHello(name string) {
	fmt.Printf("Hello from %s!\n", name)
}

// -------------------------------------------------------------------------
// 2. sync.WaitGroup — wait for a group of goroutines to finish
// JS equiv: Promise.all([...])
// -------------------------------------------------------------------------

func fetchURL(url string, wg *sync.WaitGroup) {
	defer wg.Done() // ALWAYS use defer to ensure Done() is called
	time.Sleep(10 * time.Millisecond) // simulate work
	fmt.Println("Fetched:", url)
}

// -------------------------------------------------------------------------
// 3. Channels — typed message queues
// Unbuffered channel: sender blocks until receiver is ready (synchronous)
// Buffered channel: sender can send up to N items without blocking
// -------------------------------------------------------------------------

func producer(ch chan<- int, n int) { // chan<- means send-only
	for i := 0; i < n; i++ {
		ch <- i // send to channel
	}
	close(ch) // close signals "no more values"
}

func consumer(ch <-chan int) { // <-chan means receive-only
	for v := range ch { // range on channel reads until closed
		fmt.Print(v, " ")
	}
	fmt.Println()
}

// -------------------------------------------------------------------------
// 4. select — wait on multiple channels (like switch for channels)
// JS equiv: Promise.race() or select in some async libraries
// -------------------------------------------------------------------------

func fibonacci(n int, ch chan<- int, quit <-chan bool) {
	a, b := 0, 1
	for {
		select {
		case ch <- a: // send next fibonacci number
			a, b = b, a+b
		case <-quit: // stop when told to
			fmt.Println("fibonacci stopped")
			return
		}
	}
}

// -------------------------------------------------------------------------
// 5. Mutex — protect shared state from concurrent access
// Use when goroutines share memory instead of communicating via channels
// -------------------------------------------------------------------------

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

// -------------------------------------------------------------------------
// 6. RWMutex — multiple concurrent readers, exclusive writer
// -------------------------------------------------------------------------

type ReadWriteCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

func (c *ReadWriteCache) Get(key string) (string, bool) {
	c.mu.RLock() // multiple goroutines can RLock simultaneously
	defer c.mu.RUnlock()
	v, ok := c.cache[key]
	return v, ok
}

func (c *ReadWriteCache) Set(key, value string) {
	c.mu.Lock() // exclusive lock — no readers while writing
	defer c.mu.Unlock()
	c.cache[key] = value
}

// -------------------------------------------------------------------------
// 7. sync.Once — run something exactly once (singleton init)
// -------------------------------------------------------------------------

var (
	instance *ReadWriteCache
	once     sync.Once
)

func getInstance() *ReadWriteCache {
	once.Do(func() {
		instance = &ReadWriteCache{cache: make(map[string]string)}
		fmt.Println("cache initialized (only once)")
	})
	return instance
}

// -------------------------------------------------------------------------
// 8. atomic — lockless counter for simple numeric operations
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// 9. Common patterns
// -------------------------------------------------------------------------

// Fan-out: one input, multiple workers
func fanOut(jobs <-chan int, numWorkers int) <-chan int {
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- job * job // square the input
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// Pipeline: chain of processing stages
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func main() {
	// -------------------------------------------------------------------------
	// 1. Basic goroutine
	// -------------------------------------------------------------------------

	go sayHello("goroutine 1")
	go sayHello("goroutine 2")
	// Without synchronization, main might exit before goroutines run.
	// This is the first concurrency challenge.

	// -------------------------------------------------------------------------
	// 2. WaitGroup
	// -------------------------------------------------------------------------

	var wg sync.WaitGroup
	urls := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://blog.go.dev",
	}

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg) // pass pointer — each goroutine decrements the same wg
	}
	wg.Wait() // block until all goroutines call wg.Done()
	fmt.Println("All URLs fetched")

	// -------------------------------------------------------------------------
	// 3. Unbuffered channel — synchronous handoff
	// -------------------------------------------------------------------------

	ch := make(chan int) // unbuffered

	go func() {
		ch <- 42 // blocks until someone reads
	}()

	val := <-ch // blocks until someone sends
	fmt.Println("received:", val)

	// -------------------------------------------------------------------------
	// 4. Buffered channel — async up to buffer size
	// -------------------------------------------------------------------------

	buffered := make(chan int, 3) // buffer of 3
	buffered <- 1                  // doesn't block (buffer not full)
	buffered <- 2
	buffered <- 3
	// buffered <- 4               // would block here (buffer full)

	fmt.Println(<-buffered) // 1
	fmt.Println(<-buffered) // 2
	fmt.Println(<-buffered) // 3

	// -------------------------------------------------------------------------
	// 5. range over channel with close
	// -------------------------------------------------------------------------

	numCh := make(chan int, 5)
	go producer(numCh, 5)
	consumer(numCh) // 0 1 2 3 4

	// -------------------------------------------------------------------------
	// 6. select
	// -------------------------------------------------------------------------

	fibCh := make(chan int, 10)
	quitCh := make(chan bool)

	go func() {
		for i := 0; i < 8; i++ {
			fmt.Print(<-fibCh, " ")
		}
		fmt.Println()
		quitCh <- true
	}()

	fibonacci(10, fibCh, quitCh)
	// 0 1 1 2 3 5 8 13

	// select with default (non-blocking receive)
	ch2 := make(chan string, 1)
	select {
	case msg := <-ch2:
		fmt.Println("received:", msg)
	default:
		fmt.Println("no message available (non-blocking)") // this prints
	}

	// Timeout pattern with select
	result := make(chan string, 1)
	go func() {
		time.Sleep(5 * time.Millisecond)
		result <- "done"
	}()

	select {
	case r := <-result:
		fmt.Println("got result:", r)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("timed out")
	}

	// -------------------------------------------------------------------------
	// 7. Mutex for shared state
	// -------------------------------------------------------------------------

	counter := SafeCounter{v: make(map[string]int)}
	var wg2 sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			counter.Inc("hits")
		}()
	}
	wg2.Wait()
	fmt.Println("counter:", counter.Value("hits")) // 100

	// -------------------------------------------------------------------------
	// 8. sync.Once
	// -------------------------------------------------------------------------

	// Even with 10 goroutines, init runs only once
	var wg3 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			getInstance()
		}()
	}
	wg3.Wait()

	// -------------------------------------------------------------------------
	// 9. atomic counter (lockless)
	// -------------------------------------------------------------------------

	var atomicCounter int64
	var wg4 sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}
	wg4.Wait()
	fmt.Println("atomic counter:", atomic.LoadInt64(&atomicCounter)) // 1000

	// -------------------------------------------------------------------------
	// 10. Pipeline pattern
	// -------------------------------------------------------------------------

	// generate(2, 3, 4) → square → double → print
	nums := generate(2, 3, 4)
	squared := square(nums)
	doubled := double(squared)

	for v := range doubled {
		fmt.Print(v, " ") // 8 18 32 = (2²×2, 3²×2, 4²×2)
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// 11. Fan-out
	// -------------------------------------------------------------------------

	jobs := make(chan int, 9)
	for i := 1; i <= 9; i++ {
		jobs <- i
	}
	close(jobs)

	results := fanOut(jobs, 3)
	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println("sum of squares 1..9:", sum) // 285

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Launch 5 goroutines, each sleeping a random duration (1-100ms) then
	// printing "worker N done". Use WaitGroup to wait for all and print
	// "all workers done" at the end.

	// EXERCISE 2:
	// Implement a bounded worker pool:
	// func WorkerPool(jobs <-chan int, numWorkers int) <-chan int
	// Each worker computes job² and sends to results channel.
	// Process 20 jobs with 4 workers, print total.

	// EXERCISE 3:
	// Implement a merge function:
	// func Merge(channels ...<-chan int) <-chan int
	// that reads from all input channels concurrently and sends all values
	// to a single output channel. (Fan-in pattern)

	// EXERCISE 4:
	// Implement a rate limiter that allows at most N operations per second
	// using time.Tick and channels.

	// EXERCISE 5 (Challenge):
	// Implement a concurrent map-reduce:
	// - Map phase: run n goroutines each processing a chunk of []int
	// - Reduce phase: combine the partial sums
	// Input: [1..100], Map: square each, Reduce: sum
	// Expected: 338350

	// EXERCISE 6 (sync.Once — singleton):
	// Implement a Config struct that should only be loaded once (simulating
	// reading from a file). Use sync.Once to ensure loadConfig() is only called
	// once even when called from 10 concurrent goroutines.
	// Print "loading config..." only once, then print the config from all goroutines.

	// EXERCISE 7 (sync.RWMutex — read-heavy cache):
	// Implement a thread-safe in-memory key/value cache:
	//   type Cache struct { mu sync.RWMutex; data map[string]string }
	//   func (c *Cache) Get(key string) (string, bool)
	//   func (c *Cache) Set(key, value string)
	// Launch 10 reader goroutines and 2 writer goroutines simultaneously.
	// Verify no data races with: go run -race 09_goroutines/main.go

	// EXERCISE 8 (Race condition — detect and fix):
	// The following has a data race — identify and fix it with atomic:
	//   counter := 0
	//   var wg sync.WaitGroup
	//   for i := 0; i < 1000; i++ {
	//       wg.Add(1)
	//       go func() { defer wg.Done(); counter++ }()
	//   }
	//   wg.Wait()
	//   fmt.Println(counter)  // often prints less than 1000
	// Run with -race to confirm the race. Fix with atomic.AddInt64 and verify.
}
