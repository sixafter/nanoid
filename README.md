# NanoID

[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple, fast, and efficient Go implementation of [NanoID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

- **Secure**: Uses `crypto/rand` for cryptographically secure random number generation.
- **Fast**: Optimized for performance with efficient algorithms.
- **Thread-Safe**: Safe for concurrent use in multi-threaded applications.
- **Customizable**: Specify custom ID lengths and alphabets.
- **Easy to Use**: Simple API with sensible defaults.

## Installation

To install the package, use:

```sh
go get -u github.com/sixafter/nanoid
```

## Importing the Package

To use the NanoID package in your Go project, import it as follows:

```go
import "github.com/sixafter/nanoid"
```

## Usage

### Generating a Default NanoID

Generate a NanoID using the default size (21 characters) and the default alphabet (numbers and uppercase/lowercase letters):

```go
package main

import (
    "fmt"
    "log"

    "github.com/sixafter/nanoid"
)

func main() {
    id, err := nanoid.Generate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated NanoID:", id)
}
```

### Generating a NanoID with Custom Size

Generate a NanoID with a custom length:

```go
package main

import (
    "fmt"
    "log"

    "github.com/sixafter/nanoid"
)

func main() {
    id, err := nanoid.GenerateSize(10) // Generate a 10-character NanoID
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("NanoID with custom size:", id)
}
```

### Generating a NanoID with Custom Alphabet

Generate a NanoID with a custom length and a custom set of characters:

```go
package main

import (
    "fmt"
    "log"

    "github.com/sixafter/nanoid"
)

func main() {
    alphabet := "0123456789abcdef" // Hexadecimal characters
    id, err := nanoid.GenerateCustom(16, alphabet) // Generate a 16-character NanoID
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("NanoID with custom alphabet:", id)
}
```

### Concurrency and Thread Safety

The NanoID functions are designed to be thread-safe. You can safely generate IDs from multiple goroutines concurrently without additional synchronization.

```go
package main

import (
    "fmt"
    "sync"

    "github.com/sixafter/nanoid"
)

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            id, err := nanoid.Generate()
            if err != nil {
                fmt.Println("Error generating NanoID:", err)
                return
            }
            fmt.Println("Generated NanoID:", id)
        }()
    }
    wg.Wait()
}
```

### Error Handling

All functions return an error as the second return value. Ensure you handle any potential errors:

```go
id, err := nanoid.Generate()
if err != nil {
    // Handle the error
}
```
