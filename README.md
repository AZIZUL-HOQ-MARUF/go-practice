# Go Learning Curriculum

For a JS/Angular dev targeting Go fluency.

## How to run

```bash
go run 01_basics/main.go
go run 02_control_flow/main.go
go run 03_functions/main.go
go run 04_collections/main.go
go run 05_structs/main.go
go run 06_pointers/main.go
go run 07_interfaces/main.go
go run 08_errors/main.go
go run 09_goroutines/main.go
```

## Topics

| # | File | Topic | JS Equivalent |
|---|------|-------|---------------|
| 01 | `01_basics` | Variables, types, zero values, constants, `:=` | `let`/`const`, TypeScript types |
| 02 | `02_control_flow` | `for` (all 3 forms), `range`, `if`, `switch` | `for`/`while`/`for...of`, `switch` |
| 03 | `03_functions` | Multi-return, variadic, closures, `defer` | destructuring returns, `...rest`, closures |
| 04 | `04_collections` | Arrays (fixed), slices (dynamic), maps | `[]` / `Array`, `Map`/`{}` |
| 05 | `05_structs` | Structs, methods, embedding, struct tags | `class`, `extends` → composition |
| 06 | `06_pointers` | `&`/`*`, pass-by-value vs pointer, nil | no direct equivalent in JS |
| 07 | `07_interfaces` | Implicit interfaces, type assertions, type switch | TypeScript interfaces (implicit) |
| 08 | `08_errors` | `error` type, `errors.Is/As`, wrapping, panic/recover | `try`/`catch`/`throw` |
| 09 | `09_goroutines` | Goroutines, channels, select, WaitGroup, mutex | `async`/`await`, `Promise.all` |

## Key differences from JS

- **No `undefined`** — every type has a zero value (`0`, `""`, `false`, `nil`)
- **Only `for`** — no `while`, `do-while`, `forEach`. `for` does everything.
- **Explicit errors** — no exceptions; functions return `(result, error)`
- **Static types** — compiler catches type mismatches (like strict TypeScript)
- **Value types** — structs/arrays are copied, not referenced (use `&` for sharing)
- **Exported = Capital** — `MyFunc` is public, `myFunc` is private
- **No classes** — structs + methods + interfaces replace OOP inheritance

## Each file structure

1. Concept explanation with JS comparison
2. Working examples (just run the file to see output)
3. `// EXERCISE:` comments — fill these in to practice
