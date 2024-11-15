// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"sync"
	"testing"

	"github.com/sixafter/nanoid/x/crypto/prng"
	"github.com/stretchr/testify/assert"
)

var (
	// Ensure ID implements the fmt.Stringer interface
	_ = fmt.Stringer(EmptyID)

	// Ensure ID implements the encoding.BinaryMarshaler interface
	_ = encoding.BinaryMarshaler(EmptyID)

	// Ensure ID implements the encoding.BinaryUnmarshaler interface
	_ = encoding.BinaryUnmarshaler(&EmptyID)

	// Ensure ID implements the encoding.TextMarshaler interface
	_ = encoding.TextMarshaler(EmptyID)

	// Ensure ID implements the encoding.TextUnmarshaler interface
	_ = encoding.TextUnmarshaler(&EmptyID)
)

// TestNewWithCustomLengths tests the generation of Nano IDs with custom lengths.
func TestNewWithCustomLengths(t *testing.T) {
	t.Parallel()
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

// TestNewAndMustWithCustomLengthPanic tests the must generation of a NanoID with an invalid length for panic.
func TestNewAndMustWithCustomLengthPanic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.Panics(func() {
		MustWithLength(0)
	})
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
	const idLength = 8
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	id, err := gen.New(idLength)
	is.NoError(err, "New(10) should not return an error")
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

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

// TestNewGeneratorWithInvalidAlphabet tests that the generator returns an error with invalid alphabets.
func TestNewGeneratorWithInvalidAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	lengths := []int{1, 2, 256, 257}

	mean := mean(lengths)

	// Define the alphabet types to test
	alphabetTypes := []string{"ASCII", "Unicode"}

	for _, alphabetType := range alphabetTypes {
		for _, length := range lengths {
			// New the appropriate alphabet
			var alphabet string
			if alphabetType == "ASCII" {
				alphabet = makeASCIIBasedAlphabet(length)
			} else {
				alphabet = makeUnicodeAlphabet(length)
			}
			gen, err := NewGenerator(
				WithAlphabet(alphabet),
				WithLengthHint(uint16(mean)),
			)

			alphabetRunes := []rune(alphabet)
			l := len(alphabetRunes)
			switch true {
			case l < MinAlphabetLength:
				is.Error(err, "NewGenerator() should return an error with an invalid alphabet length")
				is.Nil(gen, "Generator should be nil when initialization fails")
				is.Equal(ErrAlphabetTooShort, err, "Expected ErrAlphabetTooShort")
			case l > MaxAlphabetLength:
				is.Error(err, "NewGenerator() should return an error with an invalid alphabet length")
				is.Nil(gen, "Generator should be nil when initialization fails")
				is.Equal(ErrAlphabetTooLong, err, "Expected ErrAlphabetTooLong")
			default:
				is.NoError(err, "NewGenerator() should not return an error when initialization succeeds")
			}
		}
	}
}

// TestGetConfig tests the Config() method of the generator.
func TestInvalidUTF8Alphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a byte slice with an invalid UTF-8 sequence.
	// Here, 0x80 is a continuation byte, which by itself is not valid UTF-8.
	invalidUTF8 := []byte{0x80}

	// Convert the byte slice to a string.
	alphabet := string(invalidUTF8)

	gen, err := NewGenerator(
		WithAlphabet(alphabet),
	)

	is.Error(err, "NewGenerator() should return an error with an invalid alphabet")
	is.Nil(gen, "Generator should be nil when initialization fails")
	is.Equal(ErrNonUTF8Alphabet, err, "Expected ErrNonUTF8Alphabet")
}

