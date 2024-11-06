# nanoid <img src="https://ai.github.io/nanoid/logo.svg" align="right" alt="Nano ID logo by Anton Lovchikov" width="160" height="94">

[![CI](https://github.com/sixafter/nanoid/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid)](https://github.com/sixafter/nanoid/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)

A simple, fast, and efficient Go implementation of [Nano ID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

- **Short & Unique IDs**: Generates compact and collision-resistant identifiers.
- **Cryptographically Secure**: Utilizes Go's `crypto/rand` package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications.
- **Customizable**: 
  - Define your own set of characters for ID generation with a minimum length of 2 characters and maximum length of 256 characters.
  - Define your own random number generator.
  - Unicode and ASCII alphabets supported.
- **Concurrency Safe**: Designed to be safe for use in concurrent environments.
- **High Performance**: Optimized with buffer pooling to minimize allocations and enhance speed.
- **Optimized for Low Allocations**: Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.
    - 1 `allocs/op` for ASCII and Unicode alphabets.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.

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
  id, err := nanoid.New() 
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

### Generating a Nano ID with Custom length

Generate a NanoID with a custom length.

```go
package main

import (
  "fmt"
  "github.com/sixafter/nanoid"
)

func main() {
  id, err := nanoid.NewWithLength(10)
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

You can customize the alphabet by using the WithAlphabet option and generate an ID with a custom length.

```go
package main

import (
	"fmt"

	"github.com/sixafter/nanoid"
)

func main() {
	// Define a custom alphabet
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a new generator with custom alphabet and length hint
	gen, err := nanoid.NewGenerator(
		nanoid.WithAlphabet(alphabet),
		nanoid.WithLengthHint(10),
	)
	if err != nil {
		fmt.Println("Error creating Nano ID generator:", err)
		return
	}

	// Generate a Nano ID using the custom generator
	id, err := gen.New(10)
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

## Functions

### `New`

Generates a Nano ID with the specified length using the default generator.

```go
func New() (string, error)
```

* Returns:
    * `string`: The generated Nano ID.
    * `error`: An error if the generation fails.

### `NewWithLength`

Generates a Nano ID with the specified length using the default generator.

```go
func NewWithLength(length int) (string, error)
```

* Parameters:
    * `length` (`int`): The desired length of the Nano ID. Must be a positive integer.
* Returns:
    * `string`: The generated Nano ID.
    * `error`: An error if the generation fails.

### `Must`

Generates a Nano ID with the specified length using the default generator.

```go
func Must() string
```

* Returns:
    * `string`: The generated Nano ID.

### `MustWithLength`

Generates a Nano ID with the specified length using the default generator.

```go
func NewWithLength(length int) string
```

* Parameters:
    * `length` (`int`): The desired length of the Nano ID. Must be a positive integer.
* Returns:
    * `string`: The generated Nano ID.

### `NewGenerator`

Creates a new Nano ID generator with a custom alphabet and random source.

```go
func NewGenerator(options ...Option) (Generator, error)
```

* Parameters:
  * `options` (`Option`): Variadic Option parameters to configure the `Generator`.  Options are:
    * `WithAlphabet`: Sets a custom alphabet for the `Generator`, allowing the user to specify which characters will be used in the generated IDs.
    * `WithRandReader`: Sets a custom random reader for the `Generator`, enabling the use of a specific source of randomness.
    * `WithLengthHint`: Sets a hint for the intended length of the generated IDs, helping to optimize internal allocations based on the expected ID size.
* Returns:
  * `Generator`: A new Nano ID generator. 
  * `error`: An error if the configuration is invalid.

### `Generator` Interface

Defines the method to generate Nano IDs.

```go
// Generator defines the interface for generating Nano IDs.
type Generator interface {
    // New returns a new Nano ID of the specified length.
	New(length int) (string, error)
}
```

### `Configuration` Interface

Provides access to the generator's configuration.

```go
// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
    // Config returns the configuration of the generator.
	Config() Config
}
```

## Error Handling

The nanoid module defines several error types to handle various failure scenarios:
* `ErrDuplicateCharacters`: Returned when the alphabet contains duplicate characters.
* `ErrExceededMaxAttempts`: Returned when the generation process exceeds the maximum number of attempts.
* `ErrInvalidLength`: Returned when a non-positive Nano ID length is specified. 
* `ErrInvalidAlphabet`: Returned when an alphabet is invalid; e.g. due to length constraints.
* `ErrNonUTF8Alphabet`: Returned when an alphabet contains invalid UTF-8 characters.
* `ErrAlphabetTooShort`: Returned when alphabet length is less than 2 characters.
* `ErrAlphabetTooLong`: Returned when an alphabet length exceeds 256 characters.
* `ErrNilRandReader`: Returned when a nil random reader is provided.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
* `DefaultLength`: The default length of the generated ID: `21`
* `MinAlphabetLength`: The minimum allowed length for the alphabet: `2`
* `MaxAlphabetLength`: The maximum allowed length for the alphabet: `256`

## Performance Optimizations

### Buffer Pooling with `sync.Pool`

The nanoid generator utilizes `sync.Pool` to manage byte slice buffers efficiently. This approach minimizes memory allocations and enhances performance, especially in high-concurrency scenarios.

How It Works:
* Storing Pointers: `sync.Pool` stores pointers to `[]byte` slices (`*[]byte`) instead of the slices themselves. This avoids unnecessary allocations and aligns with best practices for using `sync.Pool`.
* Zeroing Buffers: Before returning buffers to the pool, they are zeroed out to prevent data leaks.

### Struct Optimization

The `generator` struct is optimized for memory alignment and size by ordering from largest to smallest to minimize padding and optimize memory usage.

## Execute Benchmarks:

Run the benchmarks using the `go test` command with the `bench` make target:

```shell
make bench
```

### Interpreting Results:

Sample output might look like this:

<details>
  <summary>Expand to see results</summary>

```shell
go clean
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkNanoIDAllocations-24                            3745630               319.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDAllocationsConcurrent-24                  1255897               951.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24          11057571               107.6 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          4382187               276.4 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          4246334               281.5 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          4097402               295.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          3491649               348.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         2849967               421.0 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24         11113387               106.7 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         4388538               275.2 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         4221810               282.5 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         4089571               295.7 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         3489831               345.4 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        2830932               428.5 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24         10801930               110.3 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         4352565               275.4 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         4287333               281.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         4092636               296.2 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         3496222               345.9 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        2852317               422.6 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24         10870525               107.6 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         4435015               274.3 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         4248135               281.9 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         4058250               296.3 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         3464295               346.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        2838750               423.3 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          5126839               237.1 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2622363               463.1 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         2425747               496.4 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         2068125               578.3 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24         1398754               854.4 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24        1000000              1158 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         8662704               139.0 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3546630               340.0 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        3306412               366.5 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2710264               446.1 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1901734               624.0 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24       1229500               980.4 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        8529920               140.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3551167               338.6 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       3297499               366.1 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2672764               443.5 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1918944               624.7 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24      1231447               981.3 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        8508936               139.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3554848               337.5 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       3306130               363.7 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2675566               442.3 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1914936               624.1 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24      1233818               975.2 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        8512947               139.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       3508246               339.7 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       3278722               363.5 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       2695258               445.4 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1918394               628.0 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24      1233004               975.9 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4541185               266.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2300229               514.4 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       2096215               573.2 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1679205               712.4 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24       1000000              1119 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       722259              1671 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   3491772               345.7 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1604065               759.1 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1538979               779.1 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1498453               803.1 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                  1343252               903.3 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                 1000000              1012 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  3404158               344.1 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1599379               760.2 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1549321               772.4 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1506124               799.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1343040               897.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                1000000              1005 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  3474410               345.9 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1574678               761.0 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1545774               769.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1515346               805.8 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1339198               883.2 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                1000000              1012 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  3474528               344.6 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1578373               758.2 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1545356               775.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1495804               795.4 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1356076               888.9 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                1201017              1017 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1599698               781.0 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  775138              1590 ns/op              16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  722536              1476 ns/op              24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  697624              1715 ns/op              32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  485611              2435 ns/op              64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 399118              3092 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2810716               419.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1382908               886.9 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1226542               911.0 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                1201102               998.6 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                1000000              1106 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                920350              1299 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2854525               419.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1361529               873.2 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1325887               907.0 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1204806              1005 ns/op              80 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24               1000000              1110 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               941889              1272 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2850246               417.8 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1366527               870.8 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1318562               913.8 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24               1000000              1001 ns/op              80 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24               1000000              1108 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               954940              1309 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                2857645               418.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24               1374385               872.8 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24               1317781               916.6 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24               1000000              1004 ns/op              80 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24               1000000              1107 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               917875              1297 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1464384               822.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                782992              1530 ns/op              48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                705030              1652 ns/op              48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                621445              2004 ns/op              80 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                419061              2888 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               317949              3724 ns/op             288 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24          11139302               106.0 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          4391175               272.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          4279638               280.8 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          4067823               292.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          3485182               343.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 2819616               424.4 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                 11114065               106.9 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 4399406               275.0 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 4244848               281.1 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 4102737               294.3 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 3468543               346.0 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                2820634               426.2 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                 11192881               106.2 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 4365012               274.6 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 4281464               280.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 4072840               295.6 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 3470144               345.4 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                2798563               427.1 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                 11086683               109.1 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 4346628               275.7 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 4287578               280.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 4051408               294.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 3486272               344.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                2815736               426.9 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  5074002               244.5 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2625111               461.4 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 2412589               492.9 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 2089557               576.1 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                 1366341               850.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                1000000              1155 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 8598400               139.0 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3571426               334.3 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                3299937               370.1 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2714208               443.2 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1891962               624.8 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24               1226782               973.0 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                8687305               137.5 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3578883               336.9 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               3297885               363.1 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2711613               443.8 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1922880               620.1 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24              1228870               974.6 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                8558007               139.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3528766               337.5 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               3276724               363.9 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2696563               444.7 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1920932               622.7 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24              1226420               978.3 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                8587652               138.4 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               3535434               338.3 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               3301360               364.0 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               2718532               443.3 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1927068               621.7 ns/op           144 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24              1226674               984.6 ns/op           288 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4537357               261.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2345822               513.1 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               2107111               570.4 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1695664               708.3 ns/op            80 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24               1000000              1120 ns/op             144 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               684584              1662 ns/op             288 B/op          1 allocs/op
PASS
ok      github.com/sixafter/nanoid      294.138s
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

1. Alphabet Lengths:
   * At Least Two Characters: The custom alphabet must contain at least two unique characters. An alphabet with fewer than two characters cannot produce IDs with sufficient variability or randomness.
   * Maximum Length 256 Characters: The implementation utilizes a rune-based approach, where each character in the alphabet is represented by a single rune. This allows for a broad range of unique characters, accommodating alphabets with up to 256 distinct runes. Attempting to use an alphabet with more than 256 runes will result in an error. 
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

## Determining Collisions

To determine the practical length for a NanoID for your use cases, see the collision time calculator [here](https://sixafter.github.io/nanoid/).

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.

