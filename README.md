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
go run 10_strings_strconv/main.go
go run 11_generics/main.go
go run 13_stdlib_toolkit/main.go
go run 14_packages/main.go
go run 15_context/main.go
go run 16_json_http/main.go
# 12_testing: go test ./12_testing/...
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
| 10 | `10_strings_strconv` | `strings` package, `strconv`, rune vs byte, `unicode` | `String` methods, `parseInt`/`parseFloat` |
| 11 | `11_generics` | Type parameters, constraints, generic data structures | TypeScript generics |
| 12 | `12_testing` | Table-driven tests, benchmarks, examples | Jest / describe+it |
| 13 | `13_stdlib_toolkit` | `sort`, `container/heap`, `math`, `time`, `os`, `slog` | lodash, Date, fs, winston |
| 14 | `14_packages` | Modules, packages, visibility, `init()`, `go.mod` | npm packages, `import`/`export` |
| 15 | `15_context` | `context.Context`, cancellation, deadlines, values | AbortController |
| 16 | `16_json_http` | JSON marshal/unmarshal, `net/http` server, middleware | `JSON.stringify`, Express.js |

## Key differences from JS

- **No `undefined`** — every type has a zero value (`0`, `""`, `false`, `nil`)
- **Only `for`** — no `while`, `do-while`, `forEach`. `for` does everything.
- **Explicit errors** — no exceptions; functions return `(result, error)`
- **Static types** — compiler catches type mismatches (like strict TypeScript)
- **Value types** — structs/arrays are copied, not referenced (use `&` for sharing)
- **Exported = Capital** — `MyFunc` is public, `myFunc` is private
- **No classes** — structs + methods + interfaces replace OOP inheritance

## Supplemental reading — karan99's "Learn Go:" series

Each topic pairs well with these Medium posts. Read them alongside the exercises for a second explanation angle.

| Our topic | karan99 post(s) |
|---|---|
| `01_basics` | [Variables and Data Types](https://medium.com/@karan99/learn-go-variables-and-data-types-653f0de50ca9) · [String Formatting](https://medium.com/@karan99/learn-go-string-formatting-94c8d27402b9) |
| `02_control_flow` | [Flow Control](https://medium.com/@karan99/learn-go-flow-control-9d0b663a6dfd) |
| `03_functions` | [Functions](https://medium.com/@karan99/learn-go-functions-5e2d66f67159) — note: his `init()` section maps to our `14_packages` |
| `04_collections` | [Arrays and Slices](https://medium.com/@karan99/learn-go-arrays-and-slices-0836777b06ad) · [Maps](https://medium.com/@karan99/learn-go-maps-744cc9e80166) |
| `05_structs` | [Structs](https://medium.com/@karan99/learn-go-structs-5f812e57456f) · [Methods](https://medium.com/@karan99/learn-go-methods-30102d5dc6a2) |
| `06_pointers` | [Pointers](https://medium.com/@karan99/learn-go-pointers-21cd483c72ec) |
| `07_interfaces` | [Interfaces](https://medium.com/@karan99/learn-go-interfaces-9d3d3ebc82cb) |
| `08_errors` | [Errors](https://medium.com/@karan99/learn-go-errors-57fd86542ee8) · [Panic and Recover](https://medium.com/@karan99/learn-go-panic-and-recover-8f5a23b51a62) |
| `09_goroutines` | [Concurrency](https://medium.com/@karan99/learn-go-concurrency-cfc72015db17) (theory: CSP, data race, deadlock) · [Goroutines](https://medium.com/@karan99/learn-go-goroutines-431bf6c09ef0) |
| `11_generics` | [Generics](https://medium.com/@karan99/learn-go-generics-9ff505a0c028) |
| `12_testing` | [Testing](https://medium.com/@karan99/learn-go-testing-e36763dfccee) |
| `14_packages` | [Modules](https://medium.com/@karan99/learn-go-modules-0dae76d1e250) · [Packages](https://medium.com/@karan99/learn-go-packages-433a8e7b0d24) · [Workspaces](https://medium.com/@karan99/learn-go-workspaces-4208b4e7f95e) · [Useful Commands](https://medium.com/@karan99/learn-go-useful-commands-abeab17e7f96) · [Build](https://medium.com/@karan99/learn-go-build-cd9360b19ca0) |

> `10_strings_strconv`, `13_stdlib_toolkit`, `15_context`, `16_json_http` — not covered in karan's series; our materials are the primary source for these.

## Git workflow — syncing practice branch with master

`master` holds lesson material + exercise stubs only. `practice-aziz` (or any personal branch) is where completed exercises live.

**When master gets new topics, sync like this:**

```bash
# 1. Commit your current work first
git add -A && git commit -m "wip: completed exercises up to topic X"

# 2. Preview what will change (optional dry run)
git merge master --no-commit --no-ff
git status
git merge --abort   # cancel if anything looks unexpected

# 3. Merge master into your practice branch
git merge master

# 4. Push
git push origin practice-aziz
```

**Do NOT** create a new branch from master — that discards your completed exercises. Always merge master INTO your practice branch.

If there are conflicts (only happens when master edits a file you already solved): keep your solutions, add any new exercise stubs below them, then commit.

## Each file structure

1. Concept explanation with JS comparison
2. Working examples (just run the file to see output)
3. `// EXERCISE:` comments — fill these in to practice
