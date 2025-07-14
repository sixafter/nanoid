# prng: Cryptographically Secure Pseudo-Random Number Generator (CSPRNG)

## Overview

The `prng` package provides a high-performance, cryptographically secure pseudo-random number generator (CSPRNG) 
that implements the `io.Reader` interface. Designed for concurrent use, it leverages the ChaCha20 cipher stream 
to efficiently generate random bytes.

Technically, this PRNG is not pseudo-random but is cryptographically random.

The package includes a global `Reader` and a `sync.Pool` to manage PRNG instances, ensuring low contention and 
optimized performance.

## Features

* **Cryptographic Security:** Utilizes the [ChaCha20](https://pkg.go.dev/golang.org/x/crypto/chacha20) cipher for secure random number generation. 
* **Concurrent Support:** Includes a thread-safe global `Reader` for concurrent access. 
    * Up to 98% faster when using the `prng.Reader` as a source for v4 UUID generation using Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package as compared to using the default rand reader.
    * See the benchmark results [here](#uuid-generation).
* **Efficient Resource Management:** Uses a `sync.Pool` to manage PRNG instances, reducing the overhead on `crypto/rand.Reader`. 
* **Extensible API:** Allows users to create and manage custom PRNG instances via `NewReader`.

---

## Installation

To install the package, run the following command:

```bash
go get -u github.com/sixafter/nanoid/x/crypto/prng
```

## Usage

Global Reader:

```go
package main

import (
  "fmt"
  
  "github.com/sixafter/nanoid/x/crypto/prng"
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
  "github.com/sixafter/nanoid/x/crypto/prng"
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

### NanoID Generation

Performance Benchmarks for concurrent reads of standard size Nano ID generation of 21 bytes:

* Throughput: ~3.16 `ns/op`
* Memory Usage: 0 `B/op`
* Allocations: 0 `allocs/op`

These benchmarks demonstrate the package's focus on minimizing latency, memory usage, and allocation overhead, making it suitable for high-performance applications.

<details>
  <summary>Expand to see results</summary>

```shell
 make bench-csprng
go test -bench='^BenchmarkPRNG_' -benchmem -memprofile=mem.out -cpuprofile=cpu.out ./x/crypto/prng
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/prng
cpu: Apple M4 Max
BenchmarkPRNG_ReadSerial/Serial_Read_8Bytes-16  	67554878	        17.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16Bytes-16 	51417919	        22.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_21Bytes-16 	44218303	        26.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_32Bytes-16 	34992861	        33.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_64Bytes-16 	21992018	        54.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_100Bytes-16         	14944401	        81.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_256Bytes-16         	 8358837	       142.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_512Bytes-16         	 4527169	       264.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_1000Bytes-16        	 2275684	       527.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_4096Bytes-16        	  591889	      1981 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16384Bytes-16       	  151741	      7852 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_1Goroutines-16         	635248483	         1.892 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_2Goroutines-16         	633992288	         1.963 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_4Goroutines-16         	626458201	         1.898 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_8Goroutines-16         	644301068	         1.969 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_16Goroutines-16        	644603472	         2.003 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_32Goroutines-16        	596536359	         1.911 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_64Goroutines-16        	635140470	         1.863 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_128Goroutines-16       	631339147	         1.936 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_1Goroutines-16         	537349323	         2.328 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_2Goroutines-16         	524518315	         2.258 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_4Goroutines-16         	529004313	         2.275 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_8Goroutines-16         	537154891	         2.315 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_16Goroutines-16        	530123434	         2.538 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_32Goroutines-16        	531488772	         2.314 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_64Goroutines-16        	503025224	         2.281 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_128Goroutines-16       	503538580	         2.335 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_1Goroutines-16         	407062821	         3.284 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_2Goroutines-16         	397926030	         2.979 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_4Goroutines-16         	383092478	         2.980 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_8Goroutines-16         	397549164	         3.320 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_16Goroutines-16        	404857046	         3.092 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_32Goroutines-16        	403695211	         3.157 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_64Goroutines-16        	395528656	         3.149 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_128Goroutines-16       	402283068	         3.040 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_1Goroutines-16         	252656422	         4.806 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_2Goroutines-16         	251129143	         4.756 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_4Goroutines-16         	251192641	         4.769 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_8Goroutines-16         	251188281	         4.704 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_16Goroutines-16        	251816532	         5.043 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_32Goroutines-16        	254376242	         4.795 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_64Goroutines-16        	253721534	         4.882 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_128Goroutines-16       	251037358	         4.878 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_1Goroutines-16        	168996282	         7.126 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_2Goroutines-16        	170772955	         7.162 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_4Goroutines-16        	172963831	         7.191 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_8Goroutines-16        	172869572	         7.159 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_16Goroutines-16       	170047622	         7.062 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_32Goroutines-16       	167939889	         7.045 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_64Goroutines-16       	166457977	         7.051 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_128Goroutines-16      	169524238	         7.026 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_1Goroutines-16        	96553249	        11.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_2Goroutines-16        	98563971	        11.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_4Goroutines-16        	100000000	        11.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_8Goroutines-16        	95472039	        11.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_16Goroutines-16       	100000000	        11.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_32Goroutines-16       	99647769	        11.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_64Goroutines-16       	97357482	        11.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_128Goroutines-16      	99656392	        11.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_1Goroutines-16        	53198954	        21.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_2Goroutines-16        	52842680	        21.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_4Goroutines-16        	51355572	        21.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_8Goroutines-16        	53375143	        21.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_16Goroutines-16       	51024384	        21.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_32Goroutines-16       	53242525	        22.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_64Goroutines-16       	52573258	        21.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_128Goroutines-16      	52067043	        21.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_1Goroutines-16       	27332647	        43.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_2Goroutines-16       	27003914	        43.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_4Goroutines-16       	26665702	        43.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_8Goroutines-16       	27207207	        43.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_16Goroutines-16      	26959173	        43.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_32Goroutines-16      	27306549	        43.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_64Goroutines-16      	26416467	        44.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_128Goroutines-16     	26944242	        44.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_1Goroutines-16       	 7338415	       161.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7469488	       161.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7270279	       161.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7448995	       161.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7437883	       160.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7489624	       161.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7355366	       162.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7447375	       161.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_1Goroutines-16      	 1886826	       647.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1854333	       639.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1885856	       635.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1885242	       642.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1845896	       641.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1867432	       647.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1855411	       645.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1863171	       648.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16       	  523388	      2200 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16      	  222051	      5366 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16      	  136734	      8735 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16      	   34368	     34876 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16    	    2128	    557867 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_1Goroutines-16         	 7466652	       159.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7502001	       159.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7483605	       159.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7483514	       159.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7491762	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7506014	       159.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7416632	       160.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7472210	       160.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_1Goroutines-16        	 3051492	       391.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_2Goroutines-16        	 3045686	       392.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_4Goroutines-16        	 3069856	       392.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_8Goroutines-16        	 3065576	       392.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_16Goroutines-16       	 3064304	       393.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_32Goroutines-16       	 3067126	       389.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_64Goroutines-16       	 3066285	       392.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_128Goroutines-16      	 3069032	       392.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_1Goroutines-16        	 1873540	       635.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1895378	       636.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1893477	       637.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1887652	       634.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1891016	       632.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1893213	       634.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1897372	       636.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1886792	       637.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_1Goroutines-16        	  428331	      2557 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  425432	      2566 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  427942	      2562 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  448995	      2549 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  436942	      2550 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  437209	      2548 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  442488	      2554 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  435504	      2555 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_1Goroutines-16      	   29175	     41789 ns/op	     578 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   28531	     41581 ns/op	     591 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   28742	     41162 ns/op	     587 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   28744	     41561 ns/op	     587 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   29016	     41360 ns/op	     581 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   28885	     41391 ns/op	     584 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   29019	     41337 ns/op	     581 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   28870	     41424 ns/op	     584 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                	63269031	        17.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                               	47546798	        24.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                               	40893697	        29.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                               	37978738	        31.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                               	31741378	        37.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                               	23770142	        50.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                               	18617482	        64.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                              	11167568	       107.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                              	 7656061	       157.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                              	 4072642	       295.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                             	 2088606	       570.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                             	 1000000	      1115 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                             	  544248	      2202 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_1Goroutines-16      	791767376	         1.416 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_2Goroutines-16      	823381989	         1.531 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_4Goroutines-16      	707644700	         1.696 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_8Goroutines-16      	679908354	         1.916 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_16Goroutines-16     	677823134	         1.849 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_32Goroutines-16     	649340985	         1.953 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_64Goroutines-16     	650302965	         1.834 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_128Goroutines-16    	625400376	         1.788 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_1Goroutines-16     	523881442	         2.311 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	533899465	         2.344 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	527728852	         2.322 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	518534082	         2.352 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	519025528	         2.446 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	510509312	         2.379 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	520443865	         2.322 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	507270969	         2.449 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_1Goroutines-16     	423455494	         2.856 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_2Goroutines-16     	437486252	         2.931 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_4Goroutines-16     	430021920	         2.733 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_8Goroutines-16     	442240272	         2.701 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_16Goroutines-16    	444357355	         2.734 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_32Goroutines-16    	434068728	         2.894 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_64Goroutines-16    	444295042	         2.909 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_128Goroutines-16   	428037400	         2.769 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_1Goroutines-16     	413261503	         2.918 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_2Goroutines-16     	417008079	         2.880 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_4Goroutines-16     	409639732	         3.131 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_8Goroutines-16     	415724850	         2.937 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_16Goroutines-16    	414936682	         2.860 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_32Goroutines-16    	418188345	         2.908 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_64Goroutines-16    	418807053	         3.139 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_128Goroutines-16   	398414145	         3.101 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_1Goroutines-16     	344253705	         3.434 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	356466216	         3.287 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	364713326	         3.560 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	369372076	         3.435 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	354431001	         3.492 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	348002482	         3.467 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	352474406	         3.687 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	345074368	         3.378 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_1Goroutines-16     	266253234	         4.536 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_2Goroutines-16     	271406299	         4.495 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_4Goroutines-16     	267632970	         4.482 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_8Goroutines-16     	275894203	         4.932 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_16Goroutines-16    	267850242	         4.559 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_32Goroutines-16    	262589598	         4.397 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_64Goroutines-16    	265367104	         4.816 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_128Goroutines-16   	269811342	         4.565 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_1Goroutines-16     	229454901	         5.202 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	232215346	         5.527 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	224167776	         5.267 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	232395082	         5.452 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	229815367	         5.240 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	224010520	         5.225 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	227806197	         5.460 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	230051076	         5.244 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_1Goroutines-16    	141047791	         8.865 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	140729358	         8.435 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	145064774	         8.537 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	144049452	         8.410 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	141839947	         8.720 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	139133290	         8.667 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	140278188	         8.612 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	144401906	         8.701 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_1Goroutines-16    	89568949	        12.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	90273645	        12.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	98032207	        11.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	91239426	        11.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	94771450	        12.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	90742674	        12.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	94968324	        12.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	91495381	        12.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_1Goroutines-16    	50753011	        22.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	51000891	        22.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	51776397	        22.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	51171524	        22.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	51042018	        22.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	50918116	        22.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	51752114	        22.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	50276606	        22.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_1Goroutines-16   	26980312	        43.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	27002016	        43.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	26624190	        43.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	26940007	        43.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	26724073	        43.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	27028066	        42.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	27256516	        43.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	27452074	        43.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_1Goroutines-16   	14054799	        83.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	14181148	        83.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	14237816	        83.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	14115328	        84.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	14127514	        84.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	13833184	        85.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	13917110	        85.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	13875759	        84.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_1Goroutines-16   	 7154342	       167.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 7127234	       168.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7144708	       166.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7168696	       167.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7130156	       167.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7101517	       167.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7135390	       168.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7101944	       170.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                           	     212	   5662999 ns/op	      52 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_1Goroutines-16           	    2809	    412093 ns/op	   59771 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16           	    2800	    430921 ns/op	   59955 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16           	    2803	    440070 ns/op	   59887 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16           	    2646	    445660 ns/op	   63445 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16          	    2648	    446439 ns/op	   63406 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16          	    2662	    446943 ns/op	   63065 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16          	    2685	    442157 ns/op	   62523 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16         	    2659	    438954 ns/op	   63135 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                           	      37	  28243375 ns/op	      76 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_1Goroutines-16           	     483	   2242872 ns/op	 1736968 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16           	     512	   2112440 ns/op	 1638597 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16           	     442	   2328136 ns/op	 1898103 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16           	     535	   2519526 ns/op	 1568139 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16          	     517	   2334806 ns/op	 1622729 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16          	     500	   2353781 ns/op	 1677917 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16          	     487	   2372848 ns/op	 1722710 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16         	     544	   2337972 ns/op	 1542204 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                          	      19	  56357189 ns/op	     579 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_1Goroutines-16          	     226	   4774731 ns/op	 7424002 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16          	     277	   4412455 ns/op	 6057115 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16          	     236	   5250374 ns/op	 7109408 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16          	     264	   4573370 ns/op	 6355380 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16         	     205	   5033133 ns/op	 8184571 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16         	     261	   4641640 ns/op	 6428437 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16         	     266	   4706944 ns/op	 6307581 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16        	     262	   5043935 ns/op	 6403893 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/prng	416.792s
```
</details>

### UUID Generation

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the CSPRNG-based UUID generation:

| Benchmark Scenario                   | Default ns/op | CSPRNG ns/op | % Faster (ns/op) | Default B/op | CSPRNG B/op | Default allocs/op | CSPRNG allocs/op |
|--------------------------------------|--------------:|-------------:|-----------------:|-------------:|------------:|------------------:|-----------------:|
| v4 Serial                           |      181.7    |     35.84    |      80.3%       |      16      |     16      |      1            |      1           |
| v4 Parallel                         |      447.4    |      6.011   |      98.7%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (4 goroutines)        |      465.2    |     11.47    |      97.5%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (8 goroutines)        |      485.3    |      7.770   |      98.4%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (16 goroutines)       |      448.1    |      6.047   |      98.7%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (32 goroutines)       |      515.0    |      5.882   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (64 goroutines)       |      510.7    |      5.820   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (128 goroutines)      |      512.6    |      5.722   |      98.9%       |      16      |     16      |      1            |      1           |
| v4 Concurrent (256 goroutines)      |      518.9    |      5.882   |      98.9%       |      16      |     16      |      1            |      1           |

<details>
  <summary>Expand to see results</summary>
```shell
make bench-uuid
go test -bench='^BenchmarkUUID_' -benchmem -memprofile=mem.out -cpuprofile=cpu.out ./x/crypto/prng
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/prng
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6318418	       181.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2540641	       464.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2574469	       461.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2462932	       484.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2557530	       466.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2306965	       527.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2329918	       515.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2323522	       514.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2327964	       515.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Serial-16                           	31848084	        36.45 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Parallel-16                         	202154900	         6.831 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_4-16          	100000000	        11.45 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_8-16          	154368582	         7.760 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_16-16         	203061876	         5.855 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_32-16         	200415012	         5.962 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_64-16         	198328420	         5.622 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_128-16        	202116156	         5.800 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_256-16        	203917351	         5.934 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/prng	30.172s
```
</details>

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
