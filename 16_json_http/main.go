// Topic 16: JSON + HTTP — the backbone of modern services
// Run: go run 16_json_http/main.go
//
// This is directly relevant to your migration project.
// Most migration work involves converting service endpoints —
// Go's net/http + encoding/json is the standard for REST APIs.
//
// JS analogy:
//   json.Marshal     ≈ JSON.stringify
//   json.Unmarshal   ≈ JSON.parse
//   http.HandleFunc  ≈ app.get('/path', handler) in Express
//   http.Client      ≈ fetch() or axios

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

// =========================================================================
// PART 1: JSON ENCODING/DECODING
// =========================================================================

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip,omitempty"` // omit if empty string
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`                   // never serialized
	CreatedAt time.Time `json:"created_at"`
	Address   *Address  `json:"address,omitempty"`   // omit if nil pointer
	Tags      []string  `json:"tags,omitempty"`      // omit if nil/empty
	Active    bool      `json:"active"`
	Score     *float64  `json:"score,omitempty"`     // pointer: can distinguish 0 from absent
}

// APIResponse wraps all API responses in a standard envelope.
// Common pattern in production services.
type APIResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func jsonDemo() {
	fmt.Println("=== JSON ===")

	// --- Marshal (Go → JSON) ---
	score := 95.5
	user := User{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "secret123", // will NOT appear in JSON
		CreatedAt: time.Now(),
		Tags:      []string{"admin", "user"},
		Active:    true,
		Score:     &score,
	}

	// json.MarshalIndent for pretty printing (logging, debugging)
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println(string(data))

	// Compact marshal (for network transmission)
	compact, _ := json.Marshal(user)
	fmt.Printf("  compact (%d bytes): %s...\n", len(compact), compact[:50])

	// --- Unmarshal (JSON → Go) ---
	jsonStr := `{
		"id": 2,
		"name": "Bob",
		"email": "bob@example.com",
		"active": false,
		"address": {"street": "123 Main St", "city": "Anytown"}
	}`

	var decoded User
	if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	fmt.Printf("  decoded: ID=%d Name=%s Active=%v Address=%+v\n",
		decoded.ID, decoded.Name, decoded.Active, decoded.Address)

	// --- Decode into map (when schema is unknown) ---
	var raw map[string]any
	json.Unmarshal([]byte(jsonStr), &raw)
	fmt.Printf("  raw name: %v\n", raw["name"])
	fmt.Printf("  raw address type: %T\n", raw["address"]) // map[string]interface {}

	// --- Streaming: json.Decoder (preferred for HTTP bodies and large data) ---
	// Use Decoder when reading from io.Reader (HTTP body, file)
	// Avoids loading entire JSON into memory
	reader := strings.NewReader(`{"id":3,"name":"Carol","email":"carol@example.com","active":true}`)
	var carol User
	if err := json.NewDecoder(reader).Decode(&carol); err != nil {
		fmt.Println("decode error:", err)
	}
	fmt.Printf("  streamed: %s\n", carol.Name)

	// Decode multiple JSON objects from a stream (newline-delimited JSON)
	ndJson := `{"id":1,"name":"A","email":"a@x.com","active":true}
{"id":2,"name":"B","email":"b@x.com","active":false}
{"id":3,"name":"C","email":"c@x.com","active":true}`

	dec := json.NewDecoder(strings.NewReader(ndJson))
	for dec.More() {
		var u User
		dec.Decode(&u)
		fmt.Printf("  ndjson user: %d %s\n", u.ID, u.Name)
	}

	// --- json.Encoder (stream to io.Writer) ---
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.Encode(user)
	fmt.Printf("  encoded to buffer (%d bytes)\n", buf.Len())

	// --- Handling optional fields with pointers ---
	// Without pointer: `Score float64` — can't tell if 0.0 is "set to 0" or "not provided"
	// With pointer: `Score *float64` — nil means not provided, &0.0 means explicitly 0

	noScore := User{ID: 99, Name: "Dave", Email: "dave@example.com", Active: true}
	b, _ := json.Marshal(noScore)
	fmt.Printf("  no score: %s\n", string(b)) // score field omitted

	zero := 0.0
	withZeroScore := User{ID: 100, Name: "Eve", Email: "eve@example.com", Active: true, Score: &zero}
	b2, _ := json.Marshal(withZeroScore)
	fmt.Printf("  zero score: %s\n", string(b2)) // "score":0
}

// =========================================================================
// PART 2: HTTP SERVER
// =========================================================================

// Request/response types
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// writeJSON is a helper to write JSON responses consistently.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJSON encode failed", "err", err)
	}
}

// writeError writes a standard error JSON response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// usersHandler handles /users requests.
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// In production: query DB with r.Context() for cancellation
	users := []UserResponse{
		{1, "Alice", "alice@example.com"},
		{2, "Bob", "bob@example.com"},
	}
	writeJSON(w, http.StatusOK, APIResponse[[]UserResponse]{
		Data:    users,
		Success: true,
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Limit request body size (important for production)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Validation
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if !strings.Contains(req.Email, "@") {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}

	// In production: insert into DB using r.Context()
	created := UserResponse{ID: 3, Name: req.Name, Email: req.Email}
	writeJSON(w, http.StatusCreated, APIResponse[UserResponse]{
		Data:    created,
		Success: true,
	})
}

// loggingMiddleware wraps a handler and logs each request.
// Middleware is just a function that wraps http.Handler.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
		)
		next.ServeHTTP(w, r)
		slog.Info("response",
			"method", r.Method,
			"path", r.URL.Path,
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})
}

// authMiddleware checks for a bearer token.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}
		// In production: validate token, extract user ID, add to context
		next.ServeHTTP(w, r)
	})
}

// =========================================================================
// PART 3: HTTP CLIENT
// =========================================================================

// httpClientDemo shows production-ready HTTP client usage.
func httpClientDemo(serverURL string) {
	fmt.Println("\n=== HTTP CLIENT ===")

	// Always use a custom client with timeouts — never use http.DefaultClient in production
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
		},
	}

	// GET request with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+"/users", nil)
	if err != nil {
		fmt.Println("  create request error:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("  GET error:", err)
		return
	}
	defer resp.Body.Close() // ALWAYS close the body

	fmt.Printf("  GET /users status: %d\n", resp.StatusCode)

	var result APIResponse[[]UserResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("  decode error:", err)
		return
	}
	fmt.Printf("  got %d users\n", len(result.Data))
	for _, u := range result.Data {
		fmt.Printf("    - %s (%s)\n", u.Name, u.Email)
	}

	// POST request
	body, _ := json.Marshal(CreateUserRequest{Name: "Charlie", Email: "charlie@example.com"})

	postReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/users",
		bytes.NewReader(body))
	postReq.Header.Set("Authorization", "Bearer test-token")
	postReq.Header.Set("Content-Type", "application/json")

	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("  POST error:", err)
		return
	}
	defer postResp.Body.Close()

	respBody, _ := io.ReadAll(postResp.Body)
	fmt.Printf("  POST /users status: %d body: %s\n", postResp.StatusCode, respBody)
}

func main() {
	jsonDemo()

	fmt.Println("\n=== HTTP SERVER + CLIENT ===")

	// Set up mux (router)
	mux := http.NewServeMux()
	// Go 1.22+: enhanced patterns with method and path parameters
	mux.HandleFunc("GET /users", getUsers)
	mux.HandleFunc("POST /users", createUser)
	mux.Handle("/users", loggingMiddleware(authMiddleware(http.HandlerFunc(usersHandler))))

	// Find a free port for our demo server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	serverURL := "http://" + listener.Addr().String()

	// Wrap with middleware
	handler := loggingMiddleware(mux)

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in background
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(10 * time.Millisecond)
	fmt.Println("  server running at", serverURL)

	// Demo the client
	httpClientDemo(serverURL)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	fmt.Println("  server shut down")

	// =========================================================================
	// EXERCISES
	// =========================================================================

	// EXERCISE 1:
	// Add a GET /users/{id} endpoint (Go 1.22+ pattern matching):
	// mux.HandleFunc("GET /users/{id}", getUserByID)
	// Extract the ID with r.PathValue("id"), return 404 if not found.

	// EXERCISE 2:
	// Add request body validation middleware that checks Content-Type is
	// "application/json" for POST/PUT/PATCH requests.

	// EXERCISE 3:
	// Write a generic helper: func DecodeJSON[T any](r *http.Request) (T, error)
	// that decodes the request body and validates it's not empty.

	// EXERCISE 4:
	// Write a robust HTTP client function:
	// func Get[T any](ctx context.Context, client *http.Client, url string, headers map[string]string) (T, error)
	// that handles: context cancellation, non-2xx status codes as errors,
	// and decodes the JSON response body.

	// EXERCISE 5 (Production):
	// Add graceful shutdown to the server:
	// - Listen for os.Signal (SIGTERM, SIGINT)
	// - Call server.Shutdown(ctx) with a 30s timeout
	// - Wait for in-flight requests to complete
	// This is mandatory for Google Cloud Run / Kubernetes deployments.
}
