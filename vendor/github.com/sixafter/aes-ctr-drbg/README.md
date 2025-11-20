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

## NIST SP 800-90A Compliance

For a detailed mapping between the implementation and NIST SP 800-90A requirements, see [NIST-SP-800-90A.md](docs/NIST-SP-800-90A.md).

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the release tarball, you can use Cosign to check the signature and checksums. Follow these steps:

```sh
# Fetch the latest release tag from GitHub API (e.g., "v1.14.0")
TAG=$(curl -s https://api.github.com/repos/sixafter/aes-ctr-drbg/releases/latest | jq -r .tag_name)

# Remove leading "v" for filenames (e.g., "v1.14.0" -> "1.14.0")
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
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G2-16  	1000000000	         0.6696 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G4-16  	1000000000	         0.6740 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G8-16  	1000000000	         0.6520 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G16-16 	1000000000	         0.6885 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G32-16 	1000000000	         0.6353 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G64-16 	1000000000	         0.6401 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G128-16         	1000000000	         0.6474 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16Bytes-16           	29874123	        36.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_32Bytes-16           	28411921	        42.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_64Bytes-16           	22326446	        53.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_256Bytes-16          	 9796570	       122.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_512Bytes-16          	 5619010	       213.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_4096Bytes-16         	  798786	      1470 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16384Bytes-16        	  203523	      5781 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_2Goroutines-16         	19856755	        86.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_4Goroutines-16         	19652169	        84.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_8Goroutines-16         	19535383	        88.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_16Goroutines-16        	19280488	        88.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_32Goroutines-16        	18911883	        88.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_64Goroutines-16        	20954308	        81.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_128Goroutines-16       	20513842	        73.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_2Goroutines-16         	18478191	        75.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_4Goroutines-16         	19357531	        74.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_8Goroutines-16         	20106480	        86.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_16Goroutines-16        	19378123	        74.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_32Goroutines-16        	19684783	        77.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_64Goroutines-16        	19942705	        78.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_128Goroutines-16       	20817882	        60.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_2Goroutines-16         	19866068	        75.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_4Goroutines-16         	19591770	        87.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_8Goroutines-16         	19219963	        76.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_16Goroutines-16        	19804158	        79.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_32Goroutines-16        	19630122	        67.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_64Goroutines-16        	19981904	        78.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_128Goroutines-16       	19985551	        81.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_2Goroutines-16        	10538402	       160.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_4Goroutines-16        	10828724	       161.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_8Goroutines-16        	10997553	       158.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_16Goroutines-16       	10632771	       161.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_32Goroutines-16       	10377925	       160.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_64Goroutines-16       	10851642	       148.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_128Goroutines-16      	11019211	       158.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_2Goroutines-16        	10556665	       153.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_4Goroutines-16        	10974860	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_8Goroutines-16        	11132241	       156.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_16Goroutines-16       	10818664	       154.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_32Goroutines-16       	10574911	       162.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_64Goroutines-16       	11106015	       158.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_128Goroutines-16      	11131617	       157.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7474005	       221.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 6751184	       219.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 5899460	       220.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 6126216	       220.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 6938299	       220.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7048448	       215.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7366938	       206.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1439043	       826.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1478479	       851.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1434871	       812.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1692771	       682.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1452622	       847.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1626038	       844.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1523498	       847.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_4096Bytes-16      	  778651	      1476 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_16384Bytes-16     	  208345	      5873 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_65536Bytes-16     	   51400	     23373 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_1048576Bytes-16   	    3177	    381062 ns/op	      11 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7376702	       221.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 5193856	       213.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 6114523	       224.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7188720	       205.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7138792	       217.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 5875122	       221.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 6167787	       219.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1475030	       831.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1454176	       806.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1561255	       809.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1520268	       835.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1454516	       839.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1456448	       837.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1503434	       807.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  481670	      2853 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  436726	      2932 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  467893	      2903 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  473342	      2837 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  474193	      2911 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  462038	      2924 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  475135	      2924 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   29137	     42209 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   29096	     42039 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   29031	     41666 ns/op	       4 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   28812	     41482 ns/op	       7 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   29179	     41477 ns/op	      10 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   29031	     42093 ns/op	      17 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   28680	     42144 ns/op	      20 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_16Bytes-16                                	32761209	        35.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_32Bytes-16                                	28607130	        41.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_64Bytes-16                                	22548178	        53.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_128Bytes-16                               	15760812	        76.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_256Bytes-16                               	 9770961	       122.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_512Bytes-16                               	 5631610	       213.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_1024Bytes-16                              	 3079610	       395.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1599415	       757.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_4096Bytes-16                              	  805844	      1500 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	18766542	        90.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	19131636	        87.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	18264330	        88.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	19677601	        83.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	19468838	        86.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	19406863	        86.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	18875232	        85.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	18813998	        70.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	19586839	        73.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	19415864	        87.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	18120147	        84.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	19511706	        81.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	19756026	        85.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	19666300	        84.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	18756825	        89.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	16706818	        80.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	19674012	        86.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	17465998	        87.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	19747316	        82.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	19630255	        77.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	20142564	        70.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	11142530	       143.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	10331701	       144.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	11674490	       129.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	 8749474	       124.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	11025412	       117.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	11193232	       116.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	10937000	       108.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	 9248961	       111.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	11058454	       134.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	11100536	       121.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	10285585	       115.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	10852808	       107.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	 9899952	       123.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	11235616	       130.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	10088126	       110.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	 8867764	       147.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	 8720042	       115.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11546450	       120.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	10706088	       119.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	10183741	       140.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11463069	       103.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	12553362	        95.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	12499902	        95.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	12535980	        93.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	12506155	        95.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	12727627	        93.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	12669898	        94.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12720594	        93.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	12157644	        99.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	12303175	       100.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	11785392	       100.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	11959412	       118.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	12079330	       109.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	11362237	       103.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	11637774	       120.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 6880490	       168.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7228526	       168.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7195345	       169.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7188027	       168.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7118646	       168.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7124090	       172.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7097014	       169.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     295	   3929851 ns/op	     224 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16            	    3172	    381380 ns/op	      91 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16            	    2817	    380624 ns/op	     117 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16            	    3079	    384180 ns/op	     127 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16           	    3080	    385891 ns/op	     159 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16           	    2972	    398988 ns/op	     156 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16           	    2919	    403091 ns/op	     167 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16          	    2972	    397429 ns/op	     284 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      54	  19975266 ns/op	    1222 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16            	     608	   1995378 ns/op	     352 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16            	     553	   2018336 ns/op	     428 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16            	     566	   2034392 ns/op	     511 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16           	     558	   1969891 ns/op	     547 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16           	     588	   1971865 ns/op	     505 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16           	     571	   2006916 ns/op	     914 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16          	     573	   1977134 ns/op	    1045 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      28	  39731435 ns/op	    1663 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16           	     286	   3826714 ns/op	     619 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16           	     266	   3791460 ns/op	     740 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16           	     274	   3849177 ns/op	     752 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16          	     278	   3827166 ns/op	     846 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16          	     292	   3948489 ns/op	    1099 B/op	       8 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16          	     300	   3867839 ns/op	    1226 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16         	     289	   3847689 ns/op	    1528 B/op	      18 allocs/op
BenchmarkDRBG_Read_WithKeyRotation-16                                                           	 5003805	       238.5 ns/op	     198 B/op	       1 allocs/op
BenchmarkDRBG_Read_PredictionResistance-16                                                      	 2632755	       456.1 ns/op	     634 B/op	       3 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	291.367s
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
BenchmarkUUID_v4_Default_Serial-16        	 6473760	       180.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2705866	       445.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16         	 2883284	       413.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2806682	       428.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2462146	       484.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2685201	       458.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2366074	       506.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2358429	       506.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2388648	       508.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2364384	       511.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	29120706	        40.78 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	100000000	        10.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_2-16         	52686843	        21.91 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	92968908	        12.77 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	121979662	         9.741 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	153623710	         7.668 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	154797238	         7.688 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	156757164	         7.641 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	156462766	         7.632 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	154197008	         7.795 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	33.515s
  ```
</details>

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