// TestGetConfig tests the Config() method of the generator.
func TestGetConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewGenerator()
	is.NoError(err, "NewGenerator() should not return an error with the default alphabet")

	// Assert that generator implements Configuration interface
	config, ok := gen.(Configuration)
	is.True(ok, "Generator should implement Configuration interface")

	runtimeConfig := config.Config()

	is.Equal(DefaultAlphabet, string(runtimeConfig.RuneAlphabet()), "Config.RuneAlphabet should match the default alphabet")
	is.Equal(uint16(len([]rune(DefaultAlphabet))), runtimeConfig.AlphabetLen(), "Config.AlphabetLen should match the default alphabet length")

	// Update expectedMask calculation based on RuntimeConfig
	expectedMask := uint((1 << bits.Len(uint(runtimeConfig.AlphabetLen()-1))) - 1)
	is.Equal(expectedMask, runtimeConfig.Mask(), "Config.Mask should be correctly calculated")

	is.Equal((runtimeConfig.AlphabetLen()&(runtimeConfig.AlphabetLen()-1)) == 0, runtimeConfig.IsPowerOfTwo(), "Config.IsPowerOfTwo should be correct")
	is.Positive(runtimeConfig.BaseMultiplier(), "Config.BaseMultiplier should be a positive integer")
	is.Positive(runtimeConfig.BitsNeeded(), "Config.BitsNeeded should be a positive integer")
	is.Positive(runtimeConfig.BufferMultiplier(), "Config.BufferMultiplier should be a positive integer")
	is.Positive(runtimeConfig.BufferSize(), "Config.BufferSize should be a positive integer")
	is.NotNil(runtimeConfig.ByteAlphabet(), "Config.ByteAlphabet should not be nil")
	is.Positive(runtimeConfig.BytesNeeded(), "Config.BytesNeeded should be a positive integer")
	is.Equal(true, runtimeConfig.IsASCII(), "Config.IsASCII should be true by default")
	is.Positive(runtimeConfig.LengthHint(), "Config.LengthHint should be a positive integer")
	is.Equal(prng.Reader, runtimeConfig.RandReader(), "Config.RandReader should be rand.Reader by default")
	is.NotNil(runtimeConfig.RuneAlphabet(), "Config.RuneAlphabet should not be nil")
	is.Positive(runtimeConfig.ScalingFactor(), "Config.ScalingFactor should be a positive integer")
	is.Positive(runtimeConfig.MaxBytesPerRune(), "Config.MaxBytesPerRune should be a positive integer")
}

// TestUniqueness tests that multiple generated IDs are unique.
func TestUniqueness(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	numIDs := 1000
	ids := make(map[ID]struct{}, numIDs)

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

	ids := make(chan ID, numGoroutines*numIDsPerGoroutine)
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

	idSet := make(map[ID]struct{}, numGoroutines*numIDsPerGoroutine)
	for id := range ids {
		if _, exists := idSet[id]; exists {
			is.Failf("Duplicate ID found in concurrency test", "Duplicate ID: %s", id)
		}
		idSet[id] = struct{}{}
	}
}

// isValidID checks if all characters in the ID are within the specified alphabet.
func isValidID(id ID, alphabet string) bool {
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

func TestCyclicReader(t *testing.T) {
	expected := []byte{0, 1, 2, 3, 0, 1, 2, 3}
	reader := &cyclicReader{data: []byte{0, 1, 2, 3}}
	buffer := make([]byte, len(expected))
	n, err := reader.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, len(expected), n)
	assert.Equal(t, expected, buffer)
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
	is.Equal("ABCD", string(id), "Generated ID should match the expected sequence 'ABCD'")

	// New another ID of length 4, should cycle through customBytes again
	id, err = gen.New(4)
	is.NoError(err, "New(4) should not return an error on subsequent generation")
	is.Equal("ABCD", string(id), "Generated ID should match the expected sequence 'ABCD' on subsequent generation")

	// New ID of length 8, should cycle through customBytes twice
	id, err = gen.New(8)
	is.NoError(err, "New(8) should not return an error")
	is.Equal("ABCDABCD", string(id), "Generated ID should match the expected sequence 'ABCDABCD' for length 8")
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
	is.Equal("ZYXW", string(id), "Generated ID should match the expected sequence 'ZYXW'")

	// New another ID of length 4, should cycle through customBytes again
	id, err = gen.New(4)
	is.NoError(err, "New(4) should not return an error on subsequent generation")
	is.Equal("ZYXW", string(id), "Generated ID should match the expected sequence 'ZYXW' on subsequent generation")

	// New ID of length 8, should cycle through customBytes twice
	id, err = gen.New(8)
	is.NoError(err, "New(8) should not return an error")
	is.Equal("ZYXWZYXW", string(id), "Generated ID should match the expected sequence 'ZYXWZYXW' for length 8")
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
	is.Equal("FFFF", string(id), "Generated ID should match the expected sequence 'FFFF'")

	// New ID of length 6, expecting 'FFFFFF'
	id, err = gen.New(6)
	is.NoError(err, "New(6) should not return an error")
	is.Equal("FFFFFF", string(id), "Generated ID should match the expected sequence 'FFFFFF'")
}

