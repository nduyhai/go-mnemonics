# Go Variables

## Basic Variable Declaration

```go
var name string = "John"
var age int = 30
var isActive bool = true
```

## Short Variable Declaration

```go
name := "John"
age := 30
isActive := true
```

## Zero Values

Go assigns a default "zero value" to variables declared without an explicit initial value:

- `0` for numeric types
- `false` for boolean types
- `""` (empty string) for strings
- `nil` for interfaces, slices, channels, maps, pointers and functions

## Multiple Variable Declaration

```go
var a, b, c int = 1, 2, 3
x, y, z := 1, "hello", true
```

## Constants

```go
const Pi = 3.14159
const (
    StatusOK = 200
    StatusCreated = 201
    StatusAccepted = 202
)
```

## Iota

```go
const (
    Monday = iota // 0
    Tuesday       // 1
    Wednesday     // 2
    Thursday      // 3
    Friday        // 4
)
```