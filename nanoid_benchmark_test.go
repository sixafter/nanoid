// Copyright (c) 2024-2026 Six After, Inc.
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

// Benchmark_Allocations_Serial measures allocations and nanoseconds per operation
// for generating a Nano ID with a fixed ASCII alphabet of 62 characters (A-Za-z0-9)
// and a length of 21. This benchmark is single-threaded and useful for tracking
// the baseline allocation cost and performance for the most common usage.
//
// Reports: allocs/op and ns/op.
func Benchmark_Allocations_Serial(b *testing.B) {
	b.ReportAllocs()
	const idLength = 21

	// Setup: generator uses standard ASCII alphabet and configured length.
	gen, err := NewGenerator(
		WithAlphabet(asciiAlphabet),
		WithLengthHint(idLength),
	)
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}
	b.ResetTimer()
	for b.Loop() {
		_, err = gen.NewWithLength(idLength)
	}
}

// Benchmark_Allocations_Parallel measures allocation and throughput for Nano ID generation
// under high concurrency, using the standard ASCII alphabet and ID length of 21.
// This tests thread-safety, internal contention, and real-world performance for concurrent workloads.
//
// Reports: allocs/op and ns/op under parallel execution.
func Benchmark_Allocations_Parallel(b *testing.B) {
	b.ReportAllocs()
	const idLength = 21

	gen, err := NewGenerator(
		WithAlphabet(asciiAlphabet),
		WithLengthHint(idLength),
	)
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.NewWithLength(idLength)
			if err != nil {
				b.Errorf("failed to generate ID: %v", err)
			}
		}
	})
}

// Benchmark_Read_DefaultLength benchmarks the performance and allocation profile of
// reading a Nano ID into a buffer of DefaultLength. This exercises the global Generator's
// Read method and measures its suitability for typical random byte generation scenarios.
//
// Reports: allocs/op and ns/op.
func Benchmark_Read_DefaultLength(b *testing.B) {
	buffer := make([]byte, DefaultLength)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, err := Generator.Read(buffer)
		if err != nil {
			b.Fatalf("Read returned an unexpected error: %v", err)
		}
	}
}

// Benchmark_Read_VaryingBufferSizes benchmarks the allocation and runtime behavior of
// Generator.Read across a range of buffer sizes. This is useful for analyzing whether
// performance and memory behavior remain stable as buffer size scales.
//
// Reports: allocs/op and ns/op for each buffer size.
func Benchmark_Read_VaryingBufferSizes(b *testing.B) {
	bufferSizes := []int{2, 3, 5, 13, 21, 34}
	m := mean(bufferSizes)
	b.ReportAllocs()
	gen, err := NewGenerator(WithLengthHint(uint16(m)))
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}
	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("BufferSize_%d", size), func(b *testing.B) {
			buffer := make([]byte, size)
			b.ResetTimer()
			for b.Loop() {
				_, err := gen.Read(buffer)
				if err != nil {
					b.Fatalf("Read returned an unexpected error: %v", err)
				}
			}
		})
	}
}

// Benchmark_Read_ZeroLengthBuffer ensures that reading into a zero-length buffer is a zero-cost no-op
// and does not cause allocations or panics. Verifies the contract for len(p) == 0.
func Benchmark_Read_ZeroLengthBuffer(b *testing.B) {
	buffer := make([]byte, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, err := Generator.Read(buffer)
		if err != nil {
			b.Fatalf("Read returned an unexpected error: %v", err)
		}
	}
}

// Benchmark_Read_Concurrent exercises Generator.Read in a concurrent setting, varying the number
// of goroutines (concurrency levels) to simulate realistic multi-threaded workloads. It assesses
// thread-safety, internal pool contention, and scaling behavior.
//
// Each sub-benchmark runs at a different concurrency level, measured in parallel goroutines.
func Benchmark_Read_Concurrent(b *testing.B) {
	bufferSize := DefaultLength
	concurrencyLevels := []int{1, 2, 4, 8, 16}
	b.ReportAllocs()
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
						_, err := Generator.Read(buffer)
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

// Benchmark_Alphabet_Varying_Serial systematically benchmarks Nano ID generation over
// a matrix of alphabet types (ASCII, Unicode), alphabet lengths, and ID lengths, in serial.
// This exposes the impact of alphabet and output size on allocations and runtime performance.
// Each sub-benchmark explores one configuration.
func Benchmark_Alphabet_Varying_Serial(b *testing.B) {
	b.ReportAllocs()
	idLengths := []int{8, 16, 21, 32, 64, 128}
	mn := mean(idLengths)
	alphabetLengths := []int{2, 16, 32, 64}
	alphabetTypes := []string{"ASCII", "Unicode"}
	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mn)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						b.ResetTimer()
						for b.Loop() {
							_, err := gen.NewWithLength(idLen)
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

// Benchmark_Alphabet_Varying_Parallel is the parallel/concurrent version of Benchmark_Alphabet_Varying_Serial.
// It benchmarks all combinations of alphabet type, alphabet length, and ID length under high concurrency
// to measure thread-safety, lock contention, and throughput in a parallel setting.
func Benchmark_Alphabet_Varying_Parallel(b *testing.B) {
	b.ReportAllocs()
	idLengths := []int{8, 16, 21, 32, 64, 128}
	mn := mean(idLengths)
	alphabetLengths := []int{2, 16, 32, 64}
	alphabetTypes := []string{"ASCII", "Unicode"}
	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mn)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						b.ResetTimer()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								_, err := gen.NewWithLength(idLen)
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

// Benchmark_Alphabet_Varying_Length_Varying_Serial systematically benchmarks
// the effect of both alphabet length and ID length in serial for both ASCII and Unicode alphabets.
// This allows detailed analysis of scaling behavior as both parameters increase.
func Benchmark_Alphabet_Varying_Length_Varying_Serial(b *testing.B) {
	b.ReportAllocs()
	alphabetTypes := []string{"ASCII", "Unicode"}
	alphabetLengths := []int{2, 16, 32, 64}
	idLengths := []int{8, 16, 21, 32, 64, 128}
	mn := mean(idLengths)
	for _, alphabetType := range alphabetTypes {
		for _, alphaLen := range alphabetLengths {
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(alphaLen)
			} else {
				alphabet = makeUnicodeAlphabet(alphaLen)
			}
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mn)),
			)
			if err != nil {
				b.Fatalf("Failed to create generator with %s alphabet of length %d: %v", alphabetType, alphaLen, err)
			}
			b.Run(fmt.Sprintf("%s_AlphabetLen%d", alphabetType, alphaLen), func(b *testing.B) {
				for _, idLen := range idLengths {
					b.Run(fmt.Sprintf("IDLen%d", idLen), func(b *testing.B) {
						b.ResetTimer()
						for b.Loop() {
							_, err := gen.NewWithLength(idLen)
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
