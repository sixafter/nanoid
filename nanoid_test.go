// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

// nanoid_test.go

package nanoid

import (
	"math/bits"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateDefault tests the generation of a default Nano ID.
func TestGenerateDefault(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id, err := Generate()
	is.NoError(err, "Generate() should not return an error")
	is.Equal(DefaultSize, len([]rune(id)), "Generated ID should have the default length")

	is.True(isValidID(id, DefaultAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateCustomLength tests the generation of Nano IDs with custom lengths.
func TestGenerateCustomLength(t *testing.T) {
	lengths := []int{1, 5, 10, 21, 50, 100}

	for _, length := range lengths {
		length := length // capture range variable
		t.Run("Length_"+strconv.Itoa(length), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			id, err := GenerateSize(length)
			is.NoError(err, "GenerateSize(%d) should not return an error", length)
			is.Equal(length, len([]rune(id)), "Generated ID should have the specified length")

			is.True(isValidID(id, DefaultAlphabet), "Generated ID contains invalid characters")
		})
	}
}

// TestGenerateInvalidLength tests the generator's response to invalid lengths.
func TestGenerateInvalidLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	invalidLengths := []int{0, -1, -10}

	for _, length := range invalidLengths {
		length := length // capture range variable
		t.Run("InvalidLength_"+strconv.Itoa(length), func(t *testing.T) {
			t.Parallel()
			id, err := GenerateSize(length)
			is.Error(err, "GenerateSize(%d) should return an error", length)
			is.Empty(id, "Generated ID should be empty on error")
			is.Equal(ErrInvalidLength, err, "Expected ErrInvalidLength")
		})
	}
}

// TestGenerateWithCustomAlphabet tests the generation of IDs with a custom alphabet.
func TestGenerateWithCustomAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Include Unicode characters in the custom alphabet
	customAlphabet := "abc😊🚀🌟"

	gen, err := New(customAlphabet, nil)
	is.NoError(err, "New() should not return an error with a valid custom alphabet")

	id, err := gen.Generate(10)
	is.NoError(err, "Generate(10) should not return an error")
	is.Equal(10, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, customAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithDuplicateAlphabet tests that the generator returns an error with duplicate characters.
func TestGenerateWithDuplicateAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	duplicateAlphabet := "aabbcc😊😊"
	_, err := New(duplicateAlphabet, nil)
	is.Error(err, "New() should return an error with duplicate characters in the alphabet")
	is.Equal(ErrDuplicateCharacters, err, "Expected ErrDuplicateCharacters")
}

// TestGetConfig tests the GetConfig() method of the generator.
func TestGetConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := New(DefaultAlphabet, nil)
	is.NoError(err, "New() should not return an error with the default alphabet")

	config := gen.(Configuration).GetConfig()

	is.Equal(DefaultAlphabet, string(config.Alphabet), "Config.Alphabet should match the default alphabet")
	is.Equal(uint16(len([]rune(DefaultAlphabet))), config.AlphabetLen, "Config.AlphabetLen should match the default alphabet length")

	// Update expectedMask calculation for uint32
	expectedMask := uint((1 << bits.Len(uint(config.AlphabetLen-1))) - 1)
	is.Equal(expectedMask, config.Mask, "Config.Mask should be correctly calculated")

	is.Equal((config.AlphabetLen&(config.AlphabetLen-1)) == 0, config.IsPowerOfTwo, "Config.IsPowerOfTwo should be correct")

	is.Positive(config.BitsNeeded, "Config.BitsNeeded should be a positive integer")
	is.Positive(config.BytesNeeded, "Config.BytesNeeded should be a positive integer")
}

// TestUniqueness tests that multiple generated IDs are unique.
func TestUniqueness(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	numIDs := 1000
	ids := make(map[string]struct{}, numIDs)

	for i := 0; i < numIDs; i++ {
		id, err := Generate()
		is.NoError(err, "Generate() should not return an error")
		is.NotContains(ids, id, "Duplicate ID found: %s", id)
		ids[id] = struct{}{}
	}
}

// TestConcurrency tests that the generator is safe for concurrent use.
func TestConcurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	numGoroutines := 100
	numIDsPerGoroutine := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	ids := make(chan string, numGoroutines*numIDsPerGoroutine)
	errorsChan := make(chan error, numGoroutines*numIDsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIDsPerGoroutine; j++ {
				id, err := Generate()
				if err != nil {
					errorsChan <- err
					continue
				}
				ids <- id
			}
		}()
	}

	wg.Wait()
	close(ids)
	close(errorsChan)

	for err := range errorsChan {
		is.NoError(err, "Generate() should not return an error in concurrent execution")
	}

	idSet := make(map[string]struct{}, numGoroutines*numIDsPerGoroutine)
	for id := range ids {
		if _, exists := idSet[id]; exists {
			is.Failf("Duplicate ID found in concurrency test", "Duplicate ID: %s", id)
		}
		idSet[id] = struct{}{}
	}
}

// TestInvalidAlphabetLength tests that alphabets with invalid lengths are rejected.
func TestInvalidAlphabetLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Alphabet length less than 2
	shortAlphabet := "a"
	_, err := New(shortAlphabet, nil)
	is.Error(err, "New() should return an error for alphabets shorter than 2 characters")
	is.Equal(ErrInvalidAlphabet, err, "Expected ErrInvalidAlphabet")
}

// isValidID checks if all characters in the ID are within the specified alphabet.
func isValidID(id string, alphabet string) bool {
	alphabetSet := make(map[rune]struct{}, len([]rune(alphabet)))
	for _, char := range alphabet {
		alphabetSet[char] = struct{}{}
	}

	for _, char := range id {
		if _, exists := alphabetSet[char]; !exists {
			return false
		}
	}
	return true
}
