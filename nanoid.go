// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/sixafter/prng-chacha"
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
	Interface
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
	// Validate the requested length: must be positive and non-zero. ---
	if length <= 0 {
		return EmptyID, ErrInvalidLength
	}

	// ASCII Mode: Use a stack-allocated byte buffer for maximal performance and zero pooling overhead.
	if g.config.isASCII {
		buf := make([]byte, length) // Allocate the buffer for the resulting ID.
		if err := g.newASCII(length, buf); err != nil {
			// Propagate any error (randomness failure, exceeded attempts, etc).
			return EmptyID, err
		}
		// Convert the filled byte buffer to an ID using an unsafe zero-copy conversion.
		// This is the only allocation in this fast path.
		return ID(unsafe.String(&buf[0], len(buf))), nil
	}

	// Unicode Mode: Use a pooled buffer for runes to minimize allocations and reduce GC pressure.
	idBufferPtr := g.idPool.Get().(*[]rune) // Get a pooled buffer pointer from the sync.Pool.
	idBuffer := (*idBufferPtr)[:length]     // Slice the buffer to the required length.
	defer g.idPool.Put(idBufferPtr)         // Ensure the buffer is returned to the pool after use.

	if err := g.newUnicode(length, idBuffer); err != nil {
		// Propagate errors from ID generation (randomness failure, pool issues, etc).
		return EmptyID, err
	}

	// Convert the filled rune buffer to an ID.
	// This conversion (from []rune to string) will always incur a single allocation,
	// as required by Go's memory model for Unicode string construction.
	return ID(idBuffer), nil
}

// Config holds the runtime configuration for the Nano ID generator.
//
// It is immutable after initialization and provides all the necessary
// parameters for generating unique IDs efficiently and securely.
func (g *generator) Config() Config {
	return g.config
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
func (g *generator) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if g.config.isASCII {
		// Fill ASCII directly into client buffer
		if err := g.newASCII(len(p), p); err != nil {
			return 0, err
		}
		return len(p), nil
	}

	// Unicode: Fill runes, then encode as UTF-8 into p
	runeBuf := make([]rune, g.config.lengthHint)
	if err := g.newUnicode(len(runeBuf), runeBuf); err != nil {
		return 0, err
	}
	// Encode runes into p, up to capacity
	n := 0
	for _, r := range runeBuf {
		size := utf8.RuneLen(r)
		if n+size > len(p) {
			break // buffer full
		}
		utf8.EncodeRune(p[n:], r)
		n += size
	}
	return n, nil
}

// newASCII generates a new Nano ID using the configured ASCII alphabet.
//
// This method fills the provided byte slice `idBuffer` with a random sequence of
// characters selected from the generator's ASCII alphabet. The buffer must have at
// least `length` capacity. The function performs no heap allocations beyond the use
// of internal buffer pools.
//
// Parameters:
//   - length:   The number of random ASCII characters to generate.
//   - idBuffer: The buffer to fill with generated characters.
//
// Returns:
//   - error: An error if the buffer is too small, random source fails, or max attempts are exceeded.
//
// Example:
//
//	buf := make([]byte, 21)
//	err := generator.newASCII(21, buf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(buf))
func (g *generator) newASCII(length int, idBuffer []byte) error {
	// --- Parameter validation: Ensure output buffer is large enough. ---
	if len(idBuffer) < length {
		return fmt.Errorf("buffer too small")
	}

	// --- Acquire a buffer for entropy from the pool for efficient random data handling. ---
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)
	defer g.entropyPool.Put(randomBytesPtr) // Always return buffer to the pool.

	// --- Initialize internal state for NanoID generation. ---
	cursor := 0                                   // Tracks how many characters have been written.
	maxAttempts := length * maxAttemptsMultiplier // Upper bound to prevent infinite loops.
	mask := g.config.mask                         // Bitmask for index truncation.
	bytesNeeded := g.config.bytesNeeded           // Bytes required for each index sample.
	isPowerOfTwo := g.config.isPowerOfTwo         // Optimization for binary-friendly alphabets.

	// --- Main generation loop: Fill idBuffer with valid random characters. ---
	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		// --- Determine how many random bytes are needed for this iteration. ---
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// --- Refill the entropy buffer with secure random bytes. ---
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			// Propagate any error from the random reader.
			return err
		}

		// --- Use each segment of random bytes to select characters from the alphabet. ---
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i) // Extract an integer sample.
			rnd &= mask                                 // Mask for non-power-of-two alphabet sizes.

			// --- If the random index is valid, select the character and write to output. ---
			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				idBuffer[cursor] = g.config.byteAlphabet[rnd]
				cursor++
			}
		}
	}

	// --- If unable to generate a full ID within the allowed attempts, return an error. ---
	if cursor < length {
		return ErrExceededMaxAttempts
	}

	return nil
}

