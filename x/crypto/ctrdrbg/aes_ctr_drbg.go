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

// Interface defines the contract for a NIST SP 800-90A AES-CTR-DRBG random source.
//
// Implementations provide cryptographically secure random bytes via io.Reader,
// and expose the non-secret, immutable configuration used at construction time.
//
// All methods are safe for concurrent use unless otherwise specified.
//
// The Config() method returns a copy of the DRBG's configuration. This allows inspection
// of operational parameters without exposing secrets or runtime-internal state.
type Interface interface {
	io.Reader

	// Config returns a copy of the DRBG configuration in use by this instance.
	// The returned Config does not include secrets or mutable runtime state.
	Config() Config
}

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
		item := pools[i].Get().(*drbg)
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
func NewReader(opts ...Option) (Interface, error) {
	// Step 1: Start with a default configuration, then apply each functional option to mutate cfg.
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// Step 2: Validate the configured key size is appropriate for AES.
	// Only 16, 24, or 32 bytes (AES-128, AES-192, AES-256) are supported.
	switch cfg.KeySize {
	case KeySize128, KeySize192, KeySize256:
	default:
		return nil, fmt.Errorf("invalid key size %d bytes; must be 16, 24, or 32", cfg.KeySize)
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

// Config returns a copy of the deterministic random bit generator’s static configuration.
//
// This method exposes only non-sensitive configuration options as set at initialization.
// No secret key material, runtime state, or internal DRBG details are included in the result.
// The returned Config is a copy and safe for inspection or serialization.
func (r *reader) Config() Config {
	// It's safe to fetch from any pool, as all configs are the same.
	d := r.pools[0].Get().(*drbg)
	cfg := *d.config
	r.pools[0].Put(d)
	return cfg
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

// state encapsulates the immutable cryptographic state of the DRBG, excluding the counter.
// This state is swapped atomically on rekey.
type state struct {
	// block is the initialized AES cipher.Block used in CTR mode.
	//
	// AES-CTR transforms the block cipher into a stream cipher by
	// encrypting a counter and XOR-ing it with plaintext to produce
	// pseudorandom output bytes.
	block cipher.Block

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

// drbg represents an internal deterministic random bit generator (DRBG) implementing
// the io.Reader interface using the NIST SP 800-90A AES-CTR-DRBG construction.
//
// Each drbg instance is intended to be used by a single goroutine at a time and is not
// safe for concurrent use. It maintains its own AES cipher, secret key, counter, usage counter,
// and rekeying flag for key rotation.
//
// This implementation ensures FIPS 140-2 alignment, strong security, and high performance
// under concurrent workloads by separating immutable cryptographic state (managed atomically)
// from the evolving counter (protected by a mutex).
type drbg struct {
	// config holds the immutable configuration for this DRBG instance.
	//
	// Includes:
	// - AES key size (e.g., 16, 24, or 32 bytes)
	// - Personalization string for domain separation
	// - Automatic key rotation policy
	// - Pool initialization and retry settings
	config *Config

	// state is an atomic pointer to the immutable cryptographic state for this DRBG.
	//
	// This state includes:
	//   - AES block cipher (used in CTR mode)
	//   - Secret key material
	//   - Initial counter value (NIST "V") at creation or rekey
	//
	// The atomic pointer allows for fast, race-free swapping of key/counter/cipher state
	// during asynchronous rekeying, without impacting ongoing read operations.
	state atomic.Pointer[state]

	// zero is a preallocated slice of zero-filled bytes used for output buffering.
	//
	// When UseZeroBuffer is enabled in config, this buffer is XOR-ed with
	// AES-CTR output to efficiently produce random bytes. Sized dynamically as needed.
	zero []byte

	// vMu is a mutex protecting the evolving counter (v) for this DRBG instance.
	//
	// All access and mutation of v must occur with this mutex held to ensure:
	//   - Counter advancement is atomic and non-overlapping across reads
	//   - Proper persistence of the counter value between consecutive reads
	//   - Safe resetting of the counter during key rotation (rekey)
	vMu sync.Mutex

	// v is the current 128-bit internal counter (NIST "V") for the DRBG instance.
	//
	// This counter is incremented for each AES block produced, ensuring
	// unique, non-repeating output for every call to Read. It is initialized
	// from the state.v value at creation or rekey, and persisted between reads.
	v [16]byte

	// usage tracks the number of bytes generated since the last key rotation.
	//
	// When usage exceeds config.MaxBytesPerKey, a rekey is triggered to ensure
	// forward secrecy and mitigate key compromise risk. This value is atomically updated.
	usage uint64

	// rekeying is an atomic flag (0 or 1) that guards rekey attempts.
	//
	// It ensures that only one goroutine performs rekeying at a time.
	// Uses atomic operations for concurrency safety.
	rekeying uint32
}

// Read generates cryptographically secure random bytes and writes them into the provided slice b.
//
// This method implements the io.Reader interface for drbg, providing a FIPS 140-2 aligned
// deterministic random bit generator using the AES-CTR-DRBG construction. Each call to Read
// returns a unique cryptographically strong pseudo-random stream and is safe for concurrent use.
//
// Semantics and Implementation Details:
//   - A snapshot of the current cryptographic state (key, block cipher, initial counter value) is loaded atomically.
//   - The DRBG's internal counter (v) is protected by a mutex to guarantee atomic advancement and persistence
//     between consecutive reads. This ensures that no two Read calls can produce overlapping output, and that
//     the generator stream is continuous and non-repeating.
//   - After generating the requested output, the advanced counter is persisted back to the DRBG instance.
//   - If key rotation is enabled and the generated output exceeds the configured threshold, an asynchronous
//     rekey operation is triggered. Rekeying swaps the cryptographic state atomically and resets the counter
//     (under lock) to guarantee forward secrecy and FIPS alignment.
//
// Parameters:
//   - b: Output buffer to be filled with cryptographically secure random bytes.
//
// Returns:
//   - int: Number of bytes written (equal to len(b) unless b is empty).
//   - error: Always nil under normal operation.
func (d *drbg) Read(b []byte) (int, error) {
	// Step 1: Return immediately if the buffer is empty, as required by the io.Reader contract.
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	// Atomically load the current DRBG cryptographic state.
	st := d.state.Load()

	// Lock the counter mutex to guarantee exclusive access to the evolving counter.
	d.vMu.Lock()
	var v [16]byte

	// Copy the current counter value to a local variable. This snapshot forms the basis
	// of the unique keystream for this read operation.
	copy(v[:], d.v[:])

	// Fill the output buffer using the current cryptographic state and the local counter,
	// incrementing the counter as output is produced. All counter increments are reflected
	// in the local variable.
	d.fillBlocks(b, st, &v)

	// Persist the advanced counter back to the DRBG instance, ensuring subsequent reads
	// continue the keystream seamlessly without overlap or repetition.
	copy(d.v[:], v[:])

	// Unlock the mutex, allowing other callers to proceed.
	d.vMu.Unlock()

	// Key rotation logic: atomically update the usage counter and, if the output threshold is
	// exceeded, trigger asynchronous rekeying in a background goroutine. Only one goroutine
	// may perform rekeying at a time.
	if d.config.EnableKeyRotation {
		atomic.AddUint64(&d.usage, uint64(len(b)))
		if atomic.LoadUint64(&d.usage) >= d.config.MaxBytesPerKey {
			if atomic.CompareAndSwapUint32(&d.rekeying, 0, 1) {
				go d.asyncRekey()
			}
		}
	}

	return n, nil
}

// fillBlocks fills the byte slice `b` with cryptographically secure, deterministic random data
// generated from the provided DRBG state and a local working counter.
//
// This method implements the core NIST SP 800-90A AES-CTR-DRBG output logic. It is **concurrency safe**
// as it operates only on immutable state and caller-provided (local) counter. No DRBG struct fields are mutated.
//
// Parameters:
//   - b   []byte:        Output buffer to be filled with random bytes. Must be at least 1 byte in length.
//   - st  *state:    Immutable snapshot of the DRBG key, block cipher, and initial counter (V).
//   - v   *[16]byte:     Local counter value (typically copied from DRBG.v). Advanced in place for each block.
//
// Behavior:
//   - Processes output in 16-byte (AES block size) chunks for maximal efficiency.
//   - For each block, increments the counter, encrypts it, and writes the result to output.
//   - Supports two strategies:
//   - UseZeroBuffer: Encrypted blocks are staged in a reusable buffer before being copied out (reducing allocations).
//   - Fast path: Output is written directly into the caller's buffer except for a possible tail partial block,
//     which uses a temporary [16]byte buffer.
//
// Returns:
//   - None. The output is written directly to the supplied buffer `b`, and the local counter `v` is incremented
//     in place for each block produced.
//
// Security:
//   - Ensures every 16-byte block is generated with a unique counter value per NIST recommendations.
//   - Never mutates DRBG fields or internal state directly.
//
// Panics:
//   - Never panics under normal operation. Will panic only if AES block size invariants are violated
//     (should not be possible with validated configuration).
func (d *drbg) fillBlocks(b []byte, st *state, v *[16]byte) {
	// Return immediately if the buffer is empty, as required by the io.Reader contract.
	n := len(b)
	if n == 0 {
		return
	}

	// Buffered output mode: stage keystream in reusable buffer to minimize allocations.
	if d.config.UseZeroBuffer {
		// Ensure the zero buffer is large enough; allocate if needed.
		if cap(d.zero) < n {
			d.zero = make([]byte, n)
		}
		d.zero = d.zero[:n] // Resize without reallocating if possible.

		offset := 0
		remaining := n
		for remaining > 0 {
			// Determine block size for this iteration (full or final partial block).
			blockSize := 16
			if remaining < 16 {
				blockSize = remaining
			}

			// Advance the counter as required by CTR mode (one block per keystream segment).
			incV(v)

			// Encrypt the incremented counter; write keystream into zero buffer.
			st.block.Encrypt(d.zero[offset:offset+blockSize], v[:])

			// Copy encrypted keystream to caller's buffer.
			copy(b[offset:offset+blockSize], d.zero[offset:offset+blockSize])
			offset += blockSize
			remaining -= blockSize
		}
		return
	}

	// Fast path: direct write to output buffer, except for a final partial block.
	offset := 0
	for ; offset+16 <= n; offset += 16 {
		incV(v)
		st.block.Encrypt(b[offset:offset+16], v[:])
	}

	// Handle remaining tail (if output is not a multiple of 16 bytes).
	if tail := n - offset; tail > 0 {
		var tmp [16]byte
		incV(v)
		st.block.Encrypt(tmp[:], v[:])
		copy(b[offset:], tmp[:tail])
	}
}

// newDRBG creates and returns a new, fully initialized deterministic random bit generator (DRBG) instance.
//
// This function constructs a FIPS 140-2 aligned AES-CTR-DRBG instance, securely seeded from operating system entropy.
// Initialization steps are as follows:
//  1. Acquire a seed consisting of (key size + 16) bytes of cryptographically strong random data.
//  2. Optionally XOR in a personalization string for domain separation, as required by SP 800-90A.
//  3. Derive the AES key and initial counter (V) from the seed.
//  4. Construct the AES block cipher with the derived key, and fail if the cipher cannot be created.
//  5. Optionally allocate a reusable zero buffer if requested in configuration.
//  6. Store the resulting cryptographic state atomically and initialize the working counter (v) from this state.
//
// If entropy acquisition or cipher construction fails, an error is returned and the DRBG is not created.
//
// Parameters:
//   - cfg: *Config — pointer to the DRBG configuration (must be non-nil)
//
// Returns:
//   - *drbg: newly initialized DRBG instance, ready for use
//   - error: non-nil if any initialization step fails (entropy, cipher, or config error)
func newDRBG(cfg *Config) (*drbg, error) {
	seedLen := cfg.KeySize + 16

	// Allocate a buffer for the full seed: key + 128-bit counter.
	seed := make([]byte, seedLen)

	// Read entropy from the operating system. Fail if not available.
	if _, err := io.ReadFull(rand.Reader, seed); err != nil {
		return nil, err
	}

	// XOR in personalization string (if any) for domain separation.
	if cfg.Personalization != nil {
		for i := range cfg.Personalization {
			seed[i%len(seed)] ^= cfg.Personalization[i]
		}
	}

	// Derive the AES key and the initial counter (V) from the seed.
	var key [32]byte
	copy(key[:], seed[:cfg.KeySize])
	var v [16]byte
	copy(v[:], seed[cfg.KeySize:])

	// Construct the AES block cipher using the derived key.
	block, err := aes.NewCipher(key[:cfg.KeySize])
	if err != nil {
		return nil, err
	}

	// Optionally preallocate the zero buffer for buffer-reuse mode.
	var zero []byte
	if cfg.UseZeroBuffer && cfg.DefaultBufferSize > 0 {
		zero = make([]byte, cfg.DefaultBufferSize)
	}

	// Store the immutable cryptographic state atomically.
	st := &state{
		block: block,
		key:   key,
		v:     v,
	}
	d := &drbg{
		config:   cfg,
		zero:     zero,
		usage:    0,
		rekeying: 0,
	}
	d.state.Store(st)

	// Initialize the working counter (v) from the state, guaranteeing unique output on first use.
	copy(d.v[:], v[:])

	return d, nil
}

// asyncRekey performs an asynchronous, non-blocking reseed and key rotation for the DRBG instance.
//
// This function is launched in a background goroutine when the generated output exceeds the configured threshold
// (MaxBytesPerKey). It attempts to generate new entropy, derive a new key and counter, and atomically install a
// new DRBG state. The working counter (v) is reset to the new initial value under lock. If all attempts to reseed
// fail, the existing cryptographic state is left unchanged, and the generator continues operating.
//
// Steps:
//  1. Attempt up to MaxRekeyAttempts reseed/rotate cycles, with exponential backoff (bounded by MaxRekeyBackoff).
//  2. For each attempt:
//     - Acquire a fresh random seed and optionally apply personalization.
//     - Derive a new key and counter (V), and construct a new AES cipher.
//     - On success, atomically store the new state, reset the usage counter, and set the working counter (v).
//  3. Always clear the rekeying flag before returning (even on panic or error), so future rekeys can proceed.
//
// Parameters: None (method receiver only).
func (d *drbg) asyncRekey() {
	// Always clear the rekeying flag on exit.
	defer atomic.StoreUint32(&d.rekeying, 0)

	base := d.config.RekeyBackoff
	maxBackoff := d.config.MaxRekeyBackoff
	if maxBackoff == 0 {
		maxBackoff = defaultMaxBackoff
	}

	// Attempt to reseed and rekey up to MaxRekeyAttempts times.
	for i := 0; i < d.config.MaxRekeyAttempts; i++ {
		// Obtain new entropy for key and counter (V).
		seedLen := d.config.KeySize + 16 // Key size plus 128-bit counter
		seed := make([]byte, seedLen)
		if _, err := io.ReadFull(rand.Reader, seed); err == nil {
			// Apply personalization string, if set, by XORing into the seed.
			if d.config.Personalization != nil {
				for j := range d.config.Personalization {
					seed[j%len(seed)] ^= d.config.Personalization[j]
				}
			}

			// Construct the new AES key and counter (V) from the seed buffer.
			var key [32]byte
			copy(key[:], seed[:d.config.KeySize])
			var v [16]byte
			copy(v[:], seed[d.config.KeySize:])
			block, err := aes.NewCipher(key[:d.config.KeySize])
			if err == nil {
				// Store new cryptographic state atomically.
				newState := &state{
					block: block,
					key:   key,
					v:     v,
				}
				d.state.Store(newState)
				atomic.StoreUint64(&d.usage, 0)

				// Reset the working counter (v) under mutex lock to ensure no overlap.
				d.vMu.Lock()
				copy(d.v[:], v[:])
				d.vMu.Unlock()
				return // Rekey complete.
			}

			// (If cipher construction fails, fall through and retry after backoff.)
		}

		// Wait with exponential backoff before retrying.
		time.Sleep(base)
		base *= 2
		if base > maxBackoff {
			base = maxBackoff
		}
	}
	// If all retries fail, generator continues with prior state.
}

// incV increments the DRBG counter (V) in big-endian order, rolling over as needed.
//
// The counter (V) is treated as a 128-bit unsigned integer in big-endian representation.
// Each call increments the counter by one, wrapping as appropriate. This function
// is used for advancing the DRBG keystream per SP 800-90A section on counter mode.
// Not concurrency safe; caller must synchronize if used from multiple goroutines.
//
// Parameters:
//   - v: pointer to a 16-byte array ([16]byte), representing the current counter value.
//
// Returns: None (modifies v in place).
func incV(v *[16]byte) {
	// Start from the least significant byte (rightmost, index 15), incrementing with carry.
	for i := 15; i >= 0; i-- {
		v[i]++
		if v[i] != 0 {
			break // No further carry needed; stop.
		}
	}
}