// TestGenerateWithNonPowerOfTwoAlphabetLength tests ID generation with an alphabet length that is not a power of two.
func TestGenerateWithNonPowerOfTwoAlphabetLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Alphabet length is 10, which is not a power of two
	customAlphabet := "ABCDEFGHIJ" // Length = 10
	const idLength = 16
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid non-power-of-two alphabet length")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, customAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithMinimalAlphabet tests ID generation with the minimal valid alphabet size.
func TestGenerateWithMinimalAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Minimal valid alphabet length is 2
	minimalAlphabet := "01"
	const idLength = 32
	gen, err := NewGenerator(
		WithAlphabet(minimalAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with the minimal alphabet length")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, minimalAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithMaximalAlphabet tests the generation of IDs with a large alphabet size.
func TestGenerateWithMaximalAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Generate a maximal alphabet of 256 unique runes that form a valid UTF-8 string
	var maximalAlphabet string
	for i := 0; i < MaxAlphabetLength; i++ {
		// Ensure each rune is a valid UTF-8 character
		// Runes from 0x0000 to 0x00FF are valid and can be represented in UTF-8
		maximalAlphabet += string(rune(i))
	}
	const idLength = 128
	gen, err := NewGenerator(
		WithAlphabet(maximalAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a maximal alphabet length")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, maximalAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithCustomRandReaderReturningError tests generator behavior when the custom random reader returns an error.
func TestGenerateWithCustomRandReaderReturningError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a custom random reader that always returns an error
	failingReader := &failingRandReader{}
	const idLength = 8

	// Initialize the generator with a valid alphabet and the failing random reader
	customAlphabet := "ABCDEFGH"
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(failingReader),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	// Attempt to generate an ID
	id, err := gen.New(idLength)
	is.Error(err, "gen.New() should return an error when random reader fails")
	is.Empty(id, "Generated ID should be empty on error")
	is.Equal(io.ErrUnexpectedEOF, err, "Expected io.ErrUnexpectedEOF from failingRandReader")
}

// failingRandReader is a custom io.Reader that always returns an error.
type failingRandReader struct{}

// Read implements the io.Reader interface and always returns an error.
func (f *failingRandReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// TestGenerateWithNonASCIIAlphabet tests ID generation with a Unicode alphabet when isASCII is false.
func TestGenerateWithNonASCIIAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a Unicode alphabet with emojis and special characters
	unicodeAlphabet := "αβγδε😊🚀🌟"
	const idLength = 10
	gen, err := NewGenerator(
		WithAlphabet(unicodeAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid Unicode alphabet and isASCII=false")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, unicodeAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithSpecialCharactersInAlphabet tests ID generation with an alphabet containing special characters and emojis.
func TestGenerateWithSpecialCharactersInAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Alphabet with special characters and emojis
	specialAlphabet := "!@#$%^&*()_+😊🚀"
	const idLength = 12
	gen, err := NewGenerator(
		WithAlphabet(specialAlphabet),
	)
	is.NoError(err, "NewGenerator() should not return an error with a special characters alphabet")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, specialAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithVeryLargeLength tests ID generation with a very large length.
func TestGenerateWithVeryLargeLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a standard alphabet
	standardAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const idLength = 1000 // Very large length
	gen, err := NewGenerator(
		WithAlphabet(standardAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid alphabet")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, standardAlphabet), "Generated ID contains invalid characters")
}

// TestGeneratorBufferReuse tests that buffers are correctly reused from the pool without residual data.
func TestGeneratorBufferReuse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	customAlphabet := "XYZ123"
	const idLength = 6
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	// Generate first ID
	id1, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id1.String())), "Generated ID should have the specified length")
	is.True(isValidID(id1, customAlphabet), "Generated ID contains invalid characters")

	// Generate second ID
	id2, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id2.String())), "Generated ID should have the specified length")
	is.True(isValidID(id2, customAlphabet), "Generated ID contains invalid characters")

	// Ensure that IDs are different if possible
	if id1 == id2 {
		t.Errorf("Generated IDs should be different: id1=%s, id2=%s", id1, id2)
	}
}

