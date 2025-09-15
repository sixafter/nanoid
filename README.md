# nanoid: A library for generating cryptographically random Nano IDs.

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![FIPS‑140 Mode Compatible](https://img.shields.io/badge/FIPS‑140--Mode-Compatible-brightgreen)](FIPS‑140.md)

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
- **FIPS‑140 Mode Compatible**: Designed to run in FIPS‑140 validated environments using only Go standard library crypto. 
  - For FIPS‑140 compatible random number generation, use the [aes-ctr-drbg](https://github.com/sixafter/aes-ctr-drbg) module.
  - See [FIPS‑140.md](FIPS-140.md) for details and deployment guidance.

Please see the [nanoid-cli](https://github.com/sixafter/nanoid-cli) for a command-line interface (CLI) that uses this module to generate Nano IDs.

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the `nanoid` source, first download the target version and its signature file 
from the [releases page](https://github.com/sixafter/nanoid/releases) and then run the following commands:

```sh
# Fetch the latest release tag from GitHub API (e.g., "v1.51.0")
TAG=$(curl -s https://api.github.com/repos/sixafter/nanoid/releases/latest | jq -r .tag_name)

# Remove leading "v" for filenames (e.g., "v1.51.0" -> "1.51.0")
VERSION=${TAG#v}

# Verify the release tarball
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/nanoid/main/cosign.pub \
  --signature nanoid-${VERSION}.tar.gz.sig \
  nanoid-${VERSION}.tar.gz

# Download checksums.txt and its signature from the latest release assets
curl -LO https://github.com/sixafter/nanoid/releases/download/${TAG}/checksums.txt
curl -LO https://github.com/sixafter/nanoid/releases/download/${TAG}/checksums.txt.sig

# Verify checksums.txt with cosign
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

This project integrates a cryptographically secure, high-performance random number generator (CSPRNG) from [prng-chacha](https://github.com/sixafter/prng-chacha) that can be used for UUIDv4 generation with Google’s UUID library. By replacing the default entropy source with this CSPRNG, UUIDv4 creation is significantly faster in both serial and concurrent workloads, while maintaining cryptographic quality.

For implementation details, benchmark results, and usage, see the CSPRNG [README](https://github.com/sixafter/prng-chacha).

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
make bench
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M4 Max
Benchmark_Allocations_Serial-16                        	14830308	        74.10 ns/op	      24 B/op	       1 allocs/op
Benchmark_Allocations_Parallel-16                      	87374926	        13.93 ns/op	      24 B/op	       1 allocs/op
Benchmark_Read_DefaultLength-16                        	21823753	        54.01 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_2-16      	49202928	        24.94 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_3-16      	47922843	        25.39 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_5-16      	44769436	        27.30 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_13-16     	31385779	        37.98 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_21-16     	24563296	        49.45 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_34-16     	17664668	        67.80 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_ZeroLengthBuffer-16                     	986404063	         1.212 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_1-16             	23469319	        50.41 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_2-16             	47287224	        25.17 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_4-16             	89502978	        13.22 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_8-16             	123994070	         9.912 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_16-16            	209096799	         5.656 ns/op	       0 B/op	       0 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16  	31861614	        37.86 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16 	22977042	        51.86 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16 	18545468	        64.42 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16 	15435464	        76.98 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16 	 9426519	       126.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16         	 5483312	       217.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16          	32191970	        37.13 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16         	22994911	        51.37 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16         	18733461	        64.11 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16         	15541071	        76.51 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16         	 9399879	       127.0 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16        	 5479530	       218.4 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16          	32642143	        37.36 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16         	23300292	        51.40 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16         	18884118	        63.52 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16         	15857773	        75.56 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16         	 9551781	       126.2 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16        	 5524531	       217.5 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16          	32380906	        37.33 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16         	22809477	        51.54 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16         	18773943	        63.74 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16         	15882687	        75.79 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16         	 9416259	       126.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16        	 5503002	       217.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16         	16668662	        70.04 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16        	11526090	       104.1 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16        	 9527602	       125.2 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16        	 6898809	       174.3 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16        	 3853819	       312.6 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16       	 2084970	       575.9 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16        	17133789	        69.94 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16       	11481788	       105.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16       	 9534361	       125.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16       	 6850174	       174.2 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16       	 3839973	       311.8 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16      	 2088130	       574.6 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16        	17156673	        70.04 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16       	11410390	       105.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16       	 9537810	       125.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16       	 6809893	       175.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16       	 3833170	       313.1 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16      	 2081092	       574.3 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16        	16723262	        71.30 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16       	11270502	       105.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16       	 9602179	       125.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16       	 6761265	       176.4 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16       	 3848028	       312.4 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16      	 2083189	       576.3 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen8-16         	194620867	         6.209 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen16-16        	141351794	         9.845 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen21-16        	100000000	        11.55 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen32-16        	81193770	        14.75 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen64-16        	45330201	        26.13 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen128-16       	25851924	        46.62 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen8-16        	216777411	         5.758 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen16-16       	144298010	         8.848 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen21-16       	100000000	        11.71 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen32-16       	83137520	        13.78 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen64-16       	47245101	        25.14 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen128-16      	26406150	        45.66 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen8-16        	188019464	         5.510 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen16-16       	149308290	         8.004 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen21-16       	100000000	        11.01 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen32-16       	83274057	        13.92 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen64-16       	46638768	        25.90 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen128-16      	26978290	        44.27 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen8-16        	206536862	         5.456 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen16-16       	152015960	         8.130 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen21-16       	100000000	        10.83 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen32-16       	84979563	        13.97 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen64-16       	48281804	        24.52 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen128-16      	26492161	        43.76 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen8-16       	100000000	        12.01 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen16-16      	60293805	        19.54 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen21-16      	53047172	        21.30 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen32-16      	37889799	        32.63 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen64-16      	21920443	        54.26 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen128-16     	12110583	        99.18 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen8-16      	98768138	        12.44 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen16-16     	59912625	        20.72 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen21-16     	55665618	        22.08 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen32-16     	36814191	        31.89 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen64-16     	21810398	        54.81 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen128-16    	12246813	        98.59 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen8-16      	94218027	        11.72 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen16-16     	62175088	        19.93 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen21-16     	52981692	        21.81 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen32-16     	37626205	        31.30 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen64-16     	22032514	        53.96 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen128-16    	12033945	        98.92 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen8-16      	100000000	        11.71 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen16-16     	58656163	        21.07 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen21-16     	55079770	        21.62 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen32-16     	38111671	        31.08 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen64-16     	22505924	        53.80 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen128-16    	12325257	        97.61 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16         	30982933	        38.31 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16        	22225959	        53.39 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16        	17688570	        67.18 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16        	14962710	        80.15 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16        	 8811254	       135.8 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16       	 5078203	       236.1 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16        	31699906	        38.37 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16       	22297476	        53.64 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16       	17475834	        67.39 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16       	14922480	        80.63 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16       	 8864101	       135.1 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16      	 5098405	       235.2 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16        	31704025	        38.13 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16       	22310568	        53.15 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16       	17880979	        67.22 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16       	14910853	        79.94 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16       	 8856570	       135.2 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16      	 5126726	       233.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16        	29376333	        39.24 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16       	21777496	        53.67 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16       	17666781	        66.84 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16       	15056240	        79.63 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16       	 8923198	       134.1 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16      	 5096005	       232.5 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16       	16700318	        72.12 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16      	10957123	       108.2 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16      	 9271933	       129.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16      	 6679417	       179.7 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16      	 3739922	       321.0 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16     	 2036672	       590.0 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16      	16775415	        71.37 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16     	11206695	       106.8 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16     	 9401438	       127.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16     	 6660018	       180.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16     	 3748488	       320.5 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16    	 2034178	       590.0 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16      	16434753	        72.80 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16     	11008993	       108.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16     	 9391768	       127.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16     	 6771036	       177.6 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16     	 3759814	       319.2 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16    	 2033529	       589.7 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16      	16771731	        71.23 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16     	11227098	       106.8 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16     	 9385282	       127.8 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16     	 6751688	       177.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16     	 3754378	       319.9 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16    	 2040801	       588.3 ns/op	     288 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid	199.184s
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
