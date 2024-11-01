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
	"unicode/utf8"
)

// DefaultGenerator is a global, shared instance of a Nano ID generator. It is safe for concurrent use.
var DefaultGenerator Generator

// Generate returns a new Nano ID using `DefaultLength`.
func Generate() (string, error) {
	return DefaultGenerator.Generate(DefaultLength)
}

// GenerateWithLength returns a new Nano ID of the specified length.
func GenerateWithLength(length int) (string, error) {
	return DefaultGenerator.Generate(length)
}

// MustGenerate returns a new Nano ID using `DefaultLength` if err is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustGenerate() string {
	return DefaultGenerator.MustGenerate(DefaultLength)
}

// MustGenerateWithLength returns a new Nano ID of the specified length if err is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustGenerateWithLength(length int) string {
	return DefaultGenerator.MustGenerate(length)
}

func init() {
	var err error
	DefaultGenerator, err = New(
		WithAlphabet(DefaultAlphabet),
		WithDefaultLength(DefaultLength),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DefaultGenerator: %v", err))
	}
}

var (
	ErrDuplicateCharacters = errors.New("duplicate characters in alphabet")
	ErrExceededMaxAttempts = errors.New("exceeded maximum attempts")
	ErrInvalidLength       = errors.New("invalid length")
	ErrInvalidAlphabet     = errors.New("invalid alphabet")
	ErrNonUTF8Alphabet     = errors.New("alphabet contains invalid UTF-8 characters")
	ErrAlphabetTooLong     = errors.New("alphabet length exceeds 256")
)

// Option defines a function type for configuring the Generator.
type Option func(*ConfigOptions)

// WithAlphabet sets a custom alphabet for the Generator.
func WithAlphabet(alphabet string) Option {
	return func(c *ConfigOptions) {
		c.Alphabet = alphabet
	}
}

// WithRandReader sets a custom random reader for the Generator.
func WithRandReader(reader io.Reader) Option {
	return func(c *ConfigOptions) {
		c.RandReader = reader
	}
}

// WithDefaultLength sets a custom default length for ID generation.
func WithDefaultLength(length int) Option {
	return func(c *ConfigOptions) {
		if length > 0 {
			c.DefaultLength = length
		}
	}
}

const (
	// DefaultAlphabet as per Nano ID specification.
	DefaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DefaultLength is the default size of the generated Nano ID: 21.
	DefaultLength = 21

	// maxAttemptsMultiplier defines the multiplier for maximum attempts based on length.
	maxAttemptsMultiplier = 10

	// bufferMultiplier defines how many characters the buffer should handle per read.
	// Adjust this value based on performance and memory considerations.
	bufferMultiplier = 128

	// MaxAlphabetLength defines the maximum allowed length for the alphabet.
	MaxAlphabetLength = 256 // Newly added constant
)

// ConfigOptions holds the configurable options for the Generator.
// It is used with the Function Options pattern.
type ConfigOptions struct {
	// RandReader is the source of randomness used for generating IDs.
	// By default, it uses crypto/rand.Reader, which provides cryptographically secure random bytes.
	RandReader io.Reader

	// Alphabet is the set of characters used to generate the Nano ID.
	// It must be a valid UTF-8 string containing between 2 and 256 unique characters.
	// Using a diverse and appropriately sized alphabet ensures the uniqueness and randomness of the generated IDs.
	Alphabet string

	// DefaultLength is the default length of the generated Nano ID when no specific length is provided during generation.
	DefaultLength int
}

// RuntimeConfig holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type RuntimeConfig struct {
	// RandReader is the source of randomness used for generating IDs.
	RandReader io.Reader

	// RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
	RuneAlphabet []rune

	// Mask is a bitmask used to obtain a random value from the character set.
	Mask uint

	// BitsNeeded represents the number of bits required to generate each character in the ID.
	BitsNeeded uint

	// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
	BytesNeeded uint

	// BufferSize is the buffer size used for random byte generation.
	BufferSize int

	// DefaultLength is the default size of the generated Nano ID.
	DefaultLength int

	// AlphabetLen is the length of the alphabet, stored as an uint16.
	AlphabetLen uint16

	// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
	IsPowerOfTwo bool
}