// TestGenerateWithMaxAttempts tests that the generator returns ErrExceededMaxAttempts when it cannot generate enough valid characters.
func TestGenerateWithMaxAttempts(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a small alphabet
	customAlphabet := "ABC" // len=3, bitsNeeded=2, mask=3

	// Define a random reader that always returns rnd=3 (>= len(alphabet)=3)
	failReader := &alwaysFailRandReader{}

	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(failReader),
		WithLengthHint(10),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet and fail reader")

	// Attempt to generate an ID
	id, err := gen.New(10)
	is.Error(err, "gen.New(10) should return an error when random reader cannot provide valid characters")
	is.Empty(id, "Generated ID should be empty on error")
	is.Equal(io.EOF, err, "Expected io.EOF")
}

// TestGeneratorWithZeroLengthHint tests the generator's behavior with LengthHint set to 0.
func TestGeneratorWithZeroLengthHint(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	customAlphabet := "ABCDEFGHijklmnopQR"

	lengthHint := uint16(0)
	_, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(lengthHint),
	)
	is.Error(err, "NewGenerator() should return an error with LengthHint=0")
}

// TestNewWithZeroLengthHintAndMaxAlphabet tests the generator with LengthHint=0 and maximum alphabet size.
func TestNewWithZeroLengthHintAndMaxAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define the maximum valid alphabet size
	maxAlphabet := make([]rune, MaxAlphabetLength)
	for i := 0; i < MaxAlphabetLength; i++ {
		maxAlphabet[i] = rune(i)
	}
	lengthHint := uint16(0)

	gen, err := NewGenerator(
		WithAlphabet(string(maxAlphabet)),
		WithLengthHint(lengthHint),
	)
	is.Error(err, "NewGenerator() should return an error with LengthHint=0 and maximum alphabet size")
	is.Nil(gen, "Generator should be nil when LengthHint is zero")
}

// TestGenerateWithCustomRandReaderReturningNoBytes tests generator behavior when the custom reader returns no bytes.
func TestGenerateWithCustomRandReaderReturningNoBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a custom random reader that always returns zero bytes read
	emptyReader := &emptyRandReader{}
	const idLength = 8

	// Initialize the generator with a valid alphabet and the empty random reader
	customAlphabet := "ABCDEFGH"
	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(emptyReader),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	// Attempt to generate an ID
	id, err := gen.New(idLength)
	is.Error(err, "gen.New() should return an error when random reader provides no bytes")
	is.Empty(id, "Generated ID should be empty on error")
	is.Equal(io.EOF, err, "Expected io.EOF from emptyRandReader")
}

// emptyRandReader is a custom io.Reader that always returns zero bytes read.
type emptyRandReader struct{}

// Read implements the io.Reader interface and always returns 0 bytes read.
func (f *emptyRandReader) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

// TestGeneratorConcurrencyWithCustomAlphabetLength tests that the generator can handle concurrent ID generation with custom alphabet lengths.
func TestGeneratorConcurrencyWithCustomAlphabetLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	numGoroutines := 50
	numIDsPerGoroutine := 20
	customAlphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	idLength := 15

	gen, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(uint16(idLength)),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid custom alphabet")

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	ids := make(chan ID, numGoroutines*numIDsPerGoroutine)
	errorsChan := make(chan error, numGoroutines*numIDsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIDsPerGoroutine; j++ {
				id, err := gen.New(idLength)
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
		is.NoError(err, "gen.New() should not return an error in concurrent execution")
	}

	idSet := make(map[ID]struct{}, numGoroutines*numIDsPerGoroutine)
	for id := range ids {
		if _, exists := idSet[id]; exists {
			is.Failf("Duplicate ID found in concurrency test", "Duplicate ID: %s", id)
		}
		idSet[id] = struct{}{}
	}
}

