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
cpu: Apple M3 Max
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-16         	10922077	       103.2 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-16        	 4486478	       267.6 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-16        	 4184334	       283.6 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-16        	 3922456	       307.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-16        	 3119569	       383.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-16       	 2274475	       526.9 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-16        	11451014	       104.2 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-16       	 4465977	       269.6 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-16       	 4216795	       285.5 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-16       	 3905011	       308.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-16       	 3131558	       382.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-16      	 2305051	       518.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-16        	11342877	       104.0 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-16       	 4466145	       269.4 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-16       	 4237448	       282.4 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-16       	 3915836	       305.8 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-16       	 3176829	       377.4 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-16      	 2318043	       518.0 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-16        	11139789	       106.3 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-16       	 4452184	       268.7 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-16       	 4200566	       285.5 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-16       	 3890450	       310.7 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-16       	 3153600	       382.6 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-16      	 2310556	       521.4 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen8-16        	 5252811	       227.5 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen16-16       	 2637884	       452.3 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen21-16       	 2444319	       492.6 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen32-16       	 2077375	       577.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen64-16       	 1380889	       869.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen95/IDLen128-16      	  957130	      1219 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-16       	 7770453	       151.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-16      	 3450578	       345.5 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-16      	 3077316	       397.3 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-16      	 2582444	       463.0 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-16      	 1767738	       671.8 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-16     	 1000000	      1083 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-16      	 7847635	       152.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-16     	 3439189	       345.5 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-16     	 3068278	       391.4 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-16     	 2591792	       463.8 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-16     	 1779342	       675.3 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-16    	 1000000	      1091 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-16      	 7762616	       152.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-16     	 3434809	       350.3 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-16     	 3035335	       393.2 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-16     	 2587806	       461.9 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-16     	 1779013	       676.4 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-16    	 1000000	      1089 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-16      	 5203952	       229.3 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-16     	 2577484	       466.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-16     	 2273259	       524.3 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-16     	 1907805	       625.1 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-16     	 1252963	       961.1 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-16    	  753727	      1551 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen8-16      	 5134542	       231.3 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen16-16     	 2566364	       467.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen21-16     	 2273648	       528.1 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen32-16     	 1909202	       626.2 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen64-16     	 1240447	       961.9 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen95/IDLen128-16    	  741745	      1558 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-16 	 3880791	       310.2 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-16         	 1954041	       612.6 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-16         	 1910630	       631.2 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-16         	 1800909	       662.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-16         	 1637233	       741.4 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-16        	 1566673	       769.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-16         	 3879736	       310.0 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-16        	 1966995	       610.5 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-16        	 1911943	       627.2 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-16        	 1817702	       656.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-16        	 1617438	       738.1 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-16       	 1569681	       763.2 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-16         	 3906577	       305.5 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-16        	 1970916	       608.4 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-16        	 1903587	       626.0 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-16        	 1810779	       659.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-16        	 1638396	       734.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-16       	 1592252	       756.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-16         	 3887870	       308.8 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-16        	 1974447	       608.9 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-16        	 1903549	       628.3 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-16        	 1825004	       661.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-16        	 1629457	       738.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-16       	 1571749	       756.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen8-16         	 1827397	       666.0 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen16-16        	 1000000	      1188 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen21-16        	  974991	      1271 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen32-16        	  840049	      1449 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen64-16        	  611839	      1989 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen95/IDLen128-16       	  501709	      2395 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-16        	 3287462	       367.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-16       	 1730535	       698.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-16       	 1620866	       740.1 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-16       	 1611837	       732.8 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-16       	 1473528	       817.2 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-16      	 1258662	       945.6 ns/op	     929 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-16       	 3298028	       364.7 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-16      	 1719994	       700.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-16      	 1615981	       740.7 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-16      	 1609618	       740.5 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-16      	 1467729	       814.0 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-16     	 1247503	       959.3 ns/op	     929 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-16       	 3283580	       363.7 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-16      	 1712774	       700.1 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-16      	 1625227	       740.4 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-16      	 1617552	       746.1 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-16      	 1466486	       816.5 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-16     	 1244221	       957.8 ns/op	     929 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-16       	 1932594	       623.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-16      	 1000000	      1115 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-16      	 1000000	      1214 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-16      	  887461	      1378 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-16      	  654630	      1843 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-16     	  504536	      2462 ns/op	     929 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen8-16       	 1937722	       618.9 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen16-16      	 1000000	      1115 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen21-16      	 1000000	      1211 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen32-16      	  881874	      1363 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen64-16      	  672013	      1827 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen95/IDLen128-16     	  513342	      2447 ns/op	     929 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-16  	10605374	       110.5 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-16 	 4123538	       286.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-16 	 3956678	       303.1 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-16 	 3664176	       326.0 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-16 	 2954824	       408.9 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-16         	 2124205	       565.1 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-16          	10605242	       112.6 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-16         	 4085768	       286.3 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-16         	 3962410	       304.2 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-16         	 3632366	       330.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-16         	 2841583	       424.1 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-16        	 2097081	       562.9 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-16          	10629051	       111.3 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-16         	 4177550	       285.5 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-16         	 3963330	       305.6 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-16         	 3565934	       336.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-16         	 2925541	       410.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-16        	 2118814	       564.1 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-16          	10565844	       109.2 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-16         	 4159087	       288.1 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-16         	 3936292	       303.1 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-16         	 3658989	       328.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-16         	 2896797	       409.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-16        	 2135020	       572.5 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen8-16          	 5089890	       231.9 ns/op	      16 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen16-16         	 2574538	       464.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen21-16         	 2351917	       511.8 ns/op	      48 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen32-16         	 1988805	       597.3 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen64-16         	 1320415	       918.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen95/IDLen128-16        	  891171	      1291 ns/op	     256 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-16         	 7467366	       159.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-16        	 3275968	       366.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-16        	 2893497	       414.9 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-16        	 2427856	       491.8 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-16        	 1676679	       714.0 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-16       	  996352	      1159 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-16        	 7467913	       159.8 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-16       	 3282781	       364.6 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-16       	 2883831	       415.7 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-16       	 2460036	       489.8 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-16       	 1686549	       718.3 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-16      	  954498	      1181 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-16        	 7356744	       160.0 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-16       	 3271298	       367.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-16       	 2848550	       416.1 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-16       	 2414373	       490.5 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-16       	 1669810	       717.8 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-16      	 1000000	      1154 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-16        	 4991953	       237.5 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-16       	 2471748	       485.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-16       	 2173778	       551.0 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-16       	 1835912	       657.9 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-16       	 1000000	      1012 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-16      	  712250	      1637 ns/op	     928 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen8-16        	 4992194	       237.5 ns/op	      64 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen16-16       	 2476051	       484.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen21-16       	 2178163	       547.9 ns/op	     176 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen32-16       	 1836604	       653.0 ns/op	     240 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen64-16       	 1000000	      1008 ns/op	     464 B/op	       2 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen95/IDLen128-16      	  715062	      1631 ns/op	     928 B/op	       2 allocs/op
PASS
ok  	github.com/sixafter/nanoid	289.055s
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

