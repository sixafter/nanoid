// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package prng provides a cryptographically secure pseudo-random number generator (PRNG)
// that implements the io.Reader interface. It is designed for high-performance, concurrent
// use in generating random bytes.
//
// This package is part of the experimental "x" modules and may be subject to change.
package prng

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mrand "math/rand/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/chacha20"
)

// Reader is a global, cryptographically secure random source.
// It is initialized at package load time and is safe for concurrent use.
// If initialization fails (e.g., crypto/rand is unavailable), the package will panic.
//
// Example usage:
//
//	buffer := make([]byte, 64)
//	n, err := Reader.Read(buffer)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Read %d bytes of random data: %x\n", n, buffer)
var Reader io.Reader

// Interface defines the contract for a ChaCha20-based cryptographically secure
// pseudorandom number generator (PRNG).
//
// Implementations of Interface provide a thread-safe source of cryptographically
// strong random bytes derived from the ChaCha20 stream cipher. Each implementation
// must also satisfy the io.Reader interface, making it compatible with standard
// Go APIs that consume randomness (e.g., encoding, crypto, and token generation).
//
// All methods are safe for concurrent use unless otherwise noted.
//
// The Config method allows callers to retrieve a copy of the immutable,
// non-secret configuration associated with the PRNG instance. This enables
// inspection of operational parameters—such as nonce, pool size, or reseed
// interval—without exposing any sensitive key material or mutable internal state.
type Interface interface {
	io.Reader

	// Config returns a copy of the PRNG configuration in effect for this source.
	//
	// The returned Config contains only non-secret, immutable parameters and
	// omits any runtime state or cryptographic keys. Callers may safely inspect
	// the returned value to determine operational behavior without risk of
	// secret exposure or race conditions.
	Config() Config
}

// init sets up the package‐level Reader by creating a new pooled PRNG instance.
// It is invoked automatically at program startup (package initialization).
// If NewReader fails (e.g., OS entropy unavailable), init will panic to prevent
// running without a secure random source.
//
// Panicking here is intentional and idiomatic for cryptographic primitives:
// it ensures that any critical failure in obtaining a secure entropy source
// is detected immediately and cannot be ignored.
func init() {
	cfg := DefaultConfig()
	pools := make([]*sync.Pool, cfg.Shards)
	for i := range pools {
		cfg := cfg // Capture the current configuration for this shard
		pools[i] = &sync.Pool{
			New: func() interface{} {
				var (
					p   *prng
					err error
				)
				for r := 0; r < cfg.MaxInitRetries; r++ {
					if p, err = newPRNG(&cfg); err == nil {
						return p
					}
				}
				// If initialization fails after all retries, panic.
				panic(fmt.Sprintf("prng pool init failed after %d retries: %v", cfg.MaxInitRetries, err))
			},
		}

		// Eagerly test the pool initialization to ensure that any catastrophic
		// failure is caught immediately, not deferred to the first use.
		item := pools[i].Get().(*prng)
		pools[i].Put(item)
	}

	Reader = &reader{pools: pools}
}

// reader wraps a sync.Pool of prng instances to provide an io.Reader
// that efficiently reuses ChaCha20-based PRNG objects.
// Each call to Read() pulls a prng from the pool, uses it to fill the
// provided buffer, and then returns it to the pool for future reuse.
//
// The Pool’s New function is responsible for creating and initializing
// each prng (including seeding and atomic cipher setup). This design
// minimizes allocations and contention on crypto/rand while ensuring
// each goroutine can obtain a fresh or recycled PRNG instance quickly.
type reader struct {
	config *Config
	pools  []*sync.Pool
}

