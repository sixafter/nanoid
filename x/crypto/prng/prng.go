// Copyright (c) 2024 Six After, Inc
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

// init sets up the package‐level Reader by creating a new pooled PRNG instance.
// It is invoked automatically at program startup (package initialization).
// If NewReader fails (e.g., OS entropy unavailable), init will panic to prevent
// running without a secure random source.
//
// Panicking here is intentional and idiomatic for cryptographic primitives:
// it ensures that any critical failure in obtaining a secure entropy source
// is detected immediately and cannot be ignored.
func init() {
	var err error
	Reader, err = NewReader()
	if err != nil {
		panic(fmt.Sprintf("prng.init: failed to create Reader: %v", err))
	}
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
	pool *sync.Pool
}

// NewReader constructs and returns an io.Reader that produces cryptographically secure
// pseudo-random bytes using a pool of ChaCha20‐based PRNG instances. You may supply
// zero or more functional options to customize its behavior.
//
// Each PRNG in the pool is seeded with a unique key and nonce from crypto/rand,
// and automatically rotates to a fresh key/nonce pair after emitting a configurable
// number of bytes (MaxBytesPerKey). The pool will retry PRNG initialization up to
// MaxInitRetries times, and will panic if it cannot produce a valid generator.
//
// Available Options:
//
//	WithMaxBytesPerKey(n uint64)  – set the byte threshold for key rotation (default 1GiB).
//	WithMaxInitRetries(r int)     – set the number of attempts to initialize each PRNG (default 3).
//	WithMaxRekeyAttempts(r int)   – set retry count for key rotation (default 5).
//	WithRekeyBackoff(d time.Duration) – set initial back-off duration for retries (default 100ms).
//
// Example:
//
//	reader, err := prng.NewReader()
//	if err != nil {
//	    // handle error
//	}
//
//	buf := make([]byte, 64)
//	n, err := reader.Read(buf)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("Read %d bytes: %x\n", n, buf)
func NewReader(opts ...Option) (io.Reader, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	pool := &sync.Pool{
		New: func() interface{} {
			var (
				p   *prng
				err error
			)
			for i := 0; i < cfg.MaxInitRetries; i++ {
				if p, err = newPRNG(&cfg); err == nil {
					return p
				}
			}
			panic(fmt.Sprintf("prng pool init failed after %d retries: %v", cfg.MaxInitRetries, err))
		},
	}
	return &reader{pool: pool}, nil
}

// Read pulls a *prng from the pool, claims exclusive use, fills the provided buffer
// with cryptographically secure random data, then returns the instance to the pool.
// It implements the io.Reader interface and ensures that no two goroutines can
// use the same PRNG instance simultaneously (preventing internal state corruption).
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
	p := r.pool.Get().(io.Reader)

	// Ensure the instance is returned to the pool when done
	defer r.pool.Put(p)
	return p.Read(b)
}

