// Topic 06: Pointers
// Run: go run 06_pointers/main.go
//
// JS has no pointers — objects/arrays are already reference types.
// In Go, everything is passed by VALUE (including structs, arrays).
// Pointers let you share and mutate data across function calls.
//
// Two operators:
//   & = "address of" — gets a pointer to a variable
//   * = "dereference" — gets the value at a pointer address

package main

import "fmt"

// -------------------------------------------------------------------------
// 1. Why pointers exist — pass by value vs pass by reference
// -------------------------------------------------------------------------

type BigStruct struct {
	Data [1000]int
}

// WITHOUT pointer: receives a full COPY of the int — changes don't stick
func incrementByValue(n int) {
	n++ // only changes the local copy
}

// WITH pointer: receives the ADDRESS — changes the original
func incrementByPointer(n *int) {
	*n++ // dereference to get/set the actual value
}

// Without pointer — struct change doesn't persist
func resetNameByValue(p Person) {
	p.Name = "nobody" // only affects local copy
}

// With pointer — change persists
func resetNameByPointer(p *Person) {
	p.Name = "nobody" // affects the original
	// Note: p.Name is shorthand for (*p).Name — Go auto-dereferences
}

type Person struct {
	Name string
	Age  int
}

// -------------------------------------------------------------------------
// 2. new() — allocates zeroed memory and returns a pointer
// Rarely used in practice; struct literals with & are more common
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// 3. When to use pointers
// Rule: use a pointer receiver/param when you need to:
//   a) Mutate the value
//   b) Avoid copying a large struct
//   c) Allow nil as a valid value
// -------------------------------------------------------------------------

// Linked list node — Next must be a pointer (recursive type, can't be value)
type ListNode struct {
	Val  int
	Next *ListNode // pointer to same type — only way to make recursive struct
}

func buildList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	cur := head
	for _, v := range vals[1:] {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return head
}

func printList(head *ListNode) {
	for head != nil {
		fmt.Print(head.Val)
		if head.Next != nil {
			fmt.Print(" -> ")
		}
		head = head.Next
	}
	fmt.Println()
}

// -------------------------------------------------------------------------
// 4. Pointer to pointer (**T) — rare but exists
// -------------------------------------------------------------------------

func setTo42(pp **int) {
	x := 42
	*pp = &x
}

// -------------------------------------------------------------------------
// 5. Nil pointers — the zero value for any pointer type
// ALWAYS check for nil before dereferencing (or you get a panic)
// -------------------------------------------------------------------------

func safePrint(p *int) {
	if p == nil {
		fmt.Println("nil pointer — nothing to print")
		return
	}
	fmt.Println("value:", *p)
}

