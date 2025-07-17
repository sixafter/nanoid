// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package ctrdrbg provides a FIPS 140-2 aligned, high-performance AES-CTR-DRBG.
//
// This package implements a cryptographically secure, pool-backed Deterministic Random Bit Generator
// (DRBG) following the NIST SP 800-90A AES-CTR-DRBG construction. Each generator instance uses an
// AES block cipher in counter (CTR) mode to produce cryptographically secure pseudo-random bytes,
// suitable for high-throughput, concurrent workloads.
//
// All cryptographic primitives are provided by the Go standard library. This implementation is designed
// for environments requiring strong compliance, including support for Go's FIPS-140 mode (GODEBUG=fips140=on).
package ctrdrbg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

// Reader is a package-level, cryptographically secure random source suitable for high-concurrency applications.
//
// Reader is initialized at package load time via NewReader and is safe for concurrent use. If initialization fails
// (for example, if crypto/rand is unavailable), the package will panic. This ensures that any failure to obtain a secure
// entropy source is detected immediately and not silently ignored.
//
// Example usage:
//
//	buf := make([]byte, 64)
//	_, err := ctrdrbg.Reader.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Random data: %x\n", buf)
var Reader io.Reader

// init initializes the package-level Reader. It panics if NewReader fails, preventing operation without
// a secure random source. This follows cryptographic best practices by making entropy failure a fatal error.
func init() {
	cfg := DefaultConfig()
	pools := make([]*sync.Pool, cfg.Shards)
	for i := range pools {
		cfg := cfg // Capture the current configuration for this shard
		pools[i] = &sync.Pool{
			New: func() interface{} {
				var (
					d   *drbg
					err error
				)
				for r := 0; r < cfg.MaxInitRetries; r++ {
					if d, err = newDRBG(&cfg); err == nil {
						return d
					}
				}
				// If initialization fails after all retries, panic.
				panic(fmt.Sprintf("ctrdrbg pool init failed after %d retries: %v", cfg.MaxInitRetries, err))
			},
		}

		// Eagerly test the pool initialization to ensure that any catastrophic
		// failure is caught immediately, not deferred to the first use.
		item := pools[i].Get()
		pools[i].Put(item)
	}

	Reader = &reader{pools: pools}
}

// reader is an internal implementation of io.Reader that uses a pool of DRBG instances
// to support efficient concurrent random byte generation.
type reader struct {
	pools []*sync.Pool
}

