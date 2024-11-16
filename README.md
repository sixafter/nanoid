# nanoid <img src="https://ai.github.io/nanoid/logo.svg" align="right" alt="Nano ID logo by Anton Lovchikov" width="160" height="94">

[![CI](https://github.com/sixafter/nanoid/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid)](https://github.com/sixafter/nanoid/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
![CodeQL](https://github.com/sixafter/nanoid/actions/workflows/codeql-analysis.yaml/badge.svg)

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

Please see the [Nano ID CLI](https://github.com/sixafter/nanoid-cli) for a command-line interface (CLI) that uses this package to generate Nano IDs.

---

## Installation

### Using `go get`

To install the Nano ID package, run the following command:

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
	id, err := gen.New(nanoid.DefaultLength)
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
* Storing Pointers: `sync.Pool` stores pointers to `[]byte` (or `[]rune` if Unicode) slices (`*[]byte`) instead of the slices themselves. This avoids unnecessary allocations and aligns with best practices for using `sync.Pool`.
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
cpu: Apple M4 Max
BenchmarkNanoIDAllocations-16                    	14318704	        81.72 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDAllocationsConcurrent-16          	67763110	        16.15 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_DefaultLength-16         	15188494	        78.37 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_2-16         	35454968	        34.53 ns/op	       8 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_3-16         	32779533	        37.05 ns/op	       8 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_5-16         	29836489	        40.66 ns/op	       8 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_13-16        	20298000	        57.49 ns/op	      16 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_21-16        	16523661	        72.90 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_34-16        	11567829	       102.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkGenerator_Read_ZeroLengthBuffer-16                        	899398803	         1.334 ns/op	       0 B/op	       0 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_1-16                	15924032	        73.68 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_2-16                	30869953	        39.47 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_4-16                	53020219	        25.02 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_8-16                	71823493	        16.26 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_16-16               	74829878	        14.44 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-16             	26380341	        45.18 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-16            	19585760	        62.57 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-16            	16327316	        71.70 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-16            	12858915	        92.76 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-16            	 7864554	       152.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-16           	 4521963	       263.0 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-16            	26724024	        44.68 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-16           	19288508	        61.47 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-16           	16651855	        71.67 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-16           	13239339	        91.09 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-16           	 7927832	       151.5 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-16          	 4585581	       260.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-16            	27092953	        44.56 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-16           	19397608	        61.66 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-16           	16790154	        71.77 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-16           	12921525	        92.47 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-16           	 7786845	       152.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-16          	 4586020	       261.5 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-16            	26981778	        44.63 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-16           	19697626	        61.33 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-16           	16638607	        72.65 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-16           	12382740	        94.72 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-16           	 7822154	       152.6 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-16          	 4571776	       272.1 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-16           	19501108	        62.13 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-16          	12258480	        96.82 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-16          	10284886	       114.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-16          	 7772995	       154.8 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-16          	 4410819	       272.0 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-16         	 2386099	       501.8 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-16          	19307736	        63.37 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-16         	12347330	        96.57 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-16         	10350645	       118.9 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-16         	 7597382	       156.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-16         	 4367085	       273.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-16        	 2380953	       506.9 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-16          	19134736	        62.90 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-16         	12492582	        96.04 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-16         	10391311	       115.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-16         	 7725854	       156.2 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-16         	 4314513	       273.6 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-16        	 2381019	       502.9 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-16          	19200831	        63.52 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-16         	12389340	        96.67 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-16         	10312555	       116.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-16         	 7659938	       155.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-16         	 4388805	       273.1 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-16        	 2385751	       503.7 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-16     	137644742	         7.789 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-16    	100000000	        10.73 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-16    	93969322	        13.68 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-16    	62877154	        17.81 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-16    	40507353	        30.00 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-16   	23519860	        52.73 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-16    	148337209	         8.224 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-16   	100000000	        10.73 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-16   	80615584	        13.32 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-16   	67688575	        17.66 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-16   	40898517	        28.74 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-16  	24600330	        50.13 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-16    	148020730	         8.096 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-16   	100000000	        10.78 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-16   	91318411	        14.07 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-16   	71116905	        16.99 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-16   	40690780	        28.13 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-16  	24042754	        49.82 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-16    	155735688	         7.741 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-16   	100000000	        11.15 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-16   	88183687	        13.26 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-16   	67886580	        16.79 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-16   	42262760	        28.11 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-16  	24148390	        48.57 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-16   	95122332	        12.57 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-16  	59121673	        22.50 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-16  	43497828	        27.08 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-16  	32718980	        36.79 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-16  	17846673	        66.96 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-16 	 9805506	       121.8 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-16  	94081983	        13.57 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-16 	57457389	        22.33 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-16 	39253527	        31.03 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-16 	30329932	        39.87 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-16 	17305892	        71.16 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-16         	 9408576	       127.4 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-16           	90432381	        13.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-16          	44038987	        23.46 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-16          	40879766	        30.02 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-16          	30454308	        41.64 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-16          	15016314	        74.21 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-16         	 9485343	       130.4 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-16           	88650294	        13.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-16          	53983129	        21.90 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-16          	39779500	        28.94 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-16          	29728940	        38.54 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-16          	16351167	        68.58 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-16         	 9917641	       122.0 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-16      	23241456	        45.42 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-16     	18732742	        65.11 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-16     	15977692	        72.47 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-16     	12448131	        99.06 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-16     	 7415142	       164.3 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-16    	 4145812	       289.2 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-16     	25051995	        45.72 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-16    	19202598	        62.73 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-16    	16435147	        73.21 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-16    	12831844	        94.21 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-16    	 7704615	       157.2 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-16   	 4425217	       276.7 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-16     	25598953	        45.57 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-16    	19057614	        62.66 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-16    	16212391	        73.10 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-16    	12844366	        92.96 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-16    	 7733532	       159.5 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-16   	 4260783	       281.7 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-16     	24486859	        47.67 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-16    	18676666	        63.38 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-16    	16334020	        73.24 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-16    	12672018	        94.54 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-16    	 7721707	       155.3 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-16   	 4483032	       267.4 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-16    	19529527	        62.71 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-16   	12396213	        97.23 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-16   	10278896	       116.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-16   	 7664722	       156.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-16   	 4357148	       275.2 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-16  	 2366166	       507.0 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-16   	18755970	        63.63 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-16  	12390586	        97.43 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-16  	10318215	       116.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-16  	 7675764	       156.4 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-16  	 4367828	       274.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-16 	 2368156	       506.0 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-16   	19273186	        62.95 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-16  	12282088	        97.11 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-16  	10248876	       116.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-16  	 7676949	       156.4 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-16  	 4373037	       273.8 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-16 	 2372378	       505.9 ns/op	     256 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-16   	19010188	        63.21 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-16  	12517303	        96.70 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-16  	10282838	       116.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-16  	 7661167	       156.4 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-16  	 4359604	       274.6 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-16 	 2369816	       506.8 ns/op	     256 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid	211.634s
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
