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
- **Customizable Alphabet**: Define your own set of characters for ID generation with a minimum length of 2 characters. ASCII and Unicode are supported with optimizations for both.
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

### Basic Usage with Default Settings

The simplest way to generate a Nano ID is by using the default settings. This utilizes the predefined alphabet and default ID length.

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

**Output**:

```bash
Default Nano ID: -D5f3Z_0x1Gk9Qa
```

### Generating a NanoID with Custom length

Generate a NanoID with a custom length.

```go
package main

import (
  "fmt"
  "github.com/sixafter/nanoid"
)

func main() {
  id, err := nanoid.GenerateWithLength(10)
  if err != nil {
    panic(err)
  }
  fmt.Println("Generated ID:", id)
}
```

**Output**:

```bash
Default Nano ID: 1A3F5B7C9D
```

### Customizing the Alphabet and ID Length

You can customize the alphabet and the default ID length by using the WithAlphabet and WithDefaultLength options.

```go
package main

import (
	"fmt"

	"github.com/sixafter/nanoid"
)

func main() {
	// Define a custom alphabet
	customAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Define a custom default length
	customLength := 10

	// Create a new generator with custom alphabet and default length
	gen, err := nanoid.New(
		nanoid.WithAlphabet(customAlphabet),
		nanoid.WithDefaultLength(customLength),
	)
	if err != nil {
		fmt.Println("Error creating Nano ID generator:", err)
		return
	}

	// Generate a Nano ID using the custom generator
	id, err := gen.Generate(customLength)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Custom Nano ID:", id)
}
```

**Output**"

```bash
Custom Nano ID: G5J8K2M0QZ
```

### Using a Custom Random Reader for Deterministic ID Generation

For testing or deterministic ID generation, you might want to use a custom random reader. Below is an example using a deterministic byte source.

```go
package main

import (
	"fmt"
	"io"
	"sync"

	"github.com/sixafter/nanoid"
)

// cyclicReader is a helper type that cycles through a predefined set of bytes.
// It implements the io.Reader interface.
type cyclicReader struct {
	data []byte
	mu   sync.Mutex
	pos  int
}

// Read fills p with bytes from the cyclicReader's data, cycling back to the start when necessary.
func (r *cyclicReader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.data) == 0 {
		return 0, io.EOF
	}

	n := 0
	for n < len(p) {
		p[n] = r.data[r.pos]
		n++
		r.pos = (r.pos + 1) % len(r.data)
	}

	return n, nil
}

func main() {
	// Define a custom alphabet
	customAlphabet := "ABCD"

	// Define a custom byte sequence for deterministic ID generation
	customBytes := []byte{0, 1, 2, 3} // Maps to 'A', 'B', 'C', 'D'

	// Create a cyclic random reader with the custom bytes
	cycleReader := &cyclicReader{data: customBytes}

	// Initialize the generator with custom alphabet and custom random reader
	gen, err := nanoid.New(
		nanoid.WithAlphabet(customAlphabet),
		nanoid.WithRandReader(cycleReader),
	)
	if err != nil {
		fmt.Println("Error creating Nano ID generator:", err)
		return
	}

	// Generate a Nano ID of length 4
	id, err := gen.Generate(4)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Generated Nano ID (Deterministic, Length 4):", id)

	// Generate another Nano ID of length 4, should cycle through the bytes again
	id2, err := gen.Generate(4)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Generated Nano ID (Deterministic, Length 4):", id2)

	// Generate a Nano ID of length 8, cycling through the bytes twice
	id3, err := gen.Generate(8)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Generated Nano ID (Deterministic, Length 8):", id3)
}
```

**Output**:

