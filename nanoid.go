// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"sync"
	"unicode"
	"unicode/utf8"
)

// DefaultGenerator is a global, shared instance of a Nano ID generator. It is safe for concurrent use.
var DefaultGenerator Generator

// Generate generates a Nano ID using the default generator and the default size.
func Generate() (string, error) {
	return DefaultGenerator.Generate(DefaultSize)
}

// MustGenerate generates a Nano ID using the default generator and the default size if err
// is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustGenerate() string {
	id, err := DefaultGenerator.Generate(DefaultSize)
	if err != nil {
		panic(err)
	}

	return id
}

// GenerateSize generates a Nano ID using the default generator with a specified size.
func GenerateSize(length int) (string, error) {
	return DefaultGenerator.Generate(length)
}

// MustGenerateSize generates a Nano ID using the default generator with a specified size if err
// is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustGenerateSize(length int) string {
	id, err := DefaultGenerator.Generate(length)
	if err != nil {
		panic(err)
	}

	return id
}

func init() {
	var err error
	DefaultGenerator, err = New(DefaultAlphabet, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DefaultGenerator: %v", err))
	}
}

var (
	ErrInvalidLength       = errors.New("invalid length")
	ErrInvalidAlphabet     = errors.New("invalid alphabet")
	ErrDuplicateCharacters = errors.New("duplicate characters in alphabet")
	ErrNonUTF8Alphabet     = errors.New("alphabet contains invalid UTF-8 characters")
	ErrExceededMaxAttempts = errors.New("exceeded maximum attempts")
)

// DefaultAlphabet as per Nano ID specification.
const DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

// DefaultSize is the default size of the generated Nano ID: 21.
const DefaultSize = 21

// maxAttemptsMultiplier defines the multiplier for maximum attempts based on length.
const maxAttemptsMultiplier = 10

// bufferMultiplier defines how many characters the buffer should handle per read.
// Adjust this value based on performance and memory considerations.
const bufferMultiplier = 128

// Generator defines the interface for generating Nano IDs.
type Generator interface {
	// Generate creates a new Nano ID of the specified length.
	Generate(size int) (string, error)

	// MustGenerate returns creates a new Nano ID of the specified length if err
	// is nil or panics otherwise.
	// It simplifies safe initialization of global variables holding compiled UUIDs.
	MustGenerate(length int) string
}

// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
	// GetConfig returns the configuration of the generator.
	GetConfig() Config
}

// Config holds the configuration for the Nano ID generator.
type Config struct {
	// Alphabet is a slice of bytes representing the character set used to generate IDs.
	Alphabet []byte // 24 bytes (slice header)

	// RuneAlphabet is a slice of runes, allowing support for multi-byte characters in ID generation.
	RuneAlphabet []rune // 24 bytes (slice header)

	// Mask is a bitmask used to obtain a random value from the character set.
	Mask uint // 8 bytes

	// BitsNeeded represents the number of bits required to generate each character in the ID.
	BitsNeeded uint // 8 bytes

	// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
	BytesNeeded uint // 8 bytes

	// BufferSize is the buffer size used for random byte generation.
	BufferSize int // 8 bytes

	// AlphabetLen is the length of the alphabet, stored as a uint16.
	AlphabetLen uint16 // 2 bytes

	// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
	IsPowerOfTwo bool // 1 byte

	// IsASCII indicates whether the alphabet is ASCII-only, ensuring compatibility with ASCII environments.
	IsASCII bool // 1 byte
}

// generator implements the Generator interface.
type generator struct {
	config         *Config    // 8 bytes (pointer)
	randReader     io.Reader  // 16 bytes (interface)
	byteBufferPool *sync.Pool // 8 bytes (pointer)
	runeBufferPool *sync.Pool // 8 bytes (pointer)
}

// New creates a new Generator with buffer pooling enabled.
// It returns an error if the alphabet is invalid or contains invalid UTF-8 characters.
func New(alphabet string, randReader io.Reader) (Generator, error) {
	return newGenerator(alphabet, randReader)
}

// isASCII checks if all characters in a string are ASCII.
func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// isUnicode checks if any character in a string is a non-ASCII Unicode character.
func isUnicode(s string) bool {
	return !isASCII(s)
}

