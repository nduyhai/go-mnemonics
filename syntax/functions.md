# Go Functions

## Basic Function Declaration

```go
func greet(name string) string {
    return "Hello, " + name
}
```

## Multiple Return Values

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}
```

## Named Return Values

```go
func rectangle(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return // naked return uses named return values
}
```

## Variadic Functions

```go
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}
```

## Anonymous Functions

```go
func main() {
    // Anonymous function assigned to a variable
    add := func(a, b int) int {
        return a + b
    }
    
    fmt.Println(add(3, 4)) // 7
    
    // Immediately invoked function expression (IIFE)
    result := func(a, b int) int {
        return a * b
    }(5, 6)
    
    fmt.Println(result) // 30
}
```

## Higher-Order Functions

```go
func applyOperation(a, b int, operation func(int, int) int) int {
    return operation(a, b)
}

func main() {
    multiply := func(x, y int) int {
        return x * y
    }
    
    result := applyOperation(4, 5, multiply)
    fmt.Println(result) // 20
}
```

## Closures

```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    increment := counter()
    fmt.Println(increment()) // 1
    fmt.Println(increment()) // 2
    fmt.Println(increment()) // 3
}
```

## Defer Statement

```go
func readFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // Will be executed when the function returns
    
    // Read file contents...
    return nil
}
```

## Panic and Recover

```go
func safeOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    
    // This will cause a panic
    panic("something went wrong")
}
```