// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"golang.org/x/exp/constraints"
)

// Helper function to create an ASCII-based alphabet of a specified length without duplicates
func makeASCIIBasedAlphabet(length int) string {
	const (
		start = 33  // (!)
		end   = 126 // (~)
	)
	// Calculate the number of unique printable ASCII characters in the range
	rangeSize := end - start + 1

	// Ensure the length does not exceed the number of unique characters
	if length > rangeSize {
		length = rangeSize
	}

	alphabet := make([]byte, length)
	for i := 0; i < length; i++ {
		alphabet[i] = byte(start + i)
	}
	return string(alphabet)
}

// Helper function to create a Unicode alphabet of a specified length without duplicates
// The printable Unicode range is extensive and varies widely across different scripts and symbol sets, as Unicode was designed to represent characters from numerous languages, symbols, and emojis. Unlike ASCII, Unicode doesn’t have a simple, contiguous range for all printable characters. However, there are several primary ranges in Unicode where printable characters are defined:
// 1. Basic Multilingual Plane (BMP): The majority of commonly used printable characters are in the BMP, which spans 0x0020 to 0xFFFF (decimal 32 to 65,535). This plane includes:
//   - Latin characters (including ASCII, starting from 0x0020 for space).
//   - Greek, Cyrillic, Hebrew, Arabic, and other alphabets.
//   - Mathematical symbols, punctuation, and various technical symbols.
//   - Chinese, Japanese, and Korean (CJK) characters.
//   - Emojis and other miscellaneous symbols.
//
// 2. Supplementary Multilingual Plane (SMP): Includes additional printable characters, such as:
//   - Historic scripts.
//   - Musical notation.
//   - Extended emoji sets.
//   - This plane spans 0x10000 to 0x1FFFF.
//
// 3. Supplementary Ideographic Plane (SIP): Contains additional Chinese, Japanese, and Korean ideographs from 0x20000 to 0x2FFFF.
// 4. Other Supplementary Planes: These include various specialized characters, symbols, and private-use areas.
func makeUnicodeAlphabet(length int) string {
	// Greek and Coptic Block
	const (
		start = 0x0370 // (ἰ)
		end   = 0x047F // (ѫ)
	)
	// Calculate the number of unique runes in the range
	rangeSize := end - start + 1

	// Ensure the length does not exceed the number of unique characters
	if length > rangeSize {
		length = rangeSize
	}

	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteRune(rune(start + i))
	}
	return builder.String()
}

type Number interface {
	constraints.Float | constraints.Integer
}

func mean[T Number](data []T) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, d := range data {
		sum += float64(d)
	}
	return sum / float64(len(data))
}

const asciiAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// BenchmarkNanoIDAllocations benchmarks the memory allocations and performance of generating a Nano ID
// with a length of 21 and an alphabet consisting of uppercase letters, lowercase letters, and numbers.
func BenchmarkNanoIDAllocations(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	const idLength = 21

	// Initialize the generator with the specified length and alphabet
	gen, err := NewGenerator(
		WithAlphabet(asciiAlphabet),
		WithLengthHint(idLength))
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}

	// Reset the timer to ignore setup time and track only the ID generation
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = gen.New(idLength)
	}
}

// BenchmarkNanoIDAllocationsConcurrent benchmarks the memory allocations and performance of generating
// a Nano ID concurrently with a length of 21 and an alphabet consisting of uppercase letters,
// lowercase letters, and numbers.
func BenchmarkNanoIDAllocationsConcurrent(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	// Alphabet and ID length for the test
	const idLength = 21

	// Initialize the generator with the specified length and alphabet
	gen, err := NewGenerator(
		WithAlphabet(asciiAlphabet),
		WithLengthHint(idLength))
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}

	// Reset the timer to ignore setup time and track only the ID generation
	b.ResetTimer()

	// Run the benchmark in parallel
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.New(idLength)
			if err != nil {
				b.Errorf("failed to generate ID: %v", err)
			}
		}
	})
}

