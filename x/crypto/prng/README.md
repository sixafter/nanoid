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
go test -bench='^BenchmarkPRNG_' -benchmem -memprofile=mem.out -cpuprofile=cpu.out ./x/crypto/prng
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/prng
cpu: Apple M4 Max
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G1-16 	1000000000	         0.5263 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G2-16 	1000000000	         0.5816 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G4-16 	1000000000	         0.5669 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G8-16 	1000000000	         0.5588 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G16-16         	1000000000	         0.5588 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G32-16         	1000000000	         0.5115 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G64-16         	1000000000	         0.5162 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_Concurrent_SyncPool_Baseline/G128-16        	1000000000	         0.5149 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_8Bytes-16            	66802125	        17.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16Bytes-16           	49537222	        24.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_21Bytes-16           	42370824	        28.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_32Bytes-16           	33351436	        35.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_64Bytes-16           	20887303	        57.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_100Bytes-16          	14131015	        84.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_256Bytes-16          	 8184606	       147.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_512Bytes-16          	 4359367	       276.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_1000Bytes-16         	 2194410	       546.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_4096Bytes-16         	  583491	      2051 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16384Bytes-16        	  148093	      8136 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_1Goroutines-16         	650460120	         1.847 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_2Goroutines-16         	629682855	         1.887 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_4Goroutines-16         	655235089	         1.888 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_8Goroutines-16         	610236463	         1.857 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_16Goroutines-16        	645084555	         1.863 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_32Goroutines-16        	651985818	         1.871 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_64Goroutines-16        	654493987	         1.909 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_128Goroutines-16       	651597124	         1.957 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_1Goroutines-16         	539495166	         2.220 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_2Goroutines-16         	550199494	         2.204 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_4Goroutines-16         	544748642	         2.289 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_8Goroutines-16         	533738281	         2.263 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_16Goroutines-16        	539562180	         2.279 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_32Goroutines-16        	536735340	         2.292 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_64Goroutines-16        	521125108	         2.475 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_128Goroutines-16       	519770217	         2.314 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_1Goroutines-16         	420640477	         3.037 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_2Goroutines-16         	412556794	         2.970 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_4Goroutines-16         	416960445	         3.135 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_8Goroutines-16         	402117763	         2.925 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_16Goroutines-16        	402546553	         2.963 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_32Goroutines-16        	409780027	         2.929 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_64Goroutines-16        	414022539	         2.954 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_128Goroutines-16       	408222045	         2.996 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_1Goroutines-16         	254156876	         4.775 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_2Goroutines-16         	255359896	         4.698 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_4Goroutines-16         	253389480	         4.694 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_8Goroutines-16         	259646895	         4.783 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_16Goroutines-16        	256014086	         5.022 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_32Goroutines-16        	258362790	         4.713 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_64Goroutines-16        	257810678	         4.780 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_128Goroutines-16       	258595888	         4.749 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_1Goroutines-16        	176195185	         6.910 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_2Goroutines-16        	174068412	         6.901 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_4Goroutines-16        	176464705	         7.146 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_8Goroutines-16        	175301032	         6.975 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_16Goroutines-16       	173302622	         7.092 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_32Goroutines-16       	173771790	         7.153 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_64Goroutines-16       	175646030	         6.889 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_128Goroutines-16      	176137524	         7.230 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_1Goroutines-16        	100000000	        11.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_2Goroutines-16        	99491836	        11.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_4Goroutines-16        	100000000	        11.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_8Goroutines-16        	100000000	        11.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_16Goroutines-16       	99598492	        11.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_32Goroutines-16       	99088934	        11.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_64Goroutines-16       	100000000	        11.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_128Goroutines-16      	95314687	        11.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_1Goroutines-16        	54776073	        21.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_2Goroutines-16        	55136816	        21.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_4Goroutines-16        	54264060	        21.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_8Goroutines-16        	54337674	        21.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_16Goroutines-16       	54327219	        21.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_32Goroutines-16       	54248727	        21.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_64Goroutines-16       	54337674	        21.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_128Goroutines-16      	54390421	        21.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_1Goroutines-16       	27425724	        43.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_2Goroutines-16       	27443574	        43.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_4Goroutines-16       	27448674	        43.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_8Goroutines-16       	26781976	        42.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_16Goroutines-16      	27391527	        42.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_32Goroutines-16      	27442762	        42.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_64Goroutines-16      	27823237	        42.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_128Goroutines-16     	27857820	        42.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_1Goroutines-16       	 7507477	       159.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7449465	       161.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7480952	       160.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7518315	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7456206	       160.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7523054	       161.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7440681	       161.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7505826	       160.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_1Goroutines-16      	 1895702	       636.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1892054	       635.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1894464	       635.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1891182	       635.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1895898	       638.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1880988	       635.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1884595	       636.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1883289	       640.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16       	  540337	      2205 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16      	  222061	      5375 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16      	  137257	      8738 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16      	   34378	     34912 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16    	    2143	    558389 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_1Goroutines-16         	 7461139	       161.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7446165	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7489989	       160.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7467496	       159.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7474303	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7512214	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7508348	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7465855	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_1Goroutines-16        	 3085351	       392.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_2Goroutines-16        	 3065079	       392.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_4Goroutines-16        	 3080301	       391.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_8Goroutines-16        	 3062896	       392.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_16Goroutines-16       	 3066435	       389.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_32Goroutines-16       	 3076599	       392.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_64Goroutines-16       	 3077004	       391.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_128Goroutines-16      	 3070671	       390.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_1Goroutines-16        	 1891634	       639.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1896906	       636.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1894928	       636.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1894447	       641.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1897276	       637.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1886478	       636.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1886680	       637.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1891440	       638.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_1Goroutines-16        	  459158	      2530 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  460588	      2530 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  463558	      2531 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  458139	      2535 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  460410	      2540 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  458904	      2527 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  461967	      2533 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  461701	      2533 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_1Goroutines-16      	   29361	     40618 ns/op	     575 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   29277	     40764 ns/op	     576 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   28735	     41415 ns/op	     587 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   28964	     41151 ns/op	     583 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   29313	     41115 ns/op	     576 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   29301	     41104 ns/op	     575 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   29269	     41225 ns/op	     577 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   29089	     41407 ns/op	     579 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                	62868368	        18.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                               	49340076	        24.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                               	41391501	        28.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                               	38601810	        30.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                               	32784271	        36.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                               	24270248	        49.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                               	20143902	        59.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                              	12015544	        99.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                              	 7659739	       156.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                              	 4098468	       292.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                             	 2123703	       564.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                             	 1000000	      1108 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                             	  541821	      2217 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_1Goroutines-16      	879156032	         1.396 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_2Goroutines-16      	830060163	         1.397 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_4Goroutines-16      	810751011	         1.454 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_8Goroutines-16      	834182959	         1.451 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_16Goroutines-16     	833878632	         1.552 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_32Goroutines-16     	831751417	         1.454 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_64Goroutines-16     	826555861	         1.551 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_128Goroutines-16    	818561101	         1.536 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_1Goroutines-16     	616296154	         2.004 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	623762895	         1.971 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	621838849	         2.163 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	617387428	         2.169 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	597100453	         1.960 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	600222457	         1.997 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	604787016	         2.014 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	593647965	         2.072 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_1Goroutines-16     	516760257	         2.339 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_2Goroutines-16     	521778453	         2.360 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_4Goroutines-16     	502881790	         2.351 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_8Goroutines-16     	502118924	         2.315 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_16Goroutines-16    	518531374	         2.407 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_32Goroutines-16    	518055212	         2.409 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_64Goroutines-16    	517874584	         2.557 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_128Goroutines-16   	509005688	         2.367 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_1Goroutines-16     	470248149	         2.534 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_2Goroutines-16     	472335412	         2.741 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_4Goroutines-16     	456038014	         2.566 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_8Goroutines-16     	479753166	         2.550 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_16Goroutines-16    	480415159	         2.655 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_32Goroutines-16    	478641534	         2.551 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_64Goroutines-16    	467480557	         2.607 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_128Goroutines-16   	461649064	         2.517 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_1Goroutines-16     	403111800	         3.049 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	396507921	         2.987 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	407606041	         3.032 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	407094813	         2.998 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	399461836	         2.983 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	406319395	         2.958 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	407329959	         3.033 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	407482048	         3.210 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_1Goroutines-16     	293915040	         4.162 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_2Goroutines-16     	297891238	         4.182 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_4Goroutines-16     	295933624	         4.046 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_8Goroutines-16     	305446362	         3.996 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_16Goroutines-16    	295821338	         4.171 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_32Goroutines-16    	305582935	         4.097 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_64Goroutines-16    	299973564	         4.169 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_128Goroutines-16   	297425716	         4.111 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_1Goroutines-16     	249698302	         4.777 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	247477468	         4.933 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	252724820	         5.248 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	247442215	         4.816 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	252883330	         4.829 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	253638277	         5.233 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	241729533	         4.906 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	243802878	         5.088 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_1Goroutines-16    	145801368	         7.952 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	149229717	         7.926 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	149730523	         8.075 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	148953901	         8.180 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	151685190	         8.164 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	151020409	         7.986 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	153131425	         7.963 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	151474764	         7.876 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_1Goroutines-16    	95334564	        11.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	100000000	        11.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	96465286	        12.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	97019676	        11.96 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	95729727	        12.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	97000070	        12.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	96492753	        12.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	96555518	        11.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_1Goroutines-16    	52871683	        21.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	53121926	        22.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	53999628	        22.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	54194727	        22.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	52874984	        22.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	52933489	        22.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	53348545	        22.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	52576140	        22.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_1Goroutines-16   	28015428	        42.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	28118876	        42.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	28127114	        42.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	27915578	        42.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	28209379	        42.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	28064020	        42.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	28487040	        42.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	28457935	        42.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_1Goroutines-16   	14542399	        83.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	14469808	        83.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	14474594	        83.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	13892050	        83.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	14339308	        83.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	14495328	        83.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	14583733	        83.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	14219566	        83.96 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_1Goroutines-16   	 7349744	       165.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 7293321	       166.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7355757	       165.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7244097	       165.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7304551	       165.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7193750	       165.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7261617	       166.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7299075	       165.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                           	     213	   5589546 ns/op	      49 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_1Goroutines-16           	    2779	    421182 ns/op	   60417 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16           	    2808	    420584 ns/op	   59787 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16           	    2802	    427656 ns/op	   59917 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16           	    2815	    426156 ns/op	   59641 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16          	    2817	    428621 ns/op	   59597 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16          	    2768	    429181 ns/op	   60651 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16          	    2798	    428752 ns/op	   60000 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16         	    2756	    431677 ns/op	   60918 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                           	      38	  28091592 ns/op	      74 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_1Goroutines-16           	     466	   2156953 ns/op	 1800353 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16           	     488	   2352206 ns/op	 1719225 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16           	     512	   2254062 ns/op	 1638602 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16           	     481	   2375536 ns/op	 1744203 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16          	     535	   2132853 ns/op	 1568162 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16          	     547	   2220823 ns/op	 1533769 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16          	     439	   2387933 ns/op	 1911107 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16         	     472	   2383910 ns/op	 1777495 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                          	      20	  55822662 ns/op	     140 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_1Goroutines-16          	     271	   4320822 ns/op	 6191257 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16          	     265	   4802935 ns/op	 6331424 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16          	     241	   4952143 ns/op	 6961956 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16          	     246	   5297648 ns/op	 6820440 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16         	     241	   5064146 ns/op	 6961949 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16         	     255	   5068075 ns/op	 6579733 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16         	     254	   4586988 ns/op	 6605660 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16        	     226	   4678921 ns/op	 7424077 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/prng	423.279s
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
go test -bench='^BenchmarkUUID_' -benchmem -memprofile=x/crypto/prng/mem.out -cpuprofile=x/crypto/prng/cpu.out ./x/crypto/prng
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/prng
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
ok  	github.com/sixafter/nanoid/x/crypto/prng	31.142s
```
</details>

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
