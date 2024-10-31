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

<details>
  <summary>Expand to see results</summary>

```shell
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24          10418562               112.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24          4204862               285.5 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24          3918693               300.2 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24          3727174               320.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24          3060626               394.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24         2374159               510.1 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24         10628976               111.4 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24         4227380               282.1 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24         3992378               298.6 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24         3762184               321.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24         3111062               388.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24        2363734               508.1 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24         10828202               111.9 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24         4188156               283.9 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24         3966552               303.5 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24         3707356               322.4 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24         3077642               387.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24        2336438               511.8 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24         10687689               111.7 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24         4226514               285.5 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24         3965374               301.0 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24         3736278               324.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24         3092551               389.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24        2357520               511.1 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-24          4711828               244.9 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-24         2522955               474.6 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-24         2315516               522.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-24         1949098               611.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-24         1298616               924.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-24         867211              1302 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24         7062561               168.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24        3112558               386.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24        2668669               454.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24        2235410               536.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24        1511804               793.0 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24        906760              1310 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24        7125273               167.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24       3106981               386.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24       2668101               461.7 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24       2137002               547.6 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24       1461828               814.4 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24       915102              1305 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24        7122510               168.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24       3027632               387.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24       2683312               455.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24       2228652               556.0 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24       1491031               819.2 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24       852464              1336 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24        4579814               260.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24       2259762               531.3 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24       1923090               613.0 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24       1582443               730.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24       1000000              1133 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24       661370              1839 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-24        4731512               262.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-24       2314345               517.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-24       1994914               632.7 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-24       1667074               731.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-24       1000000              1156 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-24       645205              1859 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24   3147912               388.2 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24                  1486123               825.3 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24                  1423960               856.0 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24                  1368358               882.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24                  1216256               973.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24                 1000000              1118 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24                  3216572               378.5 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24                 1485900               806.6 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24                 1424578               849.5 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24                 1335081               898.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24                 1000000              1007 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24                1000000              1116 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24                  3188726               381.6 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24                 1414419               819.0 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24                 1418808               848.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24                 1368144               873.1 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24                 1219441               982.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24                1000000              1107 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24                  3159601               405.4 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24                 1455963               822.9 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24                 1433679               846.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24                 1381220               868.7 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24                 1233231               972.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24                1000000              1096 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-24                  1428054               853.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-24                  826676              1533 ns/op              32 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-24                  627172              1639 ns/op              48 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-24                  640642              1897 ns/op              64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-24                  457977              2681 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-24                 358770              3325 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24                 2435527               488.0 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24                1252327               969.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24                1000000              1038 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24                1000000              1097 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24                 962730              1289 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24                801014              1557 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24                2310380               497.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24               1237915               964.7 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24               1000000              1035 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24               1000000              1112 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24                946850              1292 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24               740536              1575 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24                2343194               490.4 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24               1230936               976.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24               1000000              1055 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24                996325              1109 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24                927045              1307 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24               779859              1593 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24                1470223               819.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24                795807              1599 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24                684703              1769 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24                598116              1954 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24                437926              2692 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24               327703              3785 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-24                1405555               834.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-24                794175              1544 ns/op             128 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-24                735001              1713 ns/op             176 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-24                592226              1979 ns/op             240 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-24                424042              2704 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-24               325590              3664 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24          10442265               115.0 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24          4181119               291.2 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24          3898227               311.7 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          3632714               330.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          2956627               393.8 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24                 2286168               525.7 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24                 10668045               113.3 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24                 4166764               289.0 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24                 3870742               308.9 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24                 3677217               336.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24                 3017631               406.2 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24                2226591               531.6 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24                 10315075               115.2 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24                 4191613               286.0 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24                 3999072               301.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24                 3673783               328.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24                 3013660               396.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24                2289802               532.2 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24                  9420152               113.2 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24                 4086746               293.2 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24                 3830166               305.7 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24                 3693453               322.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24                 3044371               395.5 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24                2283253               523.1 ns/op           256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-24                  4883360               244.7 ns/op            16 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-24                 2472624               483.5 ns/op            32 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-24                 2276542               525.8 ns/op            48 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-24                 1956292               623.2 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-24                 1298796               922.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-24                 900985              1304 ns/op             256 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24                 7071021               168.3 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24                3074877               393.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24                2647694               449.9 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24                2217471               539.5 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24                1481122               808.4 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24                897288              1332 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24                7013146               168.5 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24               3052245               386.4 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24               2607384               465.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24               2112885               569.2 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24               1471886               819.8 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24               858808              1361 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24                6907552               170.9 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24               3003462               397.9 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24               2593538               460.8 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24               2201772               544.7 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24               1487638               809.3 ns/op           464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24               876087              1354 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24                4592804               256.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24               2244604               540.0 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24               1947669               617.4 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24               1645828               724.3 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24               1000000              1134 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24               638697              1861 ns/op             929 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-24                4711143               253.8 ns/op            64 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-24               2332794               523.6 ns/op           128 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-24               1989870               613.6 ns/op           176 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-24               1664764               728.9 ns/op           240 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-24               1000000              1138 ns/op             464 B/op          2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-24               625184              1910 ns/op             929 B/op          2 allocs/op
PASS
ok      github.com/sixafter/nanoid      282.115s
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

