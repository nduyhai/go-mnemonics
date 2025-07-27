# Go Interfaces

## Basic Interface Declaration

```go
type Writer interface {
    Write([]byte) (int, error)
}
```

## Implementing Interfaces

```go
// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

## Interface Usage

```go
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    c := Circle{Radius: 5}
    r := Rectangle{Width: 3, Height: 4}
    
    PrintShapeInfo(c) // Works with Circle
    PrintShapeInfo(r) // Works with Rectangle
}
```

## Empty Interface

```go
// The empty interface can hold values of any type
func PrintAny(v interface{}) {
    fmt.Println(v)
}

func main() {
    PrintAny(42)
    PrintAny("hello")
    PrintAny(true)
    PrintAny([]string{"a", "b", "c"})
}
```

## Type Assertions

```go
func process(i interface{}) {
    // Type assertion
    str, ok := i.(string)
    if ok {
        fmt.Println("String value:", str)
        return
    }
    
    // Type switch
    switch v := i.(type) {
    case int:
        fmt.Println("Integer:", v)
    case bool:
        fmt.Println("Boolean:", v)
    case []string:
        fmt.Println("String slice:", v)
    default:
        fmt.Println("Unknown type")
    }
}
```

## Interface Composition

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Composed interface
type ReadWriter interface {
    Reader
    Writer
}
```

## Interface Values

```go
type Stringer interface {
    String() string
}

type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s, age %d", p.Name, p.Age)
}

func main() {
    var s Stringer
    p := Person{"Alice", 30}
    
    // Interface value contains both type information and value
    s = p
    fmt.Println(s.String()) // "Alice, age 30"
    
    // nil interface value
    var nilInterface Stringer // nil interface (both type and value are nil)
    fmt.Println(nilInterface == nil) // true
    
    // Interface with nil value but non-nil type
    var nilPerson *Person
    s = nilPerson // s is not nil, but its concrete value is nil
    fmt.Println(s == nil) // false
}
```