// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/sixafter/nanoid/x/crypto/prng"
)

var (
	// Generator is a global, shared instance of a Nano ID generator. It is safe for concurrent use.
	Generator Interface

	// RandReader is the default random number generator used for generating IDs.
	RandReader = prng.Reader
)

const (
	// DefaultAlphabet defines the standard set of characters used for Nano ID generation.
	// It includes uppercase and lowercase English letters, digits, and the characters
	// '_' and '-'. This selection aligns with the Nano ID specification, ensuring
	// a URL-friendly and easily readable identifier.
	//
	// Example: "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DefaultLength specifies the default number of characters in a generated Nano ID.
	// A length of 21 characters provides a high level of uniqueness while maintaining
	// brevity, making it suitable for most applications requiring unique identifiers.
	DefaultLength = 21

	// maxAttemptsMultiplier determines the maximum number of attempts the generator
	// will make to produce a valid Nano ID before failing. It is calculated as a
	// multiplier based on the desired ID length to balance between performance
	// and the probability of successful ID generation, especially when using
	// non-power-of-two alphabets.
	maxAttemptsMultiplier = 10

	// MinAlphabetLength sets the minimum permissible number of unique characters
	// in the alphabet used for Nano ID generation. An alphabet with fewer than
	// 2 characters would not provide sufficient variability for generating unique IDs,
	// making this a lower bound to ensure meaningful ID generation.
	//
	// Example: An alphabet like "AB" is acceptable, but "A" is not.
	MinAlphabetLength = 2

	// MaxAlphabetLength defines the maximum allowable number of unique characters
	// in the alphabet for Nano ID generation. This upper limit ensures that the
	// generator operates within reasonable memory and performance constraints,
	// preventing excessively large alphabets that could degrade performance or
	// complicate index calculations.
	MaxAlphabetLength = 256
)

func init() {
	var err error
	Generator, err = NewGenerator()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Generator: %v", err))
	}
}

// Interface defines the contract for generating Nano IDs.
//
// Implementations of this interface provide methods to create new IDs
// and to read random data, supporting both ID generation and direct random byte access.
type Interface interface {
	// New generates and returns a new Nano ID as a string with configured length hint.
	// Returns an error if the ID generation fails due to issues like insufficient randomness.
	//
	// Usage:
	//   id, err := generator.New()
	//   if err != nil {
	//       // handle error
	//   }
	//   fmt.Println("Generated ID:", id)
	New() (ID, error)

	// NewWithLength generates and returns a new Nano ID as a string with the specified length.
	// The 'length' parameter determines the number of characters in the generated ID.
	// Returns an error if the ID generation fails due to issues like insufficient randomness.
	//
	// Usage:
	//   id, err := generator.NewWithLength(21)
	//   if err != nil {
	//       // handle error
	//   }
	//   fmt.Println("Generated ID:", id)
	NewWithLength(length int) (ID, error)

	// Read fills the provided byte slice 'p' with random data, reading up to len(p) bytes.
	// Returns the number of bytes read and any error encountered during the read operation.
	//
	// Implements the io.Reader interface, allowing the Interface to be used wherever an io.Reader is accepted.
	// This can be useful for directly obtaining random bytes or integrating with other components that consume random data.
	//
	// Usage:
	//   buffer := make([]byte, 21)
	//   n, err := generator.Read(buffer)
	//   if err != nil {
	//       // handle error
	//   }
	//   fmt.Printf("Read %d random bytes\n", n)
	Read(b []byte) (n int, err error)

	// Config returns the current configuration of the generator.
	Config() Config
}

type generator struct {
	config      *runtimeConfig
	entropyPool *sync.Pool
	idPool      *sync.Pool
	Configuration
}

// New generates a new Nano ID using the default length specified by `DefaultLength`.
// It returns the generated ID as a string and any error encountered during the generation.
//
// Usage:
//
//	id, err := nanoid.New()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func New() (ID, error) {
	return NewWithLength(DefaultLength)
}

// NewWithLength generates a new Nano ID of the specified length.
// It returns the generated ID as a string and any error encountered during the generation.
//
// Parameters:
//   - length int: The number of characters for the generated ID.
//
// Usage:
//
//	id, err := nanoid.NewWithLength(21)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func NewWithLength(length int) (ID, error) {
	return Generator.NewWithLength(length)
}

