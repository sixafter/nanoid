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

// DefaultAlphabet as per Nano ID specification.
const DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

// DefaultSize is the default size of the generated Nano ID: 21.
const DefaultSize = 21

// maxAttemptsMultiplier defines the multiplier for maximum attempts based on length.
const maxAttemptsMultiplier = 10

// bufferMultiplier defines how many characters the buffer should handle per read.
// Adjust this value based on performance and memory considerations.
const bufferMultiplier = 64

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
	Alphabet     []rune // 24 bytes
	Mask         uint   // 8 bytes
	BitsNeeded   uint   // 8 bytes
	BytesNeeded  uint   // 8 bytes
	AlphabetLen  uint16 // 2 bytes
	IsPowerOfTwo bool   // 1 byte
	// Padding: 5 bytes to make the struct size a multiple of 8
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

	// Check for duplicate characters using a map
	seen := make(map[rune]bool)
	for _, r := range alphabetRunes {
		if seen[r] {
			return nil, ErrDuplicateCharacters
		}
		seen[r] = true
	}

	// Calculate bitsNeeded and mask
	bitsNeeded := uint(bits.Len(uint(alphabetLen - 1)))
	if bitsNeeded == 0 {
		return nil, ErrInvalidAlphabet
	}
	mask := uint((1 << bitsNeeded) - 1)
	bytesNeeded := (bitsNeeded + 7) / 8

	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	// Dynamic bufferSize Calculation: Calculate bufferSize based on bytesNeeded and bufferMultiplier
	bufferSize := int(bytesNeeded) * bufferMultiplier

	bufferPool := &sync.Pool{
		New: func() interface{} {
			buffer := make([]byte, bufferSize)
			return &buffer
		},
	}

	return &generator{
		config: Config{
			Alphabet:     alphabetRunes,
			AlphabetLen:  uint16(alphabetLen),
			Mask:         mask,
			BitsNeeded:   bitsNeeded,
			BytesNeeded:  bytesNeeded,
			IsPowerOfTwo: isPowerOfTwo,
		},
		bufferPool: bufferPool,
		randReader: randReader,
	}, nil
}

// Generate creates a new Nano ID of the specified length.
func (g *generator) Generate(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	id := make([]rune, length)
	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	attempts := 0

	mask := g.config.Mask
	bytesNeeded := g.config.BytesNeeded

	// Retrieve a buffer from the pool
	randomBytesPtr := g.bufferPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	defer g.bufferPool.Put(randomBytesPtr)

	bufferSize := len(randomBytes)
	step := int(bytesNeeded)
	if step <= 0 {
		return "", ErrInvalidAlphabet
	}

	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		// Calculate how many random bytes we need
		neededBytes := (length - cursor) * step
		if neededBytes > bufferSize {
			neededBytes = bufferSize
		}

		// Read random bytes
		_, err := io.ReadFull(g.randReader, randomBytes[:neededBytes])
		if err != nil {
			return "", err
		}

		// Process random bytes
		for i := 0; i < neededBytes; i += step {
			rnd := uint(0)
			for j := 0; j < step; j++ {
				rnd = (rnd << 8) | uint(randomBytes[i+j])
			}
			rnd &= mask

			if g.config.IsPowerOfTwo {
				// Index is guaranteed to be within bounds
				id[cursor] = g.config.Alphabet[rnd]
				cursor++
			} else {
				if int(rnd) < int(g.config.AlphabetLen) {
					id[cursor] = g.config.Alphabet[rnd]
					cursor++
				}
			}

			if cursor == length {
				break
			}
		}
	}

	return string(id), nil
}

// GetConfig returns the configuration for the generator.
// It implements the Configuration interface.
func (g *generator) GetConfig() Config {
	return g.config
}
