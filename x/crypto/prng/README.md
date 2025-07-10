# prng: Cryptographically Secure Pseudo-Random Number Generator (CSPRNG)

## Overview

The `prng` package provides a high-performance, cryptographically secure pseudo-random number generator (CSPRNG) 
that implements the `io.Reader` interface. Designed for concurrent use, it leverages the ChaCha20 cipher stream 
to efficiently generate random bytes.

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
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid/x/crypto/prng
cpu: Apple M4 Max
BenchmarkPRNG_ReadSerial/Serial_Read_8Bytes-16  	66476008	        18.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16Bytes-16 	50425110	        23.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_21Bytes-16 	41878482	        28.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_32Bytes-16 	34911459	        34.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_64Bytes-16 	21626848	        55.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_100Bytes-16         	14644287	        81.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_256Bytes-16         	 8462792	       141.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_512Bytes-16         	 4485656	       267.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_1000Bytes-16        	 2291336	       523.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_4096Bytes-16        	  589884	      1977 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16384Bytes-16       	  151071	      7846 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_10Goroutines-16         	445045048	         2.411 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_100Goroutines-16        	479350718	         2.299 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_1000Goroutines-16       	509940837	         2.318 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_10Goroutines-16         	421898340	         2.395 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_100Goroutines-16        	524359978	         3.256 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_1000Goroutines-16       	530904837	         2.876 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_10Goroutines-16         	396200329	         4.524 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_100Goroutines-16        	371367224	         3.592 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_1000Goroutines-16       	398330661	         3.366 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_10Goroutines-16         	223768764	         4.918 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_100Goroutines-16        	262933449	         4.894 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_1000Goroutines-16       	228344895	         4.794 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_10Goroutines-16        	159867238	         7.320 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_100Goroutines-16       	165507139	         9.083 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_1000Goroutines-16      	168286864	         7.132 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_10Goroutines-16        	98790507	        12.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_100Goroutines-16       	89193912	        12.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_1000Goroutines-16      	92149676	        12.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_10Goroutines-16        	53631382	        21.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_100Goroutines-16       	52055937	        21.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_1000Goroutines-16      	53502958	        21.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_10Goroutines-16       	27251127	        43.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_100Goroutines-16      	26952336	        42.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_1000Goroutines-16     	27398458	        43.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_10Goroutines-16       	 7457576	       160.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_100Goroutines-16      	 7394361	       159.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_1000Goroutines-16     	 7453174	       161.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_10Goroutines-16      	 1887352	       636.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_100Goroutines-16     	 1862580	       642.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_1000Goroutines-16    	 1889932	       635.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16        	  518890	      2205 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16       	  222415	      5370 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16       	  136910	      8739 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16       	   34401	     34989 ns/op	       2 B/op	       0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16     	    2139	    557500 ns/op	     491 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_10Goroutines-16         	 7449819	       160.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_100Goroutines-16        	 7255634	       161.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_1000Goroutines-16       	 7248561	       160.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_10Goroutines-16        	 3055440	       392.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_100Goroutines-16       	 3057937	       391.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_1000Goroutines-16      	 3070320	       389.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_10Goroutines-16        	 1890914	       634.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_100Goroutines-16       	 1890457	       634.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_1000Goroutines-16      	 1887400	       633.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_10Goroutines-16        	  437454	      2525 ns/op	       5 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_100Goroutines-16       	  432376	      2530 ns/op	       5 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_1000Goroutines-16      	  440277	      2528 ns/op	       4 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_10Goroutines-16      	   29372	     40590 ns/op	    1144 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_100Goroutines-16     	   29516	     40706 ns/op	    1139 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_1000Goroutines-16    	   29493	     40617 ns/op	    1140 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                 	57730353	        18.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                                	47622432	        24.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                                	40288116	        29.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                                	37421308	        31.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                                	32296952	        37.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                                	23610117	        50.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                                	19637991	        59.96 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                               	11935350	        99.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                               	 7618376	       157.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                               	 4081744	       293.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                              	 2121154	       565.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1000000	      1108 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                              	  554908	      2207 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_10Goroutines-16      	724226222	         1.921 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_100Goroutines-16     	707387110	         1.953 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_1000Goroutines-16    	770943771	         2.067 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_10Goroutines-16     	585415827	         2.215 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_100Goroutines-16    	521241872	         2.175 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_1000Goroutines-16   	596724973	         3.298 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_10Goroutines-16     	490648880	         2.764 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_100Goroutines-16    	412268361	         2.718 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_1000Goroutines-16   	499436485	         2.830 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_10Goroutines-16     	432565579	         2.923 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_100Goroutines-16    	435805737	         2.844 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_1000Goroutines-16   	451968037	         3.065 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_10Goroutines-16     	379548853	         3.769 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_100Goroutines-16    	383046519	         3.563 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_1000Goroutines-16   	368127094	         3.396 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_10Goroutines-16     	279099705	         4.296 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_100Goroutines-16    	266774709	         4.313 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_1000Goroutines-16   	285190446	         6.579 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_10Goroutines-16     	249749946	         5.204 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_100Goroutines-16    	247784841	         5.265 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_1000Goroutines-16   	230481091	         5.134 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_10Goroutines-16    	148439942	         8.048 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_100Goroutines-16   	143405518	         9.013 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_1000Goroutines-16  	148282646	         7.868 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_10Goroutines-16    	90628738	        11.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_100Goroutines-16   	95038521	        11.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_1000Goroutines-16  	100000000	        11.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_10Goroutines-16    	51148168	        22.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_100Goroutines-16   	50505404	        21.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_1000Goroutines-16  	52335469	        21.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_10Goroutines-16   	27907598	        41.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_100Goroutines-16  	28232719	        41.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_1000Goroutines-16 	27732253	        41.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_10Goroutines-16   	14434743	        82.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_100Goroutines-16  	14406388	        81.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_1000Goroutines-16 	14531326	        81.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_10Goroutines-16   	 7382282	       161.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_100Goroutines-16  	 7379265	       160.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_1000Goroutines-16 	 7390212	       161.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     213	   5596168 ns/op	   49257 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_10Goroutines-16           	    2910	    416219 ns/op	  115338 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_100Goroutines-16          	    2756	    422856 ns/op	  121784 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      39	  27941823 ns/op	 1344682 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_10Goroutines-16           	     495	   2240291 ns/op	 3389523 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_100Goroutines-16          	     531	   2185563 ns/op	 3159690 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      20	  55834806 ns/op	 5243571 B/op	       0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_10Goroutines-16          	     232	   4781582 ns/op	14463506 B/op	       1 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_100Goroutines-16         	     249	   4448588 ns/op	13476002 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/nanoid/x/crypto/prng	182.010s
```
</details>

### UUID Generation

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the CSPRNG-based UUID generation:

| Benchmark Scenario                       | Default ns/op | CSPRNG ns/op | % Faster (ns/op) | Default B/op | CSPRNG B/op | Default allocs/op | CSPRNG allocs/op |
|------------------------------------------|--------------:|-------------:|-----------------:|-------------:|------------:|------------------:|-----------------:|
| v4 Serial                               |      184.4    |      36.00   |    80.5%         |      16      |     16      |      1            |      1           |
| v4 Parallel                             |      455.8    |       6.68   |    98.5%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (1 goroutine)             |      185.4    |      36.90   |    80.1%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (2 goroutines)            |      371.5    |      19.83   |    94.7%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (4 goroutines)            |      461.5    |      11.72   |    97.5%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (8 goroutines)            |      481.1    |       8.54   |    98.2%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (16 goroutines)           |      453.8    |       6.62   |    98.5%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (32 goroutines)           |      505.1    |       6.50   |    98.7%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (64 goroutines)           |      510.2    |       6.46   |    98.7%         |      16      |     16      |      1            |      1           |
| v4 Concurrent (128 goroutines)          |      511.7    |       6.41   |    98.7%         |      16      |     16      |      1            |      1           |

<details>
  <summary>Expand to see results</summary>
```shell
go test -bench='^BenchmarkUUID_' -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16               6512481               184.4 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16             2641772               455.8 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_1-16              6466557               185.4 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16              3225741               371.5 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16              2658018               461.5 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16              2490762               481.1 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16             2624019               453.8 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16             2373811               505.1 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16             2358780               510.2 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16            2350291               511.7 ns/op            16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Serial-16                               31947392                36.00 ns/op           16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Parallel-16                             177618490                6.675 ns/op          16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_1-16              30928698                36.90 ns/op           16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_2-16              60665552                19.83 ns/op           16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_4-16              91901493                11.72 ns/op           16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_8-16              140070646                8.545 ns/op          16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_16-16             182767742                6.625 ns/op          16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_32-16             184506772                6.502 ns/op          16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_64-16             185731188                6.458 ns/op          16 B/op          1 allocs/op
BenchmarkUUID_v4_CSPRNG_Concurrent/Goroutines_128-16            186703047                6.412 ns/op          16 B/op          1 allocs/op
PASS
ok      github.com/sixafter/nanoid      32.519s
```
</details>

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
