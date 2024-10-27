# NanoID

[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid.svg)](https://pkg.go.dev/github.com/sixafter/nanoid)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid)](https://goreportcard.com/report/github.com/sixafter/nanoid)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple, fast, and efficient Go implementation of [NanoID](https://github.com/ai/nanoid), a tiny, secure, URL-friendly, unique string ID generator.

## Features

* **Stateless Design**: Each function operates independently without relying on global state or caches, eliminating the need for synchronization primitives like mutexes. This design ensures predictable behavior and simplifies usage in various contexts. 
* **Cryptographically Secure**: Utilizes Go's crypto/rand package for generating cryptographically secure random numbers. This guarantees that the generated IDs are both unpredictable and suitable for security-sensitive applications. 
* **High Performance**: Optimized algorithms and efficient memory management techniques ensure rapid ID generation. Whether you're generating a few IDs or millions, the library maintains consistent speed and responsiveness. 
* **Memory Efficient**: Implements sync.Pool to reuse byte slices, minimizing memory allocations and reducing garbage collection overhead. This approach significantly enhances performance, especially in high-throughput scenarios. 
* **Thread-Safe**: Designed for safe concurrent use in multi-threaded applications. Multiple goroutines can generate IDs simultaneously without causing race conditions or requiring additional synchronization. 
* **Customizable**: Offers flexibility to specify custom ID lengths and alphabets. Whether you need short, compact IDs or longer, more complex ones, the library can accommodate your specific requirements. 
* **User-Friendly API**: Provides a simple and intuitive API with sensible defaults, making integration straightforward. Developers can start generating IDs with minimal configuration and customize as needed. 
* **Zero External Dependencies**: Relies solely on Go's standard library, ensuring ease of use, compatibility, and minimal footprint within your projects. 
* **Comprehensive Testing**: Includes a robust suite of unit tests and concurrency tests to ensure reliability, correctness, and thread safety. This commitment to quality guarantees consistent performance across different use cases. 
* **Detailed Documentation**: Accompanied by clear and thorough documentation, including examples and usage guidelines. New users can quickly understand how to implement and customize the library to fit their needs. 
* **Efficient Error Handling**: Employs predefined errors to avoid unnecessary allocations, enhancing both performance and clarity in error management. 
* **Optimized for Low Allocations**: Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.

## Installation

To install the package, use:

```sh
go get -u github.com/sixafter/nanoid
```

## Importing the Package

To use the NanoID package in your Go project, import it as follows:

```go
import "github.com/sixafter/nanoid"
```

## Usage

### Generate a Nano ID with Default Settings

Generate a Nano ID using the default size (21 characters) and default alphabet:

```go
id, err := nanoid.New()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID:", id)
```

### Generating a NanoID with Custom Size

Generate a NanoID with a custom length:

```go
id, err := nanoid.NewSize(32)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID of size 32:", id)
```

### Generate a Nano ID with Custom Alphabet

Generate a Nano ID using a custom alphabet:

```go
alphabet := "abcdef123456"
id, err := nanoid.NewCustom(16, alphabet)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID with custom alphabet:", id)
```

### Generate a Nano ID with Custom Random Source

Generate a Nano ID using a custom random source that implements io.Reader:

```go
// Example custom random source (for demonstration purposes)
var myRandomSource io.Reader = myCustomRandomReader{}

id, err := nanoid.NewCustomReader(21, nanoid.DefaultAlphabet, myRandomSource)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated Nano ID with custom random source:", id)
```

**Note:** Replace `myCustomRandomReader{}` with your actual implementation of `io.Reader`.

## Thread Safety

All functions provided by this package are safe for concurrent use by multiple goroutines. Here's an example of generating Nano IDs concurrently:

```go
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/sixafter/nanoid"
)

func main() {
	const numGoroutines = 10
	const idSize = 21

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			id, err := nanoid.New()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Generated Nano ID:", id)
		}()
	}

	wg.Wait()
}
```

## Functions

* `func New() (string, error)`: Generates a Nano ID with the default size (21 characters) and default alphabet.
* `func NewSize(size int) (string, error)`: Generates a Nano ID with a specified size using the default alphabet.
* `func NewCustom(size int, alphabet string) (string, error)`: Generates a Nano ID with a specified size and custom alphabet.
* `func NewCustomReader(size int, alphabet string, rnd io.Reader) (string, error)`: Generates a Nano ID with a specified size, custom alphabet, and custom random source.

## Constants

* `DefaultAlphabet`: The default alphabet used for ID generation: `-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz`
* `DefaultSize`: The default size of the generated ID: `21`

## Performance

The package is optimized for performance and low memory consumption:
* **Efficient Random Byte Consumption**: Uses bitwise operations to extract random bits efficiently. 
* **Avoids `math/big`**: Does not use `math/big`, relying on built-in integer types for calculations. 
* **Minimized System Calls**: Reads random bytes in batches to reduce the number of system calls.

## Execute Benchmarks:

Run the benchmarks using the go test command with the `bench` make target:

```shell
make bench
```

### Interpreting Results:

Sample output might look like this:

```shell
go test -bench=. -benchmem ./...
goos: darwin
goarch: arm64
pkg: github.com/sixafter/nanoid
cpu: Apple M3 Max
BenchmarkNew-16                     	 6329498	       189.2 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewSize/Size10-16          	11600679	       102.4 ns/op	      24 B/op	       2 allocs/op
BenchmarkNewSize/Size21-16          	 6384469	       186.7 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewSize/Size50-16          	 2680179	       448.2 ns/op	     104 B/op	       6 allocs/op
BenchmarkNewSize/Size100-16         	 1387914	       863.3 ns/op	     192 B/op	      11 allocs/op
BenchmarkNewCustom/Size10_CustomASCIIAlphabet-16         	 9306187	       128.8 ns/op	      24 B/op	       2 allocs/op
BenchmarkNewCustom/Size21_CustomASCIIAlphabet-16         	 5062975	       239.4 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom/Size50_CustomASCIIAlphabet-16         	 2322037	       515.3 ns/op	     101 B/op	       5 allocs/op
BenchmarkNewCustom/Size100_CustomASCIIAlphabet-16        	 1235755	       972.0 ns/op	     182 B/op	       9 allocs/op
BenchmarkNew_Concurrent/Concurrency1-16                  	 2368245	       513.1 ns/op	      40 B/op	       3 allocs/op
BenchmarkNew_Concurrent/Concurrency2-16                  	 1940826	       609.5 ns/op	      40 B/op	       3 allocs/op
BenchmarkNew_Concurrent/Concurrency4-16                  	 1986049	       585.6 ns/op	      40 B/op	       3 allocs/op
BenchmarkNew_Concurrent/Concurrency8-16                  	 1999959	       602.2 ns/op	      40 B/op	       3 allocs/op
BenchmarkNew_Concurrent/Concurrency16-16                 	 2018793	       595.6 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom_Concurrent/Concurrency1-16            	 1960315	       611.7 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom_Concurrent/Concurrency2-16            	 1790460	       673.7 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom_Concurrent/Concurrency4-16            	 1766841	       670.7 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom_Concurrent/Concurrency8-16            	 1768189	       677.4 ns/op	      40 B/op	       3 allocs/op
BenchmarkNewCustom_Concurrent/Concurrency16-16           	 1765303	       689.5 ns/op	      40 B/op	       3 allocs/op
PASS
ok  	github.com/sixafter/nanoid	33.279s
```

* `ns/op` (Nanoseconds per Operation):
  * Indicates the average time taken per operation. 
  * Lower values signify better CPU performance. 
* `B/op` (Bytes Allocated per Operation):
  * Shows the average number of bytes allocated per operation. 
  * `0 B/op` indicates no heap allocations, which is optimal. 
* `allocs/op` (Allocations per Operation):
  * Represents the average number of memory allocations per operation. 
  * `0 allocs/op` is ideal as it indicates no heap allocations.

## Contributing

Contributions are welcome. For larger or more material changes, please create a issue so we can discuss.

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See [LICENSE](LICENSE) file.

