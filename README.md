# Nano ID

[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid)](https://github.com/sixafter/nanoid/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple, fast, and efficient Go implementation of [Nano ID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

- **Short & Unique IDs**: Generates compact and collision-resistant identifiers.
- **Cryptographically Secure**: Utilizes Go's `crypto/rand` package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications.
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
func Generate() (string, error)
```

* Returns:
  * `string`: The generated Nano ID.
  * `error`: An error if the generation fails.


### `GenerateSize`

Generates a Nano ID with the specified length using the default generator.

```go
func GenerateSize(length int) (string, error)
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
* `ErrInvalidAlphabet`: Returned when an alphabet is invalid.
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

The `generator` struct is optimized for memory alignment and size by ordering from largest to smallest to minimize padding and optimize memory usage.

## Execute Benchmarks:

Run the benchmarks using the go test command with the `bench` make target:

```shell
make bench
```

### Interpreting Results:

Sample output might look like this:

```shell
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkDefaultGenerate-24                      4362076               274.2 ns/op            48 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size8-24            4553710               258.6 ns/op            16 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size16-24           4399778               274.1 ns/op            32 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size21-24           4256088               285.9 ns/op            48 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size32-24           4204159               284.3 ns/op            64 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size64-24           3791118               322.7 ns/op           128 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size128-24          3157098               377.9 ns/op           256 B/op          2 allocs/op
BenchmarkDefaultGenerateParallel-24              1516693               828.8 ns/op            48 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size8-24            1560518               786.8 ns/op            16 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size16-24           1364043               860.7 ns/op            32 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size21-24           1499124               786.4 ns/op            48 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size32-24           1490713               814.6 ns/op            64 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size64-24           1375413               861.1 ns/op           128 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size128-24          1220480               990.9 ns/op           256 B/op          2 allocs/op
BenchmarkGeneratorGenerate-24                            4365238               276.7 ns/op            48 B/op          2 allocs/op
BenchmarkGeneratorGenerateParallel-24                    1507068               799.6 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabet-24                               4211919               285.5 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetParallel-24                       1513590               803.8 ns/op            48 B/op          2 allocs/op
BenchmarkNewGenerator-24                                 6683338               177.0 ns/op           176 B/op          3 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength2-24        4331044               276.9 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength16-24       4350872               285.4 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength32-24       4273230               287.2 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength64-24       3908989               291.1 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength95-24       3480422               346.8 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength2-24                1496703               804.5 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength16-24               1511534               848.3 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength32-24               1507296               799.4 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength64-24               1520260               800.6 ns/op            48 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength95-24               1337623               885.6 ns/op            48 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      51.470s
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

1. **Random Byte Generation**: Nano ID generates a sequence of random bytes using a secure random source (e.g., `crypto/rand.Reader`). 
2. **Mapping to Alphabet**: Each random byte is mapped to a character in a predefined alphabet to form the final ID. 
3. **Uniform Distribution**: To ensure that each character in the alphabet has an equal probability of being selected, Nano ID employs techniques to avoid bias, especially when the alphabet size isn't a power of two.

## Custom Alphabet Constraints

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
  
## Determining Collisions

To determine the practical length for a NanoID, see the collision time calculator [here](https://sixafter.github.io/nanoid/).

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

