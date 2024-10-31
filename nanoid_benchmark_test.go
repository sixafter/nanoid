// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"fmt"
	"strings"
	"testing"
)

// Helper function to create an ASCII-based alphabet of a specified length
func makeASCIIBasedAlphabet(length int) string {
	const (
		start = 33  // '!'
		end   = 126 // '~'
	)
	// Ensure the length is within the printable ASCII range
	if length < 2 {
		length = 2
	}
	if length > end-start+1 {
		length = end - start + 1
	}
	alphabet := make([]byte, length)
	for i := 0; i < length; i++ {
		alphabet[i] = byte(start + i)
	}
	return string(alphabet)
}

// Helper function to create a Unicode alphabet of a specified length
func makeUnicodeAlphabet(length int) string {
	const (
		start = 0x0905 // 'अ'
		end   = 0x0939 // 'ह'
	)
	// Ensure the length is within the specified Unicode range
	if length < 2 {
		length = 2
	}
	if length > end-start+1 {
		length = end - start + 1
	}
	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteRune(rune(start + i))
	}
	return builder.String()
}

// BenchmarkNanoIDGeneration benchmarks Nano ID generation for varying alphabet types, alphabet lengths, and ID lengths
func BenchmarkNanoIDGeneration(b *testing.B) {
	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64, 95}

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// Generate the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator
			g, err := New(alphabet, nil)
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
							_, err := g.Generate(idLen)
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
	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64, 95}

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// Generate the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator
			g, err := New(alphabet, nil)
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
								_, err := g.Generate(idLen)
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
	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	// Define the alphabet lengths to test
	alphabetLengths := []int{2, 16, 32, 64, 95}

	// Define the Nano ID lengths to test
	idLengths := []int{8, 16, 21, 32, 64, 128}

	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			// Generate the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}

			// Initialize the generator
			g, err := New(alphabet, nil)
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
							_, err := g.Generate(idLen)
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