// TestGenerateWithAllPrintableASCII tests the generation of IDs using all printable ASCII characters.
func TestGenerateWithAllPrintableASCII(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define an alphabet with all printable ASCII characters
	var asciiAlphabet string
	for i := 32; i <= 126; i++ {
		asciiAlphabet += string(rune(i))
	}
	const idLength = 20
	gen, err := NewGenerator(
		WithAlphabet(asciiAlphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with all printable ASCII characters")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, asciiAlphabet), "Generated ID contains invalid characters")
}

// TestGenerateWithSpecialUTF8Characters tests the generation of IDs with an alphabet containing special UTF-8 characters.
func TestGenerateWithSpecialUTF8Characters(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Alphabet with special UTF-8 characters
	specialUTF8Alphabet := "äöüß😊✨💖"
	const idLength = 15
	gen, err := NewGenerator(
		WithAlphabet(specialUTF8Alphabet),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a special UTF-8 characters alphabet")

	id, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id)), "Generated ID should have the specified length")

	is.True(isValidID(id, specialUTF8Alphabet), "Generated ID contains invalid characters")
}

// TestGeneratorWithInvalidLengthHint tests that the generator returns an error when LengthHint is invalid.
func TestGeneratorWithInvalidLengthHint(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	customAlphabet := "ABCDEFG"
	lengthHint := uint16(0)

	_, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithLengthHint(lengthHint),
	)
	is.Error(err, "NewGenerator() should return an error when LengthHint is zero")
	is.Equal(ErrInvalidLength, err, "Expected ErrInvalidLength for LengthHint=0")
}

// TestGenerateWithMaxAttemptsExceeded tests the generator's behavior when it exceeds the maximum number of attempts.
func TestGenerateWithMaxAttemptsExceeded(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define a small alphabet
	smallAlphabet := "AB"
	const idLength = 100
	failReader := &alwaysFailRandReader{}

	gen, err := NewGenerator(
		WithAlphabet(smallAlphabet),
		WithRandReader(failReader),
		WithLengthHint(idLength),
	)
	is.NoError(err, "NewGenerator() should not return an error with a valid small alphabet")

	id, err := gen.New(idLength)
	is.Error(err, "gen.New(%d) should return an error when random reader cannot provide valid characters", idLength)
	is.Empty(id, "Generated ID should be empty on error")
	is.Equal(io.EOF, err, "Expected io.EOF when maximum attempts are exceeded")
}

// alwaysFailRandReader is a custom io.Reader that always returns an error.
type alwaysFailRandReader struct{}

// Read implements the io.Reader interface and always returns an error.
func (f *alwaysFailRandReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 3 // Assuming len(customAlphabet)=2, rnd=3 >= 2
	}
	return len(p), io.EOF
}

// TestGenerateWithEmptyAlphabet tests the generator's behavior when an empty alphabet is provided.
func TestGenerateWithEmptyAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	customAlphabet := ""
	_, err := NewGenerator(
		WithAlphabet(customAlphabet),
	)
	is.Error(err, "NewGenerator() should return an error when alphabet is empty")
	is.Equal(ErrInvalidAlphabet, err, "Expected ErrInvalidAlphabet when the alphabet is empty")
}

// TestGenerateWithNilRandReader tests the generator's behavior when a nil random reader is provided.
func TestGenerateWithNilRandReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	customAlphabet := "ABCDEFG"
	_, err := NewGenerator(
		WithAlphabet(customAlphabet),
		WithRandReader(io.Reader(nil)),
	)
	is.Error(err, "NewGenerator() should return an error when RandReader is nil")
}

// TestProcessRandomBytes tests the processRandomBytes function to ensure coverage for different cases.
func TestProcessRandomBytes(t *testing.T) {
	is := assert.New(t)
	gen, err := NewGenerator()
	is.NoError(err, "NewGenerator() should not return an error with default configuration")

	randomBytes := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}

	// Case 1: bytesNeeded = 1
	gen.(*generator).config.bytesNeeded = 1
	result := gen.(*generator).processRandomBytes(randomBytes, 0)
	is.Equal(uint(0x12), result, "Expected result to be 0x12 for bytesNeeded=1")

	// Case 2: bytesNeeded = 2
	gen.(*generator).config.bytesNeeded = 2
	result = gen.(*generator).processRandomBytes(randomBytes, 0)
	is.Equal(uint(0x1234), result, "Expected result to be 0x1234 for bytesNeeded=2")

	// Case 3: bytesNeeded = 4
	gen.(*generator).config.bytesNeeded = 4
	result = gen.(*generator).processRandomBytes(randomBytes, 0)
	is.Equal(uint(0x12345678), result, "Expected result to be 0x12345678 for bytesNeeded=4")

	// Case 4: bytesNeeded > 4 (default case)
	gen.(*generator).config.bytesNeeded = 6
	result = gen.(*generator).processRandomBytes(randomBytes, 0)
	is.Equal(uint(0x123456789ABC), result, "Expected result to be 0x123456789ABC for bytesNeeded=6")
}

