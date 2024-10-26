// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := New()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateSize(b *testing.B) {
	sizes := []int{10, 21, 50, 100}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := NewSize(size)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGenerateCustom(b *testing.B) {
	sizes := []int{10, 21, 50, 100}
	customAlphabet := "abcdef123456"
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := NewCustom(size, customAlphabet)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGenerateCustomUnicodeAlphabet(b *testing.B) {
	sizes := []int{10, 21, 50, 100}
	unicodeAlphabet := "あいうえお漢字🙂🚀"
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := NewCustom(size, unicodeAlphabet)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGenerateCustomReader(b *testing.B) {
	size := 21
	customAlphabet := "abcdef123456"
	randomData := make([]byte, 1024)
	for i := 0; i < len(randomData); i++ {
		randomData[i] = byte(i % 256)
	}
	rnd := bytes.NewReader(randomData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewCustomReader(size, customAlphabet, rnd)
		if err != nil {
			b.Fatal(err)
		}
		rnd.Seek(0, io.SeekStart) // Reset the reader for the next iteration
	}
}
