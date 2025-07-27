# Go Goroutines

Goroutines are lightweight threads managed by the Go runtime. They enable concurrent execution with minimal resources.

## Basic Goroutine

```go
func main() {
    // Start a goroutine
    go sayHello()
    
    // Main goroutine continues execution
    fmt.Println("Main function")
    
    // Wait to see the goroutine output
    time.Sleep(100 * time.Millisecond)
}

func sayHello() {
    fmt.Println("Hello from goroutine!")
}
```

## Multiple Goroutines

```go
func main() {
    for i := 0; i < 5; i++ {
        // Each iteration launches a new goroutine
        go func(id int) {
            fmt.Printf("Goroutine %d running\n", id)
        }(i)
    }
    
    // Wait for goroutines to finish
    time.Sleep(100 * time.Millisecond)
}
```

## Common Pitfalls

### Closure Variable Capture

```go
func main() {
    // INCORRECT: All goroutines will likely print the same value
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i) // Captures i by reference
        }()
    }
    
    // CORRECT: Pass the variable as a parameter
    for i := 0; i < 5; i++ {
        go func(val int) {
            fmt.Println(val) // Gets its own copy of i
        }(i)
    }
    
    time.Sleep(100 * time.Millisecond)
}
```

### Proper Synchronization

```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1) // Increment counter before launching goroutine
        
        go func(id int) {
            defer wg.Done() // Decrement counter when goroutine completes
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    
    // Wait for all goroutines to complete
    wg.Wait()
    fmt.Println("All workers completed")
}
```

## Limiting Concurrency

```go
func main() {
    const maxConcurrency = 3
    jobs := make([]int, 10)
    
    // Initialize jobs
    for i := range jobs {
        jobs[i] = i + 1
    }
    
    // Create a semaphore using a buffered channel
    semaphore := make(chan struct{}, maxConcurrency)
    var wg sync.WaitGroup
    
    for _, job := range jobs {
        wg.Add(1)
        
        // Acquire semaphore
        semaphore <- struct{}{}
        
        go func(jobID int) {
            defer func() {
                // Release semaphore
                <-semaphore
                wg.Done()
            }()
            
            // Process job
            fmt.Printf("Processing job %d\n", jobID)
            time.Sleep(100 * time.Millisecond)
        }(job)
    }
    
    wg.Wait()
    fmt.Println("All jobs completed")
}
```

## Worker Pool Pattern

```go
func main() {
    const numJobs = 10
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for a := 1; a <= numJobs; a++ {
        <-results
    }
}

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(100 * time.Millisecond)
        results <- j * 2
    }
}
```

## Best Practices

1. **Use sync primitives** (Mutex, WaitGroup, etc.) for proper synchronization
2. **Avoid goroutine leaks** by ensuring they can exit
3. **Pass variables as parameters** to avoid closure-related issues
4. **Limit concurrency** for resource-intensive operations
5. **Consider using worker pools** for processing many similar tasks
6. **Use context for cancellation** to gracefully stop goroutines