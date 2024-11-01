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
  gen, err := nanoid.New(alphabet, nil) // nil uses crypto/rand as the default
  if err != nil {
    panic(err)
  }

  // Generate a Nano ID
  id, err := gen.Generate(10) // Custom length: 10
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
// Generator defines the interface for generating Nano IDs.
type Generator interface {
    // Generate creates a new Nano ID of the specified length.
    Generate(size int) (string, error)
    
    // MustGenerate returns creates a new Nano ID of the specified length if err
    // is nil or panics otherwise.
    // It simplifies safe initialization of global variables holding compiled UUIDs.
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

### `Config` Struct

Holds the configuration details for the generator.

```go
// Config holds the configuration for the Nano ID generator.
type Config struct {
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

<details>
  <summary>Expand to see results</summary>

```shell
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24           9900825               119.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          3839343               310.7 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          3697426               324.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          3403065               355.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          2721128               443.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         1986505               611.1 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24          9938446               119.8 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         3804894               310.2 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         3759706               318.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         3507075               343.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         2773866               429.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        2020968               591.2 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24         10467040               116.9 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         3967908               303.9 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         3671436               317.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         3470918               345.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         2795991               430.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        2038209               591.9 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24         10072525               116.4 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         3944463               304.5 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         3778465               329.4 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         3322444               346.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         2778516               430.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        2024299               588.7 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          4929242               245.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2445019               489.8 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         2239215               533.0 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         1903724               624.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24         1246633               964.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24         880270              1365 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         7129393               166.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3088347               387.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        2676687               447.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2260320               529.4 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1503513               775.8 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24        918639              1273 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        7141119               167.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3094470               390.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       2659474               449.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2250259               533.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1525422               778.6 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24       939098              1275 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        7075434               168.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3078884               390.8 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       2653651               450.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2250336               534.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1535886               779.7 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24       905656              1273 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        4705986               254.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       2362225               511.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       2024662               589.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       1689366               705.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1000000              1098 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24       647395              1781 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4692651               251.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2358103               510.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       2025705               591.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1684885               707.9 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24       1000000              1102 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       642933              1797 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   2989003               390.5 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1466362               825.1 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1404768               857.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1333918               920.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                  1000000              1026 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                 1000000              1152 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  3085243               396.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1466227               825.1 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1410675               851.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1333658               897.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1000000              1018 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                1000000              1149 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  3039732               394.5 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1443295               821.3 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1413276               852.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1332478               914.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1000000              1029 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                1000000              1155 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  3071056               390.1 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1457083               826.4 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1398894               849.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1346475               886.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1000000              1026 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                1000000              1142 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1430025               834.6 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  740119              1543 ns/op              32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  671708              1677 ns/op              48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  605883              1925 ns/op              64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  443612              2725 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 344932              3489 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2391144               500.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1237058               966.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1000000              1044 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                 985527              1120 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                 954957              1294 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                727497              1577 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2411503               498.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1225498               969.1 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1000000              1041 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1000000              1118 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24                884580              1299 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               734211              1596 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2419585               497.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1240014               974.8 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1000000              1037 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24               1000000              1114 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24                946449              1307 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               719616              1600 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                1448635               831.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24                885998              1527 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24                702578              1714 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24                600270              1952 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24                443122              2687 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               329334              3664 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1452666               827.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                785490              1506 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                704905              1707 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                608680              1966 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                453290              2708 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               322802              3677 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24          10134021               116.7 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          3946090               305.3 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          3794966               313.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          3502370               339.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          2760170               434.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 1978986               604.9 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                 10140744               117.2 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 3913432               306.8 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 3650574               325.3 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 3444016               344.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 2784463               438.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                1968044               625.7 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                  9712027               119.9 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 3830364               314.2 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 3632488               334.7 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 3422242               345.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 2769238               433.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                1976390               604.3 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                 10145660               116.9 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 3978235               303.7 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 3778012               318.5 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 3502041               344.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 2775818               431.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                1977230               603.6 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  4833774               247.3 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2449020               491.1 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 2243475               537.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 1885855               632.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                 1235704               973.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                 900153              1364 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 7047174               170.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3047263               390.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                2634684               461.1 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2223579               541.8 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1492694               801.0 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24                906471              1319 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                6993999               169.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3055308               390.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               2628690               454.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2230983               540.8 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1488814               799.8 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24               897660              1327 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                7030861               171.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3014774               391.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               2646554               457.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2230081               539.6 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1493530               802.2 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24               924261              1326 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                4675761               258.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               2298976               524.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               2016050               592.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               1683283               712.1 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1000000              1139 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24               624578              1836 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4680818               253.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2349648               515.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               2006119               596.1 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1638229               731.4 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24               1000000              1135 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               636212              1885 ns/op             929 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      279.780s
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