// Must generates a new Nano ID using the default length specified by `DefaultLength`.
// It returns the generated ID as a string.
// If an error occurs during ID generation, it panics.
// This function simplifies safe initialization of global variables holding pre-generated Nano IDs.
//
// Usage:
//
//	id := nanoid.Must()
//	fmt.Println("Generated ID:", id)
func Must() ID {
	return MustWithLength(DefaultLength)
}

// MustWithLength generates a new Nano ID of the specified length.
// It returns the generated ID as a string.
// If an error occurs during ID generation, it panics.
// The 'length' parameter specifies the number of characters in the generated ID.
// This function simplifies safe initialization of global variables holding pre-generated Nano IDs.
//
// Parameters:
//   - length int: The number of characters for the generated ID.
//
// Usage:
//
//	id := nanoid.MustWithLength(30)
//	fmt.Println("Generated ID:", id)
func MustWithLength(length int) ID {
	id, err := NewWithLength(length)
	if err != nil {
		panic(err)
	}

	return id
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// Reader is the interface that wraps the basic Read method.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// If len(p) == 0, Read should always return n == 0. It may return a
// non-nil error if some error condition is known, such as EOF.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func Read(b []byte) (n int, err error) {
	return Generator.Read(b)
}

// NewGenerator creates a new Interface with buffer pooling enabled.
// It accepts variadic Option parameters to configure the Interface's behavior.
// The function initializes the configuration with default values, applies any provided options,
// validates the configuration, constructs the runtime configuration, initializes buffer pools,
// and returns a configured Interface or an error if the configuration is invalid.
//
// Parameters:
//   - options ...Option: A variadic list of Option functions to customize the Interface's configuration.
//
// Returns:
//   - Interface: An instance of the Interface interface configured with the specified options.
//   - error: An error object if the Interface could not be created due to invalid configuration.
//
// Error Conditions:
//   - ErrInvalidLength: Returned if the provided LengthHint is less than 1.
//   - ErrNilRandReader: Returned if the provided RandReader is nil.
//   - ErrInvalidAlphabet: Returned if the alphabet is invalid or contains invalid UTF-8 characters.
//   - ErrNonUTF8Alphabet: Returned if the alphabet contains non-UTF-8 characters.
//   - ErrDuplicateCharacters: Returned if the alphabet contains duplicate characters.
func NewGenerator(options ...Option) (Interface, error) {
	// Initialize ConfigOptions with default values.
	// These defaults include the default alphabet, the default random reader,
	// and the default length hint for ID generation.
	configOpts := &ConfigOptions{
		Alphabet:   DefaultAlphabet,
		RandReader: RandReader,
		LengthHint: DefaultLength,
	}

	// Apply provided options to customize the configuration.
	// Each Option function modifies the ConfigOptions accordingly.
	for _, opt := range options {
		opt(configOpts)
	}

	// Ensure LengthHint is within valid bounds.
	// LengthHint must be at least 1 to generate meaningful IDs.
	if configOpts.LengthHint < 1 {
		return nil, ErrInvalidLength
	}

	// Ensure RandReader is not nil.
	// A valid randomness source is essential for generating secure IDs.
	if configOpts.RandReader == nil {
		return nil, ErrNilRandReader
	}

	// Validate and construct RuntimeConfig based on the current ConfigOptions.
	// buildRuntimeConfig performs validation on the alphabet and computes necessary
	// parameters for efficient ID generation.
	config, err := buildRuntimeConfig(configOpts)
	if err != nil {
		return nil, err
	}

	// Initialize a pool of byte slices for random data generation.
	// The pool helps in reusing memory buffers, reducing garbage collection overhead.
	entropyPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, config.bufferSize*config.bufferMultiplier)
			return &buf
		},
	}

	var idPool *sync.Pool
	if config.isASCII {
		idPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, config.bufferSize*config.bufferMultiplier)
				return &buf
			},
		}
	} else {
		idPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]rune, config.bufferSize*config.bufferMultiplier)
				return &buf
			},
		}
	}

	// Return the configured Interface instance.
	// The generator holds references to the runtime configuration and buffer pools,
	// facilitating efficient and thread-safe ID generation.
	return &generator{
		config:      config,
		entropyPool: entropyPool,
		idPool:      idPool,
	}, nil
}

