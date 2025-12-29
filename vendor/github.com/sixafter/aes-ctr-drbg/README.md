# AES-CTR-DRBG

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/aes-ctr-drbg)](https://goreportcard.com/report/github.com/sixafter/aes-ctr-drbg)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/aes-ctr-drbg)](https://img.shields.io/github/go-mod/go-version/sixafter/aes-ctr-drbg)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/aes-ctr-drbg.svg)](https://pkg.go.dev/github.com/sixafter/aes-ctr-drbg)
[![FIPS‑140 Mode Compatible](https://img.shields.io/badge/FIPS‑140--Mode-Compatible-brightgreen)](FIPS‑140.md)

---

## Status

### Build & Test

[![CI](https://github.com/sixafter/aes-ctr-drbg/workflows/ci/badge.svg)](https://github.com/sixafter/aes-ctr-drbg/actions)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/aes-ctr-drbg)](https://github.com/sixafter/aes-ctr-drbg/issues)
![GitHub last commit](https://img.shields.io/github/last-commit/sixafter/aes-ctr-drbg)

### Quality

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_aes-ctr-drbg&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_aes-ctr-drbg)
![CodeQL](https://github.com/sixafter/aes-ctr-drbg/actions/workflows/codeql-analysis.yaml/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=six-after_aes-ctr-drbg&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=six-after_aes-ctr-drbg)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/11487/badge)](https://www.bestpractices.dev/projects/11487)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sixafter/aes-ctr-drbg/badge)](https://scorecard.dev/viewer/?uri=github.com/sixafter/aes-ctr-drbg)

### Package and Deploy

[![Release](https://github.com/sixafter/aes-ctr-drbg/workflows/release/badge.svg)](https://github.com/sixafter/aes-ctr-drbg/actions)

---
## Overview 

AES-CTR-DRBG (Deterministic Random Bit Generator based on AES in Counter mode) is a cryptographically secure pseudo-random number generator (CSPRNG) defined by [NIST SP 800-90A Rev. 1](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final). It's widely used in high-assurance systems, including those requiring [FIPS 140-2](https://csrc.nist.gov/pubs/fips/140-2/upd2/final) compliance. 
AES-CTR-DRBG is designed for environments requiring deterministic, reproducible, and FIPS‑140-compatible random bit generation. This module is suitable for any application that needs strong cryptographic assurance or must comply with regulated environments (e.g., FedRAMP, FIPS, PCI, HIPAA). 

The module uses only Go standard library crypto primitives (`crypto/aes` and `crypto/cipher`), making it safe for use in FIPS 140-validated Go runtimes. No third-party, homegrown, or experimental ciphers are used.

Please see the [godoc](https://pkg.go.dev/github.com/sixafter/aes-ctr-drbg) for detailed documentation.

---

## FIPS‑140 Mode

See [FIPS‑140.md](FIPS-140.md) for compliance, deployment, and configuration guidance.

---

## Features

* **Standards-Compliant Implementation:**
  Implements NIST SP 800-90A, Revision 1 AES-CTR-DRBG using the Go standard library (`crypto/aes`, `crypto/cipher`). Supports 128-, 192-, and 256-bit keys. State and counter management strictly adhere to the specification.

* **FIPS 140-2 Alignment:**
  Designed for use in FIPS 140-2 validated environments and compatible with Go’s FIPS 140 mode (`GODEBUG=fips140=on`). See [FIPS-140.md](FIPS-140.md) for platform guidance.

* **Zero-Allocation Output Path:**
  The DRBG is engineered for `0 allocs/op` in its standard `io.Reader` output path, enabling predictable resource usage and high throughput.

* **Asynchronous Key Rotation:**
  Supports automatic key rotation after a configurable number of bytes have been generated (`MaxBytesPerKey`). Rekeying occurs asynchronously with exponential backoff and configurable retry limits, reducing long-term key exposure.

* **Prediction Resistance Mode:**
  Supports NIST SP 800-90A prediction resistance. When enabled, the DRBG reseeds from system entropy before every output, as required for state compromise resilience.

* **Sharded Pooling for Concurrency:**
  Internal state pooling can be sharded across multiple `sync.Pool` instances. The number of shards is configurable, allowing improved performance under concurrent workloads.

* **Extensive Functional Configuration:**
  Exposes a comprehensive set of functional options, including:

  * AES key size (128/192/256-bit)
  * Maximum output per key (rekey threshold)
  * Personalization string (domain separation)
  * Shard/pool count
  * Reseed interval and request count
  * Buffer size controls
  * Key rotation and rekey backoff parameters
  * Prediction resistance
  * Fork detection and reseeding

* **Thread-Safe and Deterministic:**
  All DRBG instances are safe for concurrent use. Output is deterministic for a given seed and personalization.

* **io.Reader Compatibility:**
  Implements Go’s `io.Reader` interface for drop-in use as a secure random source.

* **No External Dependencies:**
  Depends exclusively on the Go standard library for cryptographic operations.

* **UUID Generation:**
  Can be used as a cryptographically secure `io.Reader` with the [`google/uuid`](https://pkg.go.dev/github.com/google/uuid) package and similar libraries.

* **Comprehensive Testing and Fuzzing:**
  Includes property-based, fuzz, concurrency, and allocation tests to validate correctness, robustness, and allocation characteristics.

* **Fork-Safety:**
  Automatic detection and reseeding on process fork. This library automatically detects process forks and reseeds in the child process to prevent random stream duplication. No manual action is required.

For an example of how this library can be consumed in practice, see [sixafter/nanoid](https://github.com/sixafter/nanoid).

NanoID uses [sixafter/prng-chacha](https://github.com/sixafter/prng-chacha) as its default high-performance RNG, but
can also use **`aes-ctr-drbg`** when FIPS 140-2/3 alignment or deterministic AES-CTR-based randomness is required. This provides a clear, real-world example of integrating this DRBG into an ID-generation workflow,
including optional FIPS-centric operation. There is also a [WithAutoRandReader()](https://github.com/sixafter/nanoid/blob/d8efbc63e5a5696a33f34b9fb5d24f5d7805a7ed/config.go#L229) option that automatically selects 
between `prng-chacha` and `aes-ctr-drbg` based on the runtime FIPS mode.

## NIST SP 800-90A Compliance

For a detailed mapping between the implementation and NIST SP 800-90A requirements, see [NIST-SP-800-90A.md](docs/NIST-SP-800-90A.md).

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the release tarball, you can use Cosign to check the signature and checksums. Follow these steps:

```sh
# Fetch the latest release tag from GitHub API (e.g., "v1.15.0")
TAG=$(curl -s https://api.github.com/repos/sixafter/aes-ctr-drbg/releases/latest | jq -r .tag_name)

# Remove leading "v" for filenames (e.g., "v1.15.0" -> "1.15.0")
VERSION=${TAG#v}

# ---------------------------------------------------------------------
# Verify the source archive using Sigstore bundles
# ---------------------------------------------------------------------

# Download the release tarball and its signature bundle
curl -LO "https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/aes-ctr-drbg-${VERSION}.tar.gz"
curl -LO "https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/aes-ctr-drbg-${VERSION}.tar.gz.sigstore.json"

# Verify the tarball with Cosign using the published public key
cosign verify-blob \
  --key "https://raw.githubusercontent.com/sixafter/aes-ctr-drbg/main/cosign.pub" \
  --bundle "aes-ctr-drbg-${VERSION}.tar.gz.sigstore.json" \
  "aes-ctr-drbg-${VERSION}.tar.gz"

# ---------------------------------------------------------------------
# Verify the checksums manifest using Sigstore bundles
# ---------------------------------------------------------------------

curl -LO "https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/checksums.txt"
curl -LO "https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/checksums.txt.sigstore.json"

cosign verify-blob \
  --key "https://raw.githubusercontent.com/sixafter/aes-ctr-drbg/main/cosign.pub" \
  --bundle "checksums.txt.sigstore.json" \
  "checksums.txt"

# ---------------------------------------------------------------------
# Confirm local artifact integrity
# ---------------------------------------------------------------------

shasum -a 256 -c checksums.txt

```

If valid, Cosign will output:

```shell
Verified OK
```

---

## Verify Go module

To validate that the Go module archive served by GitHub, go mod download, and the Go 
proxy are all consistent, run the `module-verify` target. This performs a full cross-check 
of the tag archive and module ZIPs to confirm they match byte-for-byte.

---

## Installation

```bash
go get -u github.com/sixafter/aes-ctr-drbg
```

---

## Usage

### Basic Usage: Generate Secure Random Bytes With Reader

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	buf := make([]byte, 64)
	n, err := ctrdrbg.Reader.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

### Basic Usage: Generate Secure Random Bytes with NewReader

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	// Example: AES-256 (32 bytes) key
	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(ctrdrbg.KeySize256))
	if err != nil {
		log.Fatalf("failed to create ctrdrbg.Reader: %v", err)
	}

	buf := make([]byte, 64)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

### Using Personalization and Additional Input

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	r, err := ctrdrbg.NewReader(
		ctrdrbg.WithPersonalization([]byte("service-id-1")),
		ctrdrbg.WithKeySize(ctrdrbg.KeySize256), // AES-256
	)
	if err != nil {
		log.Fatalf("failed to create ctrdrbg.Reader: %v", err)
	}

	buf := make([]byte, 64)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

---

## Performance Benchmarks

### Raw Random Byte Generation

These `ctrdrbg.Reader` benchmarks demonstrate the package's focus on minimizing latency, memory usage, and allocation overhead, making it suitable for high-performance applications.

<details>
  <summary>Expand to see results</summary>

```shell
make bench
go test -bench='^BenchmarkDRBG_' -run=^$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .
goos: darwin
goarch: arm64
pkg: github.com/sixafter/aes-ctr-drbg
cpu: Apple M4 Max
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G2-16  	1000000000	         0.6112 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G4-16  	1000000000	         0.6134 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G8-16  	1000000000	         0.5971 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G16-16 	1000000000	         0.5937 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G32-16 	1000000000	         0.5700 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G64-16 	1000000000	         0.5587 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G128-16         	1000000000	         0.5644 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16Bytes-16           	34062441	        34.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_32Bytes-16           	29895489	        39.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_64Bytes-16           	24042694	        50.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_256Bytes-16          	10366684	       115.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_512Bytes-16          	 5998489	       200.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_4096Bytes-16         	  871158	      1377 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16384Bytes-16        	  222878	      5376 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_2Goroutines-16         	19962072	        85.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_4Goroutines-16         	19596702	        84.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_8Goroutines-16         	19851349	        80.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_16Goroutines-16        	19654436	        64.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_32Goroutines-16        	19521876	        81.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_64Goroutines-16        	19818331	        82.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_128Goroutines-16       	20155250	        82.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_2Goroutines-16         	18829485	        88.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_4Goroutines-16         	19666596	        88.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_8Goroutines-16         	19722001	        83.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_16Goroutines-16        	19440373	        85.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_32Goroutines-16        	19804376	        84.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_64Goroutines-16        	19879699	        84.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_128Goroutines-16       	19940580	        82.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_2Goroutines-16         	18446626	        80.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_4Goroutines-16         	19594609	        84.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_8Goroutines-16         	19475618	        84.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_16Goroutines-16        	19886919	        85.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_32Goroutines-16        	19558533	        81.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_64Goroutines-16        	19842116	        82.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_128Goroutines-16       	20000374	        83.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_2Goroutines-16        	10465203	       161.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_4Goroutines-16        	10905505	       158.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_8Goroutines-16        	10895772	       154.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_16Goroutines-16       	10901005	       156.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_32Goroutines-16       	10952331	       151.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_64Goroutines-16       	11080442	       156.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_128Goroutines-16      	11141888	       156.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_2Goroutines-16        	10665208	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_4Goroutines-16        	10822701	       161.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_8Goroutines-16        	10816222	       155.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_16Goroutines-16       	11038309	       148.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_32Goroutines-16       	10907218	       156.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_64Goroutines-16       	11146635	       156.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_128Goroutines-16      	11120167	       154.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7884975	       205.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7952276	       203.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 6754220	       200.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 6395343	       191.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 6827542	       203.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7638505	       197.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7684468	       187.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1792022	       730.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1782008	       752.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1811017	       811.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1553781	       773.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1644704	       775.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1553227	       758.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1627238	       789.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_4096Bytes-16      	  724966	      1396 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_16384Bytes-16     	  217694	      5471 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_65536Bytes-16     	   56064	     21448 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7686088	       190.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7127140	       193.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 6847946	       203.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7394517	       202.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7078538	       198.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7311069	       213.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7093018	       199.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1553928	       787.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1505232	       806.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1586379	       824.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1567017	       834.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1515153	       826.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1498590	       835.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1637242	       823.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  500164	      2770 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  498936	      2767 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  508104	      2763 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  486643	      2768 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  493526	      2734 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  511357	      2723 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  512506	      2734 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_16Bytes-16                                	34138875	        33.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_32Bytes-16                                	30298119	        39.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_64Bytes-16                                	24156391	        50.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_128Bytes-16                               	16411648	        71.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_256Bytes-16                               	10339018	       116.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_512Bytes-16                               	 6003164	       200.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_1024Bytes-16                              	 3275638	       368.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1683925	       717.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_4096Bytes-16                              	  828991	      1396 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	19667361	        81.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	19478792	        82.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	20362839	        77.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	19621669	        79.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	19656086	        85.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	20029543	        73.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	20244778	        80.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	18923464	        86.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	19499445	        83.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	19745245	        85.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	19773741	        81.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	19825248	        84.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	19899246	        82.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	19991892	        83.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	19272760	        85.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	19945840	        84.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	19978730	        86.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	19746584	        82.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	20115018	        83.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	19668866	        80.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	20351356	        82.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	10806562	       156.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	10851708	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	10542328	       155.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	10863722	       155.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	10945759	       147.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	11071387	       153.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	11136016	       155.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	10442144	       158.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	10719852	       156.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	10798916	       156.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	10926900	       157.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	10761366	       143.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	11546500	       122.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	12631794	       114.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	10984488	       156.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	11054596	       153.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	 9367392	       148.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11085548	       149.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	10982490	       151.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	11193729	       145.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11057155	       153.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	12433735	       140.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	12563049	       142.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	12526951	       137.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	12434395	       137.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	12830301	       136.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	12661705	       146.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12534420	       147.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	12307370	       124.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	12567220	       105.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	12071736	       100.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	11801464	       106.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	12380242	       115.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	12019501	       110.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	12071012	       115.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 6966630	       220.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7283623	       222.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7271518	       216.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7153774	       217.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7220374	       215.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 6965928	       197.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 6695487	       211.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     328	   3549165 ns/op	     163 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16            	    3732	    339168 ns/op	      91 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16            	    3270	    346786 ns/op	     113 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16            	    3276	    346405 ns/op	     133 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16           	    3267	    349331 ns/op	     149 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16           	    3325	    349511 ns/op	     177 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16           	    3267	    354763 ns/op	     262 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16          	    3253	    355026 ns/op	     277 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      61	  17701644 ns/op	    1019 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16            	     717	   1717468 ns/op	     457 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16            	     645	   1733391 ns/op	     575 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16            	     656	   1740158 ns/op	     608 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16           	     642	   1738901 ns/op	     763 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16           	     625	   1764412 ns/op	     754 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16           	     660	   1753616 ns/op	    1038 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16          	     627	   1769634 ns/op	    1249 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      31	  35179063 ns/op	    1601 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16           	     338	   3426444 ns/op	     973 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16           	     336	   3494822 ns/op	    1085 B/op	       6 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16           	     314	   3478031 ns/op	    1293 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16          	     338	   3507727 ns/op	    1410 B/op	       8 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16          	     320	   3539974 ns/op	    1687 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16          	     296	   3483667 ns/op	    2022 B/op	      15 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16         	     289	   3470399 ns/op	    2329 B/op	      23 allocs/op
BenchmarkDRBG_Read_WithKeyRotation-16                                                           	 5215118	       227.6 ns/op	     196 B/op	       1 allocs/op
BenchmarkDRBG_Read_PredictionResistance-16                                                      	 2658322	       450.1 ns/op	     633 B/op	       3 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	293.330s
```

</details>

### UUID Generation with Google UUID and ctrdrbg

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the ctrdrbg-based UUID generation:

| Benchmark Scenario                  | Default ns/op | CTRDRBG ns/op | % Faster (ns/op) | Default B/op | CTRDRBG B/op | Default allocs/op | CTRDRBG allocs/op |
|-------------------------------------|--------------:|--------------:|-----------------:|-------------:|-------------:|------------------:|------------------:|
| v4 Serial                           |        180.6  |        40.78  |         77.4%    |         16   |         16   |                1  |                1  |
| v4 Parallel                         |        445.4  |        10.56  |         97.6%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (2 goroutines)        |        413.2  |        21.91  |         94.7%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (4 goroutines)        |        428.5  |        12.77  |         97.0%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (8 goroutines)        |        484.6  |         9.74  |         98.0%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (16 goroutines)       |        458.2  |         7.67  |         98.3%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (32 goroutines)       |        506.3  |         7.69  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (64 goroutines)       |        506.9  |         7.64  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (128 goroutines)      |        508.2  |         7.63  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (256 goroutines)      |        511.8  |         7.79  |         98.5%    |         16   |         16   |                1  |                1  |

Notes:
- "Default" refers to the baseline Go `crypto/rand` source.
- "CTRDRBG" refers to this AES-CTR-DRBG implementation.
- "% Faster (ns/op)" is computed as `100 * (Default - CTRDRBG) / Default`, rounded.

<details>
  <summary>Expand to see results</summary>

```shell
make bench-uuid
go test -bench='^BenchmarkUUID_' -run=^$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .
goos: darwin
goarch: arm64
pkg: github.com/sixafter/aes-ctr-drbg
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 5764681	       207.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2523262	       467.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16         	 3328860	       360.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2731453	       436.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2447664	       491.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2558335	       473.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2366698	       510.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2340867	       511.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2345401	       510.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2329904	       515.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	23044700	        50.41 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	94280326	        12.50 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_2-16         	43251036	        26.58 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	80453218	        14.85 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	100000000	        10.95 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	139321384	         8.189 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	146813631	         8.137 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	144867888	         8.253 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	147125661	         8.219 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	143012491	         8.396 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	32.857s
  ```
</details>

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
