# nanoid <img src="docs/nanoid.svg" align="right" alt="Nano ID logo by Anton Lovchikov" width="160" height="94"/>

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
---

## Status

### Build & Test

[![CI](https://github.com/sixafter/nanoid/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid/actions)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid)](https://github.com/sixafter/nanoid/issues)

### Quality

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
![CodeQL](https://github.com/sixafter/nanoid/actions/workflows/codeql-analysis.yaml/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sixafter/nanoid/badge)](https://scorecard.dev/viewer/?uri=github.com/sixafter/nanoid)

### Package and Deploy

[![Release](https://github.com/sixafter/nanoid/workflows/release/badge.svg)](https://github.com/sixafter/nanoid/actions)

---
## Overview 

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
    - 1 `allocs/op` for ASCII and Unicode alphabets regardless of alphabet size or generated ID length.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library.
- **Supports `io.Reader` Interface**: 
  - The Nano ID generator satisfies the `io.Reader` interface, allowing it to be used interchangeably with any `io.Reader` implementations. 
  - Developers can utilize the Nano ID generator in contexts such as streaming data processing, pipelines, and other I/O-driven operations.

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
	id, err := gen.NewWithLength(10)
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

	// Generate a Nano ID using the custom generator with 
	// the default length.
	id, err := gen.New()
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
BenchmarkNanoIDAllocations-16                    	14587914	        80.82 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDAllocationsConcurrent-16          	69323420	        16.73 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_DefaultLength-16         	16424012	        70.80 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_2-16         	31998825	        38.05 ns/op	       2 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_3-16         	30143905	        40.72 ns/op	       3 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_5-16         	26432930	        45.19 ns/op	       5 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_13-16        	20775577	        57.36 ns/op	      16 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_21-16        	16875480	        69.66 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_VaryingBufferSizes/BufferSize_34-16        	12903283	        92.09 ns/op	      48 B/op	       1 allocs/op
BenchmarkGenerator_Read_ZeroLengthBuffer-16                        	964289800	         1.239 ns/op	       0 B/op	       0 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_1-16                	16791974	        70.12 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_2-16                	31188610	        37.83 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_4-16                	55412426	        26.19 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_8-16                	74801612	        31.20 ns/op	      24 B/op	       1 allocs/op
BenchmarkGenerator_Read_Concurrent/Concurrency_16-16               	72736273	        22.31 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen8-16             	23004424	        51.19 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen16-16            	19190353	        62.04 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen21-16            	16964116	        71.15 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen32-16            	13652666	        87.30 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen64-16            	 8390854	       143.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen2/IDLen128-16           	 5008268	       238.7 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen8-16            	24265525	        50.30 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen16-16           	19466930	        61.27 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen21-16           	17229444	        69.10 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen32-16           	13794391	        85.86 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen64-16           	 8595162	       139.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen16/IDLen128-16          	 5110944	       235.3 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen8-16            	24191109	        49.59 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen16-16           	19380535	        61.39 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen21-16           	17132138	        69.37 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen32-16           	13725036	        86.52 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen64-16           	 8456020	       141.1 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen32/IDLen128-16          	 5058157	       235.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen8-16            	24069922	        49.86 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen16-16           	19728499	        60.49 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen21-16           	17394328	        68.44 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen32-16           	13914199	        85.25 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen64-16           	 8510121	       140.6 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/ASCII_AlphabetLen64/IDLen128-16          	 5073615	       235.5 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen8-16           	16820199	        71.43 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen16-16          	11366964	       106.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen21-16          	 9410085	       125.6 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen32-16          	 6823558	       175.7 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen64-16          	 3750666	       319.4 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen2/IDLen128-16         	 2064969	       579.9 ns/op	     289 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen8-16          	17105458	        70.94 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen16-16         	11293284	       105.8 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen21-16         	 9442479	       127.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen32-16         	 6702454	       179.0 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen64-16         	 3743170	       318.6 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen16/IDLen128-16        	 2066145	       580.4 ns/op	     289 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen8-16          	16834544	        70.83 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen16-16         	11361748	       105.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen21-16         	 9452881	       125.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen32-16         	 6732062	       177.2 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen64-16         	 3751677	       319.9 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen32/IDLen128-16        	 2060302	       581.6 ns/op	     289 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen8-16          	16610605	        70.14 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen16-16         	11277501	       105.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen21-16         	 9512638	       126.1 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen32-16         	 6695277	       179.7 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen64-16         	 3771606	       319.6 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGeneration/Unicode_AlphabetLen64/IDLen128-16        	 2061250	       580.2 ns/op	     289 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen8-16     	163857794	         7.435 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen16-16    	120517556	        10.65 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen21-16    	90712375	        13.42 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen32-16    	70485098	        15.83 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen64-16    	41104804	        27.88 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen2/IDLen128-16   	24777241	        48.36 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen8-16    	136375516	         8.372 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen16-16   	94073989	        11.88 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen21-16   	82488160	        12.94 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen32-16   	73308183	        16.81 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen64-16   	43310809	        27.55 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen16/IDLen128-16  	24440457	        48.51 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen8-16    	141769581	         9.905 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen16-16   	126553687	        10.62 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen21-16   	90315254	        11.96 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen32-16   	74982036	        15.14 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen64-16   	44936035	        26.16 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen32/IDLen128-16  	25386526	        45.59 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen8-16    	143119507	         9.058 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen16-16   	100000000	        10.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen21-16   	86324469	        12.66 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen32-16   	73032126	        18.05 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen64-16   	43514587	        27.39 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/ASCII_AlphabetLen64/IDLen128-16  	25724988	        46.07 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen8-16   	95645796	        13.73 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen16-16  	57979766	        21.32 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen21-16  	51356856	        23.51 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen32-16  	35601486	        31.63 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen64-16  	21151725	        54.50 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen2/IDLen128-16 	12294115	        97.65 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen8-16  	94587804	        14.08 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen16-16 	56783565	        21.85 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen21-16 	54984812	        22.29 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen32-16 	36905144	        33.82 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen64-16 	21909837	        55.46 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen16/IDLen128-16         	12305467	        96.95 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen8-16           	87714758	        14.20 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen16-16          	56825020	        22.30 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen21-16          	53398698	        24.17 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen32-16          	34583626	        32.74 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen64-16          	21834276	        54.27 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen32/IDLen128-16         	12407594	        97.09 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen8-16           	94156719	        12.21 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen16-16          	56045520	        21.33 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen21-16          	51659749	        23.73 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen32-16          	34311201	        31.90 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen64-16          	21373562	        54.47 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDGenerationParallel/Unicode_AlphabetLen64/IDLen128-16         	12741946	        94.92 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen8-16      	22329752	        51.44 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen16-16     	19332024	        62.58 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen21-16     	16850302	        71.79 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen32-16     	13191100	        90.20 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen64-16     	 8065251	       149.0 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen2/IDLen128-16    	 4788142	       250.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen8-16     	23545724	        51.30 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen16-16    	18922980	        63.89 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen21-16    	16443630	        83.21 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen32-16    	13306922	        89.36 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen64-16    	 8110489	       148.2 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen16/IDLen128-16   	 4772100	       250.6 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen8-16     	23818345	        50.91 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen16-16    	18987817	        62.50 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen21-16    	16761151	        70.74 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen32-16    	13572120	        89.22 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen64-16    	 8131558	       146.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen32/IDLen128-16   	 4810249	       249.9 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen8-16     	24198812	        50.42 ns/op	       8 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen16-16    	19178608	        62.21 ns/op	      16 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen21-16    	16748401	        70.96 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen32-16    	13391131	        89.29 ns/op	      32 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen64-16    	 8085642	       147.5 ns/op	      64 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/ASCII_AlphabetLen64/IDLen128-16   	 4821776	       248.3 ns/op	     128 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen8-16    	16344978	        72.49 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen16-16   	11122821	       107.6 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen21-16   	 9301333	       129.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen32-16   	 6457276	       183.4 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen64-16   	 3651783	       328.0 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen2/IDLen128-16  	 2014346	       595.8 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen8-16   	16438815	        72.38 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen16-16  	11068940	       108.0 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen21-16  	 9255571	       129.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen32-16  	 6613152	       182.9 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen64-16  	 3529232	       328.1 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen16/IDLen128-16 	 2010649	       596.3 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen8-16   	16384964	        72.73 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen16-16  	11098354	       108.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen21-16  	 9297795	       129.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen32-16  	 6673686	       178.2 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen64-16  	 3660606	       327.3 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen32/IDLen128-16 	 2015034	       597.1 ns/op	     288 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen8-16   	16504117	        72.08 ns/op	      24 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen16-16  	11093815	       108.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen21-16  	 9267514	       129.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen32-16  	 6490016	       183.8 ns/op	      80 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen64-16  	 3678355	       326.1 ns/op	     144 B/op	       1 allocs/op
BenchmarkNanoIDWithVaryingAlphabetLengths/Unicode_AlphabetLen64/IDLen128-16 	 2017965	       593.8 ns/op	     288 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid	217.846s
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
