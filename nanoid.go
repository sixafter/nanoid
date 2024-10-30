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
	Alphabet     []byte // 24 bytes
	AlphabetLen  uint16 // 2 bytes
	Step         uint16 // 2 bytes
	Mask         byte   // 1 byte
	IsPowerOfTwo bool   // 1 byte
	// Padding          // 2 bytes (to align to 8 bytes)
}

type generator struct {
	randReader io.Reader  // 16 bytes
	bufferPool *sync.Pool // 8 bytes
	config     Config     // 32 bytes (from optimized `Config` struct)
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

	alphabetBytes := []byte(alphabet)
	alphabetLen := len(alphabetBytes)

	if alphabetLen < 2 || alphabetLen > 256 {
		return nil, ErrInvalidAlphabet
	}

	// NOTE: Alternatively, a []bool slice can track seen characters. While slightly less memory-efficient than
	// using a bitmask, it's straightforward and still performant.
	//
	// Check for duplicate characters using a boolean slice
	// seen := make([]bool, 256)
	// for _, b := range alphabetBytes {
	//     if seen[b] {
	//         return nil, ErrDuplicateCharacters
	//     }
	//     seen[b] = true
	// }

	// Check for duplicate characters using a bitmask with multiple uint32s
	// A uint32 array can represent 256 bits (32 bits per uint32 × 8 = 256). This allows us to track each
	// possible byte value without the limitations of a single uint64
	var seen [8]uint32 // 8 * 32 = 256 bits
	for _, b := range alphabetBytes {
		idx := b / 32
		bit := b % 32
		if (seen[idx] & (1 << bit)) != 0 {
			return nil, ErrDuplicateCharacters
		}
		seen[idx] |= 1 << bit
	}

	// Calculate mask using power-of-two approach
	k := bits.Len(uint(alphabetLen - 1))
	if k == 0 {
		return nil, ErrInvalidAlphabet
	}
	mask := byte((1 << k) - 1)

	// Calculate step based on mask
	step := (8 * 128) / bits.OnesCount8(mask)

	// Initialize buffer pool to store pointers to byte arrays
	bufferPool := &sync.Pool{
		New: func() interface{} {
			var buffer [128]byte // Using a fixed-size array to avoid dynamic allocation
			return &buffer
		},
	}

	// Determine if alphabet length is a power of two
	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	return &generator{
		config: Config{
			Alphabet:     alphabetBytes,
			AlphabetLen:  uint16(alphabetLen),
			Mask:         mask,
			Step:         uint16(step),
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

	id := make([]byte, length)
	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	attempts := 0

	// Retrieve a pointer to the buffer from the pool
	bufferPtr := g.bufferPool.Get().(*[128]byte)
	buffer := bufferPtr[:]
	defer g.bufferPool.Put(bufferPtr) // Return the pointer to the pool

	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		// Read full buffer
		_, err := io.ReadFull(g.randReader, buffer)
		if err != nil {
			return "", err
		}

		if g.config.IsPowerOfTwo {
			for _, rnd := range buffer {
				rnd &= g.config.Mask
				// Since alphabet length is a power of two, rnd is guaranteed to be within range
				id[cursor] = g.config.Alphabet[rnd]
				cursor++
				if cursor == length {
					break
				}
			}
		} else {
			for _, rnd := range buffer {
				rnd &= g.config.Mask
				if int(rnd) < int(g.config.AlphabetLen) {
					id[cursor] = g.config.Alphabet[rnd]
					cursor++
					if cursor == length {
						break
					}
				}
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
