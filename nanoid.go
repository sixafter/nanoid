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

// New returns a new Nano ID using `DefaultLength`.
func New() (string, error) {
	return NewWithLength(DefaultLength)
}

// NewWithLength returns a new Nano ID of the specified length.
func NewWithLength(length int) (string, error) {
	return DefaultGenerator.New(length)
}

// Must returns a new Nano ID using `DefaultLength` if err is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func Must() string {
	return MustWithLength(DefaultLength)
}

// MustWithLength returns a new Nano ID of the specified length if err is nil or panics otherwise.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustWithLength(length int) string {
	id, err := NewWithLength(length)
	if err != nil {
		panic(err)
	}

	return id
}

func init() {
	var err error
	DefaultGenerator, err = NewGenerator(
		WithAlphabet(DefaultAlphabet),
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

	// MinAlphabetLength defines the minimum allowed length for the alphabet.
	MinAlphabetLength = 2

	// MaxAlphabetLength defines the maximum allowed length for the alphabet.
	MaxAlphabetLength = 256
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
}

// Config holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type Config interface {
	// RandReader is the source of randomness used for generating IDs.
	RandReader() io.Reader

	// RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
	RuneAlphabet() []rune

	// Mask is a bitmask used to obtain a random value from the character set.
	Mask() uint

	// BitsNeeded represents the number of bits required to generate each character in the ID.
	BitsNeeded() uint

	// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
	BytesNeeded() uint

	// BufferSize is the buffer size used for random byte generation.
	BufferSize() int

	// AlphabetLen is the length of the alphabet, stored as an uint16.
	AlphabetLen() uint16

	// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
	IsPowerOfTwo() bool
}

// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
	// Config returns the runtime configuration of the generator.
	Config() Config
}

// runtimeConfig holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type runtimeConfig struct {
	// RandReader is the source of randomness used for generating IDs.
	randReader io.Reader

	// RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
	runeAlphabet []rune

	// Mask is a bitmask used to obtain a random value from the character set.
	mask uint

	// BitsNeeded represents the number of bits required to generate each character in the ID.
	bitsNeeded uint

	// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
	bytesNeeded uint

	// BufferSize is the buffer size used for random byte generation.
	bufferSize int

	// AlphabetLen is the length of the alphabet, stored as an uint16.
	alphabetLen uint16

	// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
	isPowerOfTwo bool
}

// Generator defines the interface for generating Nano IDs.
type Generator interface {
	// New returns a new Nano ID of the specified length.
	New(length int) (string, error)
}

// generator implements the Generator interface.
type generator struct {
	config *runtimeConfig
	buffer *sync.Pool
}

// NewGenerator creates a new Generator with buffer pooling enabled.
// It accepts variadic Option parameters to configure the Generator.
// It returns an error if the alphabet is invalid or contains invalid UTF-8 characters.
func NewGenerator(options ...Option) (Generator, error) {
	// Initialize ConfigOptions with default values
	configOpts := &ConfigOptions{
		Alphabet:   DefaultAlphabet,
		RandReader: rand.Reader,
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
	pool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, runtimeConfig.bufferSize*utf8.UTFMax) // Max bytes per rune
			return &buf
		},
	}

	return &generator{
		config: runtimeConfig,
		buffer: pool,
	}, nil
}

// buildRuntimeConfig constructs the RuntimeConfig from ConfigOptions.
func buildRuntimeConfig(opts *ConfigOptions) (*runtimeConfig, error) {
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

	// NewGenerator check for maximum alphabet length
	if alphabetLen > MaxAlphabetLength {
		return nil, ErrAlphabetTooLong
	}

	if alphabetLen < MinAlphabetLength {
		return nil, ErrInvalidAlphabet
	}

	// Represents how many bits are required to generate an index for selecting a character from the alphabet.
	bitsNeeded := uint(bits.Len(uint(alphabetLen - 1)))
	if bitsNeeded == 0 {
		return nil, ErrInvalidAlphabet
	}

	// Ensure that only the lowest bitsNeeded bits are used from the random value
	mask := uint((1 << bitsNeeded) - 1)

	// Ensures that any fractional number of bits rounds up to the nearest whole byte.
	bytesNeeded := (bitsNeeded + 7) / 8

	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	// Calculate bufferSize dynamically based on bytesNeeded and bufferMultiplier
	bufferSize := int(bytesNeeded) * bufferMultiplier

	return &runtimeConfig{
		runeAlphabet: alphabetRunes,
		mask:         mask,
		bitsNeeded:   bitsNeeded,
		bytesNeeded:  bytesNeeded,
		bufferSize:   bufferSize,
		alphabetLen:  uint16(alphabetLen),
		isPowerOfTwo: isPowerOfTwo,
		randReader:   opts.RandReader,
	}, nil
}

// New creates a new Nano ID of the specified length.
func (g *generator) New(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	id := make([]rune, length)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	attempts := 0
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded

	// Use rune buffer pool
	randomBytesPtr := g.buffer.Get().(*[]byte)
	defer g.buffer.Put(randomBytesPtr)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)
	step := int(bytesNeeded)
	if step <= 0 {
		return "", ErrInvalidAlphabet
	}

	// New ID
	for cursor < length {
		if attempts >= maxAttempts {
			return "", ErrExceededMaxAttempts
		}
		attempts++

		neededBytes := (length - cursor) * step
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		_, err := io.ReadFull(g.config.randReader, randomBytes[:neededBytes])
		if err != nil {
			return "", err
		}

		for i := 0; i < neededBytes; i += step {
			var rnd uint
			for j := 0; j < step; j++ {
				rnd = (rnd << 8) | uint(randomBytes[i+j])
			}
			rnd &= mask

			if g.config.isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				id[cursor] = g.config.runeAlphabet[rnd]
				cursor++
			}

			if cursor == length {
				break
			}
		}
	}

	return string(id), nil
}

// Config returns the runtime configuration for the generator.
// It implements the Configuration interface.
func (g *generator) Config() Config {
	return g.config
}

// RandReader is the source of randomness used for generating IDs.
func (r *runtimeConfig) RandReader() io.Reader {
	return r.randReader
}

// RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
func (r *runtimeConfig) RuneAlphabet() []rune {
	return r.runeAlphabet
}

// Mask is a bitmask used to obtain a random value from the character set.
func (r *runtimeConfig) Mask() uint {
	return r.mask
}

// BitsNeeded represents the number of bits required to generate each character in the ID.
func (r *runtimeConfig) BitsNeeded() uint {
	return r.bitsNeeded
}

// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
func (r *runtimeConfig) BytesNeeded() uint {
	return r.bytesNeeded
}

// BufferSize is the buffer size used for random byte generation.
func (r *runtimeConfig) BufferSize() int {
	return r.bufferSize
}

// AlphabetLen is the length of the alphabet, stored as an uint16.
func (r *runtimeConfig) AlphabetLen() uint16 {
	return r.alphabetLen
}

// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
func (r *runtimeConfig) IsPowerOfTwo() bool {
	return r.isPowerOfTwo
}
