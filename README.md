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

### Generate a Nano ID with Default Settings

Generate a Nano ID using the default size (21 characters) and default alphabet:

```go
id, err := nanoid.New()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID:", id)
```

### Generating a NanoID with Custom Size

Generate a NanoID with a custom length:

```go
id, err := nanoid.NewSize(32)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID of size 32:", id)
```

### Generate a Nano ID with Custom Alphabet

Generate a Nano ID using a custom alphabet:

```go
customAlphabet := "abcdef123456"
id, err := nanoid.NewCustom(16, customAlphabet)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID with custom alphabet:", id)
```

### Generate a Nano ID with Unicode Alphabet

Generate a Nano ID using a Unicode alphabet:

```go
unicodeAlphabet := "あいうえお漢字🙂🚀"
id, err := nanoid.NewCustom(10, unicodeAlphabet)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID with Unicode alphabet:", id)
```

### Generate a Nano ID with Custom Random Source

Generate a Nano ID using a custom random source that implements io.Reader:

```go
// Example custom random source (for demonstration purposes)
var myRandomSource io.Reader = myCustomRandomReader{}

id, err := nanoid.NewCustomReader(21, nanoid.DefaultAlphabet, myRandomSource)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID with custom random source:", id)
```

**Note:** Replace `myCustomRandomReader{}` with your actual implementation of `io.Reader`.

## Thread Safety

All functions provided by this package are safe for concurrent use by multiple goroutines. Here's an example of generating Nano IDs concurrently:

```go
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/sixafter/nanoid"
)

func main() {
	const numGoroutines = 10
	const idSize = 21

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			id, err := nanoid.New()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Generated Nano ID:", id)
		}()
	}

	wg.Wait()
}
```

## Functions

* `func New() (string, error)`: Generates a Nano ID with the default size (21 characters) and default alphabet.
* `func NewSize(size int) (string, error)`: Generates a Nano ID with a specified size using the default alphabet.
* `func NewCustom(size int, alphabet string) (string, error)`: Generates a Nano ID with a specified size and custom alphabet.
* `func NewCustomReader(size int, alphabet string, rnd io.Reader) (string, error)`: Generates a Nano ID with a specified size, custom alphabet, and custom random source.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`
* `DefaultSize`: The default size of the generated ID: `21`

## Unicode Support

This implementation fully supports custom alphabets containing Unicode characters, including emojis and characters from various languages. By using []rune internally, it correctly handles multi-byte Unicode characters.

## Performance

The package is optimized for performance and low memory consumption:
* **Efficient Random Byte Consumption**: Uses bitwise operations to extract random bits efficiently. 
* **Avoids `math/big`**: Does not use `math/big`, relying on built-in integer types for calculations. 
* **Minimized System Calls**: Reads random bytes in batches to reduce the number of system calls.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

* Fork the repository. 
* Create a new branch for your feature or bugfix. 
* Write tests for your changes. 
* Ensure all tests pass. 
* Submit a pull request.

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/).

