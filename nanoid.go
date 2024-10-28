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
	return GenerateSize(DefaultSize)
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
	ErrInvalidLength       = errors.New("length must be positive")
	ErrExceededMaxAttempts = errors.New("generate method exceeded maximum attempts, possibly due to invalid mask or alphabet")
	ErrEmptyAlphabet       = errors.New("alphabet must not be empty")
	ErrAlphabetTooShort    = errors.New("alphabet length must be at least 2")
	ErrAlphabetTooLong     = errors.New("alphabet length must not exceed 256")
	ErrDuplicateCharacters = errors.New("alphabet contains duplicate characters")
)

const (
	// DefaultAlphabet Default alphabet as per Nano ID specification.
	DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

	// DefaultSize Default size of the generated Nano ID: 21.
	DefaultSize = 21
)

// Generator holds the configuration for the Nano ID generator.
type Generator interface {
	Generate(size int) (string, error)
}

type Configuration interface {
	GetConfig() Config
}

type Config struct {
	Alphabet    []byte
	AlphabetLen int
	Mask        byte
	Step        int
}

type generator struct {
	randReader io.Reader
	bufferPool *sync.Pool
	config     Config
}

// New creates a new Generator with buffer pooling enabled.
func New(alphabet string, randReader io.Reader) (Generator, error) {
	if len(alphabet) == 0 {
		return nil, ErrEmptyAlphabet
	}

	if randReader == nil {
		randReader = rand.Reader // Initialize here
	}

	alphabetBytes := []byte(alphabet)
	alphabetLen := len(alphabetBytes)

	if alphabetLen < 2 {
		return nil, ErrAlphabetTooShort
	}

	if alphabetLen > 256 {
		return nil, ErrAlphabetTooLong
	}

	// Check for duplicate characters
	seen := make(map[byte]struct{}, alphabetLen)
	for _, b := range alphabetBytes {
		if _, exists := seen[b]; exists {
			return nil, ErrDuplicateCharacters
		}
		seen[b] = struct{}{}
	}

	// Calculate mask using power-of-two approach
	k := bits.Len(uint(alphabetLen - 1))
	if k == 0 {
		return nil, ErrAlphabetTooShort
	}
	mask := byte((1 << k) - 1)

	// Calculate step based on mask
	step := (8 * 128) / bits.OnesCount8(mask)

	// Initialize buffer pool as a pointer
	bufferPool := &sync.Pool{
		New: func() interface{} {
			b := make([]byte, step)
			return &b // Store pointer to slice
		},
	}

	return &generator{
		config: Config{
			Alphabet:    alphabetBytes,
			AlphabetLen: alphabetLen,
			Mask:        mask,
			Step:        step,
		},
		randReader: randReader,
		bufferPool: bufferPool, // Always assigned
	}, nil
}

// GenerateSize creates a new Nano ID of the specified length.
// It ensures that each character in the ID is selected uniformly from the alphabet.
// Pre-allocated errors are used to minimize memory allocations.
func (g *generator) Generate(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	id := make([]byte, length)
	cursor := 0
	maxAttempts := length * 10 // Prevent infinite loops
	attempts := 0

	// Retrieve a pointer to the buffer from the pool
	bufferPtr := g.bufferPool.Get().(*[]byte)
	buffer := *bufferPtr
	defer func() {
		for i := range buffer {
			buffer[i] = 0
		}
		g.bufferPool.Put(bufferPtr) // Return the pointer to the pool
	}()

	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		n, err := g.randReader.Read(buffer)
		if err != nil {
			return "", err
		}
		buffer = buffer[:n]

		for _, rnd := range buffer {
			if int(rnd&g.config.Mask) < g.config.AlphabetLen {
				id[cursor] = g.config.Alphabet[rnd&g.config.Mask]
				cursor++
				if cursor == length {
					break
				}
			}
		}
	}

	return string(id), nil
}

// GetConfig returns the configuration for the generator.
func (g *generator) GetConfig() Config {
	return g.config
}
