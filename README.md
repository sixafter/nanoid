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
make bench
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M4 Max
Benchmark_Allocations_Serial-16                        	14637210	        75.24 ns/op	      24 B/op	       1 allocs/op
Benchmark_Allocations_Parallel-16                      	79356832	        13.57 ns/op	      24 B/op	       1 allocs/op
Benchmark_Read_DefaultLength-16                        	21212094	        53.69 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_2-16      	48615308	        24.56 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_3-16      	46863328	        25.46 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_5-16      	43054796	        28.28 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_13-16     	28753189	        39.95 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_21-16     	23684326	        50.52 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_VaryingBufferSizes/BufferSize_34-16     	17648680	        68.53 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_ZeroLengthBuffer-16                     	951876057	         1.227 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_1-16             	22802019	        51.44 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_2-16             	45954488	        25.66 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_4-16             	90800460	        13.14 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_8-16             	180719409	        10.21 ns/op	       0 B/op	       0 allocs/op
Benchmark_Read_Concurrent/Concurrency_16-16            	202727172	         5.696 ns/op	       0 B/op	       0 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16  	30999007	        38.32 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16 	22725909	        51.99 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16 	18595208	        63.74 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16 	15956376	        75.03 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16 	 9386469	       127.5 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16         	 5438360	       219.4 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16          	31870428	        37.80 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16         	22902620	        51.85 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16         	18651664	        64.06 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16         	15987067	        74.65 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16         	 9520399	       125.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16        	 5537190	       216.5 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16          	32098863	        37.67 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16         	23142808	        51.57 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16         	18592316	        64.86 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16         	15840973	        75.33 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16         	 9478008	       126.4 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16        	 5482266	       218.1 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16          	31620621	        37.51 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16         	23359842	        50.77 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16         	18964536	        63.57 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16         	15964903	        74.77 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16         	 9505382	       125.4 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16        	 5510490	       217.8 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16         	16615780	        71.26 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16        	11209304	       106.2 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16        	 9361395	       127.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16        	 6517618	       183.2 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16        	 3724382	       321.8 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16       	 1996093	       600.2 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16        	16856338	        70.79 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16       	11307592	       105.6 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16       	 9405472	       126.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16       	 6687494	       181.8 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16       	 3626528	       331.9 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16      	 1967749	       608.6 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16        	16689121	        70.95 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16       	11429323	       104.1 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16       	 9488731	       125.7 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16       	 6744073	       176.3 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16       	 3742791	       320.3 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16      	 2004769	       598.2 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16        	16736294	        70.82 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16       	11313909	       105.9 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16       	 9340014	       128.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16       	 6495428	       185.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16       	 3662890	       325.1 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16      	 2018772	       594.3 ns/op	     289 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen8-16         	202263056	         5.675 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen16-16        	142484485	         8.462 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen21-16        	98075608	        14.39 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen32-16        	78095342	        15.17 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen64-16        	46862568	        25.98 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen2/IDLen128-16       	25799908	        45.78 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen8-16        	215217715	         6.473 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen16-16       	147801945	         8.325 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen21-16       	100000000	        11.30 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen32-16       	86327575	        15.07 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen64-16       	46619745	        26.43 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen16/IDLen128-16      	26700015	        44.26 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen8-16        	207114432	         6.174 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen16-16       	124926253	         8.329 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen21-16       	100000000	        10.81 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen32-16       	85035519	        13.57 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen64-16       	45686946	        24.17 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen32/IDLen128-16      	27624494	        43.51 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen8-16        	220309203	         5.491 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen16-16       	149538975	         8.673 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen21-16       	100000000	        10.72 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen32-16       	86562655	        13.29 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen64-16       	47680539	        25.03 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/ASCII_AlphabetLen64/IDLen128-16      	26801118	        43.16 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen8-16       	100000000	        11.68 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen16-16      	57968796	        19.81 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen21-16      	55751612	        23.26 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen32-16      	36476612	        31.79 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen64-16      	21735142	        55.01 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen2/IDLen128-16     	11704328	       100.9 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen8-16      	96396837	        12.89 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen16-16     	59590932	        19.36 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen21-16     	52840933	        21.38 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen32-16     	35645769	        33.03 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen64-16     	21020624	        54.77 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen16/IDLen128-16    	11921126	        99.56 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen8-16      	100000000	        12.79 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen16-16     	59407410	        20.24 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen21-16     	51632427	        22.46 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen32-16     	38363784	        31.50 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen64-16     	21653142	        55.53 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen32/IDLen128-16    	11951625	        99.53 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen8-16      	97809140	        12.03 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen16-16     	61195740	        19.06 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen21-16     	52913358	        21.22 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen32-16     	37887158	        32.39 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen64-16     	21262551	        55.12 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Parallel/Unicode_AlphabetLen64/IDLen128-16    	12071892	        97.69 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen8-16         	29391324	        38.77 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen16-16        	22316290	        53.41 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen21-16        	17824032	        67.18 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen32-16        	14985848	        79.18 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen64-16        	 8857060	       134.9 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen2/IDLen128-16       	 5089053	       235.5 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen8-16        	30639862	        38.92 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen16-16       	22222050	        53.46 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen21-16       	17878260	        66.99 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen32-16       	15094030	        79.06 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen64-16       	 8845488	       135.4 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen16/IDLen128-16      	 5080536	       235.1 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen8-16        	31241558	        38.58 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen16-16       	22127037	        53.55 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen21-16       	17975115	        67.16 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen32-16       	15091262	        79.19 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen64-16       	 8995386	       133.9 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen32/IDLen128-16      	 5138746	       232.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen8-16        	31104634	        38.82 ns/op	       8 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen16-16       	22296405	        53.28 ns/op	      16 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen21-16       	17960217	        66.64 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen32-16       	15204202	        78.50 ns/op	      32 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen64-16       	 8978506	       132.9 ns/op	      64 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/ASCII_AlphabetLen64/IDLen128-16      	 5224174	       229.6 ns/op	     128 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen8-16       	16391967	        72.58 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen16-16      	11038838	       108.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen21-16      	 9187986	       130.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen32-16      	 6346304	       183.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen64-16      	 3619990	       330.8 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen2/IDLen128-16     	 1946517	       616.0 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen8-16      	16510437	        71.95 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen16-16     	11135348	       107.4 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen21-16     	 9244778	       129.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen32-16     	 6579336	       182.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen64-16     	 3636339	       328.9 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen16/IDLen128-16    	 1961808	       611.8 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen8-16      	16424396	        72.22 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen16-16     	11143288	       107.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen21-16     	 9240367	       129.2 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen32-16     	 6551907	       181.8 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen64-16     	 3650402	       328.0 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen32/IDLen128-16    	 1955830	       611.5 ns/op	     288 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen8-16      	16809498	        71.23 ns/op	      24 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen16-16     	11243362	       106.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen21-16     	 9310128	       128.7 ns/op	      48 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen32-16     	 6514120	       183.3 ns/op	      80 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen64-16     	 3668616	       327.9 ns/op	     144 B/op	       1 allocs/op
Benchmark_Alphabet_Varying_Length_Varying_Serial/Unicode_AlphabetLen64/IDLen128-16    	 1949391	       615.2 ns/op	     288 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid	201.113s
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
