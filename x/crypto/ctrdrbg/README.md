# ctrdrbg: AES-CTR-DRBG for Deterministic Cryptographically Secure Random Number Generation

## Overview

The `ctrdrbg` package implements an [AES-CTR-DRBG (Deterministic Random Bit Generator)](https://csrc.nist.gov/publications/detail/sp/800-90a/rev-1/final) as specified in NIST SP 800-90A.  
It is designed for environments requiring deterministic, reproducible, and **FIPS‑140-compatible** random bit generation.  
This package is suitable for any application that needs strong cryptographic assurance or must comply with regulated environments (e.g., FedRAMP, FIPS, HIPAA).

The package uses only Go standard library crypto primitives (`crypto/aes` and `crypto/cipher`), making it safe for use in FIPS 140-validated Go runtimes.  
No third-party, homegrown, or experimental ciphers are used.

---

## Features

- **NIST SP 800-90A AES-CTR-DRBG Implementation:** Implements the Deterministic Random Bit Generator (DRBG) construction defined in [NIST SP 800-90A, Revision 1](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final), using the AES block cipher in counter (CTR) mode.
  - Supports 128-, 192-, and 256-bit AES keys, with correct counter and state management as specified by the standard.
- **FIPS 140-2 Alignment:** Designed to operate in FIPS 140-2 validated environments and compatible with Go’s FIPS 140 mode (`GODEBUG=fips140=on`).
  - Uses only cryptographic primitives from the Go standard library.
  - For platform-specific guidance and deployment instructions, see [FIPS‑140.md](../../../FIPS-140.md).
- **Stateless and Concurrent Operation:** Each DRBG instance is safe for concurrent use and fully encapsulates its cryptographic state. 
  - The design supports independent operation across multiple instances, enabling scalable use in high-concurrency environments.
- **Configurable Entropy and Personalization:** Accepts externally supplied entropy sources and personalization strings, enabling domain separation and deterministic output for compliance with best practices and advanced use cases.
- **io.Reader Compatibility:** Fully satisfies Go’s `io.Reader` interface, allowing seamless integration with packages and APIs expecting a secure random source.
- **No External Dependencies:** Depends solely on the Go standard library, ensuring a lightweight and portable implementation.
- **UUID Generation Source:** Can be used as the `io.Reader` source for UUID generation with the [`google/uuid`](https://pkg.go.dev/github.com/google/uuid) package and similar libraries, providing cryptographically secure, deterministic UUIDs using AES-CTR-DRBG.

---

### NIST SP 800-90A AES-CTR-DRBG: Specification Mapping

| NIST SP 800-90A Requirement                                                                 | Implementation Reference                                   | Construction Step                                                                                           |
|---------------------------------------------------------------------------------------------|-----------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| **1. Instantiate: Acquire entropy and set initial state (`Key` and `V`)**                   | `newDRBG()` uses `io.ReadFull(rand.Reader, ...)`          | - Entropy input of `KeySize + 16` bytes is split into key and counter (V)                                  |
|                                                                                             |                                                           | - Personalization string (if provided) XORs into seed                                                      |
|                                                                                             |                                                           | - AES cipher constructed using key                                                                         |
|                                                                                             |                                                           | - Initial counter (V) set from entropy                                                                     |
| **2. Generate: For each output block, increment counter and encrypt**                        | `fillBlocks()` with `incV(v)` and `st.block.Encrypt(...)` | - For each 16-byte output block:                                                                           |
|                                                                                             |                                                           |   - Increment counter (V) (big-endian)                                                                     |
|                                                                                             |                                                           |   - Encrypt V using AES-CTR (AES block cipher in counter mode)                                             |
|                                                                                             |                                                           |   - Write result to output buffer                                                                          |
| **3. Update State After Generation**                                                        | `Read()` and `fillBlocks()` copy updated counter to state | - After output, updated counter (V) is copied back to DRBG instance                                        |
|                                                                                             |                                                           | - Mutex ensures exclusive access to counter                                                                |
| **4. Rekey/Reseed (Optional/Configurable):**                                                | `asyncRekey()` and usage logic                            | - Supports automatic rekeying after configurable number of bytes generated (`MaxBytesPerKey`)               |
|                                                                                             |                                                           | - New entropy acquired for new key and V                                                                   |
|                                                                                             |                                                           | - State atomically swapped; usage counter reset                                                            |
| **5. Personalization Support (Optional):**                                                  | `newDRBG()` and rekey use personalization                 | - Personalization string incorporated at instantiation and rekey                                            |
| **6. Edge Cases and Robustness:**                                                           | Test suite; zero/overflow logic                           | - Zero-length reads are no-ops                                                                             |
|                                                                                             |                                                           | - Counter overflow (wrap) is supported                                                                     |
|                                                                                             |                                                           | - Large and unaligned read sizes are supported                                                             |
| **7. Error Handling:**                                                                      | Error returns/panics for entropy/cipher errors            | - Instantiation returns error or panics if entropy/cipher initialization fails                             |
|                                                                                             |                                                           | - Rekey silently continues with prior state if entropy is unavailable                                      |
| **8. Concurrency:**                                                                         | Mutexes and per-instance state                            | - Mutex on counter ensures thread safety per DRBG instance                                                 |
|                                                                                             |                                                           | - Pooling enables concurrent DRBG use                                                                      |
| **9. Interface and Integration:**                                                           | `io.Reader` interface                                     | - DRBG implements `io.Reader` for use with standard Go APIs                                                |
| **10. No External Dependencies:**                                                           | Go standard library only                                  | - Implementation relies solely on standard Go cryptography primitives                                      |

---

## Installation

```bash
go get -u github.com/sixafter/nanoid/x/crypto/ctrdrbg
```

---

## Usage

### Basic Usage: Generate Secure Random Bytes With Reader

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/nanoid/x/crypto/ctrdrbg"
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

	"github.com/sixafter/nanoid/x/crypto/ctrdrbg"
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

	"github.com/sixafter/nanoid/x/crypto/ctrdrbg"
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
make bench-ctrdrbg
go test -bench='^BenchmarkDRBG_' -run=^$ -benchmem -memprofile=x/crypto/ctrdrbg/mem.out -cpuprofile=x/crypto/ctrdrbg/cpu.out ./x/crypto/ctrdrbg
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/ctrdrbg
cpu: Apple M4 Max
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G2-16  	1000000000	         0.6145 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G4-16  	1000000000	         0.6139 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G8-16  	1000000000	         0.5979 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G16-16 	1000000000	         0.5801 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G32-16 	1000000000	         0.5580 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G64-16 	1000000000	         0.5529 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G128-16         	1000000000	         0.5536 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16Bytes-16           	35978232	        32.05 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_32Bytes-16           	32421807	        36.94 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_64Bytes-16           	25152178	        47.63 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_256Bytes-16          	10319107	       116.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_512Bytes-16          	 5849815	       204.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_4096Bytes-16         	  818616	      1469 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16384Bytes-16        	  210074	      5657 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_2Goroutines-16         	19930326	        60.06 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_4Goroutines-16         	20700004	        69.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_8Goroutines-16         	20524894	        72.93 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_16Goroutines-16        	19770985	        77.23 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_32Goroutines-16        	20264736	        76.93 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_64Goroutines-16        	20141832	        78.78 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_128Goroutines-16       	19272348	        78.76 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_2Goroutines-16         	19453017	        75.79 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_4Goroutines-16         	19573448	        75.58 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_8Goroutines-16         	19479385	        79.94 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_16Goroutines-16        	19759102	        81.23 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_32Goroutines-16        	19332841	        80.46 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_64Goroutines-16        	20000319	        82.66 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_128Goroutines-16       	19738478	        77.21 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_2Goroutines-16         	19972384	        63.20 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_4Goroutines-16         	19738965	        71.55 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_8Goroutines-16         	19914256	        75.55 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_16Goroutines-16        	19793854	        82.11 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_32Goroutines-16        	19620238	        79.38 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_64Goroutines-16        	19988215	        82.49 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_128Goroutines-16       	19596702	        80.33 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_2Goroutines-16        	11146946	       142.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_4Goroutines-16        	11029047	       149.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_8Goroutines-16        	11073008	       148.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_16Goroutines-16       	11125928	       149.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_32Goroutines-16       	11119866	       148.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_64Goroutines-16       	11108440	       148.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_128Goroutines-16      	11207768	       150.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_2Goroutines-16        	11245017	       148.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_4Goroutines-16        	11143082	       151.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_8Goroutines-16        	11175612	       142.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_16Goroutines-16       	11041831	       150.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_32Goroutines-16       	11197646	       147.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_64Goroutines-16       	10933528	       148.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_128Goroutines-16      	10971085	       149.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7863334	       217.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7918147	       212.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7838018	       218.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7789429	       215.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7673707	       216.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7541233	       212.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7396773	       216.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1558494	       839.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1654434	       773.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1497637	       845.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1555300	       849.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1550515	       854.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1567885	       852.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1763186	       852.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_4096Bytes-16      	  749710	      1471 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_16384Bytes-16     	  209536	      5668 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_65536Bytes-16     	   52368	     22670 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_1048576Bytes-16   	    3211	    371099 ns/op	      31 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7764166	       212.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7356000	       213.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7275108	       222.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7252506	       218.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7222454	       215.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7152732	       218.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7141456	       218.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1531964	       854.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1552718	       858.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1506757	       858.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1512374	       857.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1511331	       820.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1612539	       849.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1528228	       850.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  500517	      2869 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  512047	      2874 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  509761	      2870 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  505234	      2870 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  502057	      2868 ns/op	      17 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  512407	      2872 ns/op	      17 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  511006	      2862 ns/op	      18 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   30416	     40926 ns/op	      19 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   29584	     41068 ns/op	      19 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   30120	     40628 ns/op	      21 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   30799	     40167 ns/op	      20 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   30522	     40132 ns/op	      26 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   30514	     40283 ns/op	      32 B/op	       1 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   30471	     40819 ns/op	      34 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_16Bytes-16                                	36082304	        32.17 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_32Bytes-16                                	32125682	        37.40 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_64Bytes-16                                	24775365	        48.45 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_128Bytes-16                               	17034656	        70.31 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_256Bytes-16                               	10265596	       117.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_512Bytes-16                               	 5841603	       207.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_1024Bytes-16                              	 3145809	       389.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1635278	       734.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_4096Bytes-16                              	  824144	      1439 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	19780939	        71.41 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	20381500	        67.59 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	20477874	        75.09 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	20288347	        75.80 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	20593477	        77.60 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	20457668	        75.27 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	20431212	        79.42 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	19519957	        65.37 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	19662943	        78.61 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	20034865	        81.71 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	19871482	        80.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	19597688	        80.79 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	19745366	        77.70 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	19542010	        83.02 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	19661346	        68.48 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	19675839	        68.37 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	19936645	        82.51 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	19774161	        82.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	19321026	        79.97 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	19652545	        81.72 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	19722231	        82.27 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	11073192	       150.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	11234876	       149.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	11159658	       152.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	11113580	       147.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	11135155	       137.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	11183436	       142.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	11162668	       152.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	11070301	       155.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	11087118	       141.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	11032106	       151.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	11073544	       149.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	11151897	       147.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	11212866	       149.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	11343543	       130.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	11545977	       148.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	11439361	       147.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	11266935	       142.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11478356	       144.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	11434555	       151.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	11345871	       150.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11290964	       148.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	12070354	       121.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	12104995	       146.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	12161227	       137.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	12072469	       141.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	12044836	       141.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	12096656	       139.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12044241	       142.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	11881642	       125.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	11984848	       128.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	11853602	       128.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	11718697	       130.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	11718024	       127.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	11453346	       129.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	11259843	       131.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 7103044	       226.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7005378	       207.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 6989464	       218.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7004456	       211.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 6948662	       219.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 6894457	       221.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 6860554	       218.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     300	   3863741 ns/op	     177 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16            	    3008	    377219 ns/op	     109 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16            	    3184	    383339 ns/op	     110 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16            	    3146	    380488 ns/op	     123 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16           	    3118	    382533 ns/op	     126 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16           	    3108	    385038 ns/op	     153 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16           	    3084	    390198 ns/op	     261 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16          	    3098	    383550 ns/op	     261 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      58	  19113881 ns/op	     850 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16            	     501	   2336128 ns/op	     391 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16            	     525	   2286823 ns/op	     438 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16            	     514	   2300128 ns/op	     478 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16           	     531	   2232886 ns/op	     531 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16           	     507	   2093149 ns/op	     643 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16           	     543	   2037698 ns/op	     859 B/op	       8 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16          	     592	   1958397 ns/op	    1037 B/op	      12 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      28	  37986455 ns/op	    1528 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16           	     236	   4963334 ns/op	     656 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16           	     212	   4827528 ns/op	     784 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16           	     214	   4748568 ns/op	     870 B/op	       6 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16          	     280	   4192354 ns/op	     859 B/op	       6 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16          	     264	   4003759 ns/op	    1062 B/op	       9 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16          	     292	   4063700 ns/op	    1177 B/op	      12 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16         	     298	   4044003 ns/op	    1445 B/op	      19 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/ctrdrbg	318.239s
```

</details>

### UUID Generation with Google UUID and ctrdrbg

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the ctrdrbg-based UUID generation:

| Benchmark Scenario                         | Default ns/op | CTRDRBG ns/op | % Faster (ns/op) | Default B/op | CTRDRBG B/op | Default allocs/op | CTRDRBG allocs/op |
|--------------------------------------------|---------------:|---------------:|------------------:|--------------:|--------------:|-------------------:|-------------------:|
| v4 Serial                                   |         178.4  |         45.90 |           74.3%   |          16   |          32   |                 1  |                 2  |
| v4 Parallel                                 |         462.2  |         12.49 |           97.3%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (2 goroutines)                |         413.4  |         26.07 |           93.7%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (4 goroutines)                |         427.0  |         16.45 |           96.1%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (8 goroutines)                |         480.9  |         12.87 |           97.3%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (16 goroutines)               |         455.6  |         10.39 |           97.7%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (32 goroutines)               |         514.6  |         10.38 |           98.0%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (64 goroutines)               |         526.1  |         10.39 |           98.0%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (128 goroutines)              |         513.7  |         10.59 |           97.9%   |          16   |          32   |                 1  |                 2  |
| v4 Concurrent (256 goroutines)              |         516.5  |         10.71 |           97.9%   |          16   |          32   |                 1  |                 2  |

<details>
  <summary>Expand to see results</summary>

  ```shell
make bench-ctrdrbg-uuid
go test -bench='^BenchmarkUUID_' -run=^$ -benchmem -memprofile=x/crypto/ctrdrbg/mem.out -cpuprofile=x/crypto/ctrdrbg/cpu.out ./x/crypto/ctrdrbg
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/ctrdrbg
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6459616	       178.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2645497	       462.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16         	 2896381	       413.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2863237	       427.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2484550	       480.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2611302	       455.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2309664	       514.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2321017	       526.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2334637	       513.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2308024	       516.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	25304576	        45.90 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	106977922	        12.49 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_2-16         	45092438	        26.07 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	74931704	        16.45 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	95511933	        12.87 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	100000000	        10.39 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	100000000	        10.38 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	100000000	        10.39 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	100000000	        10.59 ns/op	      32 B/op	       2 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	100000000	        10.71 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/ctrdrbg	29.324s
  ```
</details>

---

## FIPS‑140 Mode

See [FIPS‑140.md](../../../FIPS-140.md) for compliance, deployment, and configuration guidance.

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
