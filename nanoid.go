// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
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
	ErrAlphabetTooShort    = errors.New("alphabet length is less than 2")
	ErrAlphabetTooLong     = errors.New("alphabet length exceeds 256")
	ErrNilRandReader       = errors.New("nil random reader")
)

const (
	// DefaultAlphabet as per Nano ID specification; A-Za-z0-9_-.
	DefaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DefaultLength is the default size of the generated Nano ID: 21.
	DefaultLength = 21

	// maxAttemptsMultiplier defines the multiplier for maximum attempts based on length.
	maxAttemptsMultiplier = 10

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

// WithLengthHint sets the hint of the intended length of the IDs to be generated.
func WithLengthHint(hint uint16) Option {
	return func(c *ConfigOptions) {
		c.LengthHint = hint
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

	// LengthHint specifies a typical or default length for generated IDs.
	LengthHint uint16
}

// Config holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type Config interface {
	// RandReader returns the source of randomness used for generating IDs.
	RandReader() io.Reader

	// ByteAlphabet returns the slice of bytes for ASCII alphabets.
	ByteAlphabet() []byte

	// RuneAlphabet returns the slice of runes used for ID generation, allowing support for multibyte characters.
	RuneAlphabet() []rune

	// Mask returns the bitmask used to obtain a random value from the character set.
	Mask() uint

	// BitsNeeded returns the number of bits required to generate each character in the ID.
	BitsNeeded() uint

	// BytesNeeded returns the number of bytes required from the random source to produce the entire ID.
	BytesNeeded() uint

	// BufferSize returns the calculated size of the buffer used for random byte generation.
	BufferSize() int

	// AlphabetLen returns the length of the alphabet used for ID generation.
	AlphabetLen() uint16

	// IsPowerOfTwo returns true if the length of the alphabet is a power of two, optimizing random selection for efficient bit operations.
	IsPowerOfTwo() bool

	// IsASCII returns true if the alphabet consists solely of ASCII characters.
	IsASCII() bool

	// BufferMultiplier returns the multiplier used to determine how many characters the buffer should handle per read.
	BufferMultiplier() int

	// BaseMultiplier returns the base multiplier used to determine the growth rate of buffer size, accounting for small ID lengths to achieve balance.
	BaseMultiplier() int

	// ScalingFactor returns the scaling factor used to balance the alphabet size and ID length, ensuring smoother growth in buffer size calculations.
	ScalingFactor() int

	// LengthHint returns the hint of the intended length of the IDs to be generated.
	LengthHint() uint16
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

	// byteAlphabet is a slice of bytes for ASCII alphabets.
	byteAlphabet []byte

	// runeAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
	runeAlphabet []rune

	// Mask is a bitmask used to obtain a random value from the character set.
	mask uint

	// BitsNeeded represents the number of bits required to generate each character in the ID.
	bitsNeeded uint

	// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
	bytesNeeded uint

	// BufferSize is the buffer size used for random byte generation.
	bufferSize int

	// BufferMultiplier defines the multiplier used to calculate the buffer size for reading random bytes, ensuring gradual and consistent scaling.
	bufferMultiplier int

	// ScalingFactor adjusts the balance between alphabet size and id length to achieve smoother scaling in buffer size calculations.
	scalingFactor int

	// BaseMultiplier is used to determine the growth rate of the buffer size, adjusted for small ID lengths to ensure balance.
	baseMultiplier int

	// AlphabetLen is the length of the alphabet, stored as an uint16.
	alphabetLen uint16

	// isASCII indicates whether the alphabet consists solely of ASCII characters.
	isASCII bool

	// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
	isPowerOfTwo bool

	// LengthHint the hint of the intended length of the IDs to be generated.
	lengthHint uint16
}

// Generator defines the interface for generating Nano IDs.
type Generator interface {
	// New returns a new Nano ID of the specified length.
	New(length int) (string, error)
}

// generator implements the Generator interface.
type generator struct {
	config          *runtimeConfig
	randomBytesPool *sync.Pool
	asciiIDPool     *sync.Pool
	unicodeIDPool   *sync.Pool
}

// NewGenerator creates a new Generator with buffer pooling enabled.
// It accepts variadic Option parameters to configure the Generator.
// It returns an error if the alphabet is invalid or contains invalid UTF-8 characters.
func NewGenerator(options ...Option) (Generator, error) {
	// Initialize ConfigOptions with default values
	configOpts := &ConfigOptions{
		Alphabet:   DefaultAlphabet,
		RandReader: rand.Reader,
		LengthHint: DefaultLength,
	}

	// Apply provided options
	for _, opt := range options {
		opt(configOpts)
	}

	// ensure LengthHint is within bounds
	if configOpts.LengthHint < 1 {
		return nil, ErrInvalidLength
	}

	// ensure reader is not nil
	if configOpts.RandReader == nil {
		return nil, ErrNilRandReader
	}

	// Validate and construct RuntimeConfig
	runtimeConfig, err := buildRuntimeConfig(configOpts)
	if err != nil {
		return nil, err
	}

	// Initialize buffer pools based on Rune handling
	randomBytesPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, runtimeConfig.bufferSize*runtimeConfig.bufferMultiplier)
			return &buf
		},
	}

	// Initialize ID buffer pool with *([]byte)
	asciiIDPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, runtimeConfig.bufferSize*runtimeConfig.bufferMultiplier)
			return &buf
		},
	}

	// Initialize Rune buffer pool with *[]rune
	unicodeIDPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]rune, 0, runtimeConfig.bufferSize*runtimeConfig.bufferMultiplier)
			return &buf
		},
	}

	return &generator{
		config:          runtimeConfig,
		randomBytesPool: randomBytesPool,
		asciiIDPool:     asciiIDPool,
		unicodeIDPool:   unicodeIDPool,
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

	alphabetRunes := []rune(opts.Alphabet)
	isASCII := true
	byteAlphabet := make([]byte, len(alphabetRunes))
	for i, r := range alphabetRunes {
		if r > 0x7F { // 127: highest code point in the 7-bit ASCII character set.
			isASCII = false
			break
		}
		byteAlphabet[i] = byte(r)
	}

	if !isASCII {
		// Convert to rune alphabet if non-ASCII characters are present
		byteAlphabet = nil // Clear byteAlphabet as it's not used
	}

	// Check for duplicate characters
	seenRunes := make(map[rune]bool)
	for _, r := range alphabetRunes {
		if seenRunes[r] {
			return nil, ErrDuplicateCharacters
		}
		seenRunes[r] = true
	}

	// Check alphabet length constraints
	if len(alphabetRunes) > MaxAlphabetLength {
		return nil, ErrAlphabetTooLong
	}
	if len(alphabetRunes) < MinAlphabetLength {
		return nil, ErrAlphabetTooShort
	}

	// Calculate BitsNeeded and Mask
	bitsNeeded := uint(bits.Len(uint(len(alphabetRunes) - 1)))
	if bitsNeeded == 0 {
		return nil, ErrInvalidAlphabet
	}

	mask := uint((1 << bitsNeeded) - 1)

	// TODO: mprimeaux: Scale bitsNeeded based on length hint
	//adjustedBitsNeeded := bitsNeeded + uint(math.Log2(float64(opts.LengthHint)))

	// Ensures that any fractional number of bits rounds up to the nearest whole byte.
	bytesNeeded := (bitsNeeded + 7) / 8

	isPowerOfTwo := (len(alphabetRunes) & (len(alphabetRunes) - 1)) == 0

	// Adjust the calculation for the baseMultiplier to achieve smooth growth based on id length and alphabet length
	baseMultiplier := int(math.Ceil(math.Log2(float64(opts.LengthHint) + 2.0)))

	// Modify the scaling factor to balance alphabet size and id length for smoother scaling
	scalingFactor := int(math.Max(3.0, float64(len(alphabetRunes))/math.Pow(float64(opts.LengthHint), 0.6)))

	// Refine bufferMultiplier calculation for a smooth scaling pattern
	bufferMultiplier := baseMultiplier + int(math.Ceil(float64(scalingFactor)/1.5))

	// Recalculate bufferSize to ensure consistent and smooth scaling
	bufferSize := bufferMultiplier * int(bytesNeeded) * int(math.Max(1.5, float64(opts.LengthHint)/10.0))

	return &runtimeConfig{
		randReader:       opts.RandReader,
		byteAlphabet:     byteAlphabet,
		runeAlphabet:     alphabetRunes,
		mask:             mask,
		bitsNeeded:       bitsNeeded,
		bytesNeeded:      bytesNeeded,
		bufferSize:       bufferSize,
		bufferMultiplier: bufferMultiplier,
		scalingFactor:    scalingFactor,
		baseMultiplier:   baseMultiplier,
		alphabetLen:      uint16(len(alphabetRunes)),
		isASCII:          isASCII,
		isPowerOfTwo:     isPowerOfTwo,
		lengthHint:       opts.LengthHint,
	}, nil
}

