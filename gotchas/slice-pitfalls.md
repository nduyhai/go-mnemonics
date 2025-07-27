# Go Slice Pitfalls

Slices are one of Go's most powerful and frequently used data structures, but they come with several non-obvious behaviors that can lead to bugs.

## Slice Basics

A slice in Go is a reference to a segment of an underlying array. It consists of:
- A pointer to the first element
- A length (accessible via `len()`)
- A capacity (accessible via `cap()`)

```go
// Creating slices
s1 := []int{1, 2, 3}
s2 := make([]int, 3)      // len=3, cap=3
s3 := make([]int, 3, 5)   // len=3, cap=5
```

## Common Pitfalls

### 1. Unexpected Sharing of Underlying Arrays

```go
original := []int{1, 2, 3, 4, 5}
slice1 := original[1:3]   // [2, 3]
slice2 := original[2:4]   // [3, 4]

// Modifying one slice affects others that share the same underlying array
slice1[1] = 10            // Changes slice1[1] and slice2[0]

fmt.Println(original)     // [1, 2, 10, 4, 5]
fmt.Println(slice1)       // [2, 10]
fmt.Println(slice2)       // [10, 4]
```

### 2. Append Behavior and Capacity

```go
s := []int{1, 2, 3}
fmt.Println(len(s), cap(s))  // 3, 3

// Append when capacity is reached creates a new underlying array
s = append(s, 4)
fmt.Println(len(s), cap(s))  // 4, 6 (capacity doubled)

// This can lead to unexpected behavior
original := []int{1, 2, 3, 4, 5}
slice := original[1:3]       // [2, 3]
fmt.Println(len(slice), cap(slice))  // 2, 4

// Append within capacity modifies the original array
slice = append(slice, 10)
fmt.Println(original)        // [1, 2, 3, 10, 5]

// Append beyond capacity allocates a new array
slice = append(slice, 20, 30, 40)
fmt.Println(original)        // [1, 2, 3, 10, 5] (unchanged)
fmt.Println(slice)           // [2, 3, 10, 20, 30, 40]
```

### 3. Slice as Function Parameters

```go
func modify(s []int) {
    s[0] = 999            // This will affect the original slice
    s = append(s, 888)    // This won't affect the original slice
}

func main() {
    slice := []int{1, 2, 3}
    modify(slice)
    fmt.Println(slice)    // [999, 2, 3] (not [999, 2, 3, 888])
}
```

### 4. Slicing Beyond Capacity

```go
s := []int{1, 2, 3, 4, 5}
fmt.Println(len(s), cap(s))  // 5, 5

// This works fine
slice1 := s[1:3]             // [2, 3]

// This will panic: slice bounds out of range
// slice2 := s[1:10]

// But this works (up to capacity)
slice3 := s[1:5]             // [2, 3, 4, 5]

// And this is a runtime panic (beyond capacity)
// slice4 := s[1:6]
```

### 5. Empty vs. Nil Slices

```go
var nilSlice []int           // nil slice
emptySlice := []int{}        // empty slice
zeroLenSlice := make([]int, 0) // also empty slice

fmt.Println(nilSlice == nil)       // true
fmt.Println(emptySlice == nil)     // false
fmt.Println(zeroLenSlice == nil)   // false

// All have length 0
fmt.Println(len(nilSlice))         // 0
fmt.Println(len(emptySlice))       // 0
fmt.Println(len(zeroLenSlice))     // 0

// All can be appended to
nilSlice = append(nilSlice, 1)
emptySlice = append(emptySlice, 1)
zeroLenSlice = append(zeroLenSlice, 1)
```

### 6. Slice Copy Behavior

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)  // len=3

// Copy only copies up to the minimum length of both slices
copied := copy(dst, src)
fmt.Println(copied)    // 3
fmt.Println(dst)       // [1, 2, 3]

// To copy the entire slice, ensure destination has sufficient length
dst = make([]int, len(src))
copy(dst, src)
fmt.Println(dst)       // [1, 2, 3, 4, 5]
```

### 7. Memory Leaks with Large Slices

```go
func getSubset(large []int) []int {
    // PROBLEM: This keeps the entire original array in memory
    return large[0:10]
    
    // SOLUTION: Copy to a new slice to allow GC of the original
    // subset := make([]int, 10)
    // copy(subset, large[:10])
    // return subset
}
```

## Best Practices

1. **Use `copy()` to create independent slices**
   ```go
   original := []int{1, 2, 3, 4, 5}
   // Create a completely independent copy
   independent := make([]int, len(original))
   copy(independent, original)
   ```

2. **Be aware of capacity when appending**
   ```go
   // Force a new backing array by setting capacity to exact length
   s := make([]int, len(original), len(original))
   copy(s, original)
   s = append(s, newElement) // Guaranteed to create a new backing array
   ```

3. **Use full slice expressions to limit capacity**
   ```go
   // Third index limits the capacity
   limited := original[1:3:3] // [start:end:cap]
   fmt.Println(len(limited), cap(limited)) // 2, 2
   ```

4. **Pre-allocate slices when the size is known**
   ```go
   // More efficient when size is known
   s := make([]int, 0, 1000)
   for i := 0; i < 1000; i++ {
       s = append(s, i)
   }
   ```

5. **Be careful with large slices in long-lived objects**
   ```go
   type LargeData struct {
       // Storing a small slice of a large array can prevent GC
       data []byte
   }
   
   // Better to copy only what's needed
   func NewLargeData(source []byte) *LargeData {
       copied := make([]byte, len(source))
       copy(copied, source)
       return &LargeData{data: copied}
   }
   ```