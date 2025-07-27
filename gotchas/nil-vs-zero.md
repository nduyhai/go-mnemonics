# Nil vs. Zero Values in Go

Go has a clear distinction between nil and zero values that can cause confusion for newcomers.

## Zero Values

Every type in Go has a zero value, which is the default value a variable of that type gets when it's declared without an explicit initializer.

```go
var i int       // i = 0
var f float64   // f = 0.0
var b bool      // b = false
var s string    // s = ""
var p *int      // p = nil
```

## Nil in Go

`nil` is a predeclared identifier in Go that represents the zero value for:
- Pointers
- Interfaces
- Maps
- Slices
- Channels
- Function types

```go
var p *int      // p = nil
var i interface{} // i = nil
var m map[string]int // m = nil
var s []int     // s = nil
var c chan int  // c = nil
var f func()    // f = nil
```

## Common Gotchas

### 1. Nil vs. Empty Slice

```go
// A nil slice
var s1 []int // s1 == nil evaluates to true

// An empty slice (not nil)
s2 := []int{} // s2 == nil evaluates to false

// Both have length 0
fmt.Println(len(s1)) // 0
fmt.Println(len(s2)) // 0

// Both can be appended to
s1 = append(s1, 1)
s2 = append(s2, 1)
```

### 2. Nil vs. Empty Map

```go
// A nil map
var m1 map[string]int // m1 == nil evaluates to true

// Reading from a nil map returns zero values
fmt.Println(m1["key"]) // 0, no panic

// Writing to a nil map causes a panic
// m1["key"] = 1 // panic: assignment to entry in nil map

// An empty map (not nil)
m2 := map[string]int{} // m2 == nil evaluates to false

// Writing to an empty map is fine
m2["key"] = 1 // Works fine
```

### 3. Nil Interface Values

```go
// A nil interface value has both type and value as nil
var i1 interface{} // i1 == nil evaluates to true

// An interface containing a nil pointer is not itself nil
var p *int = nil
var i2 interface{} = p // i2 == nil evaluates to false

// This distinction matters in error handling
func returnsError() error {
    var p *MyError = nil
    if condition {
        return p // This returns a non-nil error even though p is nil!
    }
    return nil // This returns a nil error
}
```

### 4. Comparing with Nil

```go
// Correct way to check if an interface value is nil
func isNil(v interface{}) bool {
    return v == nil
}

// Incorrect way to check if a value inside an interface is nil
func isValueNil(v interface{}) bool {
    // This will not work as expected for nil pointers inside interfaces
    return reflect.ValueOf(v).IsNil()
}
```

### 5. Nil Receivers

```go
type MyStruct struct {
    Value int
}

// Method with pointer receiver
func (m *MyStruct) PointerMethod() {
    if m == nil {
        fmt.Println("Receiver is nil")
        return
    }
    fmt.Println("Value:", m.Value)
}

func main() {
    var m *MyStruct = nil
    
    // This works! Methods can be called on nil receivers
    m.PointerMethod() // Prints: "Receiver is nil"
    
    // But accessing fields of a nil struct will panic
    // fmt.Println(m.Value) // panic: runtime error: invalid memory address or nil pointer dereference
}
```

## Best Practices

1. **Be explicit about nil vs. empty collections**
   ```go
   // Prefer this for an empty slice
   s := []int{}
   
   // Over this, which is nil
   var s []int
   ```

2. **Always initialize maps before use**
   ```go
   m := make(map[string]int)
   // or
   m := map[string]int{}
   ```

3. **Check for nil before using pointers**
   ```go
   if p != nil {
       // Safe to dereference
       value := *p
   }
   ```

4. **Return explicit nil for interfaces**
   ```go
   func getError() error {
       // Don't return a nil *MyError, return a nil error
       return nil
   }
   ```

5. **Use type assertions carefully with nil interfaces**
   ```go
   var i interface{} = nil
   
   // This will panic
   // s := i.(string)
   
   // This is safe
   s, ok := i.(string)
   if !ok {
       // Handle the case where i is not a string
   }
   ```