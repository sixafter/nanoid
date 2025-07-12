# nanoid: A library for generating cryptographically random Nano IDs.

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
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10826/badge)](https://www.bestpractices.dev/projects/10826)
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
    * The custom Cryptographically Secure Pseudo Random Number Generator (CSPRNG) Includes a thread-safe global `Reader` for concurrent access.
    * Up to 98% faster when using the `prng.Reader` as a source for v4 UUID generation using Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package.
    * See the benchmark results [here](x/crypto/prng/README.md#uuid-generation).
- **Customizable**: 
  - Define your own set of characters for ID generation with a minimum length of 2 characters and maximum length of 256 characters.
  - Define your own random number generator.
  - Unicode and ASCII alphabets supported.
- **Concurrency Safe**: Designed to be safe for use in concurrent environments.
- **High Performance**: Optimized with buffer pooling to minimize allocations and enhance speed.
- **Optimized for Low Allocations**: Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.
    - 1 `allocs/op` for ASCII and Unicode alphabets regardless of alphabet size or generated ID length.
    - 0 `allocs/op` for `Reader` interface across ASCII and Unicode alphabets regardless of alphabet size or generated ID length.
- **Zero Dependencies**: Lightweight implementation with no external dependencies beyond the standard library other than for tests.
- **Supports `io.Reader` Interface**: 
  - The Nano ID generator satisfies the `io.Reader` interface, allowing it to be used interchangeably with any `io.Reader` implementations. 
  - Developers can utilize the Nano ID generator in contexts such as streaming data processing, pipelines, and other I/O-driven operations.

Please see the [Nano ID CLI](https://github.com/sixafter/nanoid-cli) for a command-line interface (CLI) that uses this package to generate Nano IDs.

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the `nanoid` source, first download the target version and its signature file 
from the [releases page](https://github.com/sixafter/nanoid/releases) along with its `.sig` file; e.g., 
`nanoid-1.32.0.tar.gz` and `nanoid-1.32.0.tar.gz.sig`. Then run the following command to verify the 
signature:

```sh
# Replace <version> with the release version you downloaded, e.g., 1.32.0

cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/nanoid/main/cosign.pub \
  --signature nanoid-<version>.tar.gz.sig \
  nanoid-<version>.tar.gz

# Example with version 1.32.0:
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/nanoid/main/cosign.pub \
  --signature nanoid-1.32.0.tar.gz.sig \
  nanoid-1.32.0.tar.gz
```

The checksums are also signed (`checksums.txt` and `checksums.txt.sig`), verify with:

```sh
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/nanoid/main/cosign.pub \
  --signature checksums.txt.sig \
  checksums.txt
```

If valid, Cosign will output:

```shell
Verified OK
```

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

### Cryptographically Secure Pseudo Random Number Generator (CSPRNG)

This project integrates a cryptographically secure, high-performance random number generator (CSPRNG) from [x/crypto/prng](x/crypto/prng) that can be used for UUIDv4 generation with Google’s UUID library. By replacing the default entropy source with this CSPRNG, UUIDv4 creation is significantly faster in both serial and concurrent workloads, while maintaining cryptographic quality.

For implementation details, benchmark results, and usage, see the CSPRNG [README](x/crypto/prng).

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
Benchmark_Allocations_Serial-16                        	14198004	        79.99 ns/op	      24 B/op	       1 allocs/op
Benchmark_Allocations_Parallel-16                      	70763061	        15.95 ns/op	      24 B/op	       1 allocs/op
Benchmark_Read_DefaultLength-16                        	21281154	        56.13 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_2-16      	46700102	        25.43 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_3-16      	46027712	        26.17 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_5-16      	41190932	        29.25 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_13-16     	28199160	        42.48 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_21-16     	21307936	        56.23 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_34-16     	16708369	        72.10 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_ZeroLengthBuffer-16                     	968834155	         1.242 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_1-16             	21443377	        55.29 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_2-16             	41843678	        28.38 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_4-16             	69332600	        14.78 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_8-16             	121183014	        15.12 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_16-16            	100000000	        11.62 ns/op	       0 B/op	       0 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16  	28916445	        40.50 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16 	21362242	        55.51 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16 	17373552	        68.55 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16 	14781601	        80.56 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16 	 8970748	       132.6 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16         	 5305686	       225.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16          	30501528	        39.23 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16         	22141957	        54.19 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16         	17827231	        67.29 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16         	14888814	        79.46 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16         	 8998646	       133.4 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16        	 5313847	       225.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16          	30973369	        39.49 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16         	22161619	        54.44 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16         	17794450	        66.74 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16         	15145483	        78.76 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16         	 9064513	       132.1 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16        	 5342866	       224.7 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16          	31185504	        39.49 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16         	22000536	        54.41 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16         	17859567	        66.90 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16         	15088242	        79.00 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16         	 9096211	       131.9 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16        	 5338404	       224.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16         	16415924	        72.64 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16        	11043372	       108.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16        	 9127800	       130.7 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16        	 6501171	       185.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16        	 3594266	       334.8 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16       	 1960977	       611.5 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16        	16586402	        72.61 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16       	11054176	       107.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16       	 9207160	       130.1 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16       	 6407494	       186.8 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16       	 3615024	       335.1 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16      	 1947868	       615.0 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16        	16251949	        73.57 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16       	10929454	       108.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16       	 9203847	       130.4 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16       	 6484030	       185.5 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16       	 3631498	       330.6 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16      	 1974337	       609.9 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16        	16457274	        72.70 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16       	10989454	       108.1 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16       	 9142694	       130.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16       	 6433011	       186.4 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16       	 3615510	       330.0 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16      	 1972446	       610.1 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen8-16         	164444101	         8.541 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen16-16        	100000000	        10.04 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen21-16        	88613745	        12.91 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen32-16        	76542405	        15.52 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen64-16        	42477812	        26.90 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen128-16       	24728820	        46.90 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen8-16        	179634476	         7.510 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen16-16       	127661520	         8.998 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen21-16       	92385864	        11.94 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen32-16       	79632363	        15.24 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen64-16       	43324426	        25.59 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen128-16      	25727125	        45.49 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen8-16        	163591700	         7.590 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen16-16       	126878126	         9.677 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen21-16       	84514975	        12.20 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen32-16       	78105081	        15.55 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen64-16       	46511478	        25.63 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen128-16      	26097026	        44.97 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen8-16        	171210379	         7.526 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen16-16       	100000000	        11.13 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen21-16       	96146780	        12.94 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen32-16       	77217180	        15.20 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen64-16       	42647775	        25.62 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen128-16      	26483367	        44.45 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen8-16       	82645812	        13.17 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen16-16      	57481702	        20.72 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen21-16      	49317178	        24.90 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen32-16      	35697950	        34.65 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen64-16      	20912510	        55.98 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen128-16     	11937585	       101.0 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen8-16      	86822826	        13.41 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen16-16     	54938557	        20.75 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen21-16     	52044649	        24.39 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen32-16     	35434160	        34.41 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen64-16     	21188840	        57.43 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen128-16    	11901056	        99.58 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen8-16      	89744787	        12.83 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen16-16     	58427806	        20.44 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen21-16     	48528964	        24.11 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen32-16     	35622930	        33.96 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen64-16     	21348229	        55.30 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen128-16    	12144427	        99.09 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen8-16      	82998313	        14.10 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen16-16     	58289803	        20.54 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen21-16     	49814835	        23.07 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen32-16     	34925558	        34.41 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen64-16     	21144768	        55.28 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen128-16    	12231183	        97.42 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16         	29768426	        40.07 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16        	21614188	        55.47 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16        	17038898	        69.35 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16        	14511915	        81.69 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16        	 8608467	       138.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16       	 5049069	       237.1 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16        	30524965	        40.00 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16       	21533708	        55.33 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16       	17262978	        68.94 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16       	14712380	        81.15 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16       	 8692237	       137.8 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16      	 5060598	       236.2 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16        	30026742	        39.93 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16       	21443089	        55.78 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16       	17341530	        69.04 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16       	14586480	        81.80 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16       	 8619507	       138.4 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16      	 5086555	       235.3 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16        	30773439	        39.80 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16       	21579045	        55.23 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16       	17644464	        68.63 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16       	14548290	        81.52 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16       	 8722580	       136.8 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16      	 5150535	       232.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16       	16142203	        73.50 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16      	11006241	       108.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16      	 9083914	       131.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16      	 6424624	       186.5 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16      	 3574914	       335.6 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16     	 1943611	       616.6 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16      	16410714	        73.03 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16     	11068536	       108.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16     	 9085238	       131.4 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16     	 6414139	       187.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16     	 3584641	       335.2 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16    	 1947680	       617.6 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16      	16402704	        73.04 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16     	11021361	       108.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16     	 9090747	       131.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16     	 6445450	       186.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16     	 3582855	       335.0 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16    	 1947332	       616.9 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16      	16392835	        72.98 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16     	11017634	       108.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16     	 9082651	       131.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16     	 6420574	       186.6 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16     	 3579476	       334.9 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16    	 1947726	       616.6 ns/op	     288 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid	215.145s
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