// TestGenerator_Read_EqualLength tests reading into a buffer equal to the default ID length.
func TestGenerator_Read_EqualLength(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	buffer := make([]byte, DefaultLength)
	n, err := Read(buffer)
	is.NoError(err, "Read should not return an error")
	is.Equal(DefaultLength, n, "Number of bytes read should equal DefaultLength")

	id := ID(buffer)
	is.Equal(DefaultLength, len(id), "Generated ID length should match DefaultLength")
	is.True(isValidID(id, DefaultAlphabet), "Generated ID should contain only valid characters")
}

// TestGenerator_Read_SmallerBuffer tests reading into a buffer smaller than the default ID length.
func TestGenerator_Read_SmallerBuffer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const length = 10
	gen, err := NewGenerator(WithLengthHint(length))
	is.NoError(err, "NewGenerator() should not return an error with length hint 10")

	buffer := make([]byte, length)
	n, err := gen.Read(buffer)
	is.NoError(err, "Read should not return an error")
	is.Equal(length, n, "Number of bytes read should equal bufferSize")

	id := ID(buffer)
	is.Equal(length, len(id), "Generated ID length should match bufferSize")
	is.True(isValidID(id, DefaultAlphabet), "Generated ID should contain only valid characters")
}

// TestGenerator_Read_LargerBuffer tests reading into a buffer larger than the default ID length.
func TestGenerator_Read_LargerBuffer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const length = 10
	gen, err := NewGenerator(WithLengthHint(length))
	is.NoError(err, "NewGenerator() should not return an error with length hint 10")

	buffer := make([]byte, length)
	n, err := gen.Read(buffer)
	is.NoError(err, "Read should not return an error")
	is.Equal(length, n, "Number of bytes read should equal bufferSize")

	id := ID(buffer)
	is.Equal(length, len(id), "Generated ID length should match bufferSize")
	is.True(isValidID(id, DefaultAlphabet), "Generated ID should contain only valid characters")
}

// TestGenerator_Read_ZeroBuffer tests reading into a zero-length buffer.
func TestGenerator_Read_ZeroBuffer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewGenerator()
	is.NoError(err, "NewGenerator() should not return an error with length hint 10")

	buffer := make([]byte, 0)
	n, err := gen.Read(buffer)
	is.NoError(err, "Read should not return an error")
	is.Equal(0, n, "Number of bytes read should be 0 for zero-length buffer")
}

