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
- **Customizable**: 
  - Define your own set of characters for ID generation with a minimum length of 2 characters and maximum length of 256 characters.
  - Define your own cryptographically secure random number generator.
- **Concurrency Safe**: Designed to be safe for use in concurrent environments.
- **High Performance**: Optimized with buffer pooling to minimize allocations and enhance speed.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.
- **Optimized for Low Allocations**: Carefully structured to minimize heap allocations at 2 `allocs/op`, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.

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

	// Create a new generator with custom alphabet and default length
	gen, err := nanoid.NewGenerator(
		nanoid.WithAlphabet(alphabet),
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

### `NewGenerator`

Creates a new Nano ID generator with a custom alphabet and random source.

```go
func NewGenerator(options ...Option) (Generator, error)
```

* Parameters:
  * `options` (`Option`): Variadic Option parameters to configure the Generator.  Options are `WithAlphabet` and `WithRandReader`.
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
    // GetConfig returns the configuration of the generator.
    GetConfig() Config
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
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24           8146620               144.1 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          3399670               350.5 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          3007588               387.2 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          2634784               495.2 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          1799394               665.4 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         1000000              1044 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24          8456922               146.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         3365966               360.0 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         3142118               382.9 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         2663785               451.4 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         1788866               659.9 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        1000000              1039 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24          8376380               141.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         3477678               351.7 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         3149929               381.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         2678804               451.2 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         1812704               660.1 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        1000000              1040 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24          8377962               142.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         3456831               346.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         3145819               382.8 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         2652284               451.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         1802006               681.8 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        1000000              1119 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          4065415               290.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2117118               569.0 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         1909093               619.8 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         1586401               754.7 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24          951430              1244 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24         621328              1856 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         7340752               163.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3154357               378.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        2677156               439.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2265291               522.0 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1569890               759.4 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24       1000000              1214 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        7339876               162.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3143100               375.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       2771413               439.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2263496               528.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1581406               771.2 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24       973166              1213 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        7479402               161.4 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3148720               372.8 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       2739915               437.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2256660               518.1 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1570305               758.5 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24       989757              1221 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        4717432               248.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       2405282               512.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       2032545               580.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       1728685               704.0 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1000000              1088 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24       688076              1750 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4865292               252.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2396606               510.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       2019394               578.1 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1680956               702.4 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24       1000000              1091 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       681628              1759 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   2654766               443.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1330676               916.7 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1241898               961.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1000000              1041 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                   992610              1205 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                  798307              1442 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  2641765               450.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1333837               906.7 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1241769               970.8 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1000000              1043 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1000000              1186 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                 835503              1461 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  2629880               450.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1309928               921.6 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1219352               973.4 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1000000              1044 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1000000              1194 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                 868380              1429 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  2641048               456.5 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1305056               923.8 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1228926               978.4 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1000000              1053 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1000000              1210 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                 860490              1422 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1280924               927.3 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  777445              1652 ns/op              88 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  672199              1816 ns/op             120 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  569658              2148 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  384903              3121 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 303528              4000 ns/op             657 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2522827               483.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1255570               953.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1000000              1035 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                1000000              1095 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                 896550              1272 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                767518              1539 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2512219               481.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1261227               949.1 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1000000              1042 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1000000              1090 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24                949870              1270 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               817960              1534 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2488392               480.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1258627               963.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1000000              1050 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24               1000000              1095 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24                988399              1264 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               761431              1538 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                1494270               807.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24                853420              1494 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24                727082              1659 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24                610881              1935 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24                455446              2619 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               338713              3542 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1515350               799.6 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                807898              1501 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                728190              1674 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                609951              1910 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                446277              2629 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               349966              3535 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24           7988445               146.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          3361719               355.0 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          3048328               392.0 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          2594168               463.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          1733986               694.2 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 1000000              1104 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                  7867738               146.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 3382924               354.8 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 3071925               391.5 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 2522364               467.1 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 1744478               702.7 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                1000000              1114 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                  8063221               146.1 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 3411207               361.0 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 2999431               390.7 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 2582430               474.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 1737324               684.8 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                1000000              1086 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                  8200508               148.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 3390000               354.2 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 3025774               405.3 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 2544818               461.3 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 1735057               690.2 ns/op           336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                1000000              1088 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  4248472               281.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2196370               547.8 ns/op            88 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 1945632               611.6 ns/op           120 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 1593208               753.7 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                  960559              1222 ns/op             336 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                 651579              1866 ns/op             656 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 7307659               161.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3155149               375.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                2754608               440.5 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2237965               528.9 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1556935               773.5 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24                876754              1273 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                7270288               163.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3109839               388.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               2675790               439.9 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2271544               528.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1554565               773.5 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24               967503              1244 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                7092847               163.4 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3213897               375.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               2764551               436.9 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2292484               523.2 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1506630               774.0 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24               977115              1242 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                4826038               252.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               2244790               517.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               2067538               594.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               1721875               700.2 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1000000              1103 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24               667432              1788 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4787002               248.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2371886               517.1 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               2029287               577.9 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1717202               693.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24               1000000              1098 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               647173              1826 ns/op             929 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      274.300s
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

## Determining Collisions

To determine the practical length for a NanoID, see the collision time calculator [here](https://sixafter.github.io/nanoid/).

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