// newGenerator is an internal constructor for generator.
func newGenerator(alphabet string, randReader io.Reader) (Generator, error) {
	if len(alphabet) == 0 {
		return nil, ErrInvalidAlphabet
	}

	if randReader == nil {
		randReader = rand.Reader
	}

	// Check if the alphabet is valid UTF-8
	if !utf8.ValidString(alphabet) {
		return nil, ErrNonUTF8Alphabet
	}

	// Determine if the alphabet is ASCII-only
	isASCII := !isUnicode(alphabet)

	var (
		alphabetBytes []byte
		alphabetRunes []rune
	)

	if isASCII {
		alphabetBytes = []byte(alphabet)
	} else {
		alphabetRunes = []rune(alphabet)
	}

	// Check for duplicate characters
	if isASCII {
		seen := make(map[byte]bool)
		for _, b := range alphabetBytes {
			if seen[b] {
				return nil, ErrDuplicateCharacters
			}
			seen[b] = true
		}
	} else {
		seenRunes := make(map[rune]bool)
		for _, r := range alphabetRunes {
			if seenRunes[r] {
				return nil, ErrDuplicateCharacters
			}
			seenRunes[r] = true
		}
	}

	// Calculate BitsNeeded and Mask
	alphabetLen := 0
	if isASCII {
		alphabetLen = len(alphabetBytes)
	} else {
		alphabetLen = len(alphabetRunes)
	}

	if alphabetLen < 2 {
		return nil, ErrInvalidAlphabet
	}

	bitsNeeded := uint(bits.Len(uint(alphabetLen - 1)))
	if bitsNeeded == 0 {
		return nil, ErrInvalidAlphabet
	}
	mask := uint((1 << bitsNeeded) - 1)
	bytesNeeded := (bitsNeeded + 7) / 8

	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	// Calculate bufferSize dynamically based on bytesNeeded and bufferMultiplier
	bufferSize := int(bytesNeeded) * bufferMultiplier

	config := &Config{
		Alphabet:     alphabetBytes,
		RuneAlphabet: alphabetRunes,
		Mask:         mask,
		BitsNeeded:   bitsNeeded,
		BytesNeeded:  bytesNeeded,
		BufferSize:   bufferSize,
		AlphabetLen:  uint16(alphabetLen),
		IsPowerOfTwo: isPowerOfTwo,
		IsASCII:      isASCII,
	}

	// Initialize buffer pools
	var byteBufferPool *sync.Pool
	var runeBufferPool *sync.Pool

	if isASCII {
		byteBufferPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, bufferSize)
				return &buf
			},
		}
	} else {
		runeBufferPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, bufferSize)
				return &buf
			},
		}
	}

	return &generator{
		config:         config,
		randReader:     randReader,
		byteBufferPool: byteBufferPool,
		runeBufferPool: runeBufferPool,
	}, nil
}

// Generate creates a new Nano ID of the specified length.
func (g *generator) Generate(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	// Choose between byte or rune slices for id
	var id interface{}
	if g.config.IsASCII {
		id = make([]byte, length)
	} else {
		id = make([]rune, length)
	}

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	attempts := 0
	mask := g.config.Mask
	bytesNeeded := g.config.BytesNeeded

	// Use appropriate buffer pool
	var randomBytesPtr *[]byte
	if g.config.IsASCII {
		randomBytesPtr = g.byteBufferPool.Get().(*[]byte)
		defer g.byteBufferPool.Put(randomBytesPtr)
	} else {
		randomBytesPtr = g.runeBufferPool.Get().(*[]byte)
		defer g.runeBufferPool.Put(randomBytesPtr)
	}
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)
	step := int(bytesNeeded)
	if step <= 0 {
		return "", ErrInvalidAlphabet
	}

	// Generate ID
	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		neededBytes := (length - cursor) * step
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		_, err := io.ReadFull(g.randReader, randomBytes[:neededBytes])
		if err != nil {
			return "", err
		}

		for i := 0; i < neededBytes; i += step {
			var rnd uint
			for j := 0; j < step; j++ {
				rnd = (rnd << 8) | uint(randomBytes[i+j])
			}
			rnd &= mask

			if g.config.IsPowerOfTwo || int(rnd) < int(g.config.AlphabetLen) {
				if g.config.IsASCII {
					id.([]byte)[cursor] = g.config.Alphabet[rnd]
				} else {
					id.([]rune)[cursor] = g.config.RuneAlphabet[rnd]
				}
				cursor++
			}

			if cursor == length {
				break
			}
		}
	}

	// Convert id to string based on its type
	if g.config.IsASCII {
		return string(id.([]byte)), nil
	}
	return string(id.([]rune)), nil
}

// MustGenerate returns creates a new Nano ID of the specified length if err
// is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func (g *generator) MustGenerate(length int) string {
	id, err := g.Generate(length)
	if err != nil {
		panic(err)
	}
	return id
}

// GetConfig returns the configuration for the generator.
// It implements the Configuration interface.
func (g *generator) GetConfig() Config {
	return *g.config
}
