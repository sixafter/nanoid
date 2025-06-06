// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package prng provides a cryptographically secure pseudo-random number generator (PRNG)
// that implements the io.Reader interface. It is designed for high-performance, concurrent
// use in generating random bytes.
//
// This package is part of the experimental "x" modules and may be subject to change.

package nanoid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// FuzzNewWithLength fuzzes NewWithLength to validate length and character constraints.
func FuzzNewWithLength(f *testing.F) {
	f.Add(10) // seed
	f.Fuzz(func(t *testing.T, length int) {
		if length <= 0 || length > 256 {
			t.Skip() // skip invalid or absurd lengths
		}

		id, err := NewWithLength(length)
		is := assert.New(t)
		is.NoError(err)
		is.Equal(length, len([]rune(id)))
		is.True(isValidID(id, DefaultAlphabet))
	})
}

// FuzzCustomAlphabet fuzzes the NewGenerator with arbitrary alphabet strings.
func FuzzCustomAlphabet(f *testing.F) {
	f.Add("abc123", 16)
	f.Fuzz(func(t *testing.T, alphabet string, length int) {
		is := assert.New(t)
		if len([]rune(alphabet)) < MinAlphabetLength || len([]rune(alphabet)) > MaxAlphabetLength {
			t.Skip()
		}
		if length <= 0 || length > 256 {
			t.Skip()
		}

		gen, err := NewGenerator(
			WithAlphabet(alphabet),
			WithLengthHint(uint16(length)),
		)
		if err != nil {
			return // skip invalid generator configs
		}

		id, err := gen.NewWithLength(length)
		is.NoError(err)
		is.Equal(length, len([]rune(id)))
		is.True(isValidID(id, alphabet))
	})
}

// FuzzCustomGenerator fuzzes combinations of alphabet, length, and generation.
func FuzzCustomGenerator(f *testing.F) {
	f.Add("xyz", 32)
	f.Fuzz(func(t *testing.T, alphabet string, size int) {
		is := assert.New(t)
		runes := []rune(alphabet)

		if len(runes) < MinAlphabetLength || len(runes) > MaxAlphabetLength {
			t.Skip()
		}
		if size <= 0 || size > 256 {
			t.Skip()
		}

		gen, err := NewGenerator(
			WithAlphabet(alphabet),
			WithLengthHint(uint16(size)),
		)
		if err != nil {
			t.Skip()
		}

		id, err := gen.NewWithLength(size)
		is.NoError(err)
		is.Equal(size, len([]rune(id)))
		is.True(isValidID(id, alphabet))
	})
}

// FuzzRead fuzzes the global Generator.Read method with varying buffer sizes.
func FuzzRead(f *testing.F) {
	f.Add(21) // DefaultLength
	f.Fuzz(func(t *testing.T, size int) {
		is := assert.New(t)
		if size < 0 || size > 256 {
			t.Skip()
		}

		buf := make([]byte, size)
		n, err := Read(buf)

		is.NoError(err)
		is.Equal(size, n)
		is.True(isValidID(ID(buf), DefaultAlphabet))
	})
}