```shell
Generated Nano ID (Deterministic, Length 4): ABCD
Generated Nano ID (Deterministic, Length 4): ABCD
Generated Nano ID (Deterministic, Length 8): ABCDABCD
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


### `GenerateWithLength`

Generates a Nano ID with the specified length using the default generator.

```go
func GenerateWithLength(length int) (string, error)
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
// Generator defines the interface for generating Nano IDs.
type Generator interface {
    // Generate returns a new Nano ID of the specified length.
    Generate(length int) (string, error)
    
    // MustGenerate returns a new Nano ID of the specified length if err is nil or panics otherwise.
    MustGenerate(length int) string
}
```

### `Configuration` Interface

Provides access to the generator's configuration.

```go
// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
    // GetConfig returns the configuration of the generator.
    GetConfig() Config
}
```

### `RuntimeConfig` Struct

Holds the configuration details for the generator.

```go
// RuntimeConfig holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type RuntimeConfig struct {
    // Alphabet is a slice of bytes representing the character set used to generate IDs.
    Alphabet []byte
    
    // RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
    RuneAlphabet []rune
    
    // Mask is a bitmask used to obtain a random value from the character set.
    Mask uint
    
    // BitsNeeded represents the number of bits required to generate each character in the ID.
    BitsNeeded uint
    
    // BytesNeeded specifies the number of bytes required from a random source to produce the ID.
    BytesNeeded uint
    
    // BufferSize is the buffer size used for random byte generation.
    BufferSize int
    
    // AlphabetLen is the length of the alphabet, stored as an uint16.
    AlphabetLen uint16
    
    // IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
    IsPowerOfTwo bool
    
    // IsASCII indicates whether the alphabet is ASCII-only, ensuring compatibility with ASCII environments.
    IsASCII bool
}
```

## Error Handling

The nanoid module defines several error types to handle various failure scenarios:
* `ErrDuplicateCharacters`: Returned when the alphabet contains duplicate characters.
* `ErrExceededMaxAttempts`: Returned when the generation process exceeds the maximum number of attempts.
* `ErrInvalidLength`: Returned when a non-positive Nano ID length is specified. 
* `ErrInvalidAlphabet`: Returned when an alphabet is invalid; e.g. due to length constraints.
* `ErrNonUTF8Alphabet`: Returned when an alphabet contains invalid UTF-8 characters.
* `ErrAlphabetTooLong`: Returned when an alphabet length exceeds 256 character.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
* `DefaultLength`: The default length of the generated ID: `21`

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

<details>
  <summary>Expand to see results</summary>

```shell
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24           8287710               143.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          3458908               349.1 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          3118034               385.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          2668833               453.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          1801447               666.0 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         1000000              1052 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24          8245402               144.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         3409947               349.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         3095005               388.4 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         2632291               455.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         1802788               665.3 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        1000000              1047 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24          8208511               145.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         3414762               351.9 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         3106312               394.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         2491836               461.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         1797640               665.4 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        1000000              1049 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24          8189415               150.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         3376741               351.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         3045118               390.2 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         2488202               457.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         1750443               675.7 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        1000000              1064 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          4129323               289.7 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2211400               543.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         2004159               600.0 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         1608355               736.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24         1000000              1206 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24         650894              1825 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         7663106               156.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3271550               368.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        2790944               431.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2361598               507.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1622791               732.8 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24       1000000              1187 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        7617163               157.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3185416               376.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       2752042               437.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2334298               505.6 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1637635               731.9 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24      1000000              1180 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        7661036               157.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3113401               375.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       2801428               426.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2371736               501.8 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1640449               730.1 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24       991594              1184 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        4916432               244.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       2408094               495.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       2108852               570.9 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       1753929               682.8 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1000000              1074 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24       629872              1787 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4896202               243.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2402408               493.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       2107456               570.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1666522               685.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24       1000000              1056 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       687913              1717 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   2679631               447.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1343784               892.1 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1257314               951.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1000000              1023 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                  1000000              1203 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                  802926              1436 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  2578495               463.3 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1317555               903.5 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1252080               955.6 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1000000              1030 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1000000              1191 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                 913140              1443 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  2639306               452.1 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1345076               898.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1263411               955.0 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1000000              1020 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1000000              1186 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                 872174              1433 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  2599221               456.1 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1348374               887.3 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1251062               941.7 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1000000              1021 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1000000              1171 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                 865360              1424 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1315418               914.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  750210              1623 ns/op              88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  727192              1787 ns/op             120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  543541              2136 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  388136              3092 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 301182              3952 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2506696               482.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1294630               925.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1000000              1005 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                1000000              1080 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                 979814              1245 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                822366              1523 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2538175               479.4 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1293229               928.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1000000              1002 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1000000              1070 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24                949635              1257 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               825680              1524 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2459442               482.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1289870               922.8 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1000000              1018 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24               1000000              1068 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24                970238              1260 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               787030              1524 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                1496124               802.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24                755793              1485 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24                696702              1672 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24                635422              1908 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24                468843              2611 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               333549              3574 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1498758               800.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                804652              1529 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                744152              1668 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                631968              1902 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                455257              2609 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               336698              3552 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24           8122065               146.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          3434354               346.5 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          3112950               384.2 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          2634900               455.7 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          1758133               683.0 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 1000000              1078 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                  8283176               145.7 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 3437125               346.7 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 3078442               394.9 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 2597701               460.2 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 1763136               712.3 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                1000000              1086 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                  8244777               145.1 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 3432304               347.3 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 3090642               387.3 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 2600820               458.1 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 1772391               680.6 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                1000000              1085 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                  8314324               144.5 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 3452112               348.6 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 3076152               389.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 2621108               474.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 1761447               686.1 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                1000000              1074 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  4312852               279.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2229373               537.6 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 2003446               604.8 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 1578762               749.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                  957735              1216 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                 627285              1861 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 7527574               158.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3208182               374.1 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                2758483               434.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2320950               519.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1568509               783.7 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24                977128              1280 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                7234887               164.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3107169               387.1 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               2633612               452.2 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2252991               535.9 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1528117               785.9 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24               884773              1282 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                7222017               166.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3107734               386.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               2653728               454.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2233526               523.9 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1598574               749.7 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24               946707              1224 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                4904601               244.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               2415222               496.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               2091091               571.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               1745792               688.8 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1000000              1072 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24               670449              1763 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4887175               243.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2414070               496.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               2105001               569.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1749606               687.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24               1000000              1072 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               681955              1764 ns/op             929 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      272.993s
```
</details>

* `ns/op`: Nanoseconds per operation. Lower values indicate faster performance.
* `B/op`: Bytes allocated per operation. Lower values indicate more memory-efficient code.
* `allocs/op`: Number of memory allocations per operation. Fewer allocations generally lead to better performance.

## Nano ID Generation

Nano ID generates unique identifiers based on the following:

1. **Random Byte Generation**: Nano ID generates a sequence of random bytes using a secure random source (e.g., `crypto/rand.Reader`). 
2. **Mapping to Alphabet**: Each random byte is mapped to a character in a predefined alphabet to form the final ID. 
3. **Uniform Distribution**: To ensure that each character in the alphabet has an equal probability of being selected, Nano ID employs techniques to avoid bias, especially when the alphabet size isn't a power of two.

## Custom Alphabet Constraints

1. Minimum Alphabet Length:
   * At Least Two Unique Characters: The custom alphabet must contain at least two unique characters. An alphabet with fewer than two characters cannot produce IDs with sufficient variability or randomness.
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
* `ErrNonUTF8Alphabet`: Returned if the alphabet contains invalid UTF-8 characters.

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
		} else if errors.Is(err, nanoid.ErrNonUTF8Alphabet) {
			// Handle invalid UTF-8 alphabet
			fmt.Println("Alphabet contains invalid UTF-8 characters.")
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