// NewReader constructs and returns an io.Reader that produces cryptographically secure
// pseudo-random bytes using a pool of ChaCha20-based PRNG instances. Functional options may be
// supplied to customize pool behavior, key rotation, and other advanced settings.
//
// Each PRNG in the pool is seeded with a unique key and nonce from crypto/rand, and automatically
// rotates to a fresh key/nonce pair after emitting a configurable number of bytes (MaxBytesPerKey).
// The pool will retry PRNG initialization up to MaxInitRetries times and will panic if it cannot
// produce a valid generator after all attempts.
//
// The returned Reader is safe for concurrent use. If the pool cannot be initialized, NewReader
// returns an error.
//
// Example:
//
//	r, err := prng.NewReader()
//	if err != nil {
//	    // handle error
//	}
//	buf := make([]byte, 32)
//	n, err := r.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Read %d bytes: %x\n", n, buf)
func NewReader(opts ...Option) (Interface, error) {
	// Step 1: Start with a default configuration and apply each functional option to allow caller customization.
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// If n <= 0, the number of shards defaults to runtime.GOMAXPROCS(0),
	// which is useful in containerized environments.
	// See: https://github.com/golang/go/issues/73193
	if cfg.Shards <= 0 {
		cfg.Shards = runtime.GOMAXPROCS(0)
	}

	// Step 2: Construct a sync.Pool for managing reusable prng instances.
	// The pool's New function attempts to construct a new *prng,
	// retrying up to cfg.MaxInitRetries times in case of failure (e.g., low entropy).
	// If all attempts fail, the function returns nil, which is caught during eager initialization below.
	pools := make([]*sync.Pool, cfg.Shards)
	for i := range pools {
		cfg := cfg // Capture the current configuration for this shard
		pools[i] = &sync.Pool{
			New: func() interface{} {
				var (
					p   *prng
					err error
				)
				for r := 0; r < cfg.MaxInitRetries; r++ {
					if p, err = newPRNG(&cfg); err == nil {
						return p
					}
				}
				// If initialization fails after all retries, return nil instead of panicking.
				// The eager initialization step below will detect and return this as an error.
				return nil
			},
		}

		// Step 3: Eagerly test the pool initialization to ensure that any catastrophic
		// failure is caught immediately, not deferred to the first use.
		// This triggers pool.New, which may return nil on failure. Any nil value is converted to an error.
		var initErr error
		item := pools[i].Get()
		if item == nil {
			initErr = fmt.Errorf("prng pool initialization failed after %d retries", cfg.MaxInitRetries)
		} else {
			pools[i].Put(item)
		}

		// Step 4: If initialization failed, return it as an error.
		if initErr != nil {
			return nil, initErr
		}
	}

	// Step 5: Return a new reader that wraps the initialized pool. This is safe for concurrent use.
	return &reader{
		pools:  pools,
		config: &cfg,
	}, nil
}

