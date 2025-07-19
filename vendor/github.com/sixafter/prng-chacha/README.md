# prng-chacha

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/prng-chacha)](https://goreportcard.com/report/github.com/sixafter/prng-chacha)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/prng-chacha)](https://img.shields.io/github/go-mod/go-version/sixafter/prng-chacha)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/prng-chacha.svg)](https://pkg.go.dev/github.com/sixafter/prng-chacha)

---

## Status

### Build & Test

[![CI](https://github.com/sixafter/prng-chacha/workflows/ci/badge.svg)](https://github.com/sixafter/prng-chacha/actions)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/prng-chacha)](https://github.com/sixafter/prng-chacha/issues)

### Quality

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_prng-chacha&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_prng-chacha)
![CodeQL](https://github.com/sixafter/prng-chacha/actions/workflows/codeql-analysis.yaml/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=six-after_prng-chacha&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=six-after_prng-chacha)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sixafter/prng-chacha/badge)](https://scorecard.dev/viewer/?uri=github.com/sixafter/prng-chacha)

### Package and Deploy

[![Release](https://github.com/sixafter/prng-chacha/workflows/release/badge.svg)](https://github.com/sixafter/prng-chacha/actions)

---
## Overview 

The `prng` package provides a high-performance, cryptographically secure pseudo-random number generator (CSPRNG) that implements the `io.Reader` interface. Designed for concurrent use, it leverages the ChaCha20 cipher stream to efficiently generate random bytes.

Technically, this PRNG is not pseudo-random but is cryptographically random.

The package includes a global `Reader` and a `sync.Pool` to manage PRNG instances, ensuring low contention and optimized performance.

Please see the [godoc](https://pkg.go.dev/github.com/sixafter/prng-chacha) for detailed documentation.

---

## Features

* **Cryptographic Security:** Utilizes the [ChaCha20](https://pkg.go.dev/golang.org/x/crypto/chacha20) cipher for secure random number generation. 
* **Concurrent Support:** Includes a thread-safe global `Reader` for concurrent access. 
    * Up to 98% faster when using the `prng.Reader` as a source for v4 UUID generation using Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package as compared to using the default rand reader.
    * See the benchmark results [here](#uuid-generation).
* **Efficient Resource Management:** Uses a `sync.Pool` to manage PRNG instances, reducing the overhead on `crypto/rand.Reader`. 
* **Extensible API:** Allows users to create and manage custom PRNG instances via `NewReader`.
- **UUID Generation Source:** Can be used as the `io.Reader` source for UUID generation with the [`google/uuid`](https://pkg.go.dev/github.com/google/uuid) package and similar libraries, providing cryptographically secure, deterministic UUIDs using PRNG-CHACHA.

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the release tarball, you can use Cosign to check the signature against the public key.

```sh
# Fetch the latest release tag from GitHub API (e.g., "v1.41.0")
TAG=$(curl -s https://api.github.com/repos/sixafter/prng-chacha/releases/latest | jq -r .tag_name)

# Remove leading "v" for filenames (e.g., "v1.41.0" -> "1.41.0")
VERSION=${TAG#v}

# Verify the release tarball
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/prng-chacha/main/cosign.pub \
  --signature prng-chacha-${VERSION}.tar.gz.sig \
  prng-chacha-${VERSION}.tar.gz

# Download checksums.txt and its signature from the latest release assets
curl -LO https://github.com/sixafter/prng-chacha/releases/download/${TAG}/checksums.txt
curl -LO https://github.com/sixafter/prng-chacha/releases/download/${TAG}/checksums.txt.sig

# Verify checksums.txt with cosign
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/prng-chacha/main/cosign.pub \
  --signature checksums.txt.sig \
  checksums.txt
```

If valid, Cosign will output:

```shell
Verified OK
```

---

## Installation

```bash
go get -u github.com/sixafter/prng-chacha
```

---

## Usage

## Usage

Global Reader:

```go
package main

import (
  "fmt"
  
  "github.com/sixafter/sixafter/prng-chacha"
)

func main() {
  buffer := make([]byte, 64)
  n, err := prng.Reader.Read(buffer)
  if err != nil {
      // Handle error
  }
  fmt.Printf("Read %d bytes of random data: %x\n", n, buffer)
}
```

Replacing default random reader for UUID Generation:

```go
package main

import (
  "fmt"

  "github.com/google/uuid"
  "github.com/sixafter/sixafter/prng-chacha"
)

func main() {
  // Set the global random reader for UUID generation
  uuid.SetRand(prng.Reader)

  // Generate a new v4 UUID
  u := uuid.New()
  fmt.Printf("Generated UUID: %s\n", u)
}
```

---

## Architecture

* Global Reader: A pre-configured io.Reader (`prng.Reader`) manages a pool of PRNG instances for concurrent use. 
* PRNG Instances: Each instance uses ChaCha20, initialized with a unique key and nonce sourced from `crypto/rand.Reader`. 
* Error Handling: The `errorPRNG` ensures safe failure when initialization errors occur. 
* Resource Efficiency: A `sync.Pool` optimizes resource reuse and reduces contention on `crypto/rand.Reader`.

---

## Performance Benchmarks

### Raw Random Byte Generation

These `prng.Reader` benchmarks demonstrate the package's focus on minimizing latency, memory usage, and allocation overhead, making it suitable for high-performance applications.

<details>
  <summary>Expand to see results</summary>

```shell
make bench-csprng
goos: darwin
goarch: arm64
pkg: github.com/sixafter/prng-chacha
cpu: Apple M4 Max
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G2-16 	1000000000	         0.5689 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G4-16 	1000000000	         0.5723 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G8-16 	1000000000	         0.5711 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G16-16         	1000000000	         0.5912 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G32-16         	1000000000	         0.5396 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G64-16         	1000000000	         0.5438 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G128-16        	1000000000	         0.5442 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_8Bytes-16            	67387809	        17.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16Bytes-16           	49208646	        24.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_21Bytes-16           	40782282	        28.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_32Bytes-16           	32735828	        36.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_64Bytes-16           	20102887	        59.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_100Bytes-16          	13549210	        88.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_256Bytes-16          	 7663681	       158.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_512Bytes-16          	 4108966	       292.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_1000Bytes-16         	 2069100	       577.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_4096Bytes-16         	  546291	      2176 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16384Bytes-16        	  138847	      8614 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_2Goroutines-16         	315906483	         3.780 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_4Goroutines-16         	315760286	         3.834 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_8Goroutines-16         	319002478	         3.759 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_16Goroutines-16        	320004976	         3.766 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_32Goroutines-16        	319288269	         3.771 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_64Goroutines-16        	320027024	         3.747 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_128Goroutines-16       	322911223	         3.718 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_2Goroutines-16         	268402184	         4.528 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_4Goroutines-16         	267405596	         4.472 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_8Goroutines-16         	266455920	         4.447 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_16Goroutines-16        	267324904	         4.396 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_32Goroutines-16        	270261190	         4.386 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_64Goroutines-16        	278794970	         4.293 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_128Goroutines-16       	278603161	         4.382 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_2Goroutines-16         	218843803	         5.521 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_4Goroutines-16         	221350075	         5.433 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_8Goroutines-16         	222382471	         5.397 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_16Goroutines-16        	219519555	         5.406 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_32Goroutines-16        	227004950	         5.379 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_64Goroutines-16        	227031237	         5.277 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_128Goroutines-16       	227121474	         5.348 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_2Goroutines-16         	127397122	         9.380 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_4Goroutines-16         	100000000	        10.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_8Goroutines-16         	122956939	        10.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_16Goroutines-16        	126966850	        10.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_32Goroutines-16        	129644491	        10.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_64Goroutines-16        	129282783	        10.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_128Goroutines-16       	125461401	        10.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_2Goroutines-16        	30224193	        39.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_4Goroutines-16        	35712380	        33.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_8Goroutines-16        	33683935	        34.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_16Goroutines-16       	35639373	        33.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_32Goroutines-16       	35541601	        33.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_64Goroutines-16       	35728906	        33.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_128Goroutines-16      	35766399	        33.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_2Goroutines-16        	12077659	        99.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_4Goroutines-16        	12070339	       100.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_8Goroutines-16        	11984224	       100.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_16Goroutines-16       	11906582	       100.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_32Goroutines-16       	11994741	        99.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_64Goroutines-16       	12000150	       100.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_128Goroutines-16      	12027015	        99.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_2Goroutines-16        	 9029817	       133.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_4Goroutines-16        	 8991758	       133.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_8Goroutines-16        	 9017007	       133.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_16Goroutines-16       	 9056287	       133.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_32Goroutines-16       	 9042925	       132.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_64Goroutines-16       	 8995414	       133.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_128Goroutines-16      	 9036214	       133.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_2Goroutines-16       	 6952444	       170.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_4Goroutines-16       	 7050771	       171.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_8Goroutines-16       	 6942898	       171.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_16Goroutines-16      	 6965104	       149.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_32Goroutines-16      	 8411968	       153.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_64Goroutines-16      	 6987492	       171.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_128Goroutines-16     	 7269903	       165.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 6193324	       194.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 6222722	       192.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 6242034	       192.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 6262288	       191.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 6235977	       191.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 6252049	       191.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 6246038	       191.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1822958	       715.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1655086	       729.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1680040	       713.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1674666	       712.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1678710	       713.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1681701	       714.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1684285	       713.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16       	  523836	      2273 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16      	  217233	      5536 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16      	  132709	      8940 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16      	   33720	     35965 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16    	    2096	    573422 ns/op	       1 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 6246136	       191.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 6249640	       191.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 6237267	       191.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 6241344	       191.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 6243712	       191.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 6245071	       191.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 6246223	       191.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_2Goroutines-16        	 2660907	       450.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_4Goroutines-16        	 2659984	       450.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_8Goroutines-16        	 2666076	       450.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_16Goroutines-16       	 2662317	       451.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_32Goroutines-16       	 2657403	       450.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_64Goroutines-16       	 2654338	       454.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_128Goroutines-16      	 2644657	       453.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1671342	       718.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1665368	       719.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1666214	       718.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1670187	       717.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1665134	       719.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1667168	       673.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1673307	       717.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  439858	      2648 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  438632	      2650 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  445155	      2647 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  424273	      2648 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  443612	      2649 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  445390	      2663 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  437482	      2728 ns/op	       1 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   28538	     41800 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   28740	     41715 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   28932	     41625 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   28928	     41618 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   28845	     41723 ns/op	       4 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   28621	     41602 ns/op	      11 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   28884	     41682 ns/op	      13 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                	64308681	        17.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                               	49512098	        24.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                               	41784912	        28.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                               	38542195	        30.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                               	32833420	        36.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                               	24234753	        49.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                               	20117899	        59.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                              	11892736	       100.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                              	 7665237	       156.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                              	 4097682	       292.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                             	 2123215	       564.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                             	 1000000	      1109 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                             	  545251	      2198 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_2Goroutines-16      	367092969	         3.295 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_4Goroutines-16      	364329969	         3.245 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_8Goroutines-16      	375396412	         3.162 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_16Goroutines-16     	369868792	         3.248 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_32Goroutines-16     	393919362	         3.124 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_64Goroutines-16     	392266092	         3.061 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_128Goroutines-16    	391131579	         3.088 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	278666618	         4.332 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	288241862	         4.259 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	287041138	         4.222 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	293286248	         4.142 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	294782804	         4.168 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	288746611	         4.141 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	293468252	         4.105 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_2Goroutines-16     	235662571	         4.996 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_4Goroutines-16     	241617771	         4.966 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_8Goroutines-16     	242537506	         4.948 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_16Goroutines-16    	241134699	         4.926 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_32Goroutines-16    	244916493	         4.890 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_64Goroutines-16    	240585705	         4.941 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_128Goroutines-16   	243669334	         5.348 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_2Goroutines-16     	220361779	         5.614 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_4Goroutines-16     	219561060	         5.561 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_8Goroutines-16     	221110374	         5.543 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_16Goroutines-16    	224684581	         5.576 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_32Goroutines-16    	225895640	         5.641 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_64Goroutines-16    	224597391	         5.557 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_128Goroutines-16   	225518374	         5.489 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	175696671	         6.860 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	178425992	         6.801 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	177544383	         6.801 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	179291941	         6.799 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	179991619	         6.720 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	180617587	         6.638 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	181707596	         6.704 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_2Goroutines-16     	91803654	        11.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_4Goroutines-16     	99399460	        12.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_8Goroutines-16     	90868075	        11.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_16Goroutines-16    	99492520	        12.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_32Goroutines-16    	90354640	        11.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_64Goroutines-16    	99487020	        12.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_128Goroutines-16   	91575384	        11.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	98527893	        11.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	97946526	        12.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	96624192	        11.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	97551394	        11.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	98195654	        11.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	96930204	        11.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	98202685	        11.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	39349153	        31.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	38870016	        31.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	39067744	        30.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	38544464	        30.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	39133942	        30.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	39151446	        30.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	39215578	        30.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	10476978	       115.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	10387612	       115.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	10316472	       115.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	10372449	       115.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	10368865	       115.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	10381939	       115.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	10397126	       115.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	 8953236	       134.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	 8971179	       134.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	 8915416	       134.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	 8901182	       134.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	 9022682	       134.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	 8984938	       133.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	 8959327	       133.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	 8767714	       137.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	 8718279	       136.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	 8760200	       137.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	 8763199	       137.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	 8737177	       137.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	 8762700	       137.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	 8846889	       137.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	 9610599	       124.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	 9631668	       124.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	 9630931	       124.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	 9637422	       124.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	 9627994	       124.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	 9625603	       124.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	 9635002	       124.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 6169490	       193.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 6165784	       194.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 6195110	       193.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 6165146	       193.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 6176155	       194.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 6144906	       193.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 6158324	       193.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                           	     213	   5587094 ns/op	     225 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16           	    2432	    412737 ns/op	      15 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16           	    2569	    413217 ns/op	      15 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16           	    2936	    407516 ns/op	      27 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16          	    2586	    414763 ns/op	      47 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16          	    2446	    414047 ns/op	      49 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16          	    2479	    409569 ns/op	      99 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16         	    2547	    416399 ns/op	     130 B/op	       2 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                           	      42	  28060090 ns/op	      67 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16           	     566	   2088432 ns/op	      38 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16           	     480	   2128022 ns/op	      85 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16           	     516	   2106695 ns/op	      84 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16          	     486	   2083866 ns/op	     154 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16          	     489	   2110314 ns/op	     240 B/op	       3 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16          	     556	   2098568 ns/op	     344 B/op	       5 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16         	     483	   2110538 ns/op	     731 B/op	      12 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                          	      20	  55798444 ns/op	     140 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16          	     283	   4151640 ns/op	      76 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16          	     288	   4139970 ns/op	     133 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16          	     283	   4153726 ns/op	     232 B/op	       2 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16         	     283	   4134685 ns/op	     249 B/op	       3 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16         	     285	   4149505 ns/op	     661 B/op	       7 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16         	     285	   4124672 ns/op	     833 B/op	      11 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16        	     284	   4138645 ns/op	    1123 B/op	      18 allocs/op
PASS
ok  	github.com/sixafter/prng-chacha	377.422s
```
</details>

### UUID Generation

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the CSPRNG-based UUID generation:

| Benchmark Scenario                         | Default ns/op | CSPRNG ns/op | % Faster (ns/op) | Default B/op | CSPRNG B/op | Default allocs/op | CSPRNG allocs/op |
|--------------------------------------------|--------------:|-------------:|-----------------:|-------------:|------------:|------------------:|-----------------:|
| v4 Serial                                 |      183.6    |     37.70    |      79.5%       |      16      |     16      |      1            |      1           |
| v4 Parallel                               |      457.2    |      5.871   |      98.7%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (4 goroutines)              |      419.2    |     11.36    |      97.3%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (8 goroutines)              |      482.1    |      7.712   |      98.4%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (16 goroutines)             |      455.6    |      5.944   |      98.7%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (32 goroutines)             |      521.1    |      5.788   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (64 goroutines)             |      533.1    |      5.735   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (128 goroutines)            |      523.4    |      5.705   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (256 goroutines)            |      523.7    |      5.794   |      98.9%       |      16      |     16      |      1            |      1           |

<details>
  <summary>Expand to see results</summary>

```shell
make bench-uuid

goos: darwin
goarch: arm64
pkg: github.com/sixafter/prng-chacha
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6239547	       183.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2614206	       457.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2867928	       419.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2520130	       482.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2617567	       455.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2312065	       521.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2300226	       533.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2300107	       523.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2331600	       523.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Serial-16                           	30584091	        37.70 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Parallel-16                         	209297205	         5.871 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_4-16          	100000000	        11.36 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_8-16          	150149610	         7.712 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_16-16         	203733687	         5.944 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_32-16         	205883962	         5.788 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_64-16         	208636114	         5.735 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_128-16        	212263852	         5.705 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_256-16        	204421857	         5.794 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/prng-chacha	31.142s
```
</details>

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