// Generator defines the interface for generating Nano IDs.
type Generator interface {
	// Generate returns a new Nano ID of the specified length.
	Generate(length int) (string, error)

	// MustGenerate returns a new Nano ID of the specified length if err is nil or panics otherwise.
	MustGenerate(length int) string
}

// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
	// GetConfig returns the runtime configuration of the generator.
	GetConfig() RuntimeConfig
}

// generator implements the Generator interface.
type generator struct {
	config         *RuntimeConfig
	runeBufferPool *sync.Pool
}

// New creates a new Generator with buffer pooling enabled.
// It accepts variadic Option parameters to configure the Generator.
// It returns an error if the alphabet is invalid or contains invalid UTF-8 characters.
func New(options ...Option) (Generator, error) {
	// Initialize ConfigOptions with default values
	configOpts := &ConfigOptions{
		Alphabet:      DefaultAlphabet,
		DefaultLength: DefaultLength,
		RandReader:    rand.Reader,
	}

	// Apply provided options
	for _, opt := range options {
		opt(configOpts)
	}

	// Validate and construct RuntimeConfig
	runtimeConfig, err := buildRuntimeConfig(configOpts)
	if err != nil {
		return nil, err
	}

	// Initialize buffer pools based on Rune handling
	runePool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, runtimeConfig.BufferSize*utf8.UTFMax) // Max bytes per rune
			return &buf
		},
	}

	return &generator{
		config:         runtimeConfig,
		runeBufferPool: runePool,
	}, nil
}

// buildRuntimeConfig constructs the RuntimeConfig from ConfigOptions.
func buildRuntimeConfig(opts *ConfigOptions) (*RuntimeConfig, error) {
	if len(opts.Alphabet) == 0 {
		return nil, ErrInvalidAlphabet
	}

	// Check if the alphabet is valid UTF-8
	if !utf8.ValidString(opts.Alphabet) {
		return nil, ErrNonUTF8Alphabet
	}

	// Convert the alphabet to runes
	alphabetRunes := []rune(opts.Alphabet)

	// Check for duplicate characters
	seenRunes := make(map[rune]bool)
	for _, r := range alphabetRunes {
		if seenRunes[r] {
			return nil, ErrDuplicateCharacters
		}
		seenRunes[r] = true
	}

	// Calculate BitsNeeded and Mask
	alphabetLen := len(alphabetRunes)

	// New check for maximum alphabet length
	if alphabetLen > MaxAlphabetLength {
		return nil, ErrAlphabetTooLong
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

	return &RuntimeConfig{
		RuneAlphabet:  alphabetRunes,
		Mask:          mask,
		BitsNeeded:    bitsNeeded,
		BytesNeeded:   bytesNeeded,
		BufferSize:    bufferSize,
		AlphabetLen:   uint16(alphabetLen),
		IsPowerOfTwo:  isPowerOfTwo,
		DefaultLength: opts.DefaultLength,
		RandReader:    opts.RandReader,
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

	// Use rune buffer pool
	randomBytesPtr := g.runeBufferPool.Get().(*[]byte)
	defer g.runeBufferPool.Put(randomBytesPtr)
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

		_, err := io.ReadFull(g.config.RandReader, randomBytes[:neededBytes])
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
				id[cursor] = g.config.RuneAlphabet[rnd]
				cursor++
			}

			if cursor == length {
				break
			}
		}
	}

	return string(id), nil
}

// MustGenerate returns a new Nano ID of the specified length if err is nil or panics otherwise.
func (g *generator) MustGenerate(length int) string {
	id, err := g.Generate(length)
	if err != nil {
		panic(err)
	}
	return id
}

// GetConfig returns the runtime configuration for the generator.
// It implements the Configuration interface.
func (g *generator) GetConfig() RuntimeConfig {
	return *g.config
}