// prng implements io.Reader using a ChaCha20 cipher stream and supports
// asynchronous, nonblocking rotation of the underlying key/nonce pair.
//
// Each instance maintains its own ChaCha20 cipher (stored atomically), a
// scratch buffer for encryption, and internal counters to enforce a
// “forward secrecy” rekey after a configurable output threshold.
type prng struct {
	// cfg holds a pointer to this PRNG instance’s configuration parameters.
	// It provides tunable settings such as MaxBytesPerKey (keystream rotation threshold)
	// and MaxInitRetries (how many times to retry initialization).
	cfg *Config

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
// It implements the `io.Reader` interface and is intended for exclusive use by a single goroutine.
//
// Internally, Read does the following:
//  1. Determines the length `n` of the requested output. If `n == 0`, returns immediately with no error.
//  2. Atomically loads the current ChaCha20 cipher stream from `p.cipher`.
//  3. Prepares a zero-valued buffer of length `n` in `p.zero`, growing it if necessary.
//  4. Calls `cipher.XORKeyStream(b, p.zero)` to generate `n` bytes of output.
//  5. Atomically increments `p.usage` by `n`.
//  6. If `p.usage` has crossed `maxBytesPerKey`, attempts a single non-blocking
//     CAS to set `p.rekeying` from 0→1, and if successful, launches `p.asyncRekey()`
//     in a background goroutine to rotate the ChaCha20 key/nonce pair.
//
// Returns the number of bytes written (`n`) and any error encountered during rekey initiation
// (though key rotation errors are logged or dropped inside `asyncRekey`).
func (p *prng) Read(b []byte) (int, error) {
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	// Atomically retrieve the active cipher stream.
	stream := p.cipher.Load().(*chacha20.Cipher)

	// Ensure `p.zero` is a zero-filled buffer of length `n`.
	if cap(p.zero) < n {
		p.zero = make([]byte, n)
	} else {
		p.zero = p.zero[:n]
	}

	// XOR the zero buffer into `b`, producing random bytes.
	stream.XORKeyStream(b, p.zero)

	// Track how many bytes we've generated under the current key.
	atomic.AddUint64(&p.usage, uint64(n))

	// If we've exceeded our per-key threshold, trigger an async rekey.
	if atomic.LoadUint64(&p.usage) > p.cfg.MaxBytesPerKey {
		if atomic.CompareAndSwapUint32(&p.rekeying, 0, 1) {
			go p.asyncRekey()
		}
	}

	return n, nil
}

// newPRNG creates and returns a fully initialized prng instance.
// It generates a fresh ChaCha20 cipher, zeroes out any sensitive seed material,
// and stores the cipher in an atomic.Value for lock-free access in Read().
// Returns an error if the underlying cipher setup fails.
func newPRNG(cfg *Config) (*prng, error) {
	stream, err := newCipher()
	if err != nil {
		return nil, err
	}

	p := &prng{
		zero: make([]byte, 0),
		cfg:  cfg,
	}

	// Store the cipher for atomic.Load() in Read().
	p.cipher.Store(stream)
	return p, nil
}

// newCipher generates a new *chacha20.Cipher seeded with a cryptographically
// secure random key and nonce. It reads exactly chacha20.KeySize bytes for the key
// and chacha20.NonceSizeX bytes for the nonce from crypto/rand.Reader.
// Immediately after creating the cipher, it wipes the key and nonce buffers
// to prevent sensitive data leakage in memory.
// Returns the initialized cipher or an error if any step fails.
func newCipher() (*chacha20.Cipher, error) {
	// Allocate key+nonce buffers
	key := make([]byte, chacha20.KeySize)
	nonce := make([]byte, chacha20.NonceSizeX)

	// Fill with random data
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("newCipher: failed to read key: %w", err)
	}
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("newCipher: failed to read nonce: %w", err)
	}

	// Create the ChaCha20 stream cipher
	stream, err := chacha20.NewUnauthenticatedCipher(key, nonce)

	// Zero out the seed material immediately
	for i := range key {
		key[i] = 0
	}
	for i := range nonce {
		nonce[i] = 0
	}

	if err != nil {
		return nil, fmt.Errorf("newCipher: unable to initialize cipher: %w", err)
	}
	return stream, nil
}

// asyncRekey performs an asynchronous, non‐blocking rotation of the internal ChaCha20 cipher.
// It is invoked when the per‐key usage threshold is exceeded and runs in its own goroutine.
// The process will retry up to Config.MaxRekeyAttempts times, waiting Config.RekeyBackoff
// (doubling on each retry) between attempts. On each attempt it captures the old cipher,
// and on success (or after all retries fail) it zeroes out the old cipher struct to remove
// any residual key or counter material from memory.
func (p *prng) asyncRekey() {
	// Always clear the rekeying flag when this goroutine exits
	defer atomic.StoreUint32(&p.rekeying, 0)

	base := p.cfg.RekeyBackoff
	var old *chacha20.Cipher

	for i := 0; i < p.cfg.MaxRekeyAttempts; i++ {
		// Capture the existing cipher so we can wipe it later
		old = p.cipher.Load().(*chacha20.Cipher)

		stream, err := newCipher()
		if err == nil {
			p.cipher.Store(stream)
			atomic.StoreUint64(&p.usage, 0)
			*old = chacha20.Cipher{}
			return
		}

		// Jitter: crypto/rand 8-byte uint64, mod base
		var b [8]byte
		if _, err := rand.Read(b[:]); err == nil {
			// interpret as big-endian uint64
			rnd := binary.BigEndian.Uint64(b[:])
			// offset in [0, base)
			delay := base + time.Duration(rnd%uint64(base))
			time.Sleep(delay)
		} else {
			// fallback to fixed back-off if RNG fails
			time.Sleep(base)
		}
		base *= 2
	}

	// All retries failed: leave the existing cipher in place, then exit
}
