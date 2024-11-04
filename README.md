# Nano ID

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
    - 1 `allocs/op` for ASCII alphabets.
    - 2 `allocs/op` for Unicode alphabets.
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
BenchmarkNanoIDAllocations-24                            3356676               353.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDAllocationsConcurrent-24                  1000000              1021 ns/op              24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24          10151994               115.6 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          4163946               287.8 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          3940454               303.9 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          3665217               327.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          3009100               397.3 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         2260245               526.6 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24         10291764               117.2 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         4176050               286.8 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         3913437               304.4 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         3661736               327.4 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         3013918               397.6 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        2277447               526.1 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24         10248081               116.2 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         4153858               286.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         3974526               305.1 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         3588565               330.1 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         2999023               399.3 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        2276793               527.3 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24         10307857               116.3 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         4178316               287.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         3948546               306.1 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         3649971               327.3 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         3005053               400.5 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        2273695               530.1 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          4783996               249.2 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2490746               487.0 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         2253277               526.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         1921735               622.9 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24         1263078               948.8 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24         889077              1352 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         7531195               154.9 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3342734               362.5 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        3006205               397.4 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2443682               490.0 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1703467               699.8 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24       1000000              1108 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        7593878               153.4 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3362204               363.0 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       2986330               398.5 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2455066               489.6 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1728088               694.8 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24      1000000              1103 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        7712916               154.3 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3353942               359.4 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       3001768               400.1 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2478320               490.6 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1723248               695.5 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24      1000000              1103 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        7683140               154.4 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       3344334               362.0 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       2986064               401.6 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       2438979               491.5 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1716052               698.3 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24      1000000              1127 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4206267               285.6 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2179117               549.8 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       1953096               616.9 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1547134               775.7 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24        940594              1244 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       606660              1911 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   3285289               361.4 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1534184               793.8 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1457750               818.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1302526               872.9 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                  1217048               975.0 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                 1000000              1139 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  2782182               408.5 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1426534               820.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1469967               893.4 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1327660               859.0 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1235076               993.2 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                1000000              1095 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  3100905               380.4 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1476397               792.1 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1439024               825.0 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1408250               875.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1000000              1004 ns/op              64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                1000000              1078 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  2513728               406.9 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1435784               835.4 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1306690               889.7 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1393742               858.4 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1233390               972.2 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                1000000              1050 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1486549               873.6 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  765040              1487 ns/op              16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  753366              1560 ns/op              24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  648486              1799 ns/op              32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  463976              2603 ns/op              64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 358416              3270 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2360455               469.0 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1278476               948.3 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1000000              1006 ns/op             144 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                1000000              1044 ns/op             208 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                1000000              1209 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                765326              1514 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2191981               470.6 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1254249               953.4 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1214083               986.0 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1000000              1092 ns/op             208 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24               1000000              1249 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               807595              1494 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2620624               459.4 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1298533               924.5 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1000000              1012 ns/op             144 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24               1000000              1054 ns/op             208 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24               1000000              1227 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               802266              1482 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                2553933               461.7 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24               1270796               918.6 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24               1228171               988.2 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24               1000000              1063 ns/op             208 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24               1000000              1208 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               800295              1472 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1323894               904.1 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                644830              1689 ns/op             112 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                673209              1799 ns/op             144 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                572910              2134 ns/op             208 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                390074              3051 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               301100              3976 ns/op             801 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24          10053727               117.5 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          4126113               287.6 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          3913356               304.3 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          3667054               332.7 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          2963713               408.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 2257764               531.1 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                 10272817               117.4 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 4071330               296.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 3685614               325.2 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 3396084               344.2 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 2807998               413.0 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                2254412               530.8 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                  9955210               115.7 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 4130734               292.2 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 3648940               313.5 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 3528259               337.8 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 2929726               408.6 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                2266981               535.1 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                  9970590               118.7 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 4124941               303.8 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 3694315               320.0 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 3397144               337.1 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 2889451               416.3 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                2185933               555.9 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  4623091               257.1 ns/op             8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2364118               496.5 ns/op            16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 2184541               545.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 1923010               630.0 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                 1239230               961.0 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                 840214              1362 ns/op             128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 7558154               160.3 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3200314               375.9 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                2865019               416.5 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2287723               523.6 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1655614               717.6 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24               1000000              1184 ns/op             800 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                7104853               157.8 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3108900               372.2 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               2860558               416.3 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2378722               515.1 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1559844               740.4 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24              1000000              1194 ns/op             800 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                7326772               159.5 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3228379               366.6 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               2939876               416.0 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2255989               520.7 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1644298               744.0 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24               990902              1177 ns/op             800 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                7096801               166.8 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               3250207               380.1 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               2869855               408.7 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               2381173               511.7 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1656698               716.6 ns/op           400 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24              1000000              1147 ns/op             800 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4154898               286.2 ns/op            56 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2186154               551.6 ns/op           112 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               1882107               615.8 ns/op           144 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1541787               793.3 ns/op           208 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24                922514              1257 ns/op             400 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               628369              1935 ns/op             800 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      283.440s
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

