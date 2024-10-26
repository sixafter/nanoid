// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

// nanoid.go
package nanoid

import (
	"crypto/rand"
	"errors"
	"math/bits"
	"sync"
)

// Constants for default settings.
const (
	DefaultAlphabet = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	DefaultSize     = 21
	MaxUintSize     = 1024 // Adjust as needed
)

// Predefined errors to avoid allocations on each call.
var (
	ErrInvalidSize        = errors.New("size must be greater than zero")
	ErrSizeExceedsMaxUint = errors.New("size exceeds maximum allowed value")
	ErrEmptyAlphabet      = errors.New("alphabet must not be empty")
	ErrRandomSourceNoData = errors.New("random source returned no data")
)

// Byte pool to reuse byte slices and minimize allocations.
var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, MaxUintSize) // Non-zero length and capacity

		return &b
	},
}

// New generates a Nano ID with the default size and alphabet using crypto/rand as the random source.
func New() (string, error) {
	return NewSize(DefaultSize)
}

// NewSize generates a Nano ID with a specified size and the default alphabet using crypto/rand as the random source.
func NewSize(size int) (string, error) {
	return NewCustom(size, DefaultAlphabet)
}

// NewCustom generates a Nano ID with a specified size and custom ASCII alphabet using crypto/rand as the random source.
func NewCustom(size int, alphabet string) (string, error) {
	if size <= 0 {
		return "", ErrInvalidSize
	}
	if size > MaxUintSize {
		return "", ErrSizeExceedsMaxUint
	}
	if len(alphabet) == 0 {
		return "", ErrEmptyAlphabet
	}

	return generateASCIIID(size, alphabet)
}

// generateASCIIID generates an ID using a byte-based (ASCII) alphabet.
func generateASCIIID(size int, alphabet string) (string, error) {
	//nolint:gosec // G115: conversion from int to uint is safe due to prior bounds checking
	bitsPerChar := bits.Len(uint(len(alphabet) - 1))
	if bitsPerChar == 0 {
		bitsPerChar = 1
	}

	// Acquire a pointer to a byte slice from the pool
	bufPtr, ok := bytePool.Get().(*[]byte)
	if !ok {
		panic("bytePool.Get() did not return a *[]byte")
	}
	buf := *bufPtr
	buf = buf[:size] // Slice to desired size

	defer func() {
		// Reset the slice back to MaxUintSize before putting it back
		*bufPtr = (*bufPtr)[:MaxUintSize]
		bytePool.Put(bufPtr)
	}()

	var bitBuffer uint64
	var bitsInBuffer int

	for i := 0; i < size; {
		if bitsInBuffer < bitsPerChar {
			var b [8]byte
			n, err := rand.Read(b[:])
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", ErrRandomSourceNoData
			}
			for j := 0; j < n; j++ {
				bitBuffer |= uint64(b[j]) << bitsInBuffer
				bitsInBuffer += 8
			}
		}

		mask := uint64((1 << bitsPerChar) - 1)
		idx := bitBuffer & mask
		bitBuffer >>= bitsPerChar
		bitsInBuffer -= bitsPerChar

		//nolint:gosec // G115: conversion from int to uint is safe due to prior bounds checking
		if int(idx) < len(alphabet) {
			buf[i] = alphabet[idx]
			i++
		}
	}

	return string(buf), nil
}
