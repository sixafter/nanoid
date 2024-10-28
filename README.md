# NanoID

[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple, fast, and efficient Go implementation of [NanoID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

- **Short & Unique IDs**: Generates compact and collision-resistant identifiers.
- **Cryptographically Secure**: Utilizes Go's crypto/rand package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications.
- **Customizable Alphabet**: Define your own set of characters for ID generation.
- **Concurrency Safe**: Designed to be safe for use in concurrent environments.
- **High Performance**: Optimized with buffer pooling to minimize allocations and enhance speed.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.
- **Optimized for Low Allocations**: Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.

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
package main

import (
  "fmt"
  "github.com/sixafter/nanoid"
)

func main() {
  id, err := nanoid.Generate() 
  if err != nil {
    panic(err)
  }
  fmt.Println("Generated ID:", id)
}
```

### Generating a NanoID with Custom Size

Generate a NanoID with a custom length:

```go
package main

import (
  "fmt"
  "github.com/sixafter/nanoid"
)

func main() {
  id, err := nanoid.GenerateSize(10)
  if err != nil {
    panic(err)
  }
  fmt.Println("Generated ID:", id)
}
```

### Generate a Nano ID with Custom Alphabet

Create a custom generator with a specific alphabet and use it to generate IDs:

```go
package main

import (
  "fmt"
  "github.com/sixafter/nanoid"
)

func main() {
  // Define a custom alphabet
  alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

  // Create a new generator
  generator, err := nanoid.New(alphabet, nil) // nil uses crypto/rand as the default
  if err != nil {
    panic(err)
  }

  // Generate a Nano ID
  id, err := generator.Generate(10) // Custom length: 10
  if err != nil {
    panic(err)
  }
  fmt.Println("Generated ID:", id)
}
```

## Functions

### `Generate`

Generates a Nano ID with the specified length using the default generator.

```go
func Generate(length int) (string, error)
```

* Parameters:
  * `length` (`int`): The desired length of the Nano ID. Must be a positive integer.
* Returns:
  * `string`: The generated Nano ID.
  * `error`: An error if the generation fails.

### `New`

Creates a new Nano ID generator with a custom alphabet and random source.

```go
func New(alphabet string, randReader io.Reader) (Generator, error)
```

* Parameters:
  * `alphabet` (`string`): The set of characters to use for generating IDs. Must not be empty, too short, or contain duplicate characters. 
  * `randReader` (`io.Reader`): The source of randomness. If `nil`, `crypto/rand` is used by default. 
* Returns:
  * `Generator`: A new Nano ID generator. 
  * `error`: An error if the configuration is invalid.

### `Generator` Interface

Defines the method to generate Nano IDs.

```go
type Generator interface {
    Generate(size int) (string, error)
}
```

### `Configuration` Interface

Provides access to the generator's configuration.

```go
type Configuration interface {
    GetConfig() Config
}
```

### `Config` Struct

Holds the configuration details for the generator.

```go
type Config struct {
    Alphabet    []byte
    AlphabetLen int
    Mask        byte
    Step        int
}
```

## Error Handling

The nanoid module defines several error types to handle various failure scenarios:
* `ErrInvalidLength`: Returned when a non-positive length is specified. 
* `ErrExceededMaxAttempts`: Returned when the generation process exceeds the maximum number of attempts. 
* `ErrEmptyAlphabet`: Returned when an empty alphabet is provided. 
* `ErrAlphabetTooShort`: Returned when the alphabet is shorter than required. 
* `ErrAlphabetTooLong`: Returned when the alphabet exceeds the maximum allowed length. 
* `ErrDuplicateCharacters`: Returned when the alphabet contains duplicate characters.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`
* `DefaultSize`: The default size of the generated ID: `21`

## Performance Optimizations

### Buffer Pooling with `sync.Pool`

The nanoid generator utilizes sync.Pool to manage byte slice buffers efficiently. This approach minimizes memory allocations and enhances performance, especially in high-concurrency scenarios.

How It Works:
* Storing Pointers: `sync.Pool` stores pointers to `[]byte` slices (`*[]byte`) instead of the slices themselves. This avoids unnecessary allocations and aligns with best practices for using `sync.Pool`.
* Zeroing Buffers: Before returning buffers to the pool, they are zeroed out to prevent data leaks.

### Struct Optimization

The `generator` struct is optimized for memory alignment and size by:

* Removing Embedded Interfaces: Interfaces like `Generator` and `Configuration` are implemented explicitly without embedding, reducing the struct's size and preventing unnecessary padding. 
* Ordering Fields by Alignment: Fields are ordered from largest to smallest alignment requirements to minimize padding and optimize memory usage.

## Execute Benchmarks:

Run the benchmarks using the go test command with the `bench` make target:

```shell
make bench
```

### Interpreting Results:

Sample output might look like this:

```shell
go test -bench=. -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkGenerateDefault-24                      3985082               300.7 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateCustomAlphabet-24               3429874               346.0 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateShortID-24                      3646383               327.2 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateLongID-24                       2557196               468.1 ns/op           128 B/op          2 allocs/op
BenchmarkGenerateMaxAlphabet-24                  4532246               263.8 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateMinAlphabet-24                  2507995               479.8 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateWithBufferPool-24               3468786               343.9 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDefaultParallel-24              1530394               790.9 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateCustomAlphabetParallel-24       1386268               861.6 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateShortIDParallel-24              1421832               842.7 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateLongIDParallel-24               1000000              1050 ns/op             128 B/op          2 allocs/op
BenchmarkGenerateExtremeConcurrency-24           1530957               785.7 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_5-24    3659472               327.7 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_10-24           3436932               346.0 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_20-24           3140282               381.1 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_50-24           2580222               470.5 ns/op           128 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_100-24          1936257               617.2 ns/op           224 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_2-24        2510594               479.6 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_6-24        3452442               346.3 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_26-24       3901122               308.0 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_38-24       3562468               336.3 ns/op            32 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      34.903s
```

* `ns/op` (Nanoseconds per Operation):
  * Indicates the average time taken per operation. 
  * Lower values signify better CPU performance. 
* `B/op` (Bytes Allocated per Operation):
  * Shows the average number of bytes allocated per operation. 
  * `0 B/op` indicates no heap allocations, which is optimal. 
* `allocs/op` (Allocations per Operation):
  * Represents the average number of memory allocations per operation. 
  * `0 allocs/op` is ideal as it indicates no heap allocations.

## Nano ID Generation

Nano ID generates unique identifiers based on the following:

1. Random Byte Generation: Nano ID generates a sequence of random bytes using a secure random source (e.g., crypto/rand.Reader). 
2. Mapping to Alphabet: Each random byte is mapped to a character in a predefined alphabet to form the final ID. 
3. Uniform Distribution: To ensure that each character in the alphabet has an equal probability of being selected, Nano ID employs techniques to avoid bias, especially when the alphabet size isn't a power of two.

## Contributing

Contributions are welcome. See [CONTRIBUTING](CODE_OF_CONDUCT)

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

