# Go Channels

Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values in another goroutine.

## Basic Channel Operations

```go
func main() {
    // Create a new channel
    ch := make(chan string)
    
    // Send a value into a channel (from a goroutine)
    go func() {
        ch <- "hello"
    }()
    
    // Receive a value from a channel
    msg := <-ch
    fmt.Println(msg) // "hello"
}
```

## Buffered Channels

```go
func main() {
    // Create a buffered channel with capacity 2
    ch := make(chan string, 2)
    
    // These sends won't block because the buffer has capacity
    ch <- "hello"
    ch <- "world"
    
    // Receiving values
    fmt.Println(<-ch) // "hello"
    fmt.Println(<-ch) // "world"
}
```

## Channel Direction

```go
// Function that only receives from a channel
func receive(ch <-chan int) {
    val := <-ch
    fmt.Println("Received:", val)
}

// Function that only sends to a channel
func send(ch chan<- int) {
    ch <- 42
}

func main() {
    ch := make(chan int)
    
    go send(ch)
    receive(ch)
}
```

## Closing Channels

```go
func main() {
    ch := make(chan int, 3)
    
    // Send values
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Close the channel
    close(ch)
    
    // Read until channel is closed
    for val := range ch {
        fmt.Println(val)
    }
    
    // Check if channel is closed
    val, ok := <-ch
    if !ok {
        fmt.Println("Channel is closed")
    }
}
```

## Select Statement

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "one"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "two"
    }()
    
    // Wait on multiple channel operations
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received", msg2)
        }
    }
}
```

## Timeout with Select

```go
func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case res := <-ch:
        fmt.Println("Received:", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

## Non-blocking Channel Operations

```go
func main() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println("Received message:", msg)
    default:
        fmt.Println("No message received")
    }
    
    select {
    case ch <- "hello":
        fmt.Println("Sent message")
    default:
        fmt.Println("No message sent")
    }
}
```

## Fan-out Pattern

```go
func main() {
    // Source channel
    src := make(chan int)
    
    // Create multiple destination channels
    dests := make([]chan int, 3)
    for i := range dests {
        dests[i] = make(chan int)
    }
    
    // Fan-out: distribute work to multiple workers
    for i := range dests {
        go func(i int, dest chan int) {
            for val := range src {
                dest <- val * (i + 1)
            }
            close(dest)
        }(i, dests[i])
    }
    
    // Send values to source
    go func() {
        for i := 1; i <= 3; i++ {
            src <- i
        }
        close(src)
    }()
    
    // Collect results from all destinations
    for i, dest := range dests {
        for val := range dest {
            fmt.Printf("Worker %d: %d\n", i+1, val)
        }
    }
}
```

## Fan-in Pattern

```go
func main() {
    // Create multiple source channels
    srcs := make([]chan int, 3)
    for i := range srcs {
        srcs[i] = make(chan int)
    }
    
    // Create destination channel
    dest := make(chan int)
    
    // Fan-in: combine multiple inputs into a single channel
    var wg sync.WaitGroup
    wg.Add(len(srcs))
    
    for i := range srcs {
        go func(i int, src chan int) {
            defer wg.Done()
            for j := 1; j <= 3; j++ {
                src <- j * (i + 1)
            }
            close(src)
        }(i, srcs[i])
    }
    
    // Merge all sources into destination
    go func() {
        for i, src := range srcs {
            for val := range src {
                dest <- val
                fmt.Printf("Source %d sent: %d\n", i+1, val)
            }
        }
        wg.Wait()
        close(dest)
    }()
    
    // Receive all values from destination
    for val := range dest {
        fmt.Println("Received:", val)
    }
}
```