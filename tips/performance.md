# Go Performance Tips

This document covers practical tips and techniques for optimizing Go code performance.

## Memory Management

### Reduce Allocations

```go
// AVOID: Creates a new slice on each iteration
for i := 0; i < 100; i++ {
    data := make([]byte, 1024)
    process(data)
}

// BETTER: Reuse the same slice
data := make([]byte, 1024)
for i := 0; i < 100; i++ {
    process(data)
}
```

### Use Object Pools for Frequently Allocated Objects

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processRequest() {
    // Get a buffer from the pool
    buf := bufferPool.Get().(*bytes.Buffer)
    
    // Reset the buffer (don't allocate a new one)
    buf.Reset()
    
    // Use the buffer...
    
    // Return it to the pool when done
    bufferPool.Put(buf)
}
```

### Preallocate Slices When Size is Known

```go
// AVOID: Multiple allocations as slice grows
var data []int
for i := 0; i < 10000; i++ {
    data = append(data, i)
}

// BETTER: Single allocation with known capacity
data := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    data = append(data, i)
}
```

## Efficient Data Structures

### Use Maps for Lookups

```go
// O(n) lookup time
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// O(1) lookup time with a map
func containsFast(items map[string]struct{}, item string) bool {
    _, exists := items[item]
    return exists
}
```

### Use Appropriate Data Structures

```go
// Use a slice for ordered data with infrequent modifications
orderedItems := []string{"a", "b", "c"}

// Use a map for fast lookups
lookupTable := map[string]int{"a": 1, "b": 2, "c": 3}

// Use a heap for priority queue operations
pq := &priorityQueue{}
heap.Init(pq)
```

## Concurrency Optimization

### Use Buffered Channels When Appropriate

```go
// Unbuffered channel - sender blocks until receiver is ready
ch := make(chan int)

// Buffered channel - sender only blocks when buffer is full
ch := make(chan int, 100)
```

### Limit Goroutine Count

```go
func processItems(items []Item) {
    // AVOID: Launching too many goroutines
    // for _, item := range items {
    //     go processItem(item)
    // }
    
    // BETTER: Limit concurrent goroutines
    const maxConcurrent = 10
    sem := make(chan struct{}, maxConcurrent)
    
    for _, item := range items {
        sem <- struct{}{} // Acquire token
        go func(item Item) {
            defer func() { <-sem }() // Release token
            processItem(item)
        }(item)
    }
    
    // Wait for all goroutines to finish
    for i := 0; i < maxConcurrent; i++ {
        sem <- struct{}{}
    }
}
```

### Use Worker Pools for Batch Processing

```go
func processWithWorkerPool(items []Item, numWorkers int) {
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // Start workers
    for w := 0; w < numWorkers; w++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(items); i++ {
        <-results
    }
}

func worker(jobs <-chan Item, results chan<- Result) {
    for job := range jobs {
        results <- processItem(job)
    }
}
```

## String Handling

### Use Strings Builder for String Concatenation

```go
// AVOID: Creates many temporary strings
func concatenateStrings(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ","
    }
    return result
}

// BETTER: Uses a single buffer
func concatenateStringsOptimized(items []string) string {
    var sb strings.Builder
    
    // Optionally preallocate the buffer
    sb.Grow(len(items) * 8) // Estimate size
    
    for _, item := range items {
        sb.WriteString(item)
        sb.WriteByte(',')
    }
    
    return sb.String()
}
```

### Avoid Unnecessary String Conversions

```go
// AVOID: Unnecessary string-to-bytes conversion
func processString(s string) {
    b := []byte(s)
    // Process bytes
    _ = b
}

// BETTER: Use strings.Reader to avoid allocation
func processStringOptimized(s string) {
    r := strings.NewReader(s)
    // Process reader
    _ = r
}
```

## Profiling and Benchmarking

### Use the Built-in Profiler

```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    // CPU profiling
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Run your program...
    
    // Memory profiling
    f2, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f2)
    f2.Close()
}
```

### Write Benchmarks

```go
// In _test.go file
func BenchmarkMyFunction(b *testing.B) {
    // Reset timer if setup is expensive
    // b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        MyFunction()
    }
}

// Run with: go test -bench=. -benchmem
```

## Compiler Optimizations

### Use Go Build Flags

```bash
# Build with optimizations
go build -gcflags="-N -l" main.go  # Disable optimizations (for debugging)
go build -ldflags="-s -w" main.go  # Strip debug info (smaller binary)
```

### Inline Small Functions

```go
// Go will automatically inline small functions
// Use //go:noinline directive to prevent inlining when needed
//go:noinline
func doNotInlineMe() {
    // ...
}
```

## General Tips

1. **Profile before optimizing** - Find actual bottlenecks first
2. **Benchmark changes** - Verify improvements with measurements
3. **Consider readability** - Don't sacrifice maintainability for small gains
4. **Use latest Go version** - Each release includes performance improvements
5. **Optimize hot paths** - Focus on frequently executed code
6. **Avoid premature optimization** - Build it right, then make it fast