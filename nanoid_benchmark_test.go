// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"strconv"
	"testing"
)

// BenchmarkGenerate benchmarks the default Generate function.
func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Generate()
	}
}

// BenchmarkGenerateSize benchmarks GenerateSize with varying sizes.
func BenchmarkGenerateSize(b *testing.B) {
	sizes := []int{8, 16, 32, 64, 128}
	for _, size := range sizes {
		b.Run("Size"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = GenerateSize(size)
			}
		})
	}
}

// BenchmarkGenerateCustom benchmarks GenerateCustom with different alphabets and sizes.
func BenchmarkGenerateCustom(b *testing.B) {
	tests := []struct {
		size     int
		alphabet string
	}{
		{21, defaultAlphabet},
		{16, "abcdef0123456789"},
		{32, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	}

	for _, test := range tests {
		label := "Size" + strconv.Itoa(test.size) + "Alphabet" + strconv.Itoa(len(test.alphabet))
		b.Run(label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = GenerateCustom(test.size, test.alphabet)
			}
		})
	}
}