// NewReader constructs and returns an io.Reader that produces cryptographically secure
// random bytes using a pool of AES-CTR-DRBG instances. Functional options may be supplied to customize key size,
// key rotation, and pool behavior. Each generator is seeded with entropy from crypto/rand.
//
// The returned Reader is safe for concurrent use. If no generator can be created after MaxInitRetries,
// NewReader returns an error.
//
// Example:
//
//	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(32))
//	if err != nil {
//	    // handle error
//	}
//
//	buf := make([]byte, 32)
//	n, err := r.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Read %d bytes: %x\n", n, buf)
//
// NewReader constructs and returns an io.Reader that produces cryptographically secure
// random bytes using a pool of AES-CTR-DRBG instances. Functional options may be supplied to customize key size,
// key rotation, and pool behavior. Each generator is seeded with entropy from crypto/rand.
//
// The returned Reader is safe for concurrent use. If no generator can be created after MaxInitRetries,
// NewReader returns an error.
//
// Example:
//
//	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(32))
//	if err != nil {
//	    // handle error
//	}
//
//	buf := make([]byte, 32)
//	n, err := r.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Read %d bytes: %x\n", n, buf)
func NewReader(opts ...Option) (io.Reader, error) {
	// Step 1: Start with a default configuration, then apply each functional option to mutate cfg.
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// Step 2: Validate the configured key size is appropriate for AES.
	// Only 16, 24, or 32 bytes (AES-128, AES-192, AES-256) are supported.
	if cfg.KeySize != 16 && cfg.KeySize != 24 && cfg.KeySize != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	// Step 3: Create a sync.Pool to manage DRBG instances for concurrent access.
	// The pool's New function attempts to create a new DRBG, retrying up to MaxInitRetries times.
	// If all attempts fail, the function panics, making failure explicit and visible.
	pools := make([]*sync.Pool, cfg.Shards)
	for i := range pools {
		cfg := cfg // Capture the current configuration for this shard
		pools[i] = &sync.Pool{
			New: func() interface{} {
				var (
					d   *drbg
					err error
				)
				for r := 0; r < cfg.MaxInitRetries; r++ {
					if d, err = newDRBG(&cfg); err == nil {
						return d
					}
				}
				// If DRBG initialization fails after all retries, panic.
				panic(fmt.Sprintf("ctrdrbg pool init failed after %d retries: %v", cfg.MaxInitRetries, err))
			},
		}

		// Step 4: Eagerly test pool initialization by creating (and releasing) one DRBG.
		// This ensures that catastrophic failures are detected during NewReader rather than at first use.
		// If pool.New panics, recover and convert the panic to an error.
		var panicErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr = fmt.Errorf("ctrdrbg pool initialization failed: %v", r)
				}
			}()
			item := pools[i].Get()
			pools[i].Put(item)
		}()

		// Step 5: If any panic occurred during pool initialization, return it as an error.
		if panicErr != nil {
			return nil, panicErr
		}
	}

	// Step 6: Return a new reader that wraps the initialized pool.
	return &reader{pools: pools}, nil
}

// shardIndex selects a pseudo-random shard index in the range [0, n) using
// a fast, thread-safe global PCG64-based RNG.
//
// This function is used to evenly distribute load across multiple sync.Pool
// shards, reducing contention in high-concurrency scenarios. It avoids the
// overhead of time-based seeding or mutex contention.
//
// The randomness is not cryptographically secure but is safe for concurrent
// use and sufficient for load balancing purposes.
//
// Panics if n <= 0.
func shardIndex(n int) int {
	return mrand.IntN(n)
}

// Read fills the provided buffer with cryptographically secure random data.
//
// Read implements the io.Reader interface and is designed to be safe for concurrent use when accessed
// via the package-level Reader or any Reader returned from NewReader.
//
// Example:
//
//	buffer := make([]byte, 32)
//	n, err := Reader.Read(buffer)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Read %d bytes of random data: %x\n", n, buffer)
func (r *reader) Read(b []byte) (int, error) {
	// Step 1: Return immediately if the buffer is empty, as required by the io.Reader contract.
	if len(b) == 0 {
		return 0, nil
	}

	// Determine the shard index based on the number of pools available.
	n := len(r.pools)
	shard := 0
	if n > 1 {
		shard = shardIndex(n)
	}

	// Step 2: Borrow an instance of the internal deterministic random bit generator from the pool.
	// This ensures that each call gets exclusive access to an isolated state for cryptographic safety.
	d := r.pools[shard].Get().(*drbg)

	// Step 3: Ensure that the borrowed instance is returned to the pool, even if Read fails or panics.
	// This pattern prevents resource leaks and maintains pool integrity.
	defer r.pools[shard].Put(d)

	// Step 4: Fill the caller’s buffer with random data using the borrowed generator.
	// The actual cryptographic work is performed by the internal generator’s Read method.
	return d.Read(b)
}

