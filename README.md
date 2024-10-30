# NanoID

[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple, fast, and efficient Go implementation of [NanoID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

- **Short & Unique IDs**: Generates compact and collision-resistant identifiers.
- **Cryptographically Secure**: Utilizes Go's crypto/rand package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications.
- **Customizable Alphabet**: Define your own set of characters for ID generation with a minimum length of 2 characters and a maximum length of 256 characters.
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
* `ErrInvalidAlphabet`: Returned when an empty alphabet is invalid.
* `ErrDuplicateCharacters`: Returned when the alphabet contains duplicate characters.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`
* `DefaultSize`: The default size of the generated ID: `21`

## Performance Optimizations

### Buffer Pooling with `sync.Pool`

The nanoid generator utilizes `sync.Pool` to manage byte slice buffers efficiently. This approach minimizes memory allocations and enhances performance, especially in high-concurrency scenarios.

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
BenchmarkGenerateDefault-24                      4318624               277.3 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateCustomAlphabet-24               4156414               288.6 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateShortID-24                      4416091               271.3 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateLongID-24                       2899802               418.8 ns/op           128 B/op          2 allocs/op
BenchmarkGenerateMaxAlphabet-24                  4463840               266.2 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateMinAlphabet-24                  4496391               269.1 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateWithBufferPool-24               3953864               288.0 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDefaultParallel-24              1730868               695.0 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateCustomAlphabetParallel-24       1692622               727.6 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateShortIDParallel-24              1856241               674.7 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateLongIDParallel-24               1297450               931.4 ns/op           128 B/op          2 allocs/op
BenchmarkGenerateExtremeConcurrency-24           1715571               685.2 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_5-24    4389192               270.5 ns/op            10 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_10-24           4158094               288.1 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_20-24           3615992               326.2 ns/op            48 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_50-24           2885797               494.5 ns/op           128 B/op          2 allocs/op
BenchmarkGenerateDifferentLengths/Length_100-24          1597360               752.7 ns/op           224 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_2-24        4523349               264.1 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_6-24        4158720               289.2 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_26-24       4235298               284.6 ns/op            32 B/op          2 allocs/op
BenchmarkGenerateDifferentAlphabets/Alphabet_38-24       3745124               327.4 ns/op            32 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      34.773s
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

### Custom Alphabet Constraints

* **Alphabet Length**: Must be between 2 and 256 unique single-byte characters. 
* **Uniqueness**: All characters in the alphabet must be unique. 
* **Character Encoding**: Only single-byte characters (byte) are supported. 
* **Error Handling**: The generator will return specific errors if the alphabet doesn't meet the constraints.

1. Length Requirements:
   * Minimum Length 2 Characters: An alphabet with fewer than two characters cannot produce diverse or secure IDs. At least two unique characters are necessary to generate a variety of IDs. 
   * Maximum Length 256 Characters: The implementation utilizes a byte-based approach where each character in the alphabet is represented by a single byte (`0-255`). This inherently limits the maximum number of unique characters to 256. Attempting to use an alphabet longer than 256 characters will result in an error.
2. Uniqueness of Characters:
   * All Characters Must Be Unique. Duplicate characters in the alphabet can introduce biases in ID generation and compromise the randomness and uniqueness of the IDs. The generator enforces uniqueness by checking for duplicates during initialization. If duplicates are detected, it will return an `ErrDuplicateCharacters` error. 
3. Character Encoding:
   * Single-Byte Characters Only: The implementation is designed to work with single-byte (`byte`) characters, which correspond to values `0-255`. Using multi-byte characters (such as UTF-8 characters beyond the basic ASCII set) can lead to unexpected behavior and is not supported. 
   * Recommended Character Sets:
     * URL-Friendly Characters: Typically, alphanumeric characters (`A-Z`, `a-z`, `0-9`) along with symbols like `-` and `_` are used to ensure that generated IDs are safe for use in URLs and file systems. 
     * Custom Sets: You can define your own set of unique single-byte characters based on your application's requirements.
4. Power-of-Two Considerations:
   * Mask Calculation: The generator calculates a mask based on the number of bits required to represent the alphabet length minus one.
    ```go
    k := bits.Len(uint(alphabetLen - 1))
    mask := byte((1 << k) - 1)
    ```
   * Implications: While the alphabet length doesn't need to be a power of two, the mask is used to efficiently reduce bias in random number generation. The implementation ensures that each character in the alphabet has an equal probability of being selected by using this mask.

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

