# ctrdrbg: AES-CTR-DRBG for Deterministic Cryptographically Secure Random Number Generation

## Overview

The `ctrdrbg` package implements an [AES-CTR-DRBG (Deterministic Random Bit Generator)](https://csrc.nist.gov/publications/detail/sp/800-90a/rev-1/final) as specified in NIST SP 800-90A.  
It is designed for environments requiring deterministic, reproducible, and **FIPS‑140-compatible** random bit generation.  
This package is suitable for any application that needs strong cryptographic assurance or must comply with regulated environments (e.g., FedRAMP, FIPS, HIPAA).

The package uses only Go standard library crypto primitives (`crypto/aes` and `crypto/cipher`), making it safe for use in FIPS 140-validated Go runtimes.  
No third-party, homegrown, or experimental ciphers are used.

## Features

- **Deterministic Random Bit Generation:** Implements AES-CTR-DRBG as specified in NIST SP 800-90A, Revision 1.
- **FIPS‑140 Mode Compatible:** Designed to run in FIPS‑140 validated environments using only Go standard library crypto.
    - For details and deployment guidance, see [FIPS‑140.md](../../../FIPS-140.md).
- **Stateless and Concurrent:** Safe for concurrent use and supports stateless operation with independent DRBG instances.
- **Customizable Entropy:** Accepts user-provided entropy and personalization strings for deterministic output.
- **io.Reader Interface:** Satisfies the `io.Reader` interface for drop-in compatibility with Go packages expecting secure random bytes.
- **No External Dependencies:** Lightweight implementation; depends only on the Go standard library.

---

## Installation

```bash
go get -u github.com/sixafter/nanoid/x/crypto/ctrdrbg
```

---

## Usage

### Basic Usage: Generate Secure Random Bytes

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/nanoid/x/crypto/ctrdrbg"
)

func main() {
	// Example: AES-256 (32 bytes) key
	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(32))
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

### Basic Usage: With Personalization

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
		ctrdrbg.WithKeySize(32), // AES-256
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

### Safe for FIPS 140 Mode