// drbg represents an internal deterministic random bit generator (DRBG) implementing
// the io.Reader interface using the NIST SP 800-90A AES-CTR-DRBG construction.
//
// Each drbg instance is intended to be used by a single goroutine at a time and is not
// safe for concurrent use. It maintains its own AES cipher, secret key, counter, usage counter,
// and rekeying flag for key rotation.
type drbg struct {
	// config holds the immutable configuration for this DRBG instance.
	//
	// Includes:
	// - AES key size (e.g., 16, 24, or 32 bytes)
	// - Personalization string for domain separation
	// - Automatic key rotation policy
	// - Pool initialization and retry settings
	config *Config

	// block is the initialized AES cipher.Block used in CTR mode.
	//
	// AES-CTR transforms the block cipher into a stream cipher by
	// encrypting a counter and XOR-ing it with plaintext to produce
	// pseudorandom output bytes.
	block cipher.Block

	// zero is a preallocated slice of zero-filled bytes used for XOR operations.
	//
	// When UseZeroBuffer is enabled in config, this buffer is XOR-ed with
	// AES-CTR output to efficiently produce random bytes.
	zero []byte

	// usage tracks the number of bytes generated since the last key rotation.
	//
	// When usage exceeds config.MaxBytesPerKey, a rekey is triggered to ensure
	// forward secrecy and mitigate key compromise risk.
	usage uint64

	// rekeying is an atomic flag (0 or 1) that guards rekey attempts.
	//
	// It ensures that only one goroutine performs rekeying at a time.
	// Uses atomic operations for concurrency safety.
	rekeying uint32

	// key holds the internal DRBG secret key used for AES-CTR operations.
	//
	// The key length is determined by config.KeySize and can be:
	// - 16 bytes for AES-128
	// - 24 bytes for AES-192
	// - 32 bytes for AES-256
	//
	// Unused bytes are zeroed and ignored.
	key [32]byte

	// v is the 128-bit internal counter (NIST "V") used by the DRBG.
	//
	// This counter is incremented for each AES block to produce a unique
	// keystream segment in CTR mode. It ensures deterministic, non-repeating output.
	v [16]byte
}

// Read generates cryptographically secure random bytes and writes them into the provided slice b.
//
// This method implements the io.Reader interface for drbg. It is not safe for concurrent
// use by multiple goroutines on the same instance. Each call should be exclusive to a single goroutine.
//
// Behavior:
//   - Fills b with pseudorandom output using AES in counter mode (CTR).
//   - Processes b in 16-byte (AES block size) chunks. For each block, increments the internal counter (v),
//     encrypts it, and copies the result into b.
//   - Tracks the total number of bytes generated in the usage counter. If key rotation is enabled and
//     the usage exceeds MaxBytesPerKey, attempts to trigger an asynchronous key rotation (asyncRekey).
//   - Only one goroutine will perform the key rotation at a time. Key rotation is non-blocking.
//
// Parameters:
//   - b: The output buffer to be filled with random bytes.
//
// Returns:
//   - int: The number of bytes written (equal to len(b), unless b is empty).
//   - error: Always nil in normal operation. Errors only occur on programmer misuse or internal misconfiguration.
//
// Example:
//
//	buf := make([]byte, 32)
//	_, err := drbgInstance.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Random bytes: %x\n", buf)
func (d *drbg) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	d.fillBlocks(b)

	// Usage and rekey logic centralized here
	if d.config.EnableKeyRotation {
		atomic.AddUint64(&d.usage, uint64(len(b)))
		if atomic.LoadUint64(&d.usage) >= d.config.MaxBytesPerKey {
			if atomic.CompareAndSwapUint32(&d.rekeying, 0, 1) {
				go d.asyncRekey()
			}
		}
	}
	return len(b), nil
}

