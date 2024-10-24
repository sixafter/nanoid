// Copyright (c) 2024 Six After, Inc.
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// TestGenerate tests the default Generate function.
func TestGenerate(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id, err := Generate()
	is.NoError(err)
	is.Equal(defaultSize, len(id))
	is.True(isValidID(id, defaultAlphabet), "ID contains invalid characters: %s", id)
}

// TestGenerateSize tests generating IDs with custom sizes.
func TestGenerateSize(t *testing.T) {
	t.Parallel()
	sizes := []int{1, 5, 10, 100, 1000}

	for _, size := range sizes {
		size := size // capture range variable
		t.Run("Size"+strconv.Itoa(size), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			id, err := GenerateSize(size)
			is.NoError(err)
			is.Equal(size, len(id))
			is.True(isValidID(id, defaultAlphabet), "ID contains invalid characters: %s", id)
		})
	}
}

// TestGenerateCustom tests generating IDs with custom sizes and alphabets.
func TestGenerateCustom(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size     int
		alphabet string
	}{
		{16, "abcdef0123456789"},
		{32, "!@#$%^&*()_+-="},
		{8, "ABCD"},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run("Size"+strconv.Itoa(test.size)+"Alphabet"+strconv.Itoa(len(test.alphabet)), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			id, err := GenerateCustom(test.size, test.alphabet)
			is.NoError(err)
			is.Equal(test.size, len(id))
			is.True(isValidID(id, test.alphabet), "ID contains invalid characters: %s", id)
		})
	}
}

// TestGenerateEdgeCases tests edge cases for invalid inputs.
func TestGenerateEdgeCases(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	_, err := GenerateCustom(0, defaultAlphabet)
	is.Error(err)

	_, err = GenerateCustom(10, "")
	is.Error(err)

	_, err = GenerateCustom(-1, defaultAlphabet)
	is.Error(err)
}

// TestGenerateCustom_SingleCharacterAlphabet tests the special case with a single-character alphabet.
func TestGenerateCustom_SingleCharacterAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 5
	alphabet := "A"
	expected := "AAAAA"

	id, err := GenerateCustom(size, alphabet)
	is.NoError(err)
	is.Equal(expected, id)
}

// Helper function to check if an ID contains only characters from the alphabet.
func isValidID(id, alphabet string) bool {
	alphabetMap := make(map[rune]struct{}, len(alphabet))
	for _, r := range alphabet {
		alphabetMap[r] = struct{}{}
	}
	for _, r := range id {
		if _, exists := alphabetMap[r]; !exists {
			return false
		}
	}
	return true
}

func TestGenerateCustom_SizeZero(t *testing.T) {
	_, err := GenerateCustom(0, "abcdef")
	if err == nil {
		t.Fatalf("Expected error for size 0, got nil")
	}
}

func TestGenerateCustom_NegativeSize(t *testing.T) {
	_, err := GenerateCustom(-5, "abcdef")
	if err == nil {
		t.Fatalf("Expected error for negative size, got nil")
	}
}
