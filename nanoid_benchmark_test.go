// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"fmt"
	"strings"
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

// BenchmarkNanoIDGeneration benchmarks Nano ID generation for varying alphabet types, alphabet lengths, and ID lengths
func BenchmarkNanoIDGeneration(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}

	mean := mean(idLengths)

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64, 95}

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
				WithLengthHint(int(mean)),
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
	alphabetLengths := []int{2, 16, 32, 64, 95}

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
				WithLengthHint(int(mean)),
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
	alphabetLengths := []int{2, 16, 32, 64, 95}

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
				WithLengthHint(int(mean)),
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
