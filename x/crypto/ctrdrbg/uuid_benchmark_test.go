// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"sync"
	"testing"

	"github.com/google/uuid"
)

// benchConcurrent is a benchmark helper for concurrent scenarios.
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

// itoa converts an integer to its decimal string representation.
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

// Baseline: uuid.New() with default random (math/rand) - serial.
func BenchmarkUUID_v4_Default_Serial(b *testing.B) {
	uuid.SetRand(nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
}

// Baseline: uuid.New() with default random (math/rand) - parallel.
func BenchmarkUUID_v4_Default_Parallel(b *testing.B) {
	uuid.SetRand(nil)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
		}
	})
}

// Baseline: uuid.New() with default random (math/rand) - concurrent.
func BenchmarkUUID_v4_Default_Concurrent(b *testing.B) {
	uuid.SetRand(nil)
	for _, gr := range []int{2, 4, 8, 16, 32, 64, 128, 256} {
		b.Run("Goroutines_"+itoa(gr), func(b *testing.B) {
			benchConcurrent(b, func() { _ = uuid.New() }, gr)
		})
	}
}

// CTRDRBG: uuid.New() with ctrdrbg.Reader - serial.
func BenchmarkUUID_v4_CTRDRBG_Serial(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
}

// CTRDRBG: uuid.New() with ctrdrbg.Reader - parallel.
func BenchmarkUUID_v4_CTRDRBG_Parallel(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
		}
	})
}

// CTRDRBG: uuid.New() with ctrdrbg.Reader - concurrent.
func BenchmarkUUID_v4_CTRDRBG_Concurrent(b *testing.B) {
	uuid.SetRand(Reader)
	defer uuid.SetRand(nil)
	for _, gr := range []int{2, 4, 8, 16, 32, 64, 128, 256} {
		b.Run("Goroutines_"+itoa(gr), func(b *testing.B) {
			benchConcurrent(b, func() { _ = uuid.New() }, gr)
		})
	}
}
