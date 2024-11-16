# prng: Cryptographically Secure Pseudo-Random Number Generator (PRNG)

## Overview

The prng package provides a high-performance, cryptographically secure pseudo-random number generator (PRNG) 
that implements the io.Reader interface. Designed for concurrent use, it leverages the ChaCha20 cipher stream 
to efficiently generate random bytes.

The package includes a global Reader and a sync.Pool to manage PRNG instances, ensuring low contention and 
optimized performance.

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
BenchmarkPRNG_ReadSerial/Serial_Read_8Bytes-16          61095649                16.73 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16Bytes-16         50855985                23.37 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_21Bytes-16         42470296                27.94 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_32Bytes-16         33680940                35.87 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_64Bytes-16         20185567                58.66 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_100Bytes-16                13742863                86.86 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_256Bytes-16                 7850715               151.8 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_512Bytes-16                 4226336               283.3 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_1000Bytes-16                2120786               568.3 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_4096Bytes-16                 570459              2125 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadSerial/Serial_Read_16384Bytes-16                142891              8472 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_10Goroutines-16            327010598                3.914 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_100Goroutines-16           625220077                2.589 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16Bytes_1000Goroutines-16          615430639                2.799 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_10Goroutines-16            460389994                3.196 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_100Goroutines-16           404379598                4.376 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_21Bytes_1000Goroutines-16          490523025                3.641 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_10Goroutines-16            385261396                4.748 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_100Goroutines-16           390760681                3.960 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_32Bytes_1000Goroutines-16          413073429                3.278 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_10Goroutines-16            245182754                7.568 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_100Goroutines-16           252511812                6.157 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_64Bytes_1000Goroutines-16          238363909                7.617 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_10Goroutines-16           149754214                9.942 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_100Goroutines-16          142018248               10.49 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_100Bytes_1000Goroutines-16         151676746               10.31 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_10Goroutines-16           100000000               12.23 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_100Goroutines-16          100000000               11.58 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_256Bytes_1000Goroutines-16         100000000               12.70 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_10Goroutines-16           54788371                21.60 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_100Goroutines-16          52693305                21.97 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_512Bytes_1000Goroutines-16         54657770                23.91 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_10Goroutines-16          27138499                43.44 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_100Goroutines-16         26837056                45.54 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_1000Bytes_1000Goroutines-16        27158050                44.98 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_10Goroutines-16           7433740               163.9 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_100Goroutines-16          7371459               163.6 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_4096Bytes_1000Goroutines-16         7314062               164.2 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_10Goroutines-16          1876988               640.2 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_100Goroutines-16         1880425               642.6 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrent/Concurrent_Read_16384Bytes_1000Goroutines-16        1891022               637.6 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_4096Bytes-16             533180              2205 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_10000Bytes-16            222049              5388 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_16384Bytes-16            136656              8759 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_65536Bytes-16             34308             35107 ns/op               1 B/op          0 allocs/op
BenchmarkPRNG_ReadSequentialLargeSizes/Serial_Read_Large_1048576Bytes-16            2116            560392 ns/op             496 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_10Goroutines-16           7386999               165.0 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_100Goroutines-16          7367784               162.0 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_4096Bytes_1000Goroutines-16         7452956               162.9 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_10Goroutines-16          3050674               397.3 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_100Goroutines-16         3032727               395.9 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_10000Bytes_1000Goroutines-16        3044510               394.2 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_10Goroutines-16          1883161               642.1 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_100Goroutines-16         1853803               646.7 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_16384Bytes_1000Goroutines-16        1856721               638.6 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_10Goroutines-16           458271              2534 ns/op               2 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_100Goroutines-16          461078              2536 ns/op               2 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_65536Bytes_1000Goroutines-16         463308              2538 ns/op               2 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_10Goroutines-16          29392             40742 ns/op             713 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_100Goroutines-16         29422             40746 ns/op             784 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentLargeSizes/Concurrent_Read_Large_1048576Bytes_1000Goroutines-16        29251             40800 ns/op             753 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_8Bytes-16                                  68608347                16.61 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_16Bytes-16                                 51152076                22.99 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_21Bytes-16                                 42930866                28.03 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_24Bytes-16                                 39859881                29.86 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_32Bytes-16                                 33484634                35.35 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_48Bytes-16                                 24490773                48.32 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_64Bytes-16                                 20496118                59.26 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_128Bytes-16                                12214537                98.06 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_256Bytes-16                                 7697662               155.9 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_512Bytes-16                                 4106547               291.4 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_1024Bytes-16                                2129442               559.6 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_2048Bytes-16                                1000000              1100 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadVariableSizes/Serial_Read_Variable_4096Bytes-16                                 553623              2171 ns/op               0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_10Goroutines-16       733742068                1.985 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_100Goroutines-16      808259511                2.588 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_8Bytes_1000Goroutines-16     686577902                3.023 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_10Goroutines-16      561635169                2.615 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_100Goroutines-16     544451948                3.381 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_16Bytes_1000Goroutines-16    555864582                3.580 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_10Goroutines-16      455553841                4.740 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_100Goroutines-16     423823033                3.833 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_21Bytes_1000Goroutines-16    450182042                4.029 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_10Goroutines-16      389610598                4.173 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_100Goroutines-16     431140887                3.930 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_24Bytes_1000Goroutines-16    440629598                3.393 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_10Goroutines-16      389705337                3.308 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_100Goroutines-16     320857830                4.472 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_32Bytes_1000Goroutines-16    326843707                4.072 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_10Goroutines-16      198504434                8.365 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_100Goroutines-16     190762611                8.312 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_48Bytes_1000Goroutines-16    199009430                7.753 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_10Goroutines-16      235403240                5.751 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_100Goroutines-16     242589153                5.827 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_64Bytes_1000Goroutines-16    244511232                5.508 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_10Goroutines-16     138769131                8.816 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_100Goroutines-16    145349035                8.788 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_128Bytes_1000Goroutines-16   151517790                8.132 ns/op           0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_10Goroutines-16     99757879                11.57 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_100Goroutines-16    94553646                11.95 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_256Bytes_1000Goroutines-16   100000000               11.82 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_10Goroutines-16     53568338                22.13 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_100Goroutines-16    54571398                22.10 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_512Bytes_1000Goroutines-16   53821010                21.66 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_10Goroutines-16    28318333                43.82 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_100Goroutines-16   27987386                42.62 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_1024Bytes_1000Goroutines-16  28246702                42.21 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_10Goroutines-16    14485872                82.88 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_100Goroutines-16   14364454                82.24 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_2048Bytes_1000Goroutines-16  14617455                83.71 ns/op            0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_10Goroutines-16     7475578               162.5 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_100Goroutines-16    7485626               163.3 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadConcurrentVariableSizes/Concurrent_Read_Variable_4096Bytes_1000Goroutines-16   7430336               166.8 ns/op             0 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                                  213           5617983 ns/op           49291 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_10Goroutines-16                2856            417124 ns/op          117495 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_100Goroutines-16               2845            429635 ns/op          117946 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                                   42          28135572 ns/op         1248426 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_10Goroutines-16                 561           2184203 ns/op         2990615 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_100Goroutines-16                510           2253316 ns/op         2261651 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                                  20          56963298 ns/op         5243134 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_10Goroutines-16                219           4637253 ns/op        11012484 B/op          0 allocs/op
BenchmarkPRNG_ReadExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_100Goroutines-16               235           4988574 ns/op        10708902 B/op          0 allocs/op
PASS
ok      github.com/sixafter/nanoid/x/crypto/prng        195.948s
```
</details>

---

## Features

* Cryptographic Security: Utilizes the [ChaCha20](https://pkg.go.dev/golang.org/x/crypto/chacha20) cipher for secure random number generation. 
* Concurrent Support: Includes a thread-safe global `Reader` for concurrent access. 
* Efficient Resource Management: Uses a `sync.Pool` to manage PRNG instances, reducing the overhead on `crypto/rand.Reader`. 
* Extensible API: Allows users to create and manage custom PRNG instances via `NewReader`.

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

---

## Architecture

* Global Reader: A pre-configured io.Reader (`prng.Reader`) manages a pool of PRNG instances for concurrent use. 
* PRNG Instances: Each instance uses ChaCha20, initialized with a unique key and nonce sourced from `crypto/rand.Reader`. 
* Error Handling: The `errorPRNG` ensures safe failure when initialization errors occur. 
* Resource Efficiency: A `sync.Pool` optimizes resource reuse and reduces contention on `crypto/rand.Reader`.

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](../../../LICENSE) file.
