// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
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
	is.Positive(runtimeConfig.BufferSize(), "Config.BufferSize should be a positive integer")
	is.Positive(runtimeConfig.BitsNeeded(), "Config.BitsNeeded should be a positive integer")
	is.Positive(runtimeConfig.BytesNeeded(), "Config.BytesNeeded should be a positive integer")
	is.Equal(rand.Reader, runtimeConfig.RandReader(), "Config.RandReader should be rand.Reader by default")
	is.Equal(true, runtimeConfig.IsASCII(), "Config.IsASCII should be true by default")
	is.NotNil(runtimeConfig.RuneAlphabet(), "Config.RuneAlphabet should not be nil")
	is.NotNil(runtimeConfig.ByteAlphabet(), "Config.ByteAlphabet should not be nil")
	is.Positive(runtimeConfig.BufferMultiplier(), "Config.BufferMultiplier should be a positive integer")
	is.Positive(runtimeConfig.LengthHint(), "Config.LengthHint should be a positive integer")
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
	is.Equal(idLength, len([]rune(id1)), "Generated ID should have the specified length")
	is.True(isValidID(id1, customAlphabet), "Generated ID contains invalid characters")

	// Generate second ID
	id2, err := gen.New(idLength)
	is.NoError(err, "gen.New(%d) should not return an error", idLength)
	is.Equal(idLength, len([]rune(id2)), "Generated ID should have the specified length")
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
	is.Equal(ErrExceededMaxAttempts, err, "Expected ErrExceededMaxAttempts")
}

// alwaysFailRandReader is a custom io.Reader that always returns rnd=3.
type alwaysFailRandReader struct{}

// Read implements the io.Reader interface and always returns 3.
func (f *alwaysFailRandReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 3 // Assuming len(customAlphabet)=3, rnd=3 >=3
	}
	return len(p), nil
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
	is.Equal(io.EOF, err, "Expected io.ErrUnexpectedEOF from emptyRandReader")
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

	ids := make(chan string, numGoroutines*numIDsPerGoroutine)
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

	idSet := make(map[string]struct{}, numGoroutines*numIDsPerGoroutine)
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
