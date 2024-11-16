// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"errors"
)

var (
	// ErrDuplicateCharacters is returned when the provided alphabet contains duplicate characters.
	ErrDuplicateCharacters = errors.New("duplicate characters in alphabet")

	// ErrExceededMaxAttempts is returned when the maximum number of attempts to perform
	// an operation, such as generating a unique ID, has been exceeded.
	ErrExceededMaxAttempts = errors.New("exceeded maximum attempts")

	// ErrInvalidLength is returned when a specified length value for an operation is invalid.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidAlphabet is returned when the provided alphabet for generating IDs is invalid.
	ErrInvalidAlphabet = errors.New("invalid alphabet")

	// ErrNonUTF8Alphabet is returned when the provided alphabet contains non-UTF-8 characters.
	ErrNonUTF8Alphabet = errors.New("alphabet contains invalid UTF-8 characters")

	// ErrAlphabetTooShort is returned when the provided alphabet has fewer than 2 characters.
	ErrAlphabetTooShort = errors.New("alphabet length is less than 2")

	// ErrAlphabetTooLong is returned when the provided alphabet exceeds 256 characters.
	ErrAlphabetTooLong = errors.New("alphabet length exceeds 256")

	// ErrNilRandReader is returned when the random number generator (rand.Reader) is nil,
	// preventing the generation of random values.
	ErrNilRandReader = errors.New("nil random reader")

	// ErrNilPointer is returned when a nil pointer is passed to a function that does not accept nil pointers.
	ErrNilPointer = errors.New("nil pointer")
)