// TestGenerator_Read_Concurrent tests concurrent reads to ensure thread safety.
func TestGenerator_Read_Concurrent(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, ok := DefaultGenerator.(*generator)
	is.True(ok, "DefaultGenerator should be of type *generator")

	numGoroutines := 10
	readsPerGoroutine := 100
	bufferSize := DefaultLength
	var wg sync.WaitGroup
	var mu sync.Mutex
	generatedIDs := make(map[ID]bool)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < readsPerGoroutine; j++ {
				buffer := make([]byte, bufferSize)
				n, err := gen.Read(buffer)
				if !is.NoError(err, "Read should not return an error") {
					return
				}
				is.Equal(bufferSize, n, "Number of bytes read should equal bufferSize")

				id := ID(buffer)
				is.Equal(bufferSize, len(id), "Generated ID length should match bufferSize")
				is.True(isValidID(id, DefaultAlphabet), "Generated ID should contain only valid characters")

				mu.Lock()
				generatedIDs[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Optionally, verify that a reasonable number of unique IDs were generated
	expectedMinUnique := numGoroutines * readsPerGoroutine / 2 // Arbitrary threshold
	is.GreaterOrEqual(len(generatedIDs), expectedMinUnique, "Expected at least %d unique IDs, but got %d", expectedMinUnique, len(generatedIDs))
}

// TestGenerator_Read_Error tests the Read method's behavior when the generator encounters an error.
func TestGenerator_Read_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a generator with a RandReader that always returns an error
	faultyReader := &errorReader{}
	gen, err := NewGenerator(WithRandReader(faultyReader))
	is.NoError(err, "NewGenerator should not return an error during creation")

	buffer := make([]byte, DefaultLength)
	n, err := gen.Read(buffer)
	is.Error(err, "Read should return an error when RandReader fails")
	is.Equal(0, n, "Number of bytes read should be 0 on error")
}

// errorReader is an io.Reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

// TestID_String tests the String() method of the ID type.
// It verifies that the String() method returns the underlying string value.
func TestID_String(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expected := expectedID.String()

	// Actual is obtained by calling String() on the ID
	actual := expectedID.String()

	// Assert that actual equals expected
	is.Equal(expected, actual, "ID.String() should return the underlying string")
}

// TestID_MarshalText tests the MarshalText() method of the ID type.
// It verifies that MarshalText() returns the correct byte slice representation of the ID.
func TestID_MarshalText(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expectedBytes := []byte(expectedID.String())

	// Actual is obtained by calling MarshalText()
	actualBytes, err := expectedID.MarshalText()

	// Assert no error occurred
	is.NoError(err, "MarshalText() should not return an error")

	// Assert that actual bytes match expected bytes
	is.Equal(expectedBytes, actualBytes, "MarshalText() should return the correct byte slice")
}

// TestID_UnmarshalText tests the UnmarshalText() method of the ID type.
// It verifies that UnmarshalText() correctly parses the byte slice and assigns the value to the ID.
func TestID_UnmarshalText(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	inputBytes := []byte(expectedID.String())

	// Initialize a zero-valued ID
	var actualID ID

	// Call UnmarshalText with inputBytes
	err := actualID.UnmarshalText(inputBytes)

	// Assert no error occurred
	is.NoError(err, "UnmarshalText() should not return an error")

	// Assert that actualID matches expectedID
	is.Equal(expectedID, actualID, "UnmarshalText() should correctly assign the input value to ID")
}

// TestID_MarshalBinary tests the MarshalBinary() method of the ID type.
// It verifies that MarshalBinary() returns the correct byte slice representation of the ID.
func TestID_MarshalBinary(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expectedBytes := []byte(expectedID.String())

	// Actual is obtained by calling MarshalBinary()
	actualBytes, err := expectedID.MarshalBinary()

	// Assert no error occurred
	is.NoError(err, "MarshalBinary() should not return an error")

	// Assert that actual bytes match expected bytes
	is.Equal(expectedBytes, actualBytes, "MarshalBinary() should return the correct byte slice")
}

// TestID_UnmarshalBinary tests the UnmarshalBinary() method of the ID type.
// It verifies that UnmarshalBinary() correctly parses the byte slice and assigns the value to the ID.
func TestID_UnmarshalBinary(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	inputBytes := []byte(expectedID.String())

	// Initialize a zero-valued ID
	var actualID ID

	// Call UnmarshalBinary with inputBytes
	err := actualID.UnmarshalBinary(inputBytes)

	// Assert no error occurred
	is.NoError(err, "UnmarshalBinary() should not return an error")

	// Assert that actualID matches expectedID
	is.Equal(expectedID, actualID, "UnmarshalBinary() should correctly assign the input value to ID")
}

// TestID_Compare tests the Compare() method of the ID type.
// It verifies that Compare() correctly compares two IDs and returns the expected result.
func TestID_Compare(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id1 := ID("FgEVN8QMTrnKGvBxFjtjw")
	id2 := ID("zTxG5Nl21ZAoM8Fabqk3H")

	// Case 1: id1 < id2
	is.Equal(-1, id1.Compare(id2), "id1 should be less than id2")

	// Case 2: id1 = id2
	is.Equal(0, id1.Compare(id1), "id1 should be equal to id1")

	// Case 3: id1 > id2
	is.Equal(1, id2.Compare(id1), "id2 should be greater than id1")
}

// TestID_IsEmpty tests the IsEmpty() method of the ID type.
// It verifies that IsEmpty() correctly returns true for an empty ID and false for a non-empty ID.
func TestID_IsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize two IDs using Must()
	id1 := Must()
	id2 := EmptyID

	// Case 1: id1 is not empty
	is.False(id1.IsEmpty(), "id1 should not be empty")

	// Case 2: id2 is empty
	is.True(id2.IsEmpty(), "id2 should be empty")
}