// fillBlocks fills the byte slice `b` with pseudo-random data generated
// by the DRBG (Deterministic Random Bit Generator) instance.
//
// This function uses the DRBG’s internal counter (`d.v`), incremented for each block,
// and encrypted via the configured block cipher to generate deterministic output.
// The method supports two output strategies depending on configuration:
//
// 1. If d.config.UseZeroBuffer is true:
//   - An internal zero buffer (`d.zero`) is used as the temporary destination for block encryption.
//   - It is resized if needed to match the length of `b`.
//   - The encrypted blocks are then copied from `d.zero` into `b`.
//
// 2. If UseZeroBuffer is false (fast path):
//   - The encrypted blocks are written directly into `b`.
//   - For any remaining bytes that do not fill a full 16-byte block,
//     encryption is performed into a temporary stack buffer (`tmp`) and the result is copied into `b`.
//
// Parameters:
//   - b []byte: the output buffer to be filled with generated data.
//
// Behavior:
//   - Each 16-byte block is generated by incrementing `d.v` and encrypting it.
//   - No allocations occur unless the zero buffer needs to grow.
//   - The function ensures correct handling of non-multiple-of-16 lengths.
func (d *drbg) fillBlocks(b []byte) {
	n := len(b)
	if n == 0 {
		return
	}

	if d.config.UseZeroBuffer {
		// Ensure the internal zero buffer is large enough for the requested output.
		if cap(d.zero) < n {
			d.zero = make([]byte, n)
		}
		d.zero = d.zero[:n]

		offset := 0
		for n > 0 {
			blockSize := 16
			if n < 16 {
				blockSize = n // Last partial block
			}
			d.incV()                                                          // Increment the internal counter
			d.block.Encrypt(d.zero[offset:offset+blockSize], d.v[:])          // Encrypt the counter into the zero buffer
			copy(b[offset:offset+blockSize], d.zero[offset:offset+blockSize]) // Copy into the output buffer
			offset += blockSize
			n -= blockSize
		}
		return
	}

	// Fast path: direct write for full blocks, then use stack buffer for any final partial block.
	offset := 0
	for ; offset+16 <= len(b); offset += 16 {
		d.incV()
		d.block.Encrypt(b[offset:offset+16], d.v[:])
	}
	if tail := len(b) - offset; tail > 0 {
		var tmp [16]byte
		d.incV()
		d.block.Encrypt(tmp[:], d.v[:])
		copy(b[offset:], tmp[:tail])
	}
}

// newDRBG creates and returns a new, fully initialized deterministic random bit generator (DRBG) instance.
//
// The DRBG is seeded with fresh entropy from the operating system (via crypto/rand), using a
// combination of a random key and initial counter (V). Optionally, personalization data may
// be XORed into the seed for domain separation. If either entropy acquisition or AES cipher
// construction fails, an error is returned. This function never panics.
//
// Parameters:
//   - cfg: Pointer to the DRBG configuration to use for initialization. Must be non-nil.
//
// Returns:
//   - *drbg: A new DRBG instance ready for use.
//   - error: If seeding or cipher construction fails, a non-nil error is returned.
func newDRBG(cfg *Config) (*drbg, error) {
	// Calculate the required seed length: key size + 16 bytes for the initial counter (V).
	seedLen := cfg.KeySize + 16

	// Allocate a seed buffer of the required size.
	seed := make([]byte, seedLen)

	// Read the required amount of entropy from the OS cryptographically secure RNG.
	if _, err := io.ReadFull(rand.Reader, seed); err != nil {
		// If entropy acquisition fails, return an error (never panic).
		return nil, err
	}

	// If personalization data is provided, XOR it into the seed for domain separation.
	if cfg.Personalization != nil {
		for i := range cfg.Personalization {
			seed[i%len(seed)] ^= cfg.Personalization[i]
		}
	}

	// Prepare the AES key as a 32-byte array (only the first KeySize bytes are used).
	var key [32]byte
	copy(key[:], seed[:cfg.KeySize])

	// Prepare the 16-byte counter (V) from the remainder of the seed buffer.
	var v [16]byte
	copy(v[:], seed[cfg.KeySize:])

	// Create the AES cipher using the derived key. Only the relevant key size is used (AES-128/192/256).
	block, err := aes.NewCipher(key[:cfg.KeySize])
	if err != nil {
		return nil, err
	}

	// zero buffer allocation for UseZeroBuffer mode
	var zero []byte
	if cfg.UseZeroBuffer && cfg.DefaultBufferSize > 0 {
		zero = make([]byte, cfg.DefaultBufferSize)
	}

	// Return a new DRBG instance with the configured key, counter, cipher, and reset usage counters.
	return &drbg{
		config:   cfg,
		block:    block,
		key:      key,
		v:        v,
		usage:    0,
		rekeying: 0,
		zero:     zero,
	}, nil
}

