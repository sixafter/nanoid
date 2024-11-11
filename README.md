# nanoid <img src="https://ai.github.io/nanoid/logo.svg" align="right" alt="Nano ID logo by Anton Lovchikov" width="160" height="94">

[![CI](https://github.com/sixafter/nanoid/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid)](https://github.com/sixafter/nanoid/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)

A simple, fast, and efficient Go implementation of [Nano ID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator. 

Please see the [godoc](https://pkg.go.dev/github.com/sixafter/nanoid) for detailed documentation.

---

## Features

- **Short & Unique IDs**: Generates compact and collision-resistant identifiers.
- **Cryptographically Secure**: Utilizes Go's `crypto/rand` and `x/crypto/chacha20` stream cypher package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications.
- **Customizable**: 
  - Define your own set of characters for ID generation with a minimum length of 2 characters and maximum length of 256 characters.
  - Define your own random number generator.
  - Unicode and ASCII alphabets supported.
- **Concurrency Safe**: Designed to be safe for use in concurrent environments.
- **High Performance**: Optimized with buffer pooling to minimize allocations and enhance speed.
- **Optimized for Low Allocations**: Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.
    - 1 `allocs/op` for ASCII and Unicode alphabets.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.
- **Supports `io.Reader` Interface**: 
  - The Nano ID generator now satisfies the `io.Reader` interface, allowing it to be used interchangeably with any `io.Reader` implementations. 
  - Developers can now utilize the Nano ID generator in contexts such as streaming data processing, pipelines, and other I/O-driven operations.

---

## Installation

To install the package, use:

```sh
go get -u github.com/sixafter/nanoid
```

To use the NanoID package in your Go project, import it as follows:

```go
import "github.com/sixafter/nanoid"
```

---

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
Generated ID: mGbzQkkPBidjL4IP_MwBM
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
Generated ID: 1A3F5B7C9D
```

### Using `io.Reader` Interface

Here's a simple example demonstrating how to use the Nano ID generator as an `io.Reader`:

```go
package main

import (
  "fmt"
  "io"
  "github.com/sixafter/nanoid"
)

func main() {
	// Nano ID default length is 21
	buf := make([]byte, nanoid.DefaultLength)

	// Read a Nano ID into the buffer
	_, err := nanoid.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Convert the byte slice to a string
	id := string(buf)
	fmt.Printf("Generated ID: %s\n", id)
}
```

**Output**:

```bash
Generated ID: 2mhTvy21bBZhZcd80ZydM
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

	fmt.Println("Generated ID:", id)
}
```

**Output**"

```bash
Generated ID: G5J8K2M0QZ
```

### Customizing the Random Number Generator

You can customize the random number generator by using the WithRandReader option and generate an ID.

```go
package main

import (
	"crypto/rand"
	"fmt"

	"github.com/sixafter/nanoid"
)

func main() {
	// Create a new generator with custom random number generator
	gen, err := nanoid.NewGenerator(
		nanoid.WithRandReader(rand.Reader),
	)
	if err != nil {
		fmt.Println("Error creating Nano ID generator:", err)
		return
	}

	// Generate a Nano ID using the custom generator
	id, err := gen.New(21)
	if err != nil {
		fmt.Println("Error generating Nano ID:", err)
		return
	}

	fmt.Println("Generated ID:", id)
}
```

**Output**"

```bash
Generated ID: A8I8K3J0QY
```

---

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
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M2 Ultra
BenchmarkNanoIDAllocations-24                           10827518               104.6 ns/op            24 B/op          1 allocs/op
BenchmarkNanoIDAllocationsConcurrent-24                 44067356                24.11 ns/op           24 B/op          1 allocs/op
BenchmarkGenerator_Read_DefaultLength-24                11847897               101.3 ns/op            24 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_2-24              24025567                47.97 ns/op            8 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_3-24              23834449                49.51 ns/op            8 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_5-24              21935820                52.17 ns/op            8 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_13-24             15598916                73.71 ns/op           16 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_21-24             12269970                99.57 ns/op           24 B/op          1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_34-24              8378222               139.8 ns/op            48 B/op          1 allocs/op
BenchmarkGenerator_Read_ZeroLengthBuffer-24                             643477209                1.833 ns/op           0 B/op          0 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_1-24                     11262013               101.5 ns/op            24 B/op          1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_2-24                     21634225                55.88 ns/op           24 B/op          1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_4-24                     29086618                51.29 ns/op           24 B/op          1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_8-24                     40623742                29.63 ns/op           24 B/op          1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_16-24                    61002217                30.19 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-24                  20715415                57.10 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-24                 14955025                79.34 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-24                 12279360                97.13 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-24                  9492841               126.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-24                  5610008               213.4 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-24                 3108288               384.6 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-24                 21041617                57.57 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-24                14756728                80.05 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-24                12446965                96.60 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-24                 9663090               128.9 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-24                 5370190               218.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-24                3104888               384.8 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-24                 20816858                57.97 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-24                14801596                80.68 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-24                12056110                97.86 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-24                 9495126               124.8 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-24                 5546528               212.5 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-24                3147789               384.0 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-24                 20705972                57.19 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-24                14748634                79.16 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-24                12383686                98.82 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-24                 9345542               125.9 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-24                 5227130               269.2 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-24                2709976               417.3 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-24                15063817                76.47 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-24                9674256               120.7 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-24                8037175               147.0 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-24                5904072               199.5 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-24                3332470               367.5 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-24               1737326               664.3 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-24               15536979                78.73 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-24               9345562               127.8 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-24               7986922               156.2 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-24               5588406               205.6 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-24               3380896               368.3 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-24              1793653               664.3 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-24               15461583                79.51 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-24               9352959               126.1 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-24               7632367               148.3 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-24               5662828               211.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-24               3163326               372.6 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-24              1719470               692.0 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-24               15438170                76.58 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-24               9705403               126.2 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-24               7993968               147.8 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-24               5979666               203.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-24               3246212               363.7 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-24              1789210               672.4 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-24          121349466               10.56 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-24         74708942                15.09 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-24         58505428                19.06 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-24         48815461                25.62 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-24         27255820                42.41 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-24        16345488                74.82 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-24         100000000               10.15 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-24        78140683                13.96 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-24        60587575                19.10 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-24        50252307                24.57 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-24        27003231                40.80 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-24       16524638                73.20 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-24         127892064                9.683 ns/op           8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-24        94048568                13.11 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-24        64258606                17.34 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-24        51897687                22.83 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-24        29772888                40.79 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-24       16770325                69.18 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-24         136429317                9.342 ns/op           8 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-24        83319828                12.83 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-24        66990139                18.20 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-24        49283337                23.55 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-24        30812617                39.21 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-24       17114200                70.27 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-24        73182445                15.25 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-24       42446697                28.16 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-24       33360979                34.14 ns/op           48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-24       24449212                45.99 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-24       13794450                88.00 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-24       7435483               157.0 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-24       75993856                14.71 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-24      44642094                26.09 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-24      34424395                34.33 ns/op           48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-24      25146447                45.21 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-24      13921240                85.48 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-24              7565599               156.6 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-24               75665024                14.61 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-24              46223847                27.86 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-24              33996219                34.00 ns/op           48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-24              26041689                46.56 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-24              14382057                86.23 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-24              7704223               162.4 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-24               74564844                15.65 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-24              46940014                25.34 ns/op           32 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-24              33220176                34.68 ns/op           48 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-24              26996548                46.71 ns/op           64 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-24              14012982                89.21 ns/op          128 B/op          1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-24              7678606               155.9 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-24          19627298                60.10 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-24         14579775                81.96 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-24         11201199                98.05 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-24          9198715               126.0 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-24          5539738               219.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-24         3102500               386.5 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-24         19726945                58.64 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-24        14528738                81.44 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-24        11859025                99.33 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-24         9354871               128.0 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-24         5454200               218.2 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-24        3035320               390.7 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-24         19795814                59.80 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-24        14687432                79.53 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-24        12302660                96.29 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-24         9485805               127.7 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-24         5519390               213.5 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-24        3087526               386.2 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-24         19896495                58.97 ns/op            8 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-24        14839662                80.72 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-24        12053421                98.48 ns/op           24 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-24         9226418               131.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-24         5591241               211.4 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-24        3081849               383.0 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-24        15504409                77.82 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-24        9529544               123.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-24        7968878               148.6 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-24        5883114               201.6 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-24        3263871               363.7 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-24       1773298               686.8 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-24       15515886                77.89 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-24       9515356               123.8 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-24       8041390               155.3 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-24       5529211               206.0 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-24       3244819               377.3 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-24      1729305               690.8 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-24       15179576                79.85 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-24       9526300               126.5 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-24       7690760               154.2 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-24       5752617               209.7 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-24       3052999               391.2 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-24      1636466               722.0 ns/op           256 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-24       14725010                80.52 ns/op           16 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-24       9553730               124.2 ns/op            32 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-24       7688282               152.3 ns/op            48 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-24       5810041               204.1 ns/op            64 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-24       3273080               367.0 ns/op           128 B/op          1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-24      1766614               688.4 ns/op           256 B/op          1 allocs/op
PASS
ok      github.com/sixafter/nanoid      218.742s
```
</details>

* `ns/op`: Nanoseconds per operation. Lower values indicate faster performance.
* `B/op`: Bytes allocated per operation. Lower values indicate more memory-efficient code.
* `allocs/op`: Number of memory allocations per operation. Fewer allocations generally lead to better performance.

---

## ID Generation

Nano ID generates unique identifiers based on the following:

1. **Random Byte Generation**: Nano ID generates a sequence of random bytes using a secure random source (e.g., `crypto/rand.Reader`). 
2. **Mapping to Alphabet**: Each random byte is mapped to a character in a predefined alphabet to form the final ID. 
3. **Uniform Distribution**: To ensure that each character in the alphabet has an equal probability of being selected, Nano ID employs techniques to avoid bias, especially when the alphabet size isn't a power of two.

---

## Custom Alphabet Constraints

1. Alphabet Lengths:
   * At Least Two Characters: The custom alphabet must contain at least two unique characters. An alphabet with fewer than two characters cannot produce IDs with sufficient variability or randomness.
   * Maximum Length 256 Characters: The implementation utilizes a rune-based approach, where each character in the alphabet is represented by a single rune. This allows for a broad range of unique characters, accommodating alphabets with up to 256 distinct runes. Attempting to use an alphabet with more than 256 runes will result in an error. 
2. Uniqueness of Characters:
   * All Characters Must Be Unique. Duplicate characters in the alphabet can introduce biases in ID generation and compromise the randomness and uniqueness of the IDs. The generator enforces uniqueness by checking for duplicates during initialization. If duplicates are detected, it will return an `ErrDuplicateCharacters` error. 
3. Character Encoding:
   * Support for ASCII and Unicode: The generator accepts alphabets containing Unicode characters, allowing you to include a wide range of symbols, emojis, or characters from various languages.

---

## Determining Collisions

To determine the practical length for a NanoID for your use cases, see the collision time calculator [here](https://sixafter.github.io/nanoid/).

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.

