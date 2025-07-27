# Go Context

The `context` package in Go is used for controlling cancellation, deadlines, and passing request-scoped values across API boundaries and between processes.

## Basic Context Usage

```go
func main() {
    // Create a background context
    ctx := context.Background()
    
    // Create a derived context with cancellation
    ctx, cancel := context.WithCancel(ctx)
    
    // Don't forget to release resources when done
    defer cancel()
    
    // Use the context
    go doSomething(ctx)
    
    // Wait a bit then cancel
    time.Sleep(100 * time.Millisecond)
    cancel()
    
    // Give doSomething time to respond to cancellation
    time.Sleep(100 * time.Millisecond)
}

func doSomething(ctx context.Context) {
    select {
    case <-time.After(1 * time.Second):
        fmt.Println("doSomething completed")
    case <-ctx.Done():
        fmt.Println("doSomething canceled")
    }
}
```

## Context with Timeout

```go
func main() {
    // Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    go func() {
        // Simulate work that takes longer than the timeout
        time.Sleep(1 * time.Second)
        fmt.Println("Work completed, but too late")
    }()
    
    select {
    case <-time.After(100 * time.Millisecond):
        fmt.Println("Doing other work...")
    case <-ctx.Done():
        fmt.Println("Context done:", ctx.Err())
    }
    
    time.Sleep(1 * time.Second)
}
```

## Context with Deadline

```go
func main() {
    // Create a context with a deadline
    deadline := time.Now().Add(200 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    go processRequest(ctx)
    
    // Wait for the context to be done
    <-ctx.Done()
    fmt.Println("Main: context is done with error:", ctx.Err())
}

func processRequest(ctx context.Context) {
    // Simulate work
    select {
    case <-time.After(300 * time.Millisecond):
        fmt.Println("processRequest: work completed")
    case <-ctx.Done():
        fmt.Println("processRequest: work canceled")
    }
}
```

## Context with Values

```go
type key string

func main() {
    // Define keys for context values
    userIDKey := key("userID")
    authTokenKey := key("authToken")
    
    // Create a context with values
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, "user-123")
    ctx = context.WithValue(ctx, authTokenKey, "token-abc")
    
    // Pass the context to a function
    processWithAuth(ctx)
}

func processWithAuth(ctx context.Context) {
    // Retrieve values from context
    userIDKey := key("userID")
    authTokenKey := key("authToken")
    
    userID, ok := ctx.Value(userIDKey).(string)
    if !ok {
        fmt.Println("userID not found or not a string")
        return
    }
    
    authToken, ok := ctx.Value(authTokenKey).(string)
    if !ok {
        fmt.Println("authToken not found or not a string")
        return
    }
    
    fmt.Printf("Processing request for user %s with token %s\n", userID, authToken)
}
```

## Propagating Context Through Layers

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    result, err := fetchUserData(ctx, "user-123")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Result:", result)
}

func fetchUserData(ctx context.Context, userID string) (string, error) {
    // Check if context is already done before starting work
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
        // Continue with the work
    }
    
    // Call a downstream function, propagating the context
    return queryDatabase(ctx, "SELECT * FROM users WHERE id = "+userID)
}

func queryDatabase(ctx context.Context, query string) (string, error) {
    // Simulate a database query that takes time
    select {
    case <-time.After(500 * time.Millisecond):
        return "User data for " + query, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

## HTTP Server with Context

```go
func main() {
    http.HandleFunc("/api", handleRequest)
    http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Get the context from the request
    ctx := r.Context()
    
    // Create a derived context with timeout
    ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()
    
    // Process the request with the context
    result, err := processWithTimeout(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintln(w, result)
}

func processWithTimeout(ctx context.Context) (string, error) {
    // Simulate work
    select {
    case <-time.After(2 * time.Second):
        return "Work completed", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

## Best Practices

1. **Always pass context as the first parameter** to functions that may block or take time
2. **Don't store contexts in structs** - pass them explicitly
3. **Use context values only for request-scoped data** that transits process or API boundaries
4. **Create a hierarchy of contexts** rather than using a single context for everything
5. **Always call cancel()** when you're done with a context, typically with `defer`
6. **Check ctx.Done() in long-running operations** to allow early cancellation
7. **Don't use context.TODO()** in production code - it's meant for temporary use during development