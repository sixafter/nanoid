// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

// nanoid_test.go

// nanoid_test.go

// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"crypto/rand"
	"io"
	"math/bits"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewWithCustomLengths tests the generation of Nano IDs with custom lengths.
func TestNewWithCustomLengths(t *testing.T) {
	lengths := []int{1, 5, 10, 21, 50, 100}

	for _, length := range lengths {
		length := length // capture range variable
		t.Run("Length_"+strconv.Itoa(length), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			id, err := NewWithLength(length)
			is.NoError(err, "NewWithLength(%d) should not return an error", length)
			is.Equal(length, len([]rune(id)), "Generated ID should have the specified length")

			is.True(isValidID(id, DefaultAlphabet), "Generated ID contains invalid characters")
		})
	}
}

// TestNewAndMustDefault tests the must generation of a default Nano ID.
func TestNewAndMustDefault(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id := Must()
	is.Equal(DefaultLength, len([]rune(id)), "Generated ID should have the default length")

	id = MustWithLength(DefaultLength)
	is.Equal(DefaultLength, len([]rune(id)), "Generated ID should have the default length")

	is.True(isValidID(id, DefaultAlphabet), "Generated ID contains invalid characters")
}

// TestNewInvalidLength tests the generator's response to invalid lengths.
func TestNewInvalidLength(t *testing.T) {
	t.Parallel()

	invalidLengths := []int{0, -1, -10}

	for _, length := range invalidLengths {
		length := length // capture range variable
		t.Run("InvalidLength_"+strconv.Itoa(length), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			id, err := NewWithLength(length)
			is.Error(err, "NewWithLength(%d) should return an error", length)
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

	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	id, err := gen.New(10)
	is.NoError(err, "New(10) should not return an error")
	is.Equal(10, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, customAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithDuplicateAlphabet tests that the generator returns an error with duplicate characters.
func TestGenerateWithDuplicateAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	duplicateAlphabet := "aabbcc😊😊"
	gen, err := NewGenerator(
		WithAlphabet(duplicateAlphabet),
	)
	is.Error(err, "NewGenerator() should return an error with duplicate characters in the alphabet")
	is.Nil(gen, "Generator should be nil when initialization fails")
	is.Equal(ErrDuplicateCharacters, err, "Expected ErrDuplicateCharacters")
}

// TestGetConfig tests the GetConfig() method of the generator.
func TestGetConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewGenerator(
		WithAlphabet(DefaultAlphabet),
	)
	is.NoError(err, "NewGenerator() should not return an error with the default alphabet")

	// Assert that generator implements Configuration interface
	config, ok := gen.(Configuration)
	is.True(ok, "Generator should implement Configuration interface")

	runtimeConfig := config.GetConfig()

	is.Equal(DefaultAlphabet, string(runtimeConfig.RuneAlphabet), "Config.RuneAlphabet should match the default alphabet")
	is.Equal(uint16(len([]rune(DefaultAlphabet))), runtimeConfig.AlphabetLen, "Config.AlphabetLen should match the default alphabet length")

	// Update expectedMask calculation based on RuntimeConfig
	expectedMask := uint((1 << bits.Len(uint(runtimeConfig.AlphabetLen-1))) - 1)
	is.Equal(expectedMask, runtimeConfig.Mask, "Config.Mask should be correctly calculated")

	is.Equal((runtimeConfig.AlphabetLen&(runtimeConfig.AlphabetLen-1)) == 0, runtimeConfig.IsPowerOfTwo, "Config.IsPowerOfTwo should be correct")

	is.Positive(runtimeConfig.BitsNeeded, "Config.BitsNeeded should be a positive integer")
	is.Positive(runtimeConfig.BytesNeeded, "Config.BytesNeeded should be a positive integer")
	is.Equal(rand.Reader, runtimeConfig.RandReader, "Config.RandReader should be rand.Reader by default")
}

// TestUniqueness tests that multiple generated IDs are unique.
func TestUniqueness(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	numIDs := 1000
	ids := make(map[string]struct{}, numIDs)

	for i := 0; i < numIDs; i++ {
		id, err := New()
		is.NoError(err, "New() should not return an error")
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
				id, err := New()
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
		is.NoError(err, "New() should not return an error in concurrent execution")
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
	gen, err := NewGenerator(
		WithAlphabet(shortAlphabet),
	)
	is.Error(err, "NewGenerator() should return an error for alphabets shorter than 2 characters")
	is.Nil(gen, "Generator should be nil when initialization fails")
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

// cyclicReader is a helper type that cycles through a predefined set of bytes.
// It implements the io.Reader interface.
type cyclicReader struct {
	data []byte
	mu   sync.Mutex
	pos  int
}

// Read fills p with bytes from the cyclicReader's data, cycling back to the start when necessary.
func (r *cyclicReader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.data) == 0 {
		return 0, io.EOF
	}

	n := 0
	for n < len(p) {
		p[n] = r.data[r.pos]
		n++
		r.pos = (r.pos + 1) % len(r.data)
	}

	return n, nil
}

// TestWithRandReader tests the WithRandReader option to ensure that the generator uses the provided random source.
func TestWithRandReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a custom alphabet
	customAlphabet := "ABCD"

	// Define a custom random source with known bytes
	// For example, bytes [0,1,2,3] should map to 'A','B','C','D'
	customBytes := []byte{0, 1, 2, 3}
	customReader := &cyclicReader{data: customBytes}

	// Initialize the generator with custom alphabet and custom random reader
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(customReader),
	)
	is.NoError(err, "NewGenerator() should not return an error with valid custom alphabet and random reader")

	// New ID of length 4
	id, err := gen.New(4)
	is.NoError(err, "New(4) should not return an error")
	is.Equal("ABCD", id, "Generated ID should match the expected sequence 'ABCD'")

	// New another ID of length 4, should cycle through customBytes again
	id, err = gen.New(4)
	is.NoError(err, "New(4) should not return an error on subsequent generation")
	is.Equal("ABCD", id, "Generated ID should match the expected sequence 'ABCD' on subsequent generation")

	// New ID of length 8, should cycle through customBytes twice
	id, err = gen.New(8)
	is.NoError(err, "New(8) should not return an error")
	is.Equal("ABCDABCD", id, "Generated ID should match the expected sequence 'ABCDABCD' for length 8")
}

// TestWithRandReaderDifferentSequence tests the WithRandReader option with a different byte sequence and alphabet.
func TestWithRandReaderDifferentSequence(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a different custom alphabet
	customAlphabet := "WXYZ"

	// Define a different custom random source with known bytes
	// For example, bytes [3,2,1,0] should map to 'Z','Y','X','W'
	customBytes := []byte{3, 2, 1, 0}
	customReader := &cyclicReader{data: customBytes}

	// Initialize the generator with custom alphabet and custom random reader
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(customReader),
	)
	is.NoError(err, "NewGenerator() should not return an error with valid custom alphabet and random reader")

	// New ID of length 4
	id, err := gen.New(4)
	is.NoError(err, "New(4) should not return an error")
	is.Equal("ZYXW", id, "Generated ID should match the expected sequence 'ZYXW'")

	// New another ID of length 4, should cycle through customBytes again
	id, err = gen.New(4)
	is.NoError(err, "New(4) should not return an error on subsequent generation")
	is.Equal("ZYXW", id, "Generated ID should match the expected sequence 'ZYXW' on subsequent generation")

	// New ID of length 8, should cycle through customBytes twice
	id, err = gen.New(8)
	is.NoError(err, "New(8) should not return an error")
	is.Equal("ZYXWZYXW", id, "Generated ID should match the expected sequence 'ZYXWZYXW' for length 8")
}

// TestWithRandReaderInsufficientBytes tests the generator's behavior when the custom reader provides insufficient bytes.
// Since cyclicReader cycles through the data, it should still work correctly.
func TestWithRandReaderInsufficientBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a custom alphabet
	customAlphabet := "EFGH"

	// Define a custom random source with a single byte
	customBytes := []byte{1} // Should map to 'F' repeatedly
	customReader := &cyclicReader{data: customBytes}

	// Initialize the generator with custom alphabet and custom random reader
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(customReader),
	)
	is.NoError(err, "NewGenerator() should not return an error with valid custom alphabet and random reader")

	// New ID of length 4, expecting 'FFFF'
	id, err := gen.New(4)
	is.NoError(err, "New(4) should not return an error")
	is.Equal("FFFF", id, "Generated ID should match the expected sequence 'FFFF'")

	// New ID of length 6, expecting 'FFFFFF'
	id, err = gen.New(6)
	is.NoError(err, "New(6) should not return an error")
	is.Equal("FFFFFF", id, "Generated ID should match the expected sequence 'FFFFFF'")
}