// New generates a new Nano ID string of the configured length hint.
//
// Returns:
//   - string: The generated Nano ID.
//   - error: An error object if the generation fails due to invalid input.
//
// Usage:
//
//	id, err := generator.New()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func (g *generator) New() (ID, error) {
	return g.NewWithLength(int(g.config.lengthHint))
}

// NewWithLength generates a new Nano ID string of the specified length.
//
// It validates the provided length to ensure it is a positive integer.
// Depending on the generator's configuration, it generates the ID using the appropriate method.
//
// Parameters:
//   - length int: The desired number of characters in the generated Nano ID.
//
// Returns:
//   - string: The generated Nano ID.
//   - error: An error object if the generation fails due to invalid input.
//
// Error Conditions:
//   - ErrInvalidLength: Returned if the provided length is less than or equal to zero.
//
// Usage:
//
//	id, err := generator.NewWithLength(21)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func (g *generator) NewWithLength(length int) (ID, error) {
	if length <= 0 {
		return EmptyID, ErrInvalidLength
	}

	if g.config.isASCII {
		return g.newASCII(length)
	}
	return g.newUnicode(length)
}

// Config holds the runtime configuration for the Nano ID generator.
//
// It is immutable after initialization and provides all the necessary
// parameters for generating unique IDs efficiently and securely.
func (g *generator) Config() Config {
	return g.config
}

// newASCII generates a new Nano ID using the ASCII alphabet.
func (g *generator) newASCII(length int) (ID, error) {
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	// Defer returning the randomBytes buffer to the pool
	defer func() {
		g.entropyPool.Put(randomBytesPtr)
	}()

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	// Retrieve the idBuffer from the pool
	idBufferPtr := g.idPool.Get().(*[]byte)
	idBuffer := (*idBufferPtr)[:length] // Ensure it has the correct length

	defer func() {
		g.idPool.Put(idBufferPtr)
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return EmptyID, err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				idBuffer[cursor] = g.config.byteAlphabet[rnd]
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return EmptyID, ErrExceededMaxAttempts
	}

	return ID(idBuffer), nil
}

// newUnicode generates a new Nano ID using the Unicode alphabet.
func (g *generator) newUnicode(length int) (ID, error) {
	// Retrieve random bytes from the pool
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	// Defer returning the randomBytes buffer to the pool
	defer func() {
		g.entropyPool.Put(randomBytesPtr)
	}()

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	// Retrieve the idBuffer from the pool
	idBufferPtr := g.idPool.Get().(*[]rune)
	idBuffer := (*idBufferPtr)[:length] // Ensure it has the correct length

	defer func() {
		g.idPool.Put(idBufferPtr)
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return EmptyID, err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				idBuffer[cursor] = g.config.runeAlphabet[rnd]
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return EmptyID, ErrExceededMaxAttempts
	}

	return ID(idBuffer), nil
}

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// If len(p) == 0, Read should always return n == 0. It may return a
// non-nil error if some error condition is known, such as EOF.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func (g *generator) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	length := len(p)
	id, err := g.NewWithLength(length)
	if err != nil {
		return 0, err
	}

	copy(p, id)
	return length, nil
}

// processRandomBytes extracts and returns an unsigned integer from the given randomBytes slice,
// starting at the specified index 'i'. The size of the returned value is determined by the
// g.config.bytesNeeded field.
//
// Parameters:
//   - randomBytes: A byte slice containing random data.
//   - i: The starting index from which to extract the required bytes from the randomBytes slice.
//
// Returns:
//   - uint: An unsigned integer constructed from the bytes, with a size defined by g.config.bytesNeeded.
//
// Behavior:
//   - If bytesNeeded is 1, a single byte is returned as an unsigned integer.
//   - If bytesNeeded is 2, the function returns a 16-bit unsigned integer (2 bytes) in Big Endian order.
//   - If bytesNeeded is 4, the function returns a 32-bit unsigned integer (4 bytes) in Big Endian order.
//   - For other values of bytesNeeded, it constructs an unsigned integer by shifting and combining each byte.
//
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

// GetConfig returns the current configuration of the generator.
func (g *generator) GetConfig() Config {
	return g.config
}