// Config returns a copy of the PRNG's configuration settings.
//
// The returned configuration describes the PRNG’s static parameters as set during initialization.
// No secret values, seeds, or internal state are included. The returned Config is a safe copy
// for inspection, logging, or diagnostics, and cannot be used to alter the PRNG’s behavior.
func (r *reader) Config() Config {
	return *r.config
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
// Read implements the io.Reader interface. It is safe for concurrent use when accessed
// via the package-level Reader or any Reader returned from NewReader. Each call to Read
// borrows an independent PRNG instance from an internal pool, ensuring safe concurrent
// usage without shared mutable state.
//
// Example usage:
//
//	buffer := make([]byte, 32)
//	n, err := Reader.Read(buffer)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Read %d bytes of random data: %x\n", n, buffer)
func (r *reader) Read(b []byte) (int, error) {
	// Step 1: If the caller provided an empty buffer, return immediately (as per io.Reader contract).
	if len(b) == 0 {
		return 0, nil
	}

	// Determine the shard index based on the number of pools available.
	n := len(r.pools)
	shard := 0
	if n > 1 {
		shard = shardIndex(n)
	}

	// Step 2: Acquire a PRNG instance from the pool for exclusive use by this call.
	// This provides thread safety and isolation of cryptographic state.
	p := r.pools[shard].Get().(*prng)

	// Step 3: Always return the PRNG instance to the pool, even if an error occurs.
	// This ensures that the pool does not leak resources and stays available for future use.
	defer r.pools[shard].Put(p)

	// Step 4: Delegate the actual generation of random bytes to the PRNG instance's Read method.
	return p.Read(b)
}

// prng implements io.Reader using a ChaCha20 cipher stream and supports
// asynchronous, nonblocking rotation of the underlying key/nonce pair.
//
// Each instance maintains its own ChaCha20 cipher (stored atomically), a
// scratch buffer for encryption, and internal counters to enforce a
// “forward secrecy” rekey after a configurable output threshold.
type prng struct {
	// config holds a pointer to this PRNG instance’s configuration parameters.
	// It provides tunable settings such as MaxBytesPerKey (keystream rotation threshold)
	// and MaxInitRetries (how many times to retry initialization).
	config *Config

	// cipher holds the active *chacha20.Cipher. We use atomic.Value so that
	// loads and stores of the cipher pointer are safe and nonblocking.
	cipher atomic.Value

	// zero is a one‐off buffer of zeros used as plaintext for XORKeyStream.
	// We grow it as needed; since each prng is single‐goroutine‐owned from the pool,
	// no synchronization around this slice is required.
	zero []byte

	// usage tracks the total number of bytes output under the current key.
	// Once usage exceeds maxBytesPerKey, we trigger an asynchronous rekey.
	// This is incremented atomically in Read().
	usage uint64

	// rekeying is a 0/1 flag (set via atomic CAS) to ensure only one
	// background goroutine at a time performs the expensive rekey operation.
	rekeying uint32
}

// Read fills the provided byte slice `b` with cryptographically secure random data.
// It implements the `io.Reader` interface and is intended for exclusive use by a single goroutine
// per PRNG instance (each pool entry is not shared between goroutines).
//
// Internally, Read does the following:
//  1. Determines the length `n` of the requested output. If `n == 0`, returns immediately with no error.
//  2. Atomically loads the current ChaCha20 cipher stream from `p.cipher`.
//  3. If `UseZeroBuffer` is true, prepares a zero-valued buffer of length `n` in `p.zero`, growing it if necessary.
//  4. Calls `cipher.XORKeyStream(b, p.zero)` if `UseZeroBuffer` is true; otherwise, calls `cipher.XORKeyStream(b, b)` for in-place output.
//  5. If key usage tracking is enabled (`EnableKeyRotation`), atomically increments `p.usage` by `n` and, if the threshold is crossed,
//     attempts a single non-blocking CAS to set `p.rekeying` from 0→1; if successful, launches `p.asyncRekey()` in a background goroutine.
//  6. Returns the number of bytes written (`n`). Errors are only expected on internal cipher malfunction, which should not occur
//     under normal operation.
func (p *prng) Read(b []byte) (int, error) {
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	// Step 1: Atomically retrieve the active cipher stream.
	stream := p.cipher.Load().(*chacha20.Cipher)

	// Step 2: Generate random output based on configuration.
	if p.config.UseZeroBuffer {
		// Ensure internal zero buffer is at least n bytes.
		if cap(p.zero) < n {
			p.zero = make([]byte, n)
		} else {
			p.zero = p.zero[:n]
		}
		// XOR the zero buffer into b, producing random bytes.
		stream.XORKeyStream(b, p.zero)
	} else {
		// XOR the buffer into itself (in-place), producing random bytes.
		stream.XORKeyStream(b, b)
	}

	// Step 3: Optionally track key usage and trigger rekeying.
	if p.config.EnableKeyRotation {
		// Atomically increment usage counter by n bytes.
		atomic.AddUint64(&p.usage, uint64(n))
		// If usage exceeds threshold, attempt async rekey.
		if atomic.LoadUint64(&p.usage) > p.config.MaxBytesPerKey {
			if atomic.CompareAndSwapUint32(&p.rekeying, 0, 1) {
				go p.asyncRekey()
			}
		}
	}

	return n, nil
}

// newPRNG creates and returns a fully initialized prng instance.
//
// This function generates a fresh ChaCha20 cipher using a cryptographically secure random key and nonce,
// securely zeroes out any sensitive seed material, and stores the cipher in an atomic.Value for lock-free
// access by Read(). If configured, it preallocates a zero buffer for optimized XORKeyStream usage.
// Returns an error if cipher setup fails.
//
// Parameters:
//   - config: Pointer to the PRNG configuration. Must not be nil.
//
// Returns:
//   - *prng: A new PRNG instance ready for random output.
//   - error: A non-nil error if cipher construction fails.
func newPRNG(config *Config) (*prng, error) {
	// Generate a fresh ChaCha20 cipher seeded with secure random key and nonce.
	stream, err := newCipher()
	if err != nil {
		// If cipher construction fails, propagate the error to caller.
		return nil, err
	}

	// Optionally preallocate a zero buffer if UseZeroBuffer is set,
	// optimizing for repeated XORKeyStream operations.
	var zero []byte
	if config.UseZeroBuffer && config.DefaultBufferSize > 0 {
		zero = make([]byte, config.DefaultBufferSize)
	} else {
		zero = make([]byte, 0)
	}

	// Initialize the PRNG instance with the selected configuration and zero buffer.
	p := &prng{
		zero:   zero,
		config: config,
	}

	// Store the cipher stream atomically for lock-free, concurrent access in Read().
	p.cipher.Store(stream)

	// Return the initialized PRNG to the caller.
	return p, nil
}

// newCipher generates and returns a new *chacha20.Cipher seeded with a cryptographically secure
// random key and nonce.
//
// The function performs the following steps:
//  1. Allocates fresh buffers for the key and nonce of the correct size.
//  2. Fills both buffers with cryptographically secure random bytes from crypto/rand.Reader.
//  3. Constructs a new ChaCha20 stream cipher instance using the generated key and nonce.
//  4. Immediately overwrites (zeroes) the key and nonce buffers in memory to prevent any
//     sensitive seed material from lingering in process memory.
//  5. If any step fails (entropy acquisition or cipher construction), returns an error with context.
//     On success, returns the initialized cipher stream.
func newCipher() (*chacha20.Cipher, error) {
	// Step 1: Allocate key and nonce buffers according to ChaCha20 specification.
	key := make([]byte, chacha20.KeySize)
	nonce := make([]byte, chacha20.NonceSizeX)

	// Step 2: Fill the key buffer with cryptographically secure random bytes.
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("newCipher: failed to read key: %w", err)
	}

	// Step 3: Fill the nonce buffer with cryptographically secure random bytes.
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("newCipher: failed to read nonce: %w", err)
	}

	// Step 4: Attempt to construct a new ChaCha20 stream cipher instance.
	stream, err := chacha20.NewUnauthenticatedCipher(key, nonce)

	// Step 5: Immediately zero out the sensitive key and nonce buffers in memory.
	for i := range key {
		key[i] = 0
	}
	for i := range nonce {
		nonce[i] = 0
	}

	// Step 6: Check for errors in cipher construction and return as needed.
	if err != nil {
		return nil, fmt.Errorf("newCipher: unable to initialize cipher: %w", err)
	}
	return stream, nil
}

