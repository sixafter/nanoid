// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package prng

import (
	"sync"
	"testing"

	"github.com/google/uuid"
)

// benchConcurrent is a benchmark helper that runs the provided function fn
// across the specified number of goroutines, distributing b.N iterations as evenly as possible.
// It is designed to test concurrent throughput or contention scenarios in benchmarks.
func benchConcurrent(b *testing.B, fn func(), goroutines int) {
	nPerG := b.N / goroutines
	rem := b.N % goroutines
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < goroutines; i++ {
		iters := nPerG
		if i < rem {
			iters++
		}
		wg.Add(1)
		go func(iters int) {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				fn()
			}
		}(iters)
	}
	wg.Wait()
}

// itoa converts an integer to its decimal string representation without heap allocations.
// It is suitable for constructing sub-benchmark names efficiently.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [12]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = '0' + byte(i%10)
		i /= 10
	}
	return string(buf[pos:])
}

// BenchmarkUUID_v4_Default_Serial measures the baseline performance of uuid.New()
// using the default (math/rand) random source in a serial (single-threaded) loop.
// This establishes a baseline for UUID v4 generation throughput and allocations.
func BenchmarkUUID_v4_Default_Serial(b *testing.B) {
	uuid.SetRand(nil)
	b.ReportAllocs()
	for b.Loop() {
		_ = uuid.New()
	}
}

// BenchmarkUUID_v4_Default_Parallel benchmarks uuid.New() using the default random source
// with the built-in Go testing RunParallel helper. It measures performance and allocation
// characteristics under parallel workload, simulating many goroutines calling uuid.New() simultaneously.
func BenchmarkUUID_v4_Default_Parallel(b *testing.B) {
	uuid.SetRand(nil)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
		}
	})
}

// BenchmarkUUID_v4_Default_Concurrent benchmarks uuid.New() using the default random source
// across a range of goroutine counts (4 to 256). It uses benchConcurrent to distribute work evenly.
// This measures scalability and contention as goroutines increase.
func BenchmarkUUID_v4_Default_Concurrent(b *testing.B) {
	uuid.SetRand(nil)
	for _, gr := range []int{2, 4, 8, 16, 32, 64, 128, 256} {
		b.Run("Goroutines_"+itoa(gr), func(b *testing.B) {
			benchConcurrent(b, func() { _ = uuid.New() }, gr)
		})
	}
}

// BenchmarkUUID_v4_CSPRNG_Serial measures the performance of uuid.New() using a custom CSPRNG-based
// Reader (e.g., cryptographic PRNG) in a serial (single-threaded) loop. This allows comparison
// of cryptographic vs. default random sources for UUID v4 generation.
func BenchmarkUUID_v4_CSPRNG_Serial(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	b.ReportAllocs()
	for b.Loop() {
		_ = uuid.New()
	}
}

// BenchmarkUUID_v4_CSPRNG_Parallel benchmarks uuid.New() with a CSPRNG random source using
// Go's testing RunParallel helper. It measures performance and allocations for cryptographically
// secure UUID v4 generation under parallel conditions.
func BenchmarkUUID_v4_CSPRNG_Parallel(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
		}
	})
}

// BenchmarkUUID_v4_CSPRNG_Concurrent benchmarks uuid.New() with a CSPRNG random source
// across multiple goroutine counts (4 to 256), using benchConcurrent to simulate highly concurrent use.
// This measures scalability, contention, and performance of cryptographic UUID generation under load.
func BenchmarkUUID_v4_CSPRNG_Concurrent(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	for _, gr := range []int{2, 4, 8, 16, 32, 64, 128, 256} {
		b.Run("Goroutines_"+itoa(gr), func(b *testing.B) {
			benchConcurrent(b, func() { _ = uuid.New() }, gr)
		})
	}
}
