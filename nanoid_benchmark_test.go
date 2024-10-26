// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid_test

import (
	"fmt"
	"github.com/sixafter/nanoid"
	"runtime"
	"testing"
)

// BenchmarkNew benchmarks the New function with default settings.
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := nanoid.New()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkNewSize benchmarks the NewSize function with various sizes.
func BenchmarkNewSize(b *testing.B) {
	sizes := []int{10, 21, 50, 100}
	for _, size := range sizes {
		size := size // Capture range variable
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := nanoid.NewSize(size)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkNewCustom benchmarks the NewCustom function with a custom ASCII alphabet and various sizes.
func BenchmarkNewCustom(b *testing.B) {
	sizes := []int{10, 21, 50, 100}
	customASCIIAlphabet := "abcdef123456"
	for _, size := range sizes {
		size := size // Capture range variable
		b.Run(fmt.Sprintf("Size%d_CustomASCIIAlphabet", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := nanoid.NewCustom(size, customASCIIAlphabet)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkNew_Concurrent benchmarks the New function under concurrent load.
func BenchmarkNew_Concurrent(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8, runtime.NumCPU()}

	for _, concurrency := range concurrencyLevels {
		concurrency := concurrency // Capture range variable
		b.Run(fmt.Sprintf("Concurrency%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := nanoid.New()
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkNewCustom_Concurrent benchmarks the NewCustom function with a custom ASCII alphabet under concurrent load.
func BenchmarkNewCustom_Concurrent(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8, runtime.NumCPU()}
	customASCIIAlphabet := "abcdef123456"

	for _, concurrency := range concurrencyLevels {
		concurrency := concurrency // Capture range variable
		b.Run(fmt.Sprintf("Concurrency%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := nanoid.NewCustom(21, customASCIIAlphabet)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
