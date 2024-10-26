// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.
package nanoid

import (
	"crypto/rand"
	"errors"
	"io"
	"math/bits"
	"strings"
)

const (
	DefaultAlphabet = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	DefaultSize     = 21
)

// New generates a NanoID with the default size and alphabet using crypto/rand as the random source.
func New() (string, error) {
	return NewSize(DefaultSize)
}

// NewSize generates a NanoID with a specified size and the default alphabet using crypto/rand as the random source.
func NewSize(size int) (string, error) {
	return NewCustom(size, DefaultAlphabet)
}

// NewCustom generates a NanoID with a specified size and custom alphabet using crypto/rand as the random source.
func NewCustom(size int, alphabet string) (string, error) {
	return NewCustomReader(size, alphabet, cryptoRandReader)
}

// NewCustomReader generates a NanoID with a specified size, custom alphabet, and custom random source.
func NewCustomReader(size int, alphabet string, rnd io.Reader) (string, error) {
	if rnd == nil {
		return "", errors.New("random source cannot be nil")
	}
	if size <= 0 {
		return "", errors.New("size must be greater than zero")
	}

	// Convert alphabet to []rune to support Unicode characters
	alphabetRunes := []rune(alphabet)
	alphabetLen := len(alphabetRunes)
	if alphabetLen == 0 {
		return "", errors.New("alphabet must not be empty")
	}

	// Handle special case when alphabet length is 1
	if alphabetLen == 1 {
		return strings.Repeat(string(alphabetRunes[0]), size), nil
	}

	// Calculate the number of bits needed to represent the alphabet indices
	bitsPerChar := bits.Len(uint(alphabetLen - 1))
	if bitsPerChar == 0 {
		bitsPerChar = 1
	}

	idRunes := make([]rune, size)
	var bitBuffer uint64
	var bitsInBuffer int
	i := 0

	for i < size {
		// If we don't have enough bits, read more random bytes
		if bitsInBuffer < bitsPerChar {
			var b [8]byte // Read up to 8 bytes at once for efficiency
			n, err := rnd.Read(b[:])
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", errors.New("random source returned no data")
			}
			// Append the new random bytes to the bit buffer
			for j := 0; j < n; j++ {
				bitBuffer |= uint64(b[j]) << bitsInBuffer
				bitsInBuffer += 8
			}
		}

		// Extract bitsPerChar bits to get the index
		idx := int(bitBuffer & ((1 << bitsPerChar) - 1))
		bitBuffer >>= bitsPerChar
		bitsInBuffer -= bitsPerChar

		// Use the index if it's within the alphabet range
		if idx < alphabetLen {
			idRunes[i] = alphabetRunes[idx]
			i++
		}
		// Else discard and continue
	}

	return string(idRunes), nil
}

// cryptoRandReader is a wrapper around crypto/rand.Reader to match io.Reader interface.
var cryptoRandReader io.Reader = cryptoRandReaderType{}

type cryptoRandReaderType struct{}

func (cryptoRandReaderType) Read(p []byte) (int, error) {
	return rand.Read(p)
}