func main() {
	// -------------------------------------------------------------------------
	// Basic pointer operations
	// -------------------------------------------------------------------------

	x := 42
	p := &x  // p is of type *int — "pointer to int"

	fmt.Println(x)  // 42    — the value
	fmt.Println(p)  // 0xc...  — the memory address
	fmt.Println(*p) // 42    — dereferenced: value at that address

	*p = 100        // modify through pointer
	fmt.Println(x)  // 100  — x was changed!

	// -------------------------------------------------------------------------
	// Pass by value vs pass by pointer
	// -------------------------------------------------------------------------

	n := 5
	incrementByValue(n)
	fmt.Println(n) // 5 — unchanged

	incrementByPointer(&n) // pass the ADDRESS of n
	fmt.Println(n)          // 6 — changed!

	// Struct example
	person := Person{Name: "Alice", Age: 30}

	resetNameByValue(person)
	fmt.Println(person.Name) // Alice — unchanged

	resetNameByPointer(&person)
	fmt.Println(person.Name) // nobody — changed!

	// -------------------------------------------------------------------------
	// new() vs & — two ways to get a pointer to a new allocation
	// -------------------------------------------------------------------------

	p1 := new(int)        // *int pointing to a zeroed int
	*p1 = 99
	fmt.Println(*p1)       // 99

	p2 := &Person{Name: "Bob", Age: 25} // more common idiom
	fmt.Println(p2.Name)                  // Bob

	// -------------------------------------------------------------------------
	// Pointer to pointer
	// -------------------------------------------------------------------------

	var pp *int
	fmt.Println(pp) // <nil>
	setTo42(&pp)
	fmt.Println(*pp) // 42

	// -------------------------------------------------------------------------
	// Nil pointer — always check before use
	// -------------------------------------------------------------------------

	var nilPtr *int
	safePrint(nilPtr) // safe: "nil pointer — nothing to print"

	val := 77
	safePrint(&val)   // "value: 77"

	// Dereferencing nil causes a PANIC (runtime error, like NPE in Java)
	// Uncomment to see:
	// fmt.Println(*nilPtr) // panic: runtime error: invalid memory address

	// -------------------------------------------------------------------------
	// Linked list with pointers
	// -------------------------------------------------------------------------

	list := buildList([]int{1, 2, 3, 4, 5})
	printList(list) // 1 -> 2 -> 3 -> 4 -> 5

	// Traverse and sum
	sum := 0
	for node := list; node != nil; node = node.Next {
		sum += node.Val
	}
	fmt.Println("sum:", sum) // 15

	// -------------------------------------------------------------------------
	// Pointer gotchas in loops
	// -------------------------------------------------------------------------

	// GOTCHA: taking address of loop variable (all point to same variable)
	ptrs := make([]*int, 3)
	for i := 0; i < 3; i++ {
		i := i          // shadow to create a new variable each iteration
		ptrs[i] = &i    // now each pointer is to a different variable
	}
	for _, ptr := range ptrs {
		fmt.Print(*ptr, " ") // 0 1 2 (correct because of shadowing)
	}
	fmt.Println()

	// WITHOUT shadowing (the bug):
	// for i := 0; i < 3; i++ {
	//     ptrs[i] = &i  // all point to the SAME i
	// }
	// After loop, i == 3, so all ptrs print 3

	// -------------------------------------------------------------------------
	// Slices and maps are already reference types — no pointer needed
	// -------------------------------------------------------------------------

	// Slice passed to function — modifications to ELEMENTS are visible
	// (because the slice header has a pointer to the backing array)
	s := []int{1, 2, 3}
	doubleElements(s)
	fmt.Println(s) // [2 4 6] — elements changed

	// But appending inside the function does NOT affect the caller
	// (unless you pass a *[]int or return the new slice)
	appendToSlice(s)
	fmt.Println(s) // [2 4 6] — unchanged! append created new backing array

	// -------------------------------------------------------------------------
	// EXERCISES
	// -------------------------------------------------------------------------

	// EXERCISE 1:
	// Write a function swap(a, b *int) that swaps two integers using pointers.
	// Verify x and y are actually swapped after the call.

	// EXERCISE 2:
	// Write a function doubleAll(nums []int) that doubles every element
	// in-place (modifying the original slice, not returning a new one).

	// EXERCISE 3:
	// Implement reverse(*ListNode) *ListNode that reverses a linked list.
	// Input: 1 -> 2 -> 3 -> 4 -> 5
	// Output: 5 -> 4 -> 3 -> 2 -> 1

	// EXERCISE 4:
	// Write a function findMiddle(*ListNode) *ListNode using the
	// slow/fast pointer technique (Floyd's algorithm).
	// For [1,2,3,4,5] return the node with Val=3.

	// EXERCISE 5 (Challenge):
	// Implement a simple doubly-linked list:
	// type DNode struct { Val int; Prev, Next *DNode }
	// Methods: PushFront, PushBack, Delete, Print (forward), Print (backward)
}

func doubleElements(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

func appendToSlice(s []int) {
	s = append(s, 99) // local re-assignment, caller's slice unchanged
	fmt.Println("inside appendToSlice:", s) // [2 4 6 99]
}