// BenchmarkGenerator_Read_DefaultLength benchmarks reading into a buffer equal to DefaultLength.
func BenchmarkGenerator_Read_DefaultLength(b *testing.B) {
	gen, ok := DefaultGenerator.(*generator)
	if !ok {
		b.Fatal("DefaultGenerator is not of type *generator")
	}

	buffer := make([]byte, DefaultLength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Read(buffer)
		if err != nil {
			b.Fatalf("Read returned an unexpected error: %v", err)
		}
	}
}

// BenchmarkGenerator_Read_VaryingBufferSizes benchmarks reading into buffers of varying sizes.
func BenchmarkGenerator_Read_VaryingBufferSizes(b *testing.B) {

	bufferSizes := []int{2, 3, 5, 13, 21, 34}
	m := mean(bufferSizes)

	gen, err := NewGenerator(WithLengthHint(uint16(m)))
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("BufferSize_%d", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := gen.Read(buffer)
				if err != nil {
					b.Fatalf("Read returned an unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkGenerator_Read_ZeroLengthBuffer benchmarks reading into a zero-length buffer.
func BenchmarkGenerator_Read_ZeroLengthBuffer(b *testing.B) {
	gen := DefaultGenerator
	buffer := make([]byte, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Read(buffer)
		if err != nil {
			b.Fatalf("Read returned an unexpected error: %v", err)
		}
	}
}

// BenchmarkGenerator_Read_Concurrent benchmarks concurrent reads to assess thread safety and performance.
func BenchmarkGenerator_Read_Concurrent(b *testing.B) {
	gen, err := NewGenerator()
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}

	bufferSize := DefaultLength
	concurrencyLevels := []int{1, 2, 4, 8, 16}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			var wg sync.WaitGroup
			b.SetParallelism(concurrency)
			b.ResetTimer()

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					buffer := make([]byte, bufferSize)
					for j := 0; j < b.N/concurrency; j++ {
						_, err := gen.Read(buffer)
						if err != nil {
							b.Errorf("Read returned an unexpected error: %v", err)
							return
						}
					}
				}()
			}
			wg.Wait()
		})
	}
}

// BenchmarkNanoIDGeneration benchmarks Nano ID generation for varying alphabet types, alphabet lengths, and ID lengths
func BenchmarkNanoIDGeneration(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}

	mean := mean(idLengths)

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64}

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// New the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator without passing 'nil'
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mean)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}

			// Create a sub-benchmark for each alphabet configuration
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					// Create a nested sub-benchmark for each Nano ID length
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						// Reset the timer to exclude setup time
						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							_, err := gen.New(idLen)
							if err != nil {
								b.Fatalf("Failed to generate Nano ID: %v", err)
							}
						}
					})
				}
			})
		}
	}
}

// BenchmarkNanoIDGenerationParallel benchmarks Nano ID generation in parallel for varying configurations
func BenchmarkNanoIDGenerationParallel(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}
	mean := mean(idLengths)

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64}

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// New the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator without passing 'nil'
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mean)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}

			// Create a sub-benchmark for each alphabet configuration
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					// Create a nested sub-benchmark for each Nano ID length
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						// Reset the timer to exclude setup time
						b.ResetTimer()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								_, err := gen.New(idLen)
								if err != nil {
									b.Fatalf("Failed to generate Nano ID: %v", err)
								}
							}
						})
					})
				}
			})
		}
	}
}

// BenchmarkNanoIDWithVaryingAlphabetLengths benchmarks how different alphabet lengths affect Nano ID generation
func BenchmarkNanoIDWithVaryingAlphabetLengths(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64}

	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}
	mean := mean(idLengths)

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// New the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator without passing 'nil'
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mean)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}

			// Create a sub-benchmark for each alphabet configuration
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					// Create a nested sub-benchmark for each Nano ID length
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						// Reset the timer to exclude setup time
						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							_, err := gen.New(idLen)
							if err != nil {
								b.Fatalf("Failed to generate Nano ID: %v", err)
							}
						}
					})
				}
			})
		}
	}
}
