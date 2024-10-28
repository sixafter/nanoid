// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"strconv"
	"testing"
)

// BenchmarkGenerateDefault benchmarks the default generator with default alphabet and ID length.
func BenchmarkGenerateDefault(b *testing.B) {
	b.ReportAllocs()
	gen, err := New(DefaultAlphabet, nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := DefaultSize // Default Nano ID length
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateCustomAlphabet benchmarks a custom alphabet generator with ID length 10.
func BenchmarkGenerateCustomAlphabet(b *testing.B) {
	b.ReportAllocs()
	alphabet := "ABCDEF"
	gen, err := New(alphabet, nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 10
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateShortID benchmarks the generator with a short ID length.
func BenchmarkGenerateShortID(b *testing.B) {
	b.ReportAllocs()
	gen, err := New("abcdef", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 5
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateLongID benchmarks the generator with a long ID length.
func BenchmarkGenerateLongID(b *testing.B) {
	b.ReportAllocs()
	gen, err := New("abcdef", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 50
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateMaxAlphabet benchmarks the generator with the maximum allowed alphabet size (256 characters).
func BenchmarkGenerateMaxAlphabet(b *testing.B) {
	b.ReportAllocs()
	// Create an alphabet of 256 unique characters
	alphabet := make([]byte, 256)
	for i := 0; i < 256; i++ {
		alphabet[i] = byte(i)
	}
	gen, err := New(string(alphabet), nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 10
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateMinAlphabet benchmarks the generator with the minimum allowed alphabet size (2 characters).
func BenchmarkGenerateMinAlphabet(b *testing.B) {
	b.ReportAllocs()
	alphabet := "AB"
	gen, err := New(alphabet, nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 10
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateWithBufferPool benchmarks the generator with buffer pooling enabled.
func BenchmarkGenerateWithBufferPool(b *testing.B) {
	b.ReportAllocs()
	gen, err := New("abcdef", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 10
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(idLength)
		if err != nil {
			b.Fatalf("GenerateSize failed: %v", err)
		}
	}
}

// BenchmarkGenerateDefaultParallel benchmarks the default generator in a parallel/concurrent setting.
func BenchmarkGenerateDefaultParallel(b *testing.B) {
	gen, err := New(DefaultAlphabet, nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := DefaultSize // Default Nano ID length

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate(idLength)
			if err != nil {
				b.Fatalf("GenerateSize failed: %v", err)
			}
		}
	})
}

// BenchmarkGenerateCustomAlphabetParallel benchmarks a custom alphabet generator in parallel.
func BenchmarkGenerateCustomAlphabetParallel(b *testing.B) {
	gen, err := New("ABCDEF", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 10

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate(idLength)
			if err != nil {
				b.Fatalf("GenerateSize failed: %v", err)
			}
		}
	})
}

// BenchmarkGenerateShortIDParallel benchmarks the generator with short IDs in parallel.
func BenchmarkGenerateShortIDParallel(b *testing.B) {
	gen, err := New("abcdef", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 5

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate(idLength)
			if err != nil {
				b.Fatalf("GenerateSize failed: %v", err)
			}
		}
	})
}

// BenchmarkGenerateLongIDParallel benchmarks the generator with long IDs in parallel.
func BenchmarkGenerateLongIDParallel(b *testing.B) {
	gen, err := New("abcdef", nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 50

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate(idLength)
			if err != nil {
				b.Fatalf("GenerateSize failed: %v", err)
			}
		}
	})
}

// BenchmarkGenerateExtremeConcurrency benchmarks the generator under extreme concurrency.
func BenchmarkGenerateExtremeConcurrency(b *testing.B) {
	gen, err := New(DefaultAlphabet, nil)
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	idLength := 21 // Default Nano ID length

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate(idLength)
			if err != nil {
				b.Fatalf("GenerateSize failed: %v", err)
			}
		}
	})
}

// BenchmarkGenerateDifferentLengths benchmarks the generator with varying ID lengths.
func BenchmarkGenerateDifferentLengths(b *testing.B) {
	alphabet := "abcdef"
	lengths := []int{5, 10, 20, 50, 100}

	for _, length := range lengths {
		// Correctly format the benchmark name
		name := "Length_" + strconv.Itoa(length)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			gen, err := New(alphabet, nil)
			if err != nil {
				b.Fatalf("Failed to create generator: %v", err)
			}

			for i := 0; i < b.N; i++ {
				_, err := gen.Generate(length)
				if err != nil {
					b.Fatalf("GenerateSize failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkGenerateDifferentAlphabets benchmarks the generator with different alphabet sizes.
func BenchmarkGenerateDifferentAlphabets(b *testing.B) {
	alphabets := []string{
		"AB",                                     // 2 characters
		"ABCDEF",                                 // 6 characters
		"abcdefghijklmnopqrstuvwxyz",             // 26 characters
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", // 64 characters
	}

	idLength := 10
	for _, alphabet := range alphabets {
		// Correctly format the benchmark name
		name := "Alphabet_" + strconv.Itoa(len(alphabet))
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			gen, err := New(alphabet, nil)
			if err != nil {
				b.Fatalf("Failed to create generator: %v", err)
			}

			for i := 0; i < b.N; i++ {
				_, err := gen.Generate(idLength)
				if err != nil {
					b.Fatalf("GenerateSize failed: %v", err)
				}
			}
		})
	}
}