// Buffer is a type constraint that allows either []byte or []rune.
type Buffer interface {
	~[]byte | ~[]rune
}

// processRandomBytes processes the random bytes into n rnd value.
// This function is kept small to encourage inlining by the compiler.
func (g *generator) processRandomBytes(randomBytes []byte, i int) uint {
	switch g.config.bytesNeeded {
	case 1:
		return uint(randomBytes[i])
	case 2:
		return uint(binary.BigEndian.Uint16(randomBytes[i : i+2]))
	case 4:
		return uint(binary.BigEndian.Uint32(randomBytes[i : i+4]))
	default:
		var rnd uint
		for j := 0; j < int(g.config.bytesNeeded); j++ {
			rnd = (rnd << 8) | uint(randomBytes[i+j])
		}
		return rnd
	}
}

// New creates a new Nano ID of the specified length.
func (g *generator) New(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	if g.config.isASCII {
		return g.newASCII(length)
	}
	return g.newUnicode(length)
}

// newASCII generates a new Nano ID using the ASCII alphabet.
func (g *generator) newASCII(length int) (string, error) {
	// Retrieve a buffer from the pool
	idPtr := g.asciiIDPool.Get().(*[]byte)

	// Ensure the buffer has enough capacity
	var id []byte
	if cap(*idPtr) >= length {
		id = (*idPtr)[:length]
	} else {
		id = make([]byte, length)
	}

	// Retrieve random bytes from the pool
	randomBytesPtr := g.randomBytesPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	// Flag to indicate successful generation
	success := false

	// Ensure resources are returned to the pool
	defer func() {
		g.randomBytesPool.Put(randomBytesPtr)
		if success {
			g.asciiIDPool.Put(idPtr)
		} else {
			// If a new buffer was created (not from the pool), do not return it
			if cap(*idPtr) >= length {
				g.asciiIDPool.Put(idPtr)
			}
		}
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return "", err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				if rnd >= uint(len(g.config.byteAlphabet)) {
					panic(fmt.Sprintf("rnd value %d exceeds byteAlphabet length %d", rnd, len(g.config.byteAlphabet)))
				}
				id[cursor] = g.config.byteAlphabet[rnd]
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return "", ErrExceededMaxAttempts
	}

	success = true
	return string(id[:cursor]), nil
}

