// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrDuplicateCharacters ensures that the generator returns ErrDuplicateCharacters
// when the provided alphabet contains duplicate characters.
func TestErrDuplicateCharacters(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Alphabet with duplicate characters
	alphabet := "abcabc"

	_, err := NewGenerator(WithAlphabet(alphabet))
	is.Equal(ErrDuplicateCharacters, err)
}

// TestErrInvalidLength verifies that the generator returns ErrInvalidLength
// when the LengthHint is set to zero, which is invalid.
func TestErrInvalidLength(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// LengthHint set to 0, which is invalid
	_, err := NewGenerator(WithLengthHint(0))
	is.Equal(ErrInvalidLength, err)
}

// TestErrInvalidAlphabet checks that the generator returns ErrInvalidAlphabet
// when an empty alphabet is provided.
func TestErrInvalidAlphabet(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Empty alphabet string, which is invalid
	alphabet := ""

	_, err := NewGenerator(WithAlphabet(alphabet))
	is.Equal(ErrInvalidAlphabet, err)
}

// TestErrNonUTF8Alphabet ensures that the generator returns ErrNonUTF8Alphabet
// when the alphabet contains invalid UTF-8 characters.
func TestErrNonUTF8Alphabet(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Invalid UTF-8 string (e.g., invalid byte sequences)
	// For example, byte 0xFF is invalid in UTF-8
	alphabet := string([]byte{0xFF, 0xFE, 0xFD})

	_, err := NewGenerator(WithAlphabet(alphabet))
	is.Equal(ErrNonUTF8Alphabet, err)
}

// TestErrAlphabetTooShort verifies that the generator returns ErrAlphabetTooShort
// when the alphabet has fewer than 2 characters.
func TestErrAlphabetTooShort(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Alphabet with only 1 character, which is too short
	alphabet := "A"

	_, err := NewGenerator(WithAlphabet(alphabet))
	is.Equal(ErrAlphabetTooShort, err)
}

// TestErrAlphabetTooLong checks that the generator returns ErrAlphabetTooLong
// when the alphabet exceeds 256 characters.
func TestErrAlphabetTooLong(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Create a unique alphabet with 257 unique characters
	alphabet := makeUnicodeAlphabet(257)

	// Verify the length. We expect 514 because each character is represented by 2 bytes.
	is.Equal(514, len(alphabet))

	_, err := NewGenerator(WithAlphabet(alphabet))
	is.Equal(ErrAlphabetTooLong, err)
}

// TestErrNilRandReader ensures that the generator returns ErrNilRandReader
// when the random reader (RandReader) is set to nil.
func TestErrNilRandReader(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// RandReader set to nil, which is invalid
	_, err := NewGenerator(WithRandReader(nil))
	is.Equal(ErrNilRandReader, err)
}

// alwaysInvalidRandReader is a mock implementation of io.Reader that always returns invalid indices.
// For an alphabet of "ABC" (length 3), it returns bytes with value 3, which are invalid indices.
type alwaysInvalidRandReader struct{}

// Read fills the provided byte slice with the invalid byte (3) and never returns an error.
func (a *alwaysInvalidRandReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 3 // Invalid index for alphabet "ABC"
	}
	return len(p), nil
}

// TestErrExceededMaxAttempts verifies that the generator returns ErrExceededMaxAttempts
// when it cannot produce a valid ID within the maximum number of attempts.
func TestErrExceededMaxAttempts(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	// Initialize the mockRandReader that always returns invalid indices.
	mockReader := &alwaysInvalidRandReader{}

	const length = 5
	generator, err := NewGenerator(
		WithAlphabet("ABC"),        // Non-power-of-two alphabet (length 3)
		WithRandReader(mockReader), // Mocked RandReader returning invalid indices (3)
		WithLengthHint(length),
	)
	is.NoError(err, "Expected no error when initializing generator with valid configuration")

	// Attempt to generate an ID; expect ErrExceededMaxAttempts
	_, err = generator.NewWithLength(length)
	is.Equal(ErrExceededMaxAttempts, err, "Expected ErrExceededMaxAttempts when generator cannot find valid indices")
}

// TestErrNilPointer_MarshalText ensures that MarshalText returns ErrNilPointer
// when called on a nil *ID.
func TestErrNilPointer_MarshalText(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	var id *ID = nil

	_, err := id.MarshalText()
	is.Equal(ErrNilPointer, err)
}

// TestErrNilPointer_UnmarshalText ensures that UnmarshalText returns ErrNilPointer
// when called on a nil *ID.
func TestErrNilPointer_UnmarshalText(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	var id *ID = nil

	err := id.UnmarshalText([]byte("test"))
	is.Equal(ErrNilPointer, err)
}

// TestErrNilPointer_MarshalBinary ensures that MarshalBinary returns ErrNilPointer
// when called on a nil *ID.
func TestErrNilPointer_MarshalBinary(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	var id *ID = nil

	_, err := id.MarshalBinary()
	is.Equal(ErrNilPointer, err)
}

// TestErrNilPointer_UnmarshalBinary ensures that UnmarshalBinary returns ErrNilPointer
// when called on a nil *ID.
func TestErrNilPointer_UnmarshalBinary(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	var id *ID = nil

	err := id.UnmarshalBinary([]byte("test"))
	is.Equal(ErrNilPointer, err)
}
