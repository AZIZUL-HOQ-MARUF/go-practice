// Topic 14: Packages, Modules, and Project Structure
// Run: go run 14_packages/main.go
//
// This file explains Go's package system. The real demo is in the
// subdirectory structure below — look at the subpackages.
//
// For a migration project at Google: proper package organization
// is critical. Go projects at Google follow specific patterns
// from the internal style guide (which mirrors the public one).

package main

import (
	"fmt"

	// Importing subpackages from this module (practice/golang)
	"practice/golang/14_packages/mathutil"
	"practice/golang/14_packages/stringutil"
)

// =========================================================================
// PACKAGE CONCEPTS (read this carefully — no JS equivalent)
// =========================================================================
//
// 1. PACKAGE vs MODULE
//    - Module: the top-level unit, defined by go.mod. Has a module path.
//    - Package: a directory of .go files. All files in a dir share one package name.
//    - One module can contain many packages.
//
// 2. PACKAGE NAMING
//    - Package name = directory name (convention, not enforced)
//    - Import path = module path + directory path from module root
//    - Usage: last segment of import path is the identifier
//      import "practice/golang/14_packages/mathutil" → use as mathutil.Add()
//
// 3. VISIBILITY (exported vs unexported)
//    - Capital letter = exported (public) — accessible from other packages
//    - Lowercase = unexported (private) — only accessible within the package
//    - Applies to: functions, types, variables, constants, struct fields
//
//    JS: everything is exported unless you use module.exports selectively
//    Go: you explicitly choose what's public via capitalization
//
// 4. INIT FUNCTIONS
//    - Each package can have one or more init() functions
//    - Run automatically after package-level vars are initialized
//    - Run before main()
//    - Used for: registering drivers, initializing global state
//    - Execution order: imports init() before the importing package's init()
//
// 5. BLANK IMPORT
//    import _ "package/path"
//    - Imports only for side effects (runs init())
//    - Common for: database drivers, image format decoders, plugin registration
//    - JS equiv: import 'side-effect-module' (no named exports used)
//
// 6. INTERNAL PACKAGES
//    - A directory named "internal" can only be imported by code in its parent
//    - Enforced by the compiler — great for hiding implementation details
//    - practice/golang/14_packages/internal/... can only be imported by
//      practice/golang/14_packages/...
//
// 7. GO.MOD ESSENTIALS
//    module practice/golang  ← module path
//    go 1.22                 ← minimum Go version
//
//    require (
//        github.com/stretchr/testify v1.9.0  ← external dependency
//    )
//
// 8. KEY COMMANDS
//    go mod init <module-path>  ← create go.mod
//    go get <package>           ← add dependency
//    go mod tidy                ← clean up go.mod (remove unused, add missing)
//    go mod vendor              ← copy deps to vendor/ (Google uses this)
//    go build ./...             ← build all packages
//    go test ./...              ← test all packages
//    go vet ./...               ← static analysis (run before every commit)
//    gofmt -w .                 ← format all files (mandatory at Google)
//    golangci-lint run          ← comprehensive linting
//
// 9. PACKAGE ORGANIZATION PATTERNS
//
//    Flat (small projects):
//    myapp/
//    ├── main.go
//    ├── handler.go
//    └── storage.go
//
//    Layered (medium projects):
//    myapp/
//    ├── main.go
//    ├── cmd/           ← entry points (multiple binaries)
//    │   └── server/
//    │       └── main.go
//    ├── internal/      ← private packages (not importable externally)
//    │   ├── handler/
//    │   ├── service/
//    │   └── repo/
//    └── pkg/           ← public reusable packages
//        └── mathutil/
//
//    Google monorepo style:
//    - Packages are very small and focused
//    - Deep directory hierarchies
//    - Heavy use of internal/
//    - Each binary is in its own cmd/ subdirectory

func main() {
	// Using exported functions from subpackages
	fmt.Println(mathutil.Add(3, 4))           // 7
	fmt.Println(mathutil.Multiply(3, 4))      // 12
	fmt.Println(mathutil.IsPrime(17))         // true
	fmt.Println(mathutil.Fibonacci(10))       // 55

	fmt.Println(stringutil.Reverse("hello"))      // olleh
	fmt.Println(stringutil.IsPalindrome("racecar")) // true
	fmt.Println(stringutil.WordCount("hello world hello")) // map[hello:2 world:1]

	// Accessing exported constants and types from subpackage
	fmt.Println(mathutil.Pi)        // 3.14159...
	fmt.Println(mathutil.MaxPrime)  // 1000

	calc := mathutil.NewCalculator()
	calc.Add(10)
	calc.Multiply(3)
	fmt.Println(calc.Result()) // 30

	// =========================================================================
	// VISIBILITY DEMONSTRATION
	// =========================================================================

	// mathutil.helperFunc() ← would NOT compile — unexported
	// mathutil.internalConst ← would NOT compile — unexported

	// But you can use exported types whose internal fields are unexported:
	// The Calculator type is exported but its `value` field is unexported
	// This is encapsulation — same as private class fields in JS

	// =========================================================================
	// EXERCISES
	// =========================================================================

	// EXERCISE 1:
	// Create a new package practice/golang/14_packages/geometry
	// with an exported Shape interface, Circle and Rectangle types,
	// and Area(), Perimeter() methods. Import and use it here.

	// EXERCISE 2:
	// Add an init() function to mathutil that prints "mathutil initialized"
	// and sets a package-level default precision variable.

	// EXERCISE 3:
	// Create practice/golang/14_packages/internal/validator
	// with ValidateEmail(s string) bool and ValidateAge(n int) bool.
	// Try importing it from 14_packages/main.go (should work).
	// Then try importing it from 15_context/main.go (should fail — internal).

	// EXERCISE 4:
	// Add an external dependency (e.g., github.com/stretchr/testify) with:
	//   go get github.com/stretchr/testify@latest
	// Write a test in 12_testing/ that uses testify's assert.Equal instead
	// of manual if/t.Error. Compare the two styles.

	// EXERCISE 5 (Production):
	// Design a package structure for a simple user service:
	// - cmd/userservice/main.go (entry point)
	// - internal/handler/user.go (HTTP handlers)
	// - internal/service/user.go (business logic)
	// - internal/repo/user.go (data access)
	// - pkg/model/user.go (shared User struct)
	// Draw the dependency graph — handler depends on service,
	// service depends on repo. Never the reverse.
}
