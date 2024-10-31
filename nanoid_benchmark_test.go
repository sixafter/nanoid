// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"fmt"
	"testing"
)

func BenchmarkDefaultGenerate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDefaultGenerateSize(b *testing.B) {
	lengths := []int{8, 16, 21, 32, 64, 128}
	for _, size := range lengths {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := GenerateSize(size)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDefaultGenerateParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Generate()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDefaultGenerateSizeParallel(b *testing.B) {
	lengths := []int{8, 16, 21, 32, 64, 128}
	for _, size := range lengths {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := GenerateSize(size)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkGeneratorGenerate(b *testing.B) {
	g, err := New(DefaultAlphabet, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := g.Generate(DefaultSize)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGeneratorGenerateParallel(b *testing.B) {
	g, err := New(DefaultAlphabet, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := g.Generate(DefaultSize)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCustomAlphabet(b *testing.B) {
	alphabet := "0123456789abcdef"
	g, err := New(alphabet, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := g.Generate(DefaultSize)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCustomAlphabetParallel(b *testing.B) {
	alphabet := "0123456789abcdef"
	g, err := New(alphabet, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := g.Generate(DefaultSize)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewGenerator(b *testing.B) {
	alphabet := DefaultAlphabet
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := New(alphabet, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCustomAlphabetLengths(b *testing.B) {
	lengths := []int{2, 16, 32, 64, 95}
	for _, alphaLen := range lengths {
		b.Run(fmt.Sprintf("AlphabetLength%d", alphaLen), func(b *testing.B) {
			alphabet := makeAlphabet(alphaLen)
			g, err := New(alphabet, nil)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := g.Generate(DefaultSize)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCustomAlphabetLengthsParallel(b *testing.B) {
	lengths := []int{2, 16, 32, 64, 95}
	for _, alphaLen := range lengths {
		b.Run(fmt.Sprintf("AlphabetLength%d", alphaLen), func(b *testing.B) {
			alphabet := makeAlphabet(alphaLen)
			g, err := New(alphabet, nil)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := g.Generate(DefaultSize)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// Helper function to create an alphabet of a specified length using printable ASCII characters
func makeAlphabet(length int) string {
	const (
		printableASCIIStart = 33  // '!' character
		printableASCIIEnd   = 126 // '~' character
	)

	printableASCIIRange := printableASCIIEnd - printableASCIIStart + 1

	if length < 2 {
		length = 2 // Minimum valid length
	}
	if length > printableASCIIRange {
		length = printableASCIIRange // Maximum valid length
	}

	alphabet := make([]byte, length)
	for i := 0; i < length; i++ {
		alphabet[i] = byte(printableASCIIStart + i)
	}
	return string(alphabet)
}