// newUnicode generates a new Nano ID using the Unicode alphabet.
func (g *generator) newUnicode(length int) (string, error) {
	// Retrieve a rune buffer from the pool
	idRunesPtr := g.unicodeIDPool.Get().(*[]rune)

	// Ensure the buffer has enough capacity
	var idRunes []rune
	if cap(*idRunesPtr) >= length {
		idRunes = (*idRunesPtr)[:length]
	} else {
		idRunes = make([]rune, length)
	}

	// Retrieve random bytes from the pool
	randomBytesPtr := g.randomBytesPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	// Flag to indicate successful generation
	success := false

	// Ensure resources are returned to the pool
	defer func() {
		g.randomBytesPool.Put(randomBytesPtr)
		if success {
			g.unicodeIDPool.Put(idRunesPtr)
		} else {
			// If a new buffer was created (not from the pool), do not return it
			if cap(*idRunesPtr) >= length {
				g.unicodeIDPool.Put(idRunesPtr)
			}
		}
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return "", err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				if rnd >= uint(len(g.config.runeAlphabet)) {
					panic(fmt.Sprintf("rnd value %d exceeds runeAlphabet length %d", rnd, len(g.config.runeAlphabet)))
				}
				idRunes[cursor] = g.config.runeAlphabet[rnd]
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return "", ErrExceededMaxAttempts
	}

	success = true
	return string(idRunes[:cursor]), nil
}

// Config returns the runtime configuration for the generator.
// It implements the Configuration interface.
func (g *generator) Config() Config {
	return g.config
}

// RandReader is the source of randomness used for generating IDs.
func (r runtimeConfig) RandReader() io.Reader {
	return r.randReader
}

// RuneAlphabet is a slice of runes, allowing support for multibyte characters in ID generation.
func (r runtimeConfig) RuneAlphabet() []rune {
	return r.runeAlphabet
}

// Mask is a bitmask used to obtain a random value from the character set.
func (r runtimeConfig) Mask() uint {
	return r.mask
}

// BitsNeeded represents the number of bits required to generate each character in the ID.
func (r runtimeConfig) BitsNeeded() uint {
	return r.bitsNeeded
}

// BytesNeeded specifies the number of bytes required from a random source to produce the ID.
func (r runtimeConfig) BytesNeeded() uint {
	return r.bytesNeeded
}

// BufferSize is the buffer size used for random byte generation.
func (r runtimeConfig) BufferSize() int {
	return r.bufferSize
}

// AlphabetLen is the length of the alphabet, stored as an uint16.
func (r runtimeConfig) AlphabetLen() uint16 {
	return r.alphabetLen
}

// IsPowerOfTwo indicates whether the length of the alphabet is a power of two, optimizing random selection.
func (r runtimeConfig) IsPowerOfTwo() bool {
	return r.isPowerOfTwo
}

// BufferMultiplier is the multiplier used to calculate the buffer size for reading random bytes, ensuring gradual and consistent scaling.
func (r runtimeConfig) BufferMultiplier() int {
	return r.bufferMultiplier
}

// BaseMultiplier is used to determine the growth rate of the buffer size, adjusted for small ID lengths to ensure balance.
func (r runtimeConfig) BaseMultiplier() int {
	return r.baseMultiplier
}

// ScalingFactor adjusts the balance between alphabet size and id length to achieve smoother scaling in buffer size calculations.
func (r runtimeConfig) ScalingFactor() int {
	return r.scalingFactor
}

// IsASCII indicates whether the alphabet consists solely of ASCII characters.
func (r runtimeConfig) IsASCII() bool {
	return r.isASCII
}

// ByteAlphabet returns a slice of bytes for ASCII alphabets.
func (r runtimeConfig) ByteAlphabet() []byte {
	return r.byteAlphabet
}

// LengthHint the hint of the intended length of the IDs to be generated.
func (r runtimeConfig) LengthHint() uint16 {
	return r.lengthHint
}
