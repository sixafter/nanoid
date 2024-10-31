// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.
// nanoid.go

package nanoid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"strings"
	"sync"
)

// DefaultGenerator is a global, shared instance of a Nano ID generator. It is safe for concurrent use.
var DefaultGenerator Generator

// Generate generates a Nano ID using the default generator and the default size.
func Generate() (string, error) {
	return DefaultGenerator.Generate(DefaultSize)
}

// GenerateSize generates a Nano ID using the default generator.
func GenerateSize(length int) (string, error) {
	return DefaultGenerator.Generate(length)
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
	ErrExceededMaxAttempts = errors.New("exceeded maximum attempts")
)

// DefaultAlphabet can now include multibyte Unicode characters.
const DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

// DefaultSize is the default size of the generated Nano ID: 21.
const DefaultSize = 21

// maxAttemptsMultiplier defines the multiplier for maximum attempts based on length.
const maxAttemptsMultiplier = 10

// Generator defines the interface for generating Nano IDs.
type Generator interface {
	Generate(size int) (string, error)
}

// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
	GetConfig() Config
}

// Config holds the configuration for the Nano ID generator.
type Config struct {
	Alphabet     []rune // Changed from []byte to []rune
	AlphabetLen  uint32 // Updated to uint32 to handle larger alphabets
	Step         uint32 // Updated to uint32
	Mask         uint32 // Updated to uint32
	IsPowerOfTwo bool
}

type generator struct {
	randReader io.Reader
	bufferPool *sync.Pool
	config     Config
}

// New creates a new Generator with buffer pooling enabled.
// It returns an error if the alphabet is invalid.
func New(alphabet string, randReader io.Reader) (Generator, error) {
	return newGenerator(alphabet, randReader)
}

// newGenerator is an internal constructor for generator.
func newGenerator(alphabet string, randReader io.Reader) (Generator, error) {
	if len(alphabet) == 0 {
		return nil, ErrInvalidAlphabet
	}

	if randReader == nil {
		randReader = rand.Reader
	}

	alphabetRunes := []rune(alphabet)
	alphabetLen := len(alphabetRunes)

	if alphabetLen < 2 {
		return nil, ErrInvalidAlphabet
	}

	// Preallocate map with capacity
	seen := make(map[rune]bool, alphabetLen)
	for _, r := range alphabetRunes {
		if seen[r] {
			return nil, ErrDuplicateCharacters
		}
		seen[r] = true
	}

	k := bits.Len(uint(alphabetLen - 1))
	if k == 0 {
		return nil, ErrInvalidAlphabet
	}
	mask := uint32((1 << k) - 1)

	onesCount := bits.OnesCount32(mask)
	if onesCount == 0 {
		return nil, ErrInvalidAlphabet
	}
	step := (8 * 128) / uint32(onesCount)

	const bufferSize = 256

	// Assign an anonymous function to a variable
	newBuffer := func() interface{} {
		buffer := make([]byte, bufferSize)
		return &buffer
	}

	bufferPool := &sync.Pool{
		New: newBuffer,
	}

	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	return &generator{
		config: Config{
			Alphabet:     alphabetRunes,
			AlphabetLen:  uint32(alphabetLen),
			Mask:         mask,
			Step:         uint32(step),
			IsPowerOfTwo: isPowerOfTwo,
		},
		bufferPool: bufferPool,
		randReader: randReader,
	}, nil
}

// Generate creates a new Nano ID of the specified length.
// It implements the Generator interface.
func (g *generator) Generate(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	var builder strings.Builder
	builder.Grow(length * 4) // Preallocate assuming max 4 bytes per rune (UTF-8 max)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	attempts := 0

	k := bits.Len(uint(g.config.AlphabetLen - 1))
	bytesPerIndex := (k + 7) / 8 // Number of bytes needed per index

	bufferSize := length * bytesPerIndex * 2 // Read extra to reduce number of reads

	// Retrieve a slice from the buffer pool
	bufferPtr := g.bufferPool.Get().(*[]byte)
	buffer := *bufferPtr
	if cap(buffer) < bufferSize {
		buffer = make([]byte, bufferSize)
	} else {
		buffer = buffer[:bufferSize]
	}
	defer g.bufferPool.Put(&buffer)

	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		// Read random bytes into buffer
		_, err := io.ReadFull(g.randReader, buffer)
		if err != nil {
			return "", err
		}

		for i := 0; i <= len(buffer)-bytesPerIndex; i += bytesPerIndex {
			var rnd uint32 = 0
			for j := 0; j < bytesPerIndex; j++ {
				rnd = (rnd << 8) | uint32(buffer[i+j])
			}
			rnd &= g.config.Mask
			if g.config.IsPowerOfTwo || rnd < g.config.AlphabetLen {
				builder.WriteRune(g.config.Alphabet[rnd])
				cursor++
				if cursor == length {
					break
				}
			}
		}
	}
	return builder.String(), nil
}

// GetConfig returns the configuration for the generator.
// It implements the Configuration interface.
func (g *generator) GetConfig() Config {
	return g.config
}
