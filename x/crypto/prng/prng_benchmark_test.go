// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package prng

import (
	"fmt"
	"testing"
)

// BenchmarkPRNG_ReadSerial benchmarks the Read method of prng.Reader with various buffer sizes in serial.
func BenchmarkPRNG_ReadSerial(b *testing.B) {
	// Define the buffer sizes to benchmark.
	bufferSizes := []int{8, 16, 21, 32, 64, 100, 256, 512, 1000, 4096, 16384}

	for _, size := range bufferSizes {
		size := size // Capture range variable
		b.Run(fmt.Sprintf("Serial_Read_%dBytes", size), func(b *testing.B) {
			rdr, err := NewReader()
			if err != nil {
				b.Fatalf("NewReader failed: %v", err)
			}
			buffer := make([]byte, size)
			b.ResetTimer()   // Reset the timer to exclude setup time
			b.ReportAllocs() // Report memory allocations
			for b.Loop() {
				_, err = rdr.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkPRNG_ReadConcurrent benchmarks the Read method of prng.Reader under concurrent access with varying buffer sizes and goroutine counts.
func BenchmarkPRNG_ReadConcurrent(b *testing.B) {
	// Define the buffer sizes and goroutine counts to benchmark concurrently.
	bufferSizes := []int{16, 21, 32, 64, 100, 256, 512, 1000, 4096, 16384}
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64, 128} // Varying goroutine counts

	for _, size := range bufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc // Capture range variables
			b.Run(fmt.Sprintf("Concurrent_Read_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					rdr, err := NewReader()
					if err != nil {
						b.Fatalf("NewReader failed: %v", err)
					}
					buffer := make([]byte, size)
					for pb.Next() {
						_, err = rdr.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

// BenchmarkPRNG_ReadSequentialLargeSizes benchmarks the Read method with large buffer sizes in serial.
func BenchmarkPRNG_ReadSequentialLargeSizes(b *testing.B) {
	// Define large buffer sizes to benchmark in serial.
	largeBufferSizes := []int{4096, 10000, 16384, 65536, 1048576}

	for _, size := range largeBufferSizes {
		size := size // Capture range variable
		b.Run(fmt.Sprintf("Serial_Read_Large_%dBytes", size), func(b *testing.B) {
			rdr, err := NewReader()
			if err != nil {
				b.Fatalf("NewReader failed: %v", err)
			}
			buffer := make([]byte, size)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, err = rdr.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkPRNG_ReadConcurrentLargeSizes benchmarks the Read method with large buffer sizes under concurrent access.
func BenchmarkPRNG_ReadConcurrentLargeSizes(b *testing.B) {
	// Define large buffer sizes and goroutine counts to benchmark concurrently.
	largeBufferSizes := []int{4096, 10000, 16384, 65536, 1048576}
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64, 128} // Varying goroutine counts

	for _, size := range largeBufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc // Capture range variables
			b.Run(fmt.Sprintf("Concurrent_Read_Large_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					rdr, err := NewReader()
					if err != nil {
						b.Fatalf("NewReader failed: %v", err)
					}
					buffer := make([]byte, size)
					for pb.Next() {
						_, err = rdr.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

// BenchmarkPRNG_ReadVariableSizes benchmarks the Read method with variable buffer sizes in serial.
func BenchmarkPRNG_ReadVariableSizes(b *testing.B) {
	// Define a range of buffer sizes to benchmark in serial.
	variableBufferSizes := []int{8, 16, 21, 24, 32, 48, 64, 128, 256, 512, 1024, 2048, 4096}

	for _, size := range variableBufferSizes {
		size := size // Capture range variable
		b.Run(fmt.Sprintf("Serial_Read_Variable_%dBytes", size), func(b *testing.B) {
			rdr, err := NewReader()
			if err != nil {
				b.Fatalf("NewReader failed: %v", err)
			}
			buffer := make([]byte, size)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, err = rdr.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkPRNG_ReadConcurrentVariableSizes benchmarks the Read method with variable buffer sizes under concurrent access.
func BenchmarkPRNG_ReadConcurrentVariableSizes(b *testing.B) {
	// Define a range of buffer sizes and goroutine counts to benchmark concurrently.
	variableBufferSizes := []int{8, 16, 21, 24, 32, 48, 64, 128, 256, 512, 1024, 2048, 4096}
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64, 128}

	for _, size := range variableBufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc // Capture range variables
			b.Run(fmt.Sprintf("Concurrent_Read_Variable_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					rdr, err := NewReader()
					if err != nil {
						b.Fatalf("NewReader failed: %v", err)
					}
					buffer := make([]byte, size)
					for pb.Next() {
						_, err = rdr.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

// BenchmarkPRNG_ReadExtremeSizes benchmarks the Read method with extremely large buffer sizes.
func BenchmarkPRNG_ReadExtremeSizes(b *testing.B) {
	// Define extremely large buffer sizes to benchmark in serial and concurrent contexts.
	extremeBufferSizes := []int{10485760, 52428800, 104857600} // 10MB, 50MB, 100MB

	for _, size := range extremeBufferSizes {
		size := size // Capture range variable
		// Serial Benchmark
		b.Run(fmt.Sprintf("Serial_Read_Extreme_%dBytes", size), func(b *testing.B) {
			rdr, err := NewReader()
			if err != nil {
				b.Fatalf("NewReader failed: %v", err)
			}
			buffer := make([]byte, size)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, err = rdr.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})

		// Concurrent Benchmark
		goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64, 128}
		for _, gc := range goroutineCounts {
			gc := gc // Capture range variable
			b.Run(fmt.Sprintf("Concurrent_Read_Extreme_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					rdr, err := NewReader()
					if err != nil {
						b.Fatalf("NewReader failed: %v", err)
					}
					buffer := make([]byte, size)
					for pb.Next() {
						_, err = rdr.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}