// newUnicode generates a new Nano ID using the configured Unicode alphabet.
//
// This method fills the provided rune slice `idBuffer` with a random sequence of
// runes selected from the generator's Unicode alphabet. The buffer must have at
// least `length` capacity. No heap allocations occur outside internal buffer pools.
//
// Parameters:
//   - length:   The number of random runes to generate.
//   - idBuffer: The buffer to fill with generated runes.
//
// Returns:
//   - error: An error if the buffer is too small, random source fails, or max attempts are exceeded.
//
// Example:
//
//	buf := make([]rune, 16)
//	err := generator.newUnicode(16, buf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(buf))
func (g *generator) newUnicode(length int, idBuffer []rune) error {
	// --- Parameter validation: Ensure output buffer is large enough. ---
	if len(idBuffer) < length {
		return fmt.Errorf("buffer too small")
	}

	// --- Acquire a buffer for entropy from the pool for efficient random data handling. ---
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)
	defer g.entropyPool.Put(randomBytesPtr) // Always return buffer to the pool.

	// --- Initialize internal state for NanoID generation. ---
	cursor := 0                                   // Tracks how many runes have been written.
	maxAttempts := length * maxAttemptsMultiplier // Upper bound to prevent infinite loops.
	mask := g.config.mask                         // Bitmask for index truncation.
	bytesNeeded := g.config.bytesNeeded           // Bytes required for each index sample.
	isPowerOfTwo := g.config.isPowerOfTwo         // Optimization for binary-friendly alphabets.

	// --- Main generation loop: Fill idBuffer with valid random runes. ---
	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		// --- Determine how many random bytes are needed for this iteration. ---
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// --- Refill the entropy buffer with secure random bytes. ---
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			// Propagate any error from the random reader.
			return err
		}

		// --- Use each segment of random bytes to select runes from the alphabet. ---
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i) // Extract an integer sample.
			rnd &= mask                                 // Mask for non-power-of-two alphabet sizes.

			// --- If the random index is valid, select the rune and write to output. ---
			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				idBuffer[cursor] = g.config.runeAlphabet[rnd]
				cursor++
			}
		}
	}

	// --- If unable to generate a full ID within the allowed attempts, return an error. ---
	if cursor < length {
		return ErrExceededMaxAttempts
	}

	return nil
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
		// Fast path: Use a single byte as the random value.
		// Suitable for small alphabets (≤ 256) where 8 bits is sufficient.
		return uint(randomBytes[i])
	case 2:
		// Use 2 bytes to construct a 16-bit unsigned integer in Big Endian order.
		// Suitable for alphabets up to 2^16 unique characters.
		return uint(binary.BigEndian.Uint16(randomBytes[i : i+2]))
	case 4:
		// Use 4 bytes to construct a 32-bit unsigned integer in Big Endian order.
		// Suitable for very large alphabets, and provides a fast path for power-of-two sizing.
		return uint(binary.BigEndian.Uint32(randomBytes[i : i+4]))
	default:
		// General case: Combine `bytesNeeded` bytes into a single unsigned integer.
		// Shifts the result left by 8 bits on each iteration, then adds the next byte.
		// This path is used when the alphabet size requires a custom number of bytes.
		var rnd uint
		for j := 0; j < int(g.config.bytesNeeded); j++ {
			rnd = (rnd << 8) | uint(randomBytes[i+j])
		}
		return rnd
	}
}