// asyncRekey performs an asynchronous, non-blocking rotation of the internal ChaCha20 cipher.
//
// This method is invoked when the PRNG's per-key usage threshold is exceeded. It runs in its own
// goroutine, and attempts to rekey the PRNG up to Config.MaxRekeyAttempts times, doubling the
// backoff after each failure (jittered by a random value for each attempt).
//
// On each attempt, the function captures the current cipher pointer so that it can zero out the
// old cipher (removing key/counter material) after a successful rekey. If all attempts fail,
// the function leaves the existing cipher in place and simply exits. The rekeying flag is always
// cleared before returning to allow future rekey attempts.
func (p *prng) asyncRekey() {
	// Always clear the rekeying flag when this goroutine exits, so rekey can be attempted again.
	defer atomic.StoreUint32(&p.rekeying, 0)

	// Start with the configured base backoff duration.
	base := p.config.RekeyBackoff

	// Track the previous cipher pointer for secure wiping after rotation.
	var old *chacha20.Cipher

	// Determine the maximum allowed backoff (with fallback to default).
	maxBackoff := p.config.MaxRekeyBackoff
	if maxBackoff == 0 {
		maxBackoff = maxRekeyBackoff // Use library default if unset.
	}

	for i := 0; i < p.config.MaxRekeyAttempts; i++ {
		// Capture the currently-active cipher pointer for later zeroization.
		old = p.cipher.Load().(*chacha20.Cipher)

		// Attempt to create a new ChaCha20 cipher (with new key and nonce).
		stream, err := newCipher()
		if err == nil {
			// Store the new cipher atomically.
			p.cipher.Store(stream)

			// Reset usage count for new key/nonce.
			atomic.StoreUint64(&p.usage, 0)

			// Wipe the memory of the old cipher (zero out struct fields).
			*old = chacha20.Cipher{}

			// Rekey successful; exit the function.
			return
		}

		// If cipher initialization failed, jitter the retry delay by a random amount.
		var b [8]byte
		if _, err := rand.Read(b[:]); err == nil {
			// Interpret b as a big-endian uint64 for jitter.
			rnd := binary.BigEndian.Uint64(b[:])

			// Calculate delay: base + (rnd mod base) for randomness.
			delay := base + time.Duration(rnd%uint64(base))
			time.Sleep(delay)
		} else {
			// If reading random bytes fails, fall back to fixed backoff.
			time.Sleep(base)
		}

		// Exponentially backoff for the next retry, up to the maximum allowed.
		base *= 2
		if base > maxBackoff {
			base = maxBackoff
		}
	}

	// All attempts to rekey failed; function exits with existing cipher in place.
}
