# Go Testing Tips

This document covers best practices and techniques for effective testing in Go.

## Basic Testing

### Writing Simple Tests

```go
// In file: math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}
```

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed signs", -2, 3, 1},
    }
    
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
            }
        })
    }
}
```

### Subtests

```go
func TestStrings(t *testing.T) {
    t.Run("concatenation", func(t *testing.T) {
        result := "Hello, " + "World"
        if result != "Hello, World" {
            t.Error("String concatenation failed")
        }
    })
    
    t.Run("length", func(t *testing.T) {
        s := "Hello"
        if len(s) != 5 {
            t.Errorf("Length of %q = %d; want 5", s, len(s))
        }
    })
}
```

## Test Organization

### Test Helpers

```go
func setupTestCase(t *testing.T) func(t *testing.T) {
    t.Log("Setting up test case")
    
    // Return a function to be called at the end of the test
    return func(t *testing.T) {
        t.Log("Tearing down test case")
        // Cleanup code here
    }
}

func TestWithHelper(t *testing.T) {
    teardown := setupTestCase(t)
    defer teardown(t)
    
    // Test code here
}
```

### Using testdata Directory

```go
func TestReadConfig(t *testing.T) {
    // Read test data from the testdata directory
    data, err := os.ReadFile("testdata/config.json")
    if err != nil {
        t.Fatalf("Failed to read test data: %v", err)
    }
    
    // Use the test data
    config, err := ParseConfig(data)
    if err != nil {
        t.Fatalf("ParseConfig() error = %v", err)
    }
    
    if config.Name != "test-config" {
        t.Errorf("config.Name = %q; want %q", config.Name, "test-config")
    }
}
```

## Mocking and Interfaces

### Using Interfaces for Testability

```go
// Define an interface for external dependencies
type DataStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

// Implementation that can be replaced in tests
type Service struct {
    store DataStore
}

func (s *Service) Process(key string) (string, error) {
    value, err := s.store.Get(key)
    if err != nil {
        return "", err
    }
    // Process value...
    return value + " processed", nil
}
```

### Mock Implementation

```go
// Mock implementation for testing
type MockDataStore struct {
    GetFunc func(key string) (string, error)
    SetFunc func(key, value string) error
}

func (m *MockDataStore) Get(key string) (string, error) {
    return m.GetFunc(key)
}

func (m *MockDataStore) Set(key, value string) error {
    return m.SetFunc(key, value)
}

func TestService_Process(t *testing.T) {
    // Setup mock
    mock := &MockDataStore{
        GetFunc: func(key string) (string, error) {
            if key == "test-key" {
                return "test-value", nil
            }
            return "", fmt.Errorf("key not found")
        },
    }
    
    service := &Service{store: mock}
    
    // Test successful case
    result, err := service.Process("test-key")
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if result != "test-value processed" {
        t.Errorf("result = %q; want %q", result, "test-value processed")
    }
    
    // Test error case
    _, err = service.Process("unknown-key")
    if err == nil {
        t.Error("Expected error, got nil")
    }
}
```

## HTTP Testing

### Testing HTTP Handlers

```go
func TestHandler(t *testing.T) {
    // Create a request
    req, err := http.NewRequest("GET", "/api/items", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create the handler
    handler := http.HandlerFunc(GetItemsHandler)
    
    // Serve the request
    handler.ServeHTTP(rr, req)
    
    // Check status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // Check response body
    expected := `{"items":[]}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

### Testing HTTP Clients

```go
func TestClient(t *testing.T) {
    // Start a test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/api/data" {
            t.Errorf("Expected request to '/api/data', got %q", r.URL.Path)
            http.Error(w, "Not found", http.StatusNotFound)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    }))
    defer server.Close()
    
    // Use the test server URL in your client
    client := NewAPIClient(server.URL)
    resp, err := client.FetchData()
    
    if err != nil {
        t.Fatalf("FetchData() error = %v", err)
    }
    
    if resp.Status != "ok" {
        t.Errorf("resp.Status = %q; want %q", resp.Status, "ok")
    }
}
```

## Advanced Testing

### Parallel Tests

```go
func TestParallel(t *testing.T) {
    // Mark test as capable of running in parallel with other tests
    t.Parallel()
    
    // Test code here
}
```

### Benchmarks

```go
func BenchmarkFibonacci(b *testing.B) {
    // Run the Fibonacci function b.N times
    for n := 0; n < b.N; n++ {
        Fibonacci(10)
    }
}

// Run with: go test -bench=.
```

### Fuzzing

```go
// Available in Go 1.18+
func FuzzReverse(f *testing.F) {
    // Provide seed corpus
    testcases := []string{"hello", "world", ""}
    for _, tc := range testcases {
        f.Add(tc)
    }
    
    // Fuzz test
    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)
        
        // Property: reversing twice should return original string
        if orig != doubleRev {
            t.Errorf("Reverse(Reverse(%q)) = %q, want %q", orig, doubleRev, orig)
        }
    })
}

// Run with: go test -fuzz=Fuzz
```

## Test Coverage

### Measuring Coverage

```bash
# Run tests with coverage
go test -cover

# Generate coverage profile
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Coverage in CI

```yaml
# Example GitHub Actions workflow step
- name: Run tests with coverage
  run: go test -race -coverprofile=coverage.txt -covermode=atomic ./...

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v1
  with:
    file: ./coverage.txt
```

## Best Practices

1. **Write tests before fixing bugs** - Create a test that reproduces the issue
2. **Keep tests fast** - Slow tests discourage frequent testing
3. **One assertion per test** - Makes it clear what failed
4. **Use table-driven tests** - Reduces boilerplate for multiple test cases
5. **Test exported functionality** - Focus on the public API
6. **Use subtests for organization** - Improves readability and allows running specific cases
7. **Don't test standard library** - Assume it works correctly
8. **Use t.Helper() for helper functions** - Improves error reporting
9. **Clean up test resources** - Use defer for cleanup
10. **Run tests with race detector periodically** - `go test -race`