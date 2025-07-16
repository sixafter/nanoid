// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"fmt"
	"testing"
)

// For benchmarking sync.Pool get/put only (DRBG instancing contention, not output).
func (r *reader) syncPoolGetPut() {
	dr := r.pool.Get().(*drbg)
	r.pool.Put(dr)
}

func BenchmarkDRBG_Concurrent_SyncPool_Baseline(b *testing.B) {
	rdr, _ := NewReader()
	goroutineCounts := []int{2, 4, 8, 16, 32, 64, 128}
	if r, ok := rdr.(*reader); ok {
		for _, count := range goroutineCounts {
			benchName := fmt.Sprintf("G%d", count)
			b.Run(benchName, func(b *testing.B) {
				b.SetParallelism(count)
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						r.syncPoolGetPut()
					}
				})
			})
		}
	}
}

func BenchmarkDRBG_ReadSerial(b *testing.B) {
	bufferSizes := []int{8, 16, 21, 32, 64, 100, 256, 512, 1000, 4096, 16384}
	for _, size := range bufferSizes {
		size := size
		b.Run(fmt.Sprintf("Serial_Read_%dBytes", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Reader.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkDRBG_ReadConcurrent(b *testing.B) {
	bufferSizes := []int{16, 21, 32, 64, 100, 256, 512, 1000, 4096, 16384}
	goroutineCounts := []int{2, 4, 8, 16, 32, 64, 128}
	for _, size := range bufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc
			b.Run(fmt.Sprintf("Concurrent_Read_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				buffer := make([]byte, size)
				b.SetParallelism(gc)
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, err := Reader.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

func BenchmarkDRBG_ReadSequentialLargeSizes(b *testing.B) {
	largeBufferSizes := []int{4096, 10000, 16384, 65536, 1048576}
	for _, size := range largeBufferSizes {
		size := size
		b.Run(fmt.Sprintf("Serial_Read_Large_%dBytes", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Reader.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkDRBG_ReadConcurrentLargeSizes(b *testing.B) {
	largeBufferSizes := []int{4096, 10000, 16384, 65536, 1048576}
	goroutineCounts := []int{2, 4, 8, 16, 32, 64, 128}
	for _, size := range largeBufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc
			b.Run(fmt.Sprintf("Concurrent_Read_Large_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				buffer := make([]byte, size)
				b.SetParallelism(gc)
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, err := Reader.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

func BenchmarkDRBG_ReadVariableSizes(b *testing.B) {
	variableBufferSizes := []int{8, 16, 21, 24, 32, 48, 64, 128, 256, 512, 1024, 2048, 4096}
	for _, size := range variableBufferSizes {
		size := size
		b.Run(fmt.Sprintf("Serial_Read_Variable_%dBytes", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Reader.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkDRBG_ReadConcurrentVariableSizes(b *testing.B) {
	variableBufferSizes := []int{8, 16, 21, 24, 32, 48, 64, 128, 256, 512, 1024, 2048, 4096}
	goroutineCounts := []int{2, 4, 8, 16, 32, 64, 128}
	for _, size := range variableBufferSizes {
		for _, gc := range goroutineCounts {
			size, gc := size, gc
			b.Run(fmt.Sprintf("Concurrent_Read_Variable_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				buffer := make([]byte, size)
				b.SetParallelism(gc)
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, err := Reader.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}

func BenchmarkDRBG_ReadExtremeSizes(b *testing.B) {
	extremeBufferSizes := []int{10485760, 52428800, 104857600} // 10MB, 50MB, 100MB
	for _, size := range extremeBufferSizes {
		size := size
		// Serial
		b.Run(fmt.Sprintf("Serial_Read_Extreme_%dBytes", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Reader.Read(buffer)
				if err != nil {
					b.Fatalf("Read failed: %v", err)
				}
			}
		})
		// Concurrent
		goroutineCounts := []int{2, 4, 8, 16, 32, 64, 128}
		for _, gc := range goroutineCounts {
			gc := gc
			b.Run(fmt.Sprintf("Concurrent_Read_Extreme_%dBytes_%dGoroutines", size, gc), func(b *testing.B) {
				buffer := make([]byte, size)
				b.SetParallelism(gc)
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, err := Reader.Read(buffer)
						if err != nil {
							b.Fatalf("Read failed: %v", err)
						}
					}
				})
			})
		}
	}
}
