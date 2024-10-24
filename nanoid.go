// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"crypto/rand"
	"errors"
	"math/bits"
	"strings"
)

const (
	defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultSize     = 21
	bitsPerByte     = 8 // Number of bits in a byte
)

// Generate generates a NanoID with the default size and alphabet.
func Generate() (string, error) {
	return GenerateSize(defaultSize)
}

// GenerateSize generates a NanoID with a specified size and the default alphabet.
func GenerateSize(size int) (string, error) {
	return GenerateCustom(size, defaultAlphabet)
}

// GenerateCustom generates a NanoID with a specified size and custom alphabet.
func GenerateCustom(size int, alphabet string) (string, error) {
	if size <= 0 {
		return "", errors.New("size must be greater than zero")
	}
	alphabetLen := len(alphabet)
	if alphabetLen == 0 {
		return "", errors.New("alphabet must not be empty")
	}

	alphabetBytes := []byte(alphabet)

	// Handle special case when alphabet length is 1
	if alphabetLen == 1 {
		return strings.Repeat(alphabet, size), nil
	}

	mask := (1 << bits.Len(uint(alphabetLen-1))) - 1
	bitsPerChar := bits.Len(uint(alphabetLen - 1))
	totalBits := size * bitsPerChar
	step := (totalBits + bitsPerByte - 1) / bitsPerByte // Number of random bytes needed

	id := make([]byte, size)
	bytes := make([]byte, step)

	for i := 0; i < size; {
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		for _, b := range bytes {
			idx := int(b) & mask
			if idx < alphabetLen {
				id[i] = alphabetBytes[idx]
				i++
				if i == size {
					break
				}
			}
		}
	}

	return string(id), nil
}