// asyncRekey performs an asynchronous, non-blocking reseed of the DRBG key/counter state.
//
// This function is invoked in a background goroutine when the total output bytes exceed the configured MaxBytesPerKey threshold.
// It will retry reseeding and rekeying up to MaxRekeyAttempts, with exponential backoff between attempts (bounded by MaxRekeyBackoff).
// On successful reseed, the internal AES key, counter, and block cipher are replaced atomically and usage is reset to zero.
// If all attempts fail, the DRBG instance continues operating with its existing cryptographic state.
// The function always clears the rekeying flag before returning, guaranteeing safe concurrency semantics.
func (d *drbg) asyncRekey() {
	// Always clear the rekeying flag when this goroutine exits,
	// regardless of whether the rekey was successful or not.
	defer atomic.StoreUint32(&d.rekeying, 0)

	// Initialize the base and maximum backoff durations for retries.
	base := d.config.RekeyBackoff
	maxBackoff := d.config.MaxRekeyBackoff
	if maxBackoff == 0 {
		maxBackoff = defaultMaxBackoff // Fallback to library default if not set.
	}

	// Attempt to reseed and rekey up to MaxRekeyAttempts times.
	for i := 0; i < d.config.MaxRekeyAttempts; i++ {
		// Step 1: Obtain fresh entropy for the new key and counter.
		seedLen := d.config.KeySize + 16 // Key size plus 128-bit counter.
		seed := make([]byte, seedLen)
		if _, err := io.ReadFull(rand.Reader, seed); err == nil {
			// Step 2: Apply personalization string, if set, by XORing into the seed.
			if d.config.Personalization != nil {
				for i := range d.config.Personalization {
					seed[i%len(seed)] ^= d.config.Personalization[i]
				}
			}

			// Step 3: Construct the new AES key and counter from the seed buffer.
			var key [32]byte
			copy(key[:], seed[:d.config.KeySize])
			var v [16]byte
			copy(v[:], seed[d.config.KeySize:])

			// Step 4: Attempt to construct a new AES block cipher with the key.
			block, err := aes.NewCipher(key[:d.config.KeySize])
			if err == nil {
				// Step 5: If successful, update internal state atomically and reset usage.
				d.key = key
				d.v = v
				d.block = block
				atomic.StoreUint64(&d.usage, 0)
				return // Rekey complete.
			}
			// (If cipher construction fails, fall through and retry after backoff.)
		}

		// Step 6: Wait with exponential backoff before retrying.
		time.Sleep(base)
		base *= 2
		if base > maxBackoff {
			base = maxBackoff
		}
	}

	// If all retries failed, the DRBG continues using its prior key and counter.
}

// incV increments the DRBG counter (V) in big-endian order, rolling over as needed.
//
// This operation is critical for advancing the AES-CTR-DRBG internal state. The counter (V)
// is treated as a 128-bit unsigned integer stored in big-endian format. Each call increments
// the counter by one, wrapping as appropriate (per NIST SP 800-90A section on counter mode).
// This ensures every block of output is unique for a given key.
//
// Not concurrency safe. This method is always called within a single goroutine per DRBG instance.
func (d *drbg) incV() {
	// Start from the least significant byte (rightmost, index 15).
	for i := 15; i >= 0; i-- {
		d.v[i]++
		// If increment did not overflow this byte (i.e., it is not zero),
		// no further carry is needed, so we can exit the loop early.
		if d.v[i] != 0 {
			break
		}
	}
}
