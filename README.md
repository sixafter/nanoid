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
BenchmarkDefaultGenerate-24                      3322077               360.7 ns/op           120 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size8-24            3882877               306.0 ns/op            56 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size16-24           3488208               351.4 ns/op            88 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size21-24           3297568               362.2 ns/op           120 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size32-24           2974650               411.5 ns/op           152 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size64-24           2322823               513.3 ns/op           280 B/op          2 allocs/op
BenchmarkDefaultGenerateSize/Size128-24          1548751               771.9 ns/op           536 B/op          2 allocs/op
BenchmarkDefaultGenerateParallel-24              1278450               937.9 ns/op           120 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size8-24            1435849               839.9 ns/op            56 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size16-24           1319145               903.9 ns/op            88 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size21-24           1266843               949.2 ns/op           120 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size32-24           1000000              1010 ns/op             152 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size64-24           1000000              1158 ns/op             280 B/op          2 allocs/op
BenchmarkDefaultGenerateSizeParallel/Size128-24           884023              1445 ns/op             537 B/op          2 allocs/op
BenchmarkGeneratorGenerate-24                            3182769               371.9 ns/op           120 B/op          2 allocs/op
BenchmarkGeneratorGenerateParallel-24                    1255461               945.8 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabet-24                               3274995               372.7 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetParallel-24                       1252706               950.6 ns/op           120 B/op          2 allocs/op
BenchmarkNewGenerator-24                                  365620              3198 ns/op            2416 B/op         14 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength2-24        3184551               368.2 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength16-24       3314415               374.3 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength32-24       3218119               364.4 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength64-24       3199032               368.0 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengths/AlphabetLength95-24       2844274               432.8 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength2-24                1267381               951.0 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength16-24               1269626               942.1 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength32-24               1278676               947.0 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength64-24               1254948               950.0 ns/op           120 B/op          2 allocs/op
BenchmarkCustomAlphabetLengthsParallel/AlphabetLength95-24               1000000              1020 ns/op             120 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      50.181s
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

1. Minimum Alphabet Length:
   * At Least Two Unique Characters: The custom alphabet must contain at least two unique characters. An alphabet with fewer than two characters cannot produce IDs with sufficient variability or randomness. 
   * Maximum Length 256 Characters: The implementation utilizes a byte-based approach where each character in the alphabet is represented by a single byte (`0-255`). This inherently limits the maximum number of unique characters to 256. Attempting to use an alphabet longer than 256 characters will result in an error.
2. Uniqueness of Characters:
   * All Characters Must Be Unique. Duplicate characters in the alphabet can introduce biases in ID generation and compromise the randomness and uniqueness of the IDs. The generator enforces uniqueness by checking for duplicates during initialization. If duplicates are detected, it will return an `ErrDuplicateCharacters` error. 
3. Character Encoding:
   * Support for Unicode: The generator accepts alphabets containing Unicode characters, allowing you to include a wide range of symbols, emojis, or characters from various languages.
4. Power-of-Two Considerations:
   * Mask Calculation: The generator calculates a mask based on the number of bits required to represent the alphabet length minus one.
    ```go
    k := bits.Len(uint(alphabetLen - 1))
    mask := byte((1 << k) - 1)
    ```
   * Implications: While the alphabet length doesn't need to be a power of two, the mask is used to efficiently reduce bias in random number generation. The implementation ensures that each character in the alphabet has an equal probability of being selected by using this mask.

## Error Handling

When initializing the generator with a custom alphabet, the following errors might occur:
* `ErrInvalidAlphabet`: Returned if the alphabet length is less than 2. 
* `ErrDuplicateCharacters`: Returned if duplicate characters are found in the alphabet.

Example of Handling Errors:

```go
package main

import (
	"errors"
	"fmt"
	"github.com/sixafter/nanoid"
)

func main() {
	// Define a custom alphabet (e.g., lowercase letters and digits)
	customAlphabet := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Create a new generator with the custom alphabet
	generator, err := nanoid.New(customAlphabet, nil)
	if err != nil {
		if errors.Is(err, nanoid.ErrInvalidAlphabet) {
			// Handle invalid alphabet length
			fmt.Println("Alphabet must contain at least two unique characters.")
		} else if errors.Is(err, nanoid.ErrDuplicateCharacters) {
			// Handle duplicate characters in the alphabet
			fmt.Println("Alphabet contains duplicate characters.")
		} else {
			// Handle other potential errors
			fmt.Println("Error initializing Nano ID generator:", err)
		}
		return // Exit if generator initialization fails
	}

	// Generate a Nano ID of length 15
	id, err := generator.Generate(15)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Generated ID:", id)
}
```

Output:

```bash
Generated ID: k5f3z8n2q1w9b0d
```

## Determining Collisions

To determine the practical length for a NanoID, see the collision time calculator [here](https://sixafter.github.io/nanoid/).

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