- The DRBG is implemented strictly using Go's standard library crypto.
- Fully compatible with Go’s [`GODEBUG=fips140=on`](https://go.dev/doc/security/fips140) runtime enforcement.
- No external or non-validated algorithms are used at any time.

---

## Architecture

- **DRBG Instance:** Each DRBG is independently seeded and maintains its own internal state (Key, V).
- **AES-CTR Mode:** All random data is generated using AES in CTR mode, following NIST requirements.
- **Personalization:** Personalization string is incorporated at initialization for added uniqueness.
- **Concurrency:** The package is safe for concurrent use across goroutines.

---

## Performance Benchmarks

### Raw Random Byte Generation

Performance Benchmarks for various read sizes using the `ctrdrbg.Reader`.

* Throughput: ~69.28 `ns/op`
* Memory Usage: 0 `B/op`
* Allocations: 0 `allocs/op`

These benchmarks demonstrate the package's focus on minimizing latency, memory usage, and allocation overhead, making it suitable for high-performance applications.


<details>
  <summary>Expand to see results</summary>

```shell
make bench-ctrdrbg
go test -bench='^BenchmarkDRBG_' -benchmem -memprofile=x/crypto/ctrdrbg/mem.out -cpuprofile=x/crypto/ctrdrbg/cpu.out ./x/crypto/ctrdrbg
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/ctrdrbg
cpu: Apple M4 Max
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G2-16 	1000000000	         0.5867 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G4-16 	1000000000	         0.5950 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G8-16 	1000000000	         0.5726 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G16-16         	1000000000	         0.5398 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G32-16         	1000000000	         0.5341 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G64-16         	1000000000	         0.5224 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Concurrent_SyncPool_Baseline/G128-16        	1000000000	         0.5194 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_8Bytes-16            	41356254	        27.57 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_16Bytes-16           	66431233	        17.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_21Bytes-16           	38076753	        31.59 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_32Bytes-16           	51806107	        23.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_64Bytes-16           	37704330	        31.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_100Bytes-16          	21041049	        57.55 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_256Bytes-16          	12830982	        93.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_512Bytes-16          	 7029285	       171.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_1000Bytes-16         	 3602702	       332.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_4096Bytes-16         	  955296	      1261 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSerial/Serial_Read_16384Bytes-16        	  240126	      4990 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_2Goroutines-16         	24605625	        70.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_4Goroutines-16         	23419946	        73.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_8Goroutines-16         	23912062	        74.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_16Goroutines-16        	25121748	        67.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_32Goroutines-16        	26564932	        67.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_64Goroutines-16        	27518595	        64.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16Bytes_128Goroutines-16       	28097067	        65.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_2Goroutines-16         	16891406	       190.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_4Goroutines-16         	16116306	        98.83 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_8Goroutines-16         	15318064	       186.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_16Goroutines-16        	15804871	        96.53 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_32Goroutines-16        	16165725	        93.25 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_64Goroutines-16        	15918681	        87.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_21Bytes_128Goroutines-16       	16833019	       186.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_2Goroutines-16         	23454793	        81.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_4Goroutines-16         	22897376	        76.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_8Goroutines-16         	22688366	        69.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_16Goroutines-16        	24935323	        69.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_32Goroutines-16        	26483926	        67.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_64Goroutines-16        	27455266	        64.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_32Bytes_128Goroutines-16       	27119690	        62.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_2Goroutines-16         	24394726	        79.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_4Goroutines-16         	26096364	        77.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_8Goroutines-16         	25666962	        73.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_16Goroutines-16        	23711256	        70.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_32Goroutines-16        	27163968	        66.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_64Goroutines-16        	26312206	        66.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_64Bytes_128Goroutines-16       	25379367	        66.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_2Goroutines-16        	19304656	        81.00 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_4Goroutines-16        	19495551	        80.82 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_8Goroutines-16        	20357168	        79.74 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_16Goroutines-16       	19713253	        79.61 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_32Goroutines-16       	20250174	        74.99 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_64Goroutines-16       	19995680	        78.52 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_100Bytes_128Goroutines-16      	20529811	        80.44 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_2Goroutines-16        	11332866	       146.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_4Goroutines-16        	11900329	       138.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_8Goroutines-16        	11857389	       138.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_16Goroutines-16       	11628859	       136.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_32Goroutines-16       	12352996	       135.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_64Goroutines-16       	12126044	       137.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_256Bytes_128Goroutines-16      	12430214	       134.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_2Goroutines-16        	10687228	       154.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_4Goroutines-16        	10312039	       151.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_8Goroutines-16        	11452233	       147.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_16Goroutines-16       	11786670	       130.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_32Goroutines-16       	12293763	       153.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_64Goroutines-16       	11624494	       138.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_512Bytes_128Goroutines-16      	12017239	       156.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_2Goroutines-16       	11941149	       135.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_4Goroutines-16       	12125616	       130.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_8Goroutines-16       	11931241	       152.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_16Goroutines-16      	11713725	       141.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_32Goroutines-16      	12344746	       132.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_64Goroutines-16      	11962581	       143.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_1000Bytes_128Goroutines-16     	12175143	       142.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 5850606	       216.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7575766	       196.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7576767	       195.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7287967	       207.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7672298	       216.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7684561	       210.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7592994	       208.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1492654	       757.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1551319	       748.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1545484	       707.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1637460	       794.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1565677	       826.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1522772	       745.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1547283	       829.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16       	  794491	      1294 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16      	  389359	      3069 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16      	  235744	      4999 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16      	   59410	     20097 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16    	    3661	    324805 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 6672326	       189.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7079473	       213.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7138963	       182.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7135531	       202.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7374717	       220.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 6309536	       214.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 6509306	       204.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_2Goroutines-16        	 2249860	       557.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_4Goroutines-16        	 2256259	       511.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_8Goroutines-16        	 2539240	       474.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_16Goroutines-16       	 2873290	       498.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_32Goroutines-16       	 2558970	       495.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_64Goroutines-16       	 2613232	       506.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_128Goroutines-16      	 2499842	       470.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1492027	       868.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1414190	       887.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1660220	       771.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1705810	       729.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1535983	       893.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1623381	       716.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1646668	       692.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  459470	      2692 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  477390	      2722 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  509078	      2520 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  483771	      2710 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  498073	      2879 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  479391	      2888 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  502183	      2824 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   29029	     43505 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   29434	     42218 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   28771	     40847 ns/op	       2 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   28894	     40501 ns/op	       4 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   32516	     39658 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   30337	     40232 ns/op	      14 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   31057	     39979 ns/op	      15 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                	40272176	        29.22 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                               	65559736	        18.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                               	35809492	        33.69 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                               	35501340	        33.12 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                               	51358869	        23.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                               	42948087	        27.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                               	36682944	        32.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                              	22665475	        52.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                              	12645388	        97.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                              	 6916808	       173.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                             	 3620143	       332.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                             	 1845120	       648.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                             	  936307	      1291 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_2Goroutines-16      	20338738	        82.19 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_4Goroutines-16      	21369930	        78.94 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_8Goroutines-16      	20435473	        79.23 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_16Goroutines-16     	20460459	        78.37 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_32Goroutines-16     	20304883	        78.62 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_64Goroutines-16     	21448855	        74.44 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_128Goroutines-16    	21768312	        78.61 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	25935174	        75.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	23076331	        74.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	23677356	        70.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	24741459	        68.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	27184737	        62.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	28358011	        63.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	27622189	        63.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_2Goroutines-16     	16098954	       100.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_4Goroutines-16     	16030456	        98.94 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_8Goroutines-16     	12072583	       185.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_16Goroutines-16    	16511914	        96.88 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_32Goroutines-16    	16747729	       183.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_64Goroutines-16    	16346972	        95.14 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_128Goroutines-16   	16267012	        98.29 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_2Goroutines-16     	15657304	       190.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_4Goroutines-16     	15460786	        99.60 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_8Goroutines-16     	15668736	       187.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_16Goroutines-16    	15204483	        97.79 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_32Goroutines-16    	15880639	        95.20 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_64Goroutines-16    	16000852	       170.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_128Goroutines-16   	16305999	        91.99 ns/op	      16 B/op	       1 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	22960354	        78.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	23220020	        76.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	22819490	        75.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	26700981	        65.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	26463655	        64.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	27179196	        66.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	27110244	        62.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_2Goroutines-16     	21672531	        81.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_4Goroutines-16     	22630531	        73.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_8Goroutines-16     	24540880	        72.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_16Goroutines-16    	22656826	        71.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_32Goroutines-16    	26505790	        68.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_64Goroutines-16    	26689104	        65.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_128Goroutines-16   	27369868	        63.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	22562204	        79.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	24310504	        73.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	23567419	        71.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	25485416	        68.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	26559909	        68.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	26792042	        62.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	27030907	        65.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	11948868	       151.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	11782282	       148.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	11398888	       148.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	11631193	       146.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	12177382	       145.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	12876951	       142.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	12929379	       141.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	10390746	       144.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	11094285	       140.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	12036021	       137.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	12056655	       136.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	12299145	       118.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	12289383	       134.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	12372164	       134.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	11239603	       152.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	10941400	       150.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	11149155	       132.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11680039	       157.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	11701561	       139.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	11836663	       140.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11662245	       157.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	11722641	       149.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	11204637	       146.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	11648162	       151.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	10546617	       151.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	10446212	       155.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	11780808	       153.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12127719	       152.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	11497666	       139.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	11602941	       107.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	12177105	       128.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	11557608	       122.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	12385678	       120.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	11743467	       124.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	11852905	       115.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 6338384	       189.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 5815082	       178.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 6920331	       179.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 6441996	       196.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7082898	       181.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7199071	       223.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7184109	       214.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                           	     324	   3642148 ns/op	      14 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16           	    3103	    375278 ns/op	      14 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16           	    3159	    373398 ns/op	      16 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16           	    2980	    366157 ns/op	      34 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16          	    2726	    373430 ns/op	      43 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16          	    3158	    372897 ns/op	      55 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16          	    2869	    380115 ns/op	     147 B/op	       1 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16         	    2949	    368705 ns/op	     187 B/op	       1 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                           	      56	  18037516 ns/op	      73 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16           	     568	   1860397 ns/op	      64 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16           	     561	   1895792 ns/op	      95 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16           	     542	   1885216 ns/op	     186 B/op	       1 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16          	     531	   1899703 ns/op	     206 B/op	       1 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16          	     571	   1856677 ns/op	     376 B/op	       3 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16          	     588	   1880466 ns/op	     578 B/op	       5 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16         	     576	   1884183 ns/op	     730 B/op	       9 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                          	      31	  35947733 ns/op	     133 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16          	     326	   3780643 ns/op	     116 B/op	       0 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16          	     292	   3726583 ns/op	     182 B/op	       1 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16          	     286	   3618397 ns/op	     339 B/op	       2 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16         	     304	   3571114 ns/op	     409 B/op	       3 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16         	     320	   3582016 ns/op	     595 B/op	       5 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16         	     312	   3610114 ns/op	     971 B/op	       9 allocs/op
BenchmarkDRBG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16        	     313	   3619726 ns/op	    1232 B/op	      16 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/ctrdrbg	424.104s
```

</details>

### UUID Generation with Google UUID and ctrdrbg

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the ctrdrbg-based UUID generation:

| Benchmark Scenario                         | Default ns/op | CTRDRBG ns/op | % Faster (ns/op) | Default B/op | CTRDRBG B/op | Default allocs/op | CTRDRBG allocs/op |
|--------------------------------------------|--------------:|--------------:|-----------------:|-------------:|-------------:|------------------:|------------------:|
| v4 Serial                                 |     176.0     |     31.84     |     81.9%        |      16      |     16       |      1            |      1            |
| v4 Parallel                               |     460.7     |      7.85     |     98.3%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (4 goroutines)              |     471.1     |     10.25     |     97.8%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (8 goroutines)              |     480.0     |      7.63     |     98.4%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (16 goroutines)             |     448.3     |      5.86     |     98.7%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (32 goroutines)             |     509.5     |      5.89     |     98.8%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (64 goroutines)             |     509.3     |      5.92     |     98.8%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (128 goroutines)            |     515.5     |      5.86     |     98.9%        |      16      |     16       |      1            |      1            |
| v4 Concurrent (256 goroutines)            |     513.8     |      5.96     |     98.8%        |      16      |     16       |      1            |      1            |

<details>
  <summary>Expand to see results</summary>

  ```shell
 make bench-ctrdrbg-uuid
go test -bench='^BenchmarkUUID_' -benchmem -memprofile=x/crypto/ctrdrbg/mem.out -cpuprofile=x/crypto/ctrdrbg/cpu.out ./x/crypto/ctrdrbg
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/ctrdrbg
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6100108	       176.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2648044	       460.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2529234	       471.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2481907	       480.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2672932	       448.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2332858	       509.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2362628	       509.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2354358	       515.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2295604	       513.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	36376542	        31.84 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	131925153	         7.848 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	100000000	        10.25 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	155344178	         7.628 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	203604535	         5.862 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	202456543	         5.892 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	200709396	         5.921 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	205246484	         5.864 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	200157261	         5.961 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/ctrdrbg	30.286s
  ```
</details>

---

## FIPS‑140 Mode

See [FIPS‑140.md](../../../FIPS-140.md) for compliance, deployment, and configuration guidance.

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
