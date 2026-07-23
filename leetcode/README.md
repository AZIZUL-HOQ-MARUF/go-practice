# Go Cheatsheet for Problem Solving

Go has no set, no deque, no ordered map, no ternary, no built-in abs for int.
Here's how to fill those gaps and avoid the traps.

---

## Quick Index

- [Set](#set)
- [Map](#map)
- [Slice — Push, Pop, Dequeue](#slice--push-pop-dequeue)
- [Array vs Slice](#array-vs-slice)
- [Sorting](#sorting)
- [String](#string)
- [Integer — Size, Overflow, abs](#integer--size-overflow-abs)
- [No Ternary](#no-ternary--compact-alternatives)
- [Swap](#swap)
- [Checking Membership / Default Value](#checking-membership--default-value)
- [Deque](#deque-go-has-none)
- [Nil / Zero Value Traps](#nil--zero-value-traps)
- [int32 vs int64 — When to Use Which](#int32-vs-int64--when-to-use-which)
- [Pointers in Tree / Linked List Problems](#pointers-in-tree--linked-list-problems)
- [Reference Types in Recursion](#reference-types-in-recursion)
- [Rune ↔ Int Conversions](#rune--int-conversions)
- [Heap / Priority Queue](#heap--priority-queue)
- [Bit Manipulation — A Few Tricks Worth Knowing](#bit-manipulation--a-few-tricks-worth-knowing)

---

## Set

Go has no built-in set. Two idioms:

```go
// map[T]bool — easier to read
seen := map[int]bool{}
seen[3] = true
if seen[3] { }
delete(seen, 3)

// map[T]struct{} — zero bytes per entry, idiomatic Go
seen := map[int]struct{}{}
seen[3] = struct{}{}
_, ok := seen[3]   // ok == true means present
delete(seen, 3)
```

**Gotcha:** `seen[x]` on a `map[int]bool` returns `false` (zero value) for missing keys — so `if seen[x]` silently does the right thing without an existence check. Fine for booleans, dangerous for `map[int]int` (see Map section).

---

## Map

```go
m := map[string]int{}

// ALWAYS use two-value form when existence matters
val, ok := m["key"]
if !ok { /* key is absent */ }

// Zero value trap — this is fine for counting
m["x"]++   // starts from 0 even if "x" was never inserted

// But this is a bug if you mean "only update existing keys"
m["x"] += 10  // silently inserts "x" = 10 if absent
```

**Gotcha: map iteration order is random every run.** If you need sorted output:

```go
keys := make([]string, 0, len(m))
for k := range m { keys = append(keys, k) }
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

**Gotcha: you cannot take the address of a map value.** `&m["key"]` does not compile. Copy it out first.

**Gotcha: writing to a nil map panics.**

```go
var m map[string]int
m["x"] = 1  // panic: assignment to entry in nil map
// Always initialize: m := map[string]int{}
```

---

## Slice — Push, Pop, Dequeue

```go
s := []int{1, 2, 3}

// push (append to end)
s = append(s, 4)

// pop from end  (stack)
top := s[len(s)-1]
s = s[:len(s)-1]

// pop from front  (queue) — O(n), fine for LeetCode
front := s[0]
s = s[1:]

// peek without removing
s[len(s)-1]   // top of stack
s[0]          // front of queue

// remove at index i, order preserved — O(n)
s = append(s[:i], s[i+1:]...)

// remove at index i, order NOT needed — O(1)
s[i] = s[len(s)-1]
s = s[:len(s)-1]
```

**Gotcha: subslice shares the backing array.** Mutating one mutates the other:

```go
a := []int{1, 2, 3, 4}
b := a[1:3]  // b = [2, 3], but backed by a
b[0] = 99   // a is now [1, 99, 3, 4]

// Safe copy:
b := append([]int(nil), a[1:3]...)
```

**Gotcha: `append` may or may not allocate.** If cap is enough it reuses the array; if not, it allocates a new one. You can never assume the old slice variable points to the same memory after an append.

---

## Array vs Slice

Use a **fixed-size array** when the size is a known constant — it lives on the stack, is comparable with `==`, and avoids heap allocation.

```go
// Frequency count for lowercase letters — array beats slice here
var freq [26]int
for _, ch := range word {
    freq[ch-'a']++
}
// Compare two frequency arrays directly
if freq1 == freq2 { /* anagram */ }  // works because arrays are comparable

// Same thing as a slice — NOT comparable with ==
freq := make([]int, 26)
// reflect.DeepEqual(f1, f2) or manual loop needed
```

```go
// Fixed-size DP row — array if n is a small constant
var dp [1001]int

// Dynamic size — must use slice
dp := make([]int, n+1)
```

**Rule of thumb:** if the size is a compile-time constant ≤ a few hundred, prefer an array. Otherwise use `make([]int, n)`.

**Gotcha: arrays are value types.** Passing an array to a function copies the whole thing. Slices pass a header (pointer + len + cap) — cheap.

```go
func bad(a [1000]int) { }    // copies 8000 bytes
func good(a []int) { }       // copies 24 bytes (slice header)
func alsoGood(a *[1000]int) { }  // pointer to array, no copy
```

---

## Sorting

```go
import "sort"

sort.Ints(nums)      // ascending, in-place
sort.Strings(strs)   // ascending, in-place

// Reverse sort
sort.Sort(sort.Reverse(sort.IntSlice(nums)))

// Custom comparator
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

// Multi-key sort
sort.Slice(items, func(i, j int) bool {
    if items[i].Age != items[j].Age {
        return items[i].Age < items[j].Age
    }
    return items[i].Name < items[j].Name
})
```

**Gotcha: `sort.Slice` is NOT stable.** Equal elements may swap positions. Use `sort.SliceStable` when relative order of equals matters.

**Gotcha: sorting a copy does not sort the original.**

```go
sorted := append([]int(nil), nums...)  // copy first
sort.Ints(sorted)
// nums is unchanged
```

**Shortcut — sort a string's characters:**

```go
b := []byte(s)
sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
result := string(b)
```

---

## String

Strings in Go are **immutable byte sequences**. `s[i]` gives a `byte` (0–255), not a character.

```go
s := "hello"
s[0] = 'H'          // compile error — strings are immutable
b := []byte(s)
b[0] = 'H'          // fine
s = string(b)       // convert back
```

**`len(s)` counts bytes, not characters.** For ASCII this is the same; for Unicode it is not.

```go
s := "héllo"
len(s)        // 6 (é is 2 bytes in UTF-8)
len([]rune(s)) // 5 (5 characters)
```

**Iterating characters:**

```go
for i, ch := range s { }  // ch is rune (Unicode codepoint), i is byte index
for i := 0; i < len(s); i++ { b := s[i] }  // b is byte — use for ASCII problems
```

**Gotcha: `+` in a loop is O(n²).** Use `strings.Builder`:

```go
var sb strings.Builder
for _, ch := range chars {
    sb.WriteByte(ch)     // byte
    sb.WriteRune(ch)     // rune
    sb.WriteString(str)  // string
}
result := sb.String()
```

---

## Integer

**No built-in `abs` for int:**

```go
if x < 0 { x = -x }
// or as a function
func abs(x int) int { if x < 0 { return -x }; return x }
```

**Division truncates toward zero** (not floor like Python):

```go
-7 / 2   // -3  (Python gives -4)
-7 % 2   // -1  (sign follows the dividend)

// Safe modulo (always non-negative):
((a % m) + m) % m
```

**`int` is 64-bit on 64-bit platforms** — rarely need `int64` explicitly. But when the problem says values up to 10^18 or you multiply two 10^9 values, be explicit:

```go
var a, b int64 = 1_000_000_000, 1_000_000_000
product := a * b   // fine in int64, overflows int32
```

**Useful constants:**

```go
import "math"

math.MaxInt    // 9223372036854775807
math.MinInt    // -9223372036854775368
math.MaxInt32  // 2147483647
const INF = math.MaxInt / 2  // safe to add to without overflow
```

---

## No Ternary — Compact Alternatives

Go has no `x ? a : b`. Options:

```go
// inline if — most readable
result := a
if condition { result = b }

// for simple min/max (Go 1.21+)
x := min(a, b)
x := max(a, b)

// manual min/max for older Go or custom types
func minInt(a, b int) int { if a < b { return a }; return b }
```

---

## Swap

No temp variable needed:

```go
a, b = b, a
nums[i], nums[j] = nums[j], nums[i]
```

---

## Checking Membership / Default Value

```go
// map — existence check is separate from zero value
count := map[int]int{}
count[x]++   // starts at 0 automatically — this is fine
v := count[x]  // returns 0 if x absent, no panic

// slice — no contains built-in before Go 1.21
// Go 1.21+:
import "slices"
slices.Contains(nums, target)
slices.Index(nums, target)  // returns -1 if absent
```

---

## Deque (Go has none)

For problems needing push/pop from both ends (sliding window max, palindrome checks):

```go
dq := []int{}

dq = append(dq, x)         // push back
dq = append([]int{x}, dq...) // push front — O(n), avoid in hot loops

back := dq[len(dq)-1]; dq = dq[:len(dq)-1]  // pop back
front := dq[0]; dq = dq[1:]                  // pop front — O(n)
```

For O(1) both ends, use a circular buffer or `container/ring`. For most LeetCode sliding-window problems the O(n) front-pop is fast enough since total elements popped = n.

---

## Nil / Zero Value Traps

```go
var s []int      // nil slice — len=0, cap=0
s := []int{}     // empty slice — len=0, cap=0

append(s, 1)     // works on nil slice, returns new slice
len(s)           // 0 for both

s == nil         // true for var s []int, FALSE for s := []int{}
// Prefer len(s) == 0 to check emptiness — works for both
```

```go
var m map[int]int   // nil map
m[1]                // read is fine — returns 0
m[1] = 1            // PANIC — cannot write to nil map
```

```go
var p *TreeNode     // nil pointer
p.Val               // PANIC — nil dereference
if p != nil { p.Val }
```

---

## int32 vs int64 — When to Use Which

In Go, **`int` is 64-bit on all modern platforms** (LeetCode runs 64-bit). You almost never need `int32` or `int64` explicitly — but you do need to think about overflow when multiplying.

**Quick reference by constraint:**

| Problem says max N is... | Safe for intermediate products? |
|---|---|
| N ≤ 10³ | `int` always fine |
| N ≤ 10⁴ | `int` fine, N² = 10⁸ fits int32 too |
| N ≤ 10⁵ | N² = 10¹⁰ — overflows int32, fine in `int` (64-bit) |
| N ≤ 10⁹ | N² = 10¹⁸ — fits int64, overflows if you go higher |
| N ≤ 10¹⁸ | N*2 may overflow — use modular arithmetic |

**Rule:** Go's `int` handles everything up to ~9.2 × 10¹⁸. The only time you need `int64` explicitly is when interfacing with APIs or problem structs that use `int32`, or when doing modular arithmetic where an intermediate product can exceed 9.2 × 10¹⁸.

```go
// This is the dangerous pattern — int32 cast loses information
var x int32 = 1_000_000
var y int32 = 1_000_000
product := x * y  // OVERFLOW: 10^12 doesn't fit int32 (max ~2.1 × 10^9)

// Fix: promote before multiplying
product := int64(x) * int64(y)

// In Go LeetCode, TreeNode.Val and ListNode.Val are plain `int` (64-bit)
// so this is rarely an issue unless you cast yourself
```

**Modular arithmetic — when values grow huge:**

```go
const MOD = 1_000_000_007

// Safe: intermediate result fits int64 (both operands < MOD < 10^9, product < 10^18)
result = result * base % MOD

// Unsafe if result or base can be >= sqrt(MaxInt64) ≈ 3 × 10^9
// Use int64 explicitly or reduce modulo at every step
```

---

## Pointers in Tree / Linked List Problems

### Reassigning a pointer variable doesn't affect the caller

```go
// WRONG — head inside the function is a copy of the pointer
func deleteHead(head *ListNode) {
    head = head.Next  // only changes the local variable
}

// RIGHT — return the new head
func deleteHead(head *ListNode) *ListNode {
    return head.Next
}
```

### Modifying a node's field IS visible to the caller

```go
// This works — you're dereferencing the pointer and changing the struct
func doubleVal(node *ListNode) {
    node.Val *= 2  // caller sees the change
}
```

### Dummy node — eliminates nil-head edge cases

```go
dummy := &ListNode{Next: head}
cur := dummy
// ... build or modify the list using cur
return dummy.Next  // actual head (handles case where original head was removed)
```

### Slow / fast pointer — cycle detection, midpoint

```go
slow, fast := head, head
for fast != nil && fast.Next != nil {
    slow = slow.Next
    fast = fast.Next.Next
}
// when loop ends: slow is at the midpoint (or start of cycle)
```

**Gotcha: `fast.Next.Next` panics if `fast.Next == nil`.** Always check `fast != nil && fast.Next != nil` in that order (short-circuit saves you).

### Tree — passing a pointer-to-pointer when you need to rewire

```go
// If you need to delete/replace a node from inside a recursive call,
// return the node instead of trying to modify the parent's pointer field.

func deleteNode(root *TreeNode, key int) *TreeNode {
    if root == nil { return nil }
    if key < root.Val {
        root.Left = deleteNode(root.Left, key)
    } else if key > root.Val {
        root.Right = deleteNode(root.Right, key)
    } else {
        if root.Left == nil { return root.Right }
        if root.Right == nil { return root.Left }
        // find inorder successor, etc.
    }
    return root
}
// Pattern: always return the (possibly updated) node — caller re-wires via assignment
```

### Two-pointer on linked list — don't forget to cut the connection

```go
// Finding the node before the midpoint and splitting:
slow, fast := head, head.Next
for fast != nil && fast.Next != nil {
    slow = slow.Next
    fast = fast.Next.Next
}
second := slow.Next
slow.Next = nil  // CUT — forgetting this causes infinite loops in merge/sort
```

---

## Reference Types in Recursion

**Slices, maps, and channels are already reference types** — mutations inside a recursive call are visible outside. But `append` can silently break this.

```go
// BUG — append may allocate a new backing array; caller's slice is stale
func collect(node *TreeNode, result []int) {
    if node == nil { return }
    result = append(result, node.Val)  // local reassignment, caller never sees new elements
    collect(node.Left, result)
    collect(node.Right, result)
}

// FIX 1 — return the slice (cleanest)
func collect(node *TreeNode, result []int) []int {
    if node == nil { return result }
    result = append(result, node.Val)
    result = collect(node.Left, result)
    result = collect(node.Right, result)
    return result
}

// FIX 2 — pass a pointer to the slice
func collect(node *TreeNode, result *[]int) {
    if node == nil { return }
    *result = append(*result, node.Val)
    collect(node.Left, result)
    collect(node.Right, result)
}

// FIX 3 — close over an outer variable (most common in practice)
func inorder(root *TreeNode) []int {
    result := []int{}
    var dfs func(*TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil { return }
        dfs(node.Left)
        result = append(result, node.Val)  // captures result from outer scope
        dfs(node.Right)
    }
    dfs(root)
    return result
}
```

**Maps don't have the append problem** — they're always reference types and mutations are always visible:

```go
func countFreq(node *TreeNode, freq map[int]int) {
    if node == nil { return }
    freq[node.Val]++  // visible to caller — no pointer tricks needed
    countFreq(node.Left, freq)
    countFreq(node.Right, freq)
}
```

**Gotcha: integer/bool/struct parameters are always copied.** If you need to accumulate a count or sum across recursive calls, either return it, pass a pointer, or close over an outer variable.

```go
// WRONG — count is a local copy, changes don't propagate
func countNodes(node *TreeNode, count int) {
    if node == nil { return }
    count++
    countNodes(node.Left, count)
    countNodes(node.Right, count)
}

// RIGHT — close over
count := 0
var dfs func(*TreeNode)
dfs = func(node *TreeNode) {
    if node == nil { return }
    count++
    dfs(node.Left)
    dfs(node.Right)
}
dfs(root)
```

---

## Rune ↔ Int Conversions

```go
// char → index (0-based)
idx := ch - 'a'   // 'a'=0, 'b'=1, ..., 'z'=25
idx := ch - 'A'   // 'A'=0, 'B'=1, ..., 'Z'=25
digit := ch - '0' // '0'=0, '1'=1, ..., '9'=9

// index → char (rune)
ch := rune('a' + idx)   // 0→'a', 1→'b', ...
ch := rune('0' + digit) // 0→'0', 1→'1', ...
ch := byte('a' + idx)   // same but as byte — use in []byte

// rune ↔ string
string(rune(65))   // "A"
string('a')        // "a"
[]rune("hello")    // [104 101 108 108 111]

// ASCII case tricks (only valid for a-z / A-Z)
lower := ch | 32      // 'A' → 'a'  (sets bit 5)
upper := ch &^ 32     // 'a' → 'A'  (clears bit 5)
toggle := ch ^ 32     // 'a' ↔ 'A'

// Safe version using unicode package (handles non-ASCII too)
import "unicode"
unicode.ToLower(ch)
unicode.ToUpper(ch)
unicode.IsLetter(ch)
unicode.IsDigit(ch)
unicode.IsSpace(ch)

// strconv for string ↔ number
import "strconv"
n, _ := strconv.Atoi("42")      // "42" → 42
s := strconv.Itoa(42)           // 42 → "42"
n, _ := strconv.ParseInt("ff", 16, 64)  // hex string → int64
s := strconv.FormatInt(255, 2)  // int64 → binary string "11111111"
```

**Gotcha: `string(65)` gives `"A"`, not `"65"`.** Use `strconv.Itoa(65)` for the digit string.

**Gotcha: iterating a string with `s[i]` gives a `byte`.** For multi-byte Unicode chars, use `range`:

```go
for i, ch := range s {  // ch is rune
    idx := ch - 'a'     // works correctly even with range
}

// If you only have ASCII (most LeetCode problems), index directly:
for i := 0; i < len(s); i++ {
    idx := s[i] - 'a'  // s[i] is byte, same result for ASCII
}
```

---

## Heap / Priority Queue

Go's `container/heap` requires implementing 5 methods. Here's the minimal copy-paste for the two most common cases.

### Min-Heap of ints

```go
import "container/heap"

type MinHeap []int
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any          { a := *h; x := a[len(a)-1]; *h = a[:len(a)-1]; return x }

h := &MinHeap{3, 1, 4}
heap.Init(h)
heap.Push(h, 2)
top := heap.Pop(h).(int)  // 1
peek := (*h)[0]           // peek without removing
```

**Max-Heap:** flip the `Less` comparison: `return h[i] > h[j]`

### Heap of structs (e.g., `[value, index]` pairs)

```go
type Item struct{ val, idx int }
type ItemHeap []Item
func (h ItemHeap) Len() int           { return len(h) }
func (h ItemHeap) Less(i, j int) bool { return h[i].val < h[j].val }  // min by val
func (h ItemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ItemHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *ItemHeap) Pop() any          { a := *h; x := a[len(a)-1]; *h = a[:len(a)-1]; return x }
```

**Gotcha: `heap.Pop` returns `any`, you must type-assert.** `heap.Pop(h).(int)` — forgetting the assertion is a compile error.

**Gotcha: after directly modifying a heap element, call `heap.Fix(h, i)` to restore the invariant.** If you just assign `(*h)[i] = newVal` without calling Fix, the heap is silently broken.

---

### Shortcut — sorted slice instead of heap (insert + retrieve)

If you need to **insert values and retrieve the min/max** but don't want to write the heap interface, maintain a sorted slice with a binary-search insert. Same logical result, zero boilerplate.

```go
sorted := []int{}

// INSERT — find position, shift right, place value — O(n)
i := sort.SearchInts(sorted, val)
sorted = append(sorted, 0)
copy(sorted[i+1:], sorted[i:])
sorted[i] = val

// RETRIEVE min / max — O(1), no Pop needed
min := sorted[0]
max := sorted[len(sorted)-1]

// RETRIEVE AND REMOVE min (like heap.Pop on a min-heap)
min, sorted = sorted[0], sorted[1:]

// RETRIEVE AND REMOVE max
max, sorted = sorted[len(sorted)-1], sorted[:len(sorted)-1]
```

**When to use which:**

| | Heap | Sorted slice |
|---|---|---|
| Insert | O(log n) | O(n) — shifts elements |
| Retrieve min/max | O(log n) | O(1) |
| Boilerplate | ~10 lines | 3 lines |
| Use when | n > ~1000, or streaming inserts matter | n ≤ ~1000, or you want fast code-writing |

> For most LeetCode medium problems the input is small enough that the O(n) insert never shows up in the time limit. Write the sorted slice, move on. Reach for the heap only when n is large or the problem is explicitly about priority ordering.

---

## Bit Manipulation — A Few Tricks Worth Knowing

Not a full list. These are the ones that actually show up.

### `n & (n-1)` — strip the lowest set bit

Subtracting 1 flips the lowest set bit and sets everything below it. ANDing with the original clears exactly that bit.

```
n     = 1100
n-1   = 1011
n&(n-1) = 1000   ← lowest set bit gone
```

```go
// Is n a power of 2? (exactly one bit set)
n > 0 && n&(n-1) == 0

// Count set bits manually (Brian Kernighan) — runs once per set bit, not per bit
count := 0
for n > 0 {
    n &= n - 1
    count++
}
// In practice just use: bits.OnesCount(uint(n))
```

### `n & (-n)` — isolate the lowest set bit

Negation in two's complement flips all bits then adds 1, which has the effect of isolating the rightmost 1.

```
n      = 0110 1100
-n     = 1001 0100
n&(-n) = 0000 0100   ← only the lowest set bit remains
```

```go
lowest := n & (-n)          // value of lowest set bit (power of 2)
pos    := bits.TrailingZeros(uint(n))  // index of lowest set bit

// Iterate over all set bits in n
for n > 0 {
    bit := n & (-n)   // isolate
    n &= n - 1        // clear it
    // do something with bit
}
```

### XOR — the "find the odd one out" weapon

XOR cancels pairs: `a ^ a == 0`, `a ^ 0 == a`. So XOR-ing a list where every element appears twice leaves only the one that appears once.

```go
// Find the single number (every other appears twice)
result := 0
for _, v := range nums { result ^= v }  // pairs cancel, lone survivor remains

// Find the single number where others appear an even number of times — same trick
```

XOR also tells you which bits differ between two numbers, which makes it useful for finding the position of a change:

```go
diff := a ^ b
// bit i is set in diff ↔ a and b differ at bit i
// position of the first differing bit:
pos := bits.TrailingZeros(uint(diff))
```

### Shift as fast multiply / divide (powers of 2 only)

```go
n << k   // n * 2^k
n >> k   // n / 2^k  (floor, for positive n)

// Common: check if bit i is set
(n >> i) & 1 == 1

// Set bit i
n |= 1 << i

// Clear bit i
n &^= 1 << i   // &^ is Go's AND-NOT (bit clear)

// Toggle bit i
n ^= 1 << i
```

**Gotcha: shift amount must be non-negative and < bit width.** Shifting a 64-bit int by 64 or more is undefined behavior in some languages; in Go it's well-defined (result is 0) but you should still guard against it if the shift amount comes from user input.
