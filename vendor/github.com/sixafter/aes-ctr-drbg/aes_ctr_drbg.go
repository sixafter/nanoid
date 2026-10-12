// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package ctrdrbg implements a pool-backed deterministic random bit generator (DRBG) based on the CTR_DRBG
// mechanism of NIST SP 800-90A Rev. 1, §10.2.1, using AES without a derivation function.
//
// It has two modes:
//   - Default: output is AES in counter mode under a key and counter taken directly from crypto/rand entropy.
//     The key is not changed after each request, so a later compromise of the internal state exposes earlier
//     output generated under the same key. Read performs 0 allocations per call.
//   - WithBacktrackingResistance: instantiation, reseeding, generation, additional input, prediction resistance,
//     and the per-request and reseed limits follow NIST SP 800-90A §10.2.1 and §9.3.1, known-answer self-tests
//     of these functions run before first use (§11.3), and they are tested against NIST CAVP CTR_DRBG
//     known-answer vectors. A new key and counter are derived after every request, at
//     2 allocations per Read call.
//
// WithContinuousHealthTest adds the FIPS 140-2 §4.9.2 continuous random number generator test. All
// cryptographic primitives come from the Go standard library, and the package works with Go's FIPS 140 mode
// (GODEBUG=fips140=on).
//
// Reference:
//
//	NIST Special Publication 800-90A Rev. 1, §10.2.1 (CTR_DRBG Construction)
//	https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-90Ar1.pdf
package ctrdrbg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	mrand "math/rand/v2"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// MaxBytesPerRequest is the NIST SP 800-90A maximum bytes per request for CTR_DRBG (2^19 bits = 64 KB).
	MaxBytesPerRequest = 1 << 16
)

var (
	// ErrRequestTooLarge is returned when a Read request exceeds the NIST SP 800-90A maximum (64 KB).
	ErrRequestTooLarge = errors.New("ctrdrbg: request exceeds NIST SP 800-90A max_number_of_bits_per_request (64 KB)")

	// ErrHealthTestFailed indicates that the continuous health test detected stuck output. The DRBG instance that
	// detected it stays in an error state and returns ErrHealthTestFailed for every later request.
	ErrHealthTestFailed = errors.New("ctrdrbg: continuous health test failed (stuck output detected)")

	// ErrPersonalizationTooLong is returned by NewReader when BacktrackingResistance is enabled and the
	// personalization string is longer than the NIST SP 800-90A maximum of key size + 16 bytes.
	ErrPersonalizationTooLong = errors.New("ctrdrbg: personalization string exceeds NIST SP 800-90A max_personalization_string_length (seedlen)")

	// ErrAdditionalInputTooLong is returned by ReadWithAdditionalInput when BacktrackingResistance is enabled and
	// the additional input is longer than the NIST SP 800-90A maximum of key size + 16 bytes.
	ErrAdditionalInputTooLong = errors.New("ctrdrbg: additional input exceeds NIST SP 800-90A max_additional_input_length (seedlen)")
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

	// Reseed injects new entropy and optional additional input, refreshing the DRBG state.
	//
	// In the default mode, additionalInput is combined with system entropy and personalization (if set)
	// to derive a new internal key and counter. With BacktrackingResistance, fresh system entropy combined
	// with additionalInput is mixed into the existing key and counter as specified by NIST SP 800-90A
	// §10.2.1.4.1, and personalization is not used; additionalInput must not exceed key size + 16 bytes,
	// or ErrAdditionalInputTooLong is returned.
	// Reseed can be called at any time to proactively refresh the DRBG.
	//
	// Returns an error if the reseed operation fails.
	Reseed(additionalInput []byte) error

	// ReadWithAdditionalInput generates cryptographically secure random bytes,
	// supplementing system entropy with optional per-call additional input.
	//
	// This method is NIST-compliant and allows callers to inject additional
	// entropy or context for a single Read operation.
	// If PredictionResistance is enabled, fresh entropy is always used; additionalInput is ignored in the default
	// mode and, with BacktrackingResistance, mixed into the prediction-resistance reseed (NIST SP 800-90A §9.3.1).
	//
	// When BacktrackingResistance is enabled (and PredictionResistance is not), additionalInput is mixed into
	// the key and counter before and after generation, as specified by NIST SP 800-90A §10.2.1.5.1 (steps 2
	// and 6), without obtaining new entropy. It must not exceed key size + 16 bytes, or ErrAdditionalInputTooLong
	// is returned; empty additionalInput is treated as absent.
	//
	// When EnableKeyRotation is enabled, an asynchronous key rotation that completes concurrently may replace
	// the state derived from additionalInput with a new state derived from fresh system entropy alone. Output
	// remains seeded from fresh entropy, but the additional input may not be reflected in it.
	//
	// Returns the number of bytes read (equal to len(b)) and an error (if any).
	ReadWithAdditionalInput(b []byte, additionalInput []byte) (int, error)
}

// init initializes the package-level Reader. It panics if NewReader fails, preventing operation without
// a secure random source. This follows cryptographic best practices by making entropy failure a fatal error.
func init() {
	cfg := DefaultConfig()
	pools, err := initShardPools(cfg)
	if err != nil {
		panic(fmt.Sprintf("ctrdrbg: failed to initialize pools: %v", err))
	}

	Reader = &reader{pools: pools}
}

// initShardPools creates and validates all sync.Pool shards for concurrent DRBG use.
//
// For each shard, a sync.Pool is created whose New function constructs a DRBG instance using the provided config.
// If instantiation fails, it retries up to MaxInitRetries times, then panics if unsuccessful.
// After creating each pool, it is eagerly tested by borrowing and returning an instance, to ensure failures are
// caught at construction rather than at first use.
//
// If any pool.New panics (due to repeated DRBG initialization failure), the function recovers and returns the
// error to the caller (for use in NewReader). In package-level init, a panic is allowed to abort process startup.
//
// Parameters:
//   - cfg: Config specifying DRBG options (shard count, retries, etc.)
//
// Returns:
//   - []*sync.Pool: slice of initialized pools, one per shard.
//   - error: non-nil if pool initialization panicked for any shard.
func initShardPools(cfg Config) ([]*sync.Pool, error) {
	// Create a slice of sync.Pool pointers, one for each shard.
	pools := make([]*sync.Pool, cfg.Shards)
	for i := range pools {
		// Capture the config by value for use in the pool.New closure (avoids loop variable capture bug).
		capturedCfg := cfg
		pools[i] = &sync.Pool{
			New: func() interface{} {
				var d *drbg
				var err error
				// Attempt to instantiate a DRBG instance, retrying up to MaxInitRetries times.
				for r := 0; r < capturedCfg.MaxInitRetries; r++ {
					if d, err = newDRBG(&capturedCfg); err == nil {
						return d
					}
				}
				// If all attempts fail here, we return nil. The eager initialization below
				// performs the same construction and will surface a concrete error to the caller.
				return nil
			},
		}

		// Eagerly test pool initialization to ensure catastrophic failures are caught immediately,
		// not deferred until first use. Attempt construction directly and return an error on failure.
		var (
			warm *drbg
			err  error
		)
		for r := 0; r < capturedCfg.MaxInitRetries; r++ {
			if warm, err = newDRBG(&capturedCfg); err == nil {
				pools[i].Put(warm)
				err = nil
				break
			}
		}

		// If initialization failed, return error immediately.
		if err != nil {
			return nil, fmt.Errorf("ctrdrbg pool initialization failed after %d retries: %v", capturedCfg.MaxInitRetries, err)
		}
	}

	return pools, nil
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
//	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(ctrdrbg.KeySize256))
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
	// Start with a default configuration, then apply each functional option to mutate cfg.
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// FIPS 140-2 §4.9.1: Run Known Answer Tests if enabled.
	if cfg.EnableSelfTests {
		if err := RunSelfTests(); err != nil {
			return nil, err
		}
	}

	// Validate the configured key size is appropriate for AES.
	// Only 16, 24, or 32 bytes (AES-128, AES-192, AES-256) are supported.
	switch cfg.KeySize {
	case KeySize128, KeySize192, KeySize256:
	default:
		return nil, fmt.Errorf("invalid key size %d bytes; must be 16, 24, or 32", cfg.KeySize)
	}

	if cfg.MaxInitRetries < 1 {
		return nil, fmt.Errorf("invalid MaxInitRetries: must be >= 1")
	}

	// NIST SP 800-90A §10.2.1 Table 3: without a derivation function, the personalization string
	// must not exceed key size + 16 bytes.
	if cfg.BacktrackingResistance && len(cfg.Personalization) > int(cfg.KeySize)+16 {
		return nil, ErrPersonalizationTooLong
	}

	// NIST SP 800-90A §10.2.1 Table 3: a reseed is required after at most 2^48 output requests. With
	// BacktrackingResistance, that limit applies even when ReseedRequests is not set.
	if cfg.BacktrackingResistance && cfg.ReseedRequests == 0 {
		cfg.ReseedRequests = maxReseedInterval
	}

	// NIST SP 800-90A §11.3: with BacktrackingResistance, known-answer tests of the instantiate, reseed, and
	// generate functions must pass before first use.
	if cfg.BacktrackingResistance {
		if err := runDRBGSelfTests(); err != nil {
			return nil, err
		}
	}

	// Initialize the shard pools using the validated configuration.
	pools, err := initShardPools(cfg)
	if err != nil {
		return nil, err
	}

	//  Return a new reader that wraps the initialized pool.
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

// Reseed refreshes the state of one DRBG instance per shard pool with new entropy and optional additional input.
//
// This method implements the Interface contract for NIST SP 800-90A compliant deterministic random bit generators.
// For the instance borrowed from each shard pool, Reseed obtains fresh entropy from the system entropy source and
// combines it with the provided additionalInput and, in the default mode, the configured personalization string.
// With BacktrackingResistance, the result is mixed into the existing key and counter as specified by
// NIST SP 800-90A §10.2.1.4.1, and personalization is not used.
// This operation ensures that each reseeded instance has a newly derived key and counter (V), providing domain separation
// and supporting explicit entropy injection for compliance, recovery, or defense-in-depth.
//
// Parameters:
//   - additionalInput []byte: Optional per-call entropy or domain-separation input to be mixed into the reseed.
//     May be nil if no additional input is required.
//
// Returns:
//   - error: Returns the first error encountered if any DRBG instance fails to reseed; otherwise returns nil.
//
// Security and Compliance Notes:
//   - Reseed is safe for concurrent use and may be called at any time during operation.
//   - Only the instance borrowed from each shard is reseeded. Other pooled instances, including those in use by
//     concurrent readers, keep their current state until their own reseed, rekey, or fork-detection policy applies.
//     To reseed on every output request, use WithPredictionResistance.
//   - This is required for some FIPS and high-assurance applications and is recommended for recovery from suspected
//     entropy pool compromise or for regulatory compliance triggers.
//
// Example usage:
//
//	err := reader.Reseed([]byte("audit-event-timestamp"))
//	if err != nil {
//	    log.Fatalf("reseed failed: %v", err)
//	}
func (r *reader) Reseed(additionalInput []byte) error {
	// Iterate over each sync.Pool in the shard pool array.
	for _, pool := range r.pools {
		// Borrow a DRBG instance from the pool.
		d := pool.Get().(*drbg)
		// Attempt to reseed this DRBG instance using the provided additionalInput.
		// Reseed will combine system entropy, personalization, and additionalInput as per NIST.
		err := d.Reseed(additionalInput)
		// Return the DRBG instance to the pool for future reuse, regardless of success.
		pool.Put(d)
		// If reseed fails for any DRBG, immediately return the error to the caller.
		if err != nil {
			return err
		}
	}
	// If all DRBGs reseeded successfully, return nil.
	return nil
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

// ReadWithAdditionalInput fills the provided buffer with cryptographically secure random bytes,
// optionally supplementing system entropy with caller-provided additionalInput, per NIST SP 800-90A.
//
// This method enables advanced consumers to inject per-call entropy, external event data, or
// session context into the DRBG reseed process, as specified by the NIST DRBG "additional input" feature.
// If PredictionResistance is enabled on the DRBG configuration, fresh entropy is always used; the additionalInput
// argument is ignored in the default mode and, with BacktrackingResistance, mixed into the prediction-resistance
// reseed (NIST SP 800-90A §9.3.1).
//
// Parameters:
//   - b []byte: The output buffer to fill with random bytes. Must be non-nil; may be zero-length.
//   - additionalInput []byte: Optional per-call entropy or context for NIST-compliant reseed. May be nil.
//
// Returns:
//   - int: The number of bytes written to b (always len(b) unless b is empty).
//   - error: Error returned if random generation fails; nil on success.
//
// Concurrency and Pooling:
//   - This method is safe for concurrent use. The underlying DRBG instance is selected from an internal pool
//     using a sharding strategy to maximize throughput and minimize contention.
//
// Security and Compliance Notes:
//   - additionalInput is cryptographically mixed with system entropy and personalization (if any) for this call only.
//   - If PredictionResistance is enabled, reseed is always performed from entropy; additionalInput is ignored in
//     the default mode and mixed into that reseed with BacktrackingResistance.
//   - For most use cases, use the standard Read method. This method is intended for regulatory, compliance, or
//     advanced event-driven entropy injection scenarios.
//
// Example usage:
//
//	n, err := reader.ReadWithAdditionalInput(buf, []byte("user-event-entropy"))
//	if err != nil {
//	    // handle error
//	}
func (r *reader) ReadWithAdditionalInput(b []byte, additionalInput []byte) (int, error) {
	// Determine the number of pools (shards) in the reader for load balancing.
	n := len(r.pools)
	shard := 0
	// For multiple shards, select a random shard index for this call.
	if n > 1 {
		shard = shardIndex(n)
	}
	// Borrow a DRBG instance from the selected pool for this operation.
	d := r.pools[shard].Get().(*drbg)
	// Ensure the instance is returned to the pool after use (even on error or panic).
	defer r.pools[shard].Put(d)
	// Fill the buffer using the borrowed DRBG, injecting additionalInput as specified.
	return d.ReadWithAdditionalInput(b, additionalInput)
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
	// Return immediately if the buffer is empty, as required by the io.Reader contract.
	if len(b) == 0 {
		return 0, nil
	}

	// Determine the shard index based on the number of pools available.
	n := len(r.pools)
	shard := 0
	if n > 1 {
		shard = shardIndex(n)
	}

	// Borrow an instance of the internal deterministic random bit generator from the pool.
	// This ensures that each call gets exclusive access to an isolated state for cryptographic safety.
	d := r.pools[shard].Get().(*drbg)

	// Ensure that the borrowed instance is returned to the pool, even if Read fails or panics.
	// This pattern prevents resource leaks and maintains pool integrity.
	defer r.pools[shard].Put(d)

	// Fill the caller’s buffer with random data using the borrowed generator.
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

	// lastReseedTime records the time of the last successful reseed.
	// Used to determine if the configured ReseedInterval has elapsed and
	// automatic reseeding should occur before the next output.
	lastReseedTime time.Time

	// entropySource supplies entropy for reseeds when BacktrackingResistance is enabled. It is nil, meaning
	// crypto/rand, for every instance except those created by the NIST SP 800-90A §11.3 known-answer self-test,
	// which must use fixed entropy.
	entropySource io.Reader

	// zero is a pre-allocated slice of zero-filled bytes used for output buffering.
	//
	// When UseZeroBuffer is enabled in config, this buffer is XOR-ed with
	// AES-CTR output to efficiently produce random bytes. Sized dynamically as needed.
	zero []byte

	// vMu is a mutex protecting the evolving counter (v) for this DRBG instance.
	//
	// All access and mutation of v must occur with this mutex held to ensure:
	//   - Counter advancement is atomic and non-overlapping across reads
	//   - Proper persistence of the counter-value between consecutive reads
	//   - Safe resetting of the counter during key rotation (rekey)
	vMu sync.Mutex

	// v is the current 128-bit internal counter (NIST "V") for the DRBG instance.
	//
	// This counter is incremented for each AES block produced, ensuring
	// unique, non-repeating output for every call to Read. It is initialized
	// from the state.v value at creation or rekey, and persisted between reads.
	v [16]byte

	// encV is a persistent [16]byte working buffer used as a session-local counter
	// during output generation.
	//
	// On each call to Read, the current counter-value (d.v) is copied into encV, which is
	// then incremented and used for block generation throughout the request. Only after all
	// output is produced is encV copied back into d.v, ensuring atomic and consistent counter
	// advancement. This prevents partial or inconsistent counter updates if Read exits early
	// due to errors or panics.
	encV [16]byte

	// tmp is a persistent [16]byte working buffer used during output generation.
	// It holds the encrypted output for the final (partial) block in fillBlocks.
	// This avoids repeated stack allocations and ensures maximum efficiency.
	tmp [16]byte

	// Previous output block for continuous health test
	lastOutputBlock [16]byte

	// requests counts the number of output requests (calls to Read or ReadWithAdditionalInput)
	// since the last reseed. Used to enforce the ReseedRequests limit and trigger reseeding
	// after a configured number of requests.
	requests uint64

	// forkChecks counts output requests for interval-based fork detection (ForkDetectionInterval > 0).
	// It is separate from the reseed request counter, which every reseed resets, so reseeding cannot
	// prevent the fork-detection interval from being reached.
	forkChecks uint64

	// usage tracks the number of bytes generated since the last key rotation.
	//
	// When usage exceeds config.MaxBytesPerKey, a rekey is triggered to ensure
	// forward secrecy and mitigate key compromise risk. This value is atomically updated.
	usage uint64

	// pid caches the process identifier (PID) of the operating system process in which
	// this DRBG instance was most recently initialized or reseeded.
	//
	// Purpose:
	//   - Enables robust detection of process-level forks (e.g., via fork(2) or similar system calls).
	//   - After a fork, the child process receives a new, unique PID, but the DRBG instance initially
	//     retains the PID from its parent process.
	//   - By comparing the current process PID (os.Getpid()) to this cached value, the DRBG can reliably
	//     detect fork events at runtime.
	//   - When a fork is detected, the DRBG securely reseeds its cryptographic state, preventing random
	//     stream duplication and ensuring forward and backward security in both parent and child processes.
	//
	// Security Rationale:
	//   - Eliminates the risk of duplicated random streams following a process fork—a known pitfall of
	//     userspace CSPRNGs in forking environments (Linux, macOS).
	//   - Aligns the safety guarantees of userspace DRBGs with those provided by kernel-backed CSPRNGs,
	//     such as Linux getrandom(2), which handle fork-safety internally.
	pid int

	// rekeying is an atomic flag (0 or 1) that guards rekey attempts.
	//
	// It ensures that only one goroutine performs rekeying at a time.
	// Uses atomic operations for concurrency safety.
	rekeying uint32

	// True once first output block recorded
	healthTestReady bool

	// healthTestFailed records that the continuous health test has failed. FIPS 140-2 §4.9.2 requires the
	// instance to stay in an error state, so every later output request returns ErrHealthTestFailed.
	healthTestFailed bool
}

// continuousHealthTest performs the FIPS 140-2 §4.9.2 continuous random number generator test on an output
// request: every 128-bit block generated for b is compared with the previously generated block of this instance,
// including the full block behind a final partial block, which is still in the scratch buffer. If two compared
// blocks are equal, the instance enters a permanent error state and ErrHealthTestFailed is returned.
//
// The compared blocks are output delivered to the caller, so a constant-time comparison is not required.
func (d *drbg) continuousHealthTest(b []byte) error {
	if d.healthTestFailed {
		return ErrHealthTestFailed
	}

	// Compare each block with the one before it, starting from the block saved by the previous request.
	prev := d.lastOutputBlock[:]
	full := len(b) &^ 15
	for off := 0; off < full; off += 16 {
		curr := b[off : off+16]
		if d.healthTestReady && bytes.Equal(prev, curr) {
			d.healthTestFailed = true
			return ErrHealthTestFailed
		}
		prev = curr
		d.healthTestReady = true
	}
	if full < len(b) {
		if d.healthTestReady && bytes.Equal(prev, d.tmp[:]) {
			d.healthTestFailed = true
			return ErrHealthTestFailed
		}
		prev = d.tmp[:]
		d.healthTestReady = true
	}

	// Save the last generated block for comparison with the next request.
	copy(d.lastOutputBlock[:], prev)
	return nil
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
	// Return immediately if the buffer is empty, as required by the io.Reader contract.
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	// NIST SP 800-90A §10.2.1: Validate max_number_of_bits_per_request (64 KB).
	if n > MaxBytesPerRequest {
		return 0, ErrRequestTooLarge
	}

	d.reseedIfForked()

	// Prediction Resistance
	if d.config.PredictionResistance {
		var err error
		if d.config.BacktrackingResistance {
			// NIST SP 800-90A §9.3.1 step 7.1: reseed per §10.2.1.4.1.
			err = d.reseedUpdate(nil)
		} else {
			err = d.reseed(nil)
		}
		if err != nil {
			return 0, fmt.Errorf("prediction resistance reseed failed: %w", err)
		}
	} else {
		// Optional: Reseed if the configured interval has elapsed since the last reseed.
		if d.config.ReseedInterval > 0 {
			now := time.Now()
			if now.Sub(d.lastReseedTime) >= d.config.ReseedInterval {
				if err := d.reseed(nil); err != nil {
					return 0, fmt.Errorf("interval reseed failed: %w", err)
				}
			}
		}

		// NIST-required: Reseed if the configured request count is exceeded.
		if d.config.ReseedRequests > 0 && atomic.LoadUint64(&d.requests) >= d.config.ReseedRequests {
			if err := d.reseed(nil); err != nil {
				return 0, fmt.Errorf("request-count reseed failed: %w", err)
			}
		}
	}

	// Atomically load the current DRBG cryptographic state.
	st := d.state.Load()

	// Lock the counter mutex to guarantee exclusive access to the evolving counter.
	d.vMu.Lock()

	// Copy the current counter value to a local variable. This snapshot forms the basis
	// of the unique keystream for this read operation.
	copy(d.encV[:], d.v[:])

	// Fill the output buffer using the current cryptographic state and the local counter,
	// incrementing the counter as output is produced. All counter increments are reflected
	// in the local variable.
	d.fillBlocks(b, st, &d.encV)

	if d.config.ContinuousHealthTest {
		if err := d.continuousHealthTest(b); err != nil {
			clear(b)
			d.vMu.Unlock()
			return 0, err
		}
	}

	// Persist the advanced counter back to the DRBG instance, ensuring subsequent reads
	// continue the keystream seamlessly without overlap or repetition.
	copy(d.v[:], d.encV[:])

	// NIST SP 800-90A §10.2.1.5.1 step 6: derive a new key and counter for backtracking resistance.
	if d.config.BacktrackingResistance {
		if err := d.update(st, nil); err != nil {
			d.vMu.Unlock()
			return 0, err
		}
	}

	// Unlock the mutex, allowing other callers to proceed.
	d.vMu.Unlock()

	// NIST-required: Increment the requests counter for this DRBG instance.
	if !d.config.PredictionResistance {
		atomic.AddUint64(&d.requests, 1)
	}

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

// ReadWithAdditionalInput fills the provided buffer with cryptographically secure random bytes,
// optionally reseeding the DRBG instance with caller-provided additional input per NIST SP 800-90A.
//
// This method is intended for advanced use cases where explicit entropy injection or
// per-call domain separation is required, as enabled by the NIST DRBG "additional input" feature.
// If PredictionResistance is enabled, the DRBG reseeds from system entropy for every read; additionalInput is
// ignored in the default mode and, with BacktrackingResistance, mixed into that reseed (NIST SP 800-90A §9.3.1
// step 7.1). Otherwise, in the default mode, if additionalInput is non-nil, the DRBG instance reseeds using a
// combination of system entropy, personalization, and the supplied additionalInput before generating output.
//
// Semantics and Implementation Details:
//   - If PredictionResistance is enabled, reseeds from fresh system entropy before generating output;
//     additionalInput is ignored in the default mode and mixed into that reseed with BacktrackingResistance.
//   - In the default mode, if PredictionResistance is not enabled and additionalInput is non-nil, reseeds
//     using both entropy and the caller's input.
//   - Loads the current cryptographic state (AES key, block cipher, initial counter) atomically.
//   - The DRBG's internal counter (v) is protected by a mutex to guarantee non-overlapping output across reads.
//   - Output is generated using fillBlocks and the advanced counter value is persisted for continuity.
//   - If key rotation is enabled and the usage threshold is exceeded, an asynchronous rekey is triggered.
//
// Parameters:
//   - b []byte: Output buffer to be filled with cryptographically secure random bytes.
//   - additionalInput []byte: Optional per-call entropy or context to inject during reseed. May be nil.
//
// Returns:
//   - int: Number of bytes written (equal to len(b) unless b is empty).
//   - error: Error if entropy acquisition, reseed, or output generation fails.
//
// Example:
//
//	n, err := drbg.ReadWithAdditionalInput(buf, []byte("request-entropy"))
//	if err != nil {
//	    // handle error
//	}
func (d *drbg) ReadWithAdditionalInput(b []byte, additionalInput []byte) (int, error) {
	// Return immediately if the buffer is empty, as required by the io.Reader contract.
	n := len(b)
	if n == 0 {
		return 0, nil
	}

	// NIST SP 800-90A §10.2.1: Validate max_number_of_bits_per_request (64 KB).
	if n > MaxBytesPerRequest {
		return 0, ErrRequestTooLarge
	}

	// NIST SP 800-90A §10.2.1 Table 3: without a derivation function, additional input must not exceed
	// key size + 16 bytes.
	if d.config.BacktrackingResistance && len(additionalInput) > int(d.config.KeySize)+16 {
		return 0, ErrAdditionalInputTooLong
	}

	d.reseedIfForked()

	// If PredictionResistance is enabled, always reseed from fresh entropy before output. With
	// BacktrackingResistance, the reseed follows NIST SP 800-90A §9.3.1 step 7.1 and §10.2.1.4.1 and mixes in
	// the additional input, which is then not used again for this request; in the default mode it is ignored.
	if d.config.PredictionResistance {
		var err error
		if d.config.BacktrackingResistance {
			err = d.reseedUpdate(additionalInput)
		} else {
			err = d.reseed(nil)
		}
		if err != nil {
			return 0, fmt.Errorf("prediction resistance reseed failed: %w", err)
		}
	} else {
		// Optional: Reseed if the configured interval has elapsed since the last reseed.
		if d.config.ReseedInterval > 0 {
			now := time.Now()
			if now.Sub(d.lastReseedTime) >= d.config.ReseedInterval {
				if err := d.reseed(nil); err != nil {
					return 0, fmt.Errorf("interval reseed failed: %w", err)
				}
			}
		}

		// NIST-required: Reseed if the configured request count is exceeded.
		if d.config.ReseedRequests > 0 && atomic.LoadUint64(&d.requests) >= d.config.ReseedRequests {
			if err := d.reseed(nil); err != nil {
				return 0, fmt.Errorf("request-count reseed failed: %w", err)
			}
		}

		// If additionalInput is provided and prediction resistance is not enabled,
		// reseed using both entropy and additional input for this output request. With BacktrackingResistance,
		// additional input is instead mixed into the key and counter before and after generation
		// (NIST SP 800-90A §10.2.1.5.1 steps 2 and 6).
		if additionalInput != nil && !d.config.BacktrackingResistance {
			if err := d.reseed(additionalInput); err != nil {
				return 0, fmt.Errorf("reseed with additional input failed: %w", err)
			}
		}
	}

	// Load the current cryptographic state (AES key, block cipher, initial counter) atomically.
	st := d.state.Load()

	// Lock the counter mutex to guarantee exclusive access to the evolving counter.
	d.vMu.Lock()

	// Copy the current counter value to a local variable for use in output generation.
	copy(d.encV[:], d.v[:])

	// NIST SP 800-90A §10.2.1.5.1 step 2: mix the additional input into the key and counter before generating.
	// With prediction resistance, additional input is not used by the generate steps.
	var generateInput []byte
	if d.config.BacktrackingResistance && !d.config.PredictionResistance && len(additionalInput) > 0 {
		generateInput = additionalInput
		if err := d.update(st, generateInput); err != nil {
			d.vMu.Unlock()
			return 0, err
		}
		st = d.state.Load()
		copy(d.encV[:], d.v[:])
	}

	// Fill the output buffer using the current cryptographic state and the local counter,
	// incrementing the counter as output is produced.
	d.fillBlocks(b, st, &d.encV)

	if d.config.ContinuousHealthTest {
		if err := d.continuousHealthTest(b); err != nil {
			clear(b)
			d.vMu.Unlock()
			return 0, err
		}
	}

	// Persist the advanced counter back to the DRBG instance, ensuring
	// future reads continue the keystream seamlessly.
	copy(d.v[:], d.encV[:])

	// NIST SP 800-90A §10.2.1.5.1 step 6: derive a new key and counter for backtracking resistance.
	if d.config.BacktrackingResistance {
		if err := d.update(st, generateInput); err != nil {
			d.vMu.Unlock()
			return 0, err
		}
	}

	// Unlock the mutex, allowing other callers to proceed.
	d.vMu.Unlock()

	// NIST-required: Increment the requests counter for this DRBG instance.
	if !d.config.PredictionResistance {
		atomic.AddUint64(&d.requests, 1)
	}

	// Key rotation logic: update the usage counter and, if the output threshold is
	// exceeded, trigger asynchronous rekeying in a background goroutine.
	if d.config.EnableKeyRotation {
		atomic.AddUint64(&d.usage, uint64(n))
		if atomic.LoadUint64(&d.usage) >= d.config.MaxBytesPerKey {
			if atomic.CompareAndSwapUint32(&d.rekeying, 0, 1) {
				go d.asyncRekey()
			}
		}
	}

	return n, nil
}

// Reseed injects new entropy and optional additional input, refreshing the DRBG instance's internal state.
//
// This method can be called at any time to force a rekey and counter reset. In the default mode, the
// additionalInput parameter, if non-nil, is cryptographically combined with system entropy and the
// configured personalization string to derive a new key and counter. With BacktrackingResistance, the
// reseed follows NIST SP 800-90A §10.2.1.4.1: fresh entropy combined with additionalInput is mixed into
// the existing key and counter, and personalization is not used.
//
// Reseed is suitable for explicit entropy injection after external events, regulatory triggers,
// or as a proactive measure to ensure stream independence and forward secrecy. Upon successful
// reseed, all future outputs are derived from the new state.
//
// Parameters:
//   - additionalInput []byte: Optional extra entropy or domain-separation input for this reseed operation.
//     May be nil if not required.
//
// Returns:
//   - error: Non-nil if entropy acquisition, cipher construction, or state replacement fails; nil on success.
//
// Example usage:
//
//	err := drbg.Reseed([]byte("device-boot-entropy"))
//	if err != nil {
//	    log.Fatalf("reseed failed: %v", err)
//	}
func (d *drbg) Reseed(additionalInput []byte) error {
	// Reseed the DRBG instance using system entropy and any caller-provided additional input.
	// The reseed function will cryptographically mix system entropy, personalization, and additionalInput,
	// replacing the internal key, counter, and AES state atomically. If reseed fails, the previous state is retained.
	return d.reseed(additionalInput)
}

// fillBlocks fills the byte slice `b` with cryptographically secure, deterministic random data
// generated from the provided DRBG state and a caller-provided working counter.
//
// This method implements the core NIST SP 800-90A AES-CTR-DRBG output logic. It is concurrency safe
// as it operates only on immutable state and caller-provided (session-local) counter. No DRBG struct
// fields are mutated during block generation.
//
// Parameters:
//   - b   []byte:      Output buffer to be filled with random bytes. Must be at least 1 byte in length.
//   - st  *state:      Immutable snapshot of the DRBG key, block cipher, and initial counter (V).
//   - v   *[16]byte:   Session-local working counter for this output operation. Advanced in place.
//
// Behavior:
//   - Processes output in 16-byte (AES block size) chunks for maximal efficiency.
//   - For each block, increments the session-local counter, encrypts it, and writes the result to output.
//   - Supports two strategies:
//   - UseZeroBuffer: Encrypted blocks are staged in a reusable buffer before being copied out (reducing allocations).
//   - Fast path: Output is written directly into the caller's buffer except for a possible tail partial block,
//     which uses the persistent drbg.tmp [16]byte buffer.
//
// Security:
//   - Ensures every 16-byte block is generated with a unique counter value per NIST recommendations.
//   - Never mutates DRBG fields or global internal state directly.
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

			// Advance the session-local counter as required by CTR mode (one block per keystream segment).
			incV(v)

			// Encrypt the incremented counter; write keystream into zero buffer.
			// A final partial block is encrypted into a full-block scratch buffer and then copied,
			// because Encrypt requires a full-block destination.
			if blockSize == 16 {
				st.block.Encrypt(d.zero[offset:offset+16], v[:])
			} else {
				st.block.Encrypt(d.tmp[:], v[:])
				copy(d.zero[offset:offset+blockSize], d.tmp[:blockSize])
			}

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
		incV(v)
		st.block.Encrypt(d.tmp[:], v[:])
		copy(b[offset:], d.tmp[:tail])
	}
}

// update mixes the given data into this instance's key and counter, as specified by the update function of
// NIST SP 800-90A §10.2.1.2. It encrypts successive values of the working counter under the current key to
// produce key size + 16 bytes of keystream and XORs that with the given data (shorter data is treated as
// zero-padded, and nil as all zeros, as in §10.2.1.5.1 step 6 when there is no additional input). The first
// key-size bytes become the new key and the last 16 bytes the new counter; the new state is published
// atomically and the working counter is set to the new counter.
//
// The caller must hold the counter lock. Keystream blocks are staged one at a time in a persistent
// per-instance buffer and written directly into the new state, so no stack buffers escape to the heap
// through the cipher.Block interface. The given data must not exceed key size + 16 bytes.
func (d *drbg) update(st *state, providedData []byte) error {
	keyLen := int(d.config.KeySize)
	seedLen := keyLen + 16

	newState := &state{}
	for pos := 0; pos < seedLen; pos += 16 {
		incV(&d.encV)
		st.block.Encrypt(d.tmp[:], d.encV[:])
		if pos < len(providedData) {
			m := min(16, seedLen-pos, len(providedData)-pos)
			subtle.XORBytes(d.tmp[:m], d.tmp[:m], providedData[pos:pos+m])
		}

		// Bytes [0, keyLen) of the result form the key; bytes [keyLen, seedLen) form V.
		n := 0
		if pos < keyLen {
			n = copy(newState.key[pos:keyLen], d.tmp[:])
		}
		if pos+n >= keyLen {
			copy(newState.v[pos+n-keyLen:], d.tmp[n:])
		}
	}
	clear(d.tmp[:])

	block, err := aes.NewCipher(newState.key[:keyLen])
	if err != nil {
		return err
	}
	newState.block = block

	d.state.Store(newState)
	copy(d.v[:], newState.v[:])
	return nil
}

// instantiate sets the initial key and counter from the given seed material (fresh entropy XOR the
// personalization string), as specified by NIST SP 800-90A §10.2.1.3.1: starting from an all-zero key and
// counter, the seed material is mixed in by the §10.2.1.2 update function. It is used when
// BacktrackingResistance is enabled and runs only when the instance is created.
func (d *drbg) instantiate(seedMaterial []byte) error {
	block, err := aes.NewCipher(make([]byte, d.config.KeySize))
	if err != nil {
		return err
	}

	d.vMu.Lock()
	defer d.vMu.Unlock()
	clear(d.encV[:])
	return d.update(&state{block: block}, seedMaterial)
}

// reseed refreshes the DRBG instance with new system entropy, personalization, and optional additional input.
//
// This function generates a new internal state (AES key, counter, and cipher block) using the provided DRBG
// configuration and any additionalInput supplied by the caller. It atomically installs the new cryptographic state,
// ensuring no overlap with previous keystream output. After reseed, both the byte usage counter and request count
// are reset, and the reseed timestamp is updated to support interval/request-count-based reseed policies.
//
// Parameters:
//   - additionalInput []byte: Optional, caller-supplied entropy or domain-separation material that is cryptographically
//     combined with system entropy and personalization. This may be nil for standard reseeds.
//
// Returns:
//   - error: Non-nil if entropy acquisition or state initialization fails; nil on success.
//
// Concurrency Notes:
//   - This method synchronizes counter updates using a mutex and uses atomic operations for usage/request counters.
//   - If called concurrently, it is possible for closely timed reseeds to slightly race in updating the metadata fields;
//     this does not impact cryptographic safety.
func (d *drbg) reseed(additionalInput []byte) error {
	// With BacktrackingResistance, every reseed follows NIST SP 800-90A §10.2.1.4.1.
	if d.config.BacktrackingResistance {
		return d.reseedUpdate(additionalInput)
	}

	// Generate a new DRBG state (key, counter, AES block) using system entropy,
	// personalization string, and optional additional input.
	newState, err := newDRBGState(d.config, additionalInput)
	if err != nil {
		return err
	}

	// Atomically install the new cryptographic state for this DRBG instance.
	d.state.Store(newState)

	// Acquire the counter mutex and update the working counter (v)
	// to match the new state, ensuring unique, non-overlapping output.
	d.vMu.Lock()
	copy(d.v[:], newState.v[:])
	d.vMu.Unlock()

	// Reset the usage counter, guaranteeing fresh key usage tracking.
	atomic.StoreUint64(&d.usage, 0)

	// Update reseed tracking metadata.
	d.lastReseedTime = time.Now()
	atomic.StoreUint64(&d.requests, 0)

	return nil
}

// reseedUpdate reseeds this instance as specified by NIST SP 800-90A §10.2.1.4.1 (no derivation function).
// Fresh entropy (key size + 16 bytes) is combined with the optional additional input (at most key size + 16
// bytes; shorter input is zero-padded) and mixed into the existing key and counter, so the prior state also
// contributes to the result. The personalization string is not used; it applies only when the instance is
// created. It is used when BacktrackingResistance is enabled.
func (d *drbg) reseedUpdate(additionalInput []byte) error {
	if len(additionalInput) > int(d.config.KeySize)+16 {
		return ErrAdditionalInputTooLong
	}
	if err := d.mixEntropy(additionalInput); err != nil {
		return err
	}

	// Restart byte-usage, reseed-interval, and request-count tracking from this reseed.
	atomic.StoreUint64(&d.usage, 0)
	d.lastReseedTime = time.Now()
	atomic.StoreUint64(&d.requests, 0)
	return nil
}

// mixEntropy obtains key size + 16 bytes of fresh entropy, combines it with the optional additional input,
// and mixes the result into the existing key and counter, as specified by NIST SP 800-90A §10.2.1.4.1. It
// does not change reseed bookkeeping, so it is safe to call from the background key-rotation goroutine.
func (d *drbg) mixEntropy(additionalInput []byte) error {
	src := d.entropySource
	if src == nil {
		src = rand.Reader
	}
	seedMaterial, err := readEntropy(src, d.config)
	if err != nil {
		return err
	}
	subtle.XORBytes(seedMaterial, seedMaterial, additionalInput)

	d.vMu.Lock()
	copy(d.encV[:], d.v[:])
	err = d.update(d.state.Load(), seedMaterial)
	d.vMu.Unlock()
	clear(seedMaterial)
	return err
}

// readEntropy returns key size + 16 bytes of entropy read from src.
func readEntropy(src io.Reader, cfg *Config) ([]byte, error) {
	seed := make([]byte, cfg.KeySize+16)
	if _, err := io.ReadFull(src, seed); err != nil {
		return nil, err
	}
	return seed, nil
}

// newDRBGState generates a fresh, fully initialized DRBG state (AES key and counter) per NIST SP 800-90A.
//
// This function combines system entropy, the configured personalization string, and optional additional input
// to derive a unique and unpredictable state. This process aligns with NIST recommendations for DRBG
// instantiation and reseed, ensuring domain separation and strong security.
//
// The returned state encapsulates an AES block cipher initialized with the derived key, a copy of the key
// material, and the initial 128-bit counter value (V). If entropy acquisition or cipher creation fails,
// an error is returned and no state is produced.
//
// Parameters:
//   - cfg *Config: The DRBG configuration, including key size, personalization, and related options.
//   - additionalInput []byte: Optional per-call entropy or context to further randomize the state; may be nil.
//
// Returns:
//   - *state: Newly derived DRBG state with a fresh key, cipher, and counter.
//   - error: Non-nil if entropy acquisition or cipher construction fails; nil on success.
func newDRBGState(cfg *Config, additionalInput []byte) (*state, error) {
	// Acquire fresh entropy from the operating system. This forms the basis of the DRBG seed material.
	seed, err := readEntropy(rand.Reader, cfg)
	if err != nil {
		return nil, err
	}

	// Incorporate the personalization string, if provided, by XOR-ing it into the seed for domain separation.
	// Mix in any caller-supplied additional input by XOR-ing it into the seed, further randomizing the state.
	mixSeed(seed, cfg.Personalization, additionalInput)

	// Split the seed buffer into an AES key (of the configured size) and a 128-bit counter (V).
	var key [32]byte
	copy(key[:], seed[:cfg.KeySize])
	var v [16]byte
	copy(v[:], seed[cfg.KeySize:])

	// Initialize the AES block cipher using the derived key. Return an error if cipher creation fails.
	block, err := aes.NewCipher(key[:cfg.KeySize])
	if err != nil {
		return nil, err
	}

	// Construct and return the new DRBG state, encapsulating the cipher, key, and counter.
	return &state{
		block: block,
		key:   key,
		v:     v,
	}, nil
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
//   - cfg: *Config — A pointer to the DRBG configuration (must be non-nil)
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
	mixSeed(seed, cfg.Personalization, nil)

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
		pid:      os.Getpid(),
	}
	d.state.Store(st)
	d.lastReseedTime = time.Now() // The reseed interval is measured from instantiation.

	// Initialize the working counter (v) from the state, guaranteeing unique output on first use.
	copy(d.v[:], v[:])

	// NIST SP 800-90A §10.2.1.3.1: derive the initial key and counter from the seed material.
	if cfg.BacktrackingResistance {
		if err := d.instantiate(seed); err != nil {
			return nil, err
		}
	}

	// FIPS 140-2 §4.9.2: the first block generated after initialization is not output; it is kept only for
	// comparison with the next generated block.
	if cfg.ContinuousHealthTest {
		if _, err := d.Read(d.lastOutputBlock[:]); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// mixSeed incorporates the personalization string and additional input into the DRBG seed.
//
// This function XORs both the personalization string (for domain separation) and any
// caller-provided additional input (for added entropy or context) into the given seed buffer,
// wrapping both as needed if their lengths exceed the seed length. This process is per NIST SP 800-90A.
//
// Parameters:
//   - seed: The entropy seed buffer to be mixed into (usually key+V).
//   - personalization: Optional domain-separation string; XOR-ed into the seed.
//   - additionalInput: Optional caller-supplied input; further XOR-ed into the seed.
//
// Each byte of personalization and additionalInput is XOR-ed into seed at the corresponding
// index, wrapping with modulo if necessary.
//
// Example usage:
//
//	mixSeed(seed, cfg.Personalization, additionalInput)
func mixSeed(seed []byte, personalization, additionalInput []byte) {
	// Incorporate the personalization string by XOR-ing it into the seed
	// to achieve domain separation. This ensures that DRBG instances with
	// different personalization values generate independent streams.
	for i := range personalization {
		seed[i%len(seed)] ^= personalization[i]
	}

	// Mix in any caller-supplied additional input (such as explicit entropy)
	// by XOR-ing it into the seed. This further randomizes the DRBG state,
	// enhancing uniqueness or context-sensitivity per NIST guidelines.
	for i := range additionalInput {
		seed[i%len(seed)] ^= additionalInput[i]
	}
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

	// With BacktrackingResistance, key rotation is a NIST SP 800-90A §10.2.1.4.1 reseed.
	if d.config.BacktrackingResistance {
		d.rekeyByReseed(base, maxBackoff)
		return
	}

	// Attempt to reseed and rekey up to MaxRekeyAttempts times.
	for i := 0; i < d.config.MaxRekeyAttempts; i++ {
		// Obtain new entropy for key and counter (V).
		seedLen := d.config.KeySize + 16 // Key size plus 128-bit counter
		seed := make([]byte, seedLen)
		if _, err := io.ReadFull(rand.Reader, seed); err == nil {
			// Apply personalization string, if set, by XORing into the seed.
			mixSeed(seed, d.config.Personalization, nil)

			// Construct the new AES key and counter (V) from the seed buffer.
			var key [32]byte
			copy(key[:], seed[:d.config.KeySize])
			var v [16]byte
			copy(v[:], seed[d.config.KeySize:])
			var block cipher.Block
			block, err = aes.NewCipher(key[:d.config.KeySize])
			if err == nil {
				// FIPS 140-2 §4.7.6: Zeroize old key material before replacement.
				if d.config.EnableZeroization {
					if oldState := d.state.Load(); oldState != nil {
						// Use subtle.XORBytes to prevent compiler optimization.
						subtle.XORBytes(oldState.key[:], oldState.key[:], oldState.key[:])
						subtle.XORBytes(oldState.v[:], oldState.v[:], oldState.v[:])
					}
				}

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
				// Zeroize old working counter before overwriting.
				if d.config.EnableZeroization {
					subtle.XORBytes(d.v[:], d.v[:], d.v[:])
				}
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

// rekeyByReseed performs automatic key rotation when BacktrackingResistance is enabled. Instead of replacing
// the key and counter, it mixes fresh entropy into them as specified by NIST SP 800-90A §10.2.1.4.1, retrying
// with the same exponential backoff as the default key rotation. Only the byte-usage counter is reset;
// reseed-interval and request-count tracking are owned by the reading goroutine and are left unchanged.
func (d *drbg) rekeyByReseed(backoff, maxBackoff time.Duration) {
	for i := 0; i < d.config.MaxRekeyAttempts; i++ {
		oldState := d.state.Load()
		if err := d.mixEntropy(nil); err == nil {
			// FIPS 140-2 §4.7.6: zeroize the replaced key material.
			if d.config.EnableZeroization && oldState != nil {
				subtle.XORBytes(oldState.key[:], oldState.key[:], oldState.key[:])
				subtle.XORBytes(oldState.v[:], oldState.v[:], oldState.v[:])
			}
			atomic.StoreUint64(&d.usage, 0)
			return
		}

		time.Sleep(backoff)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

// incV increments the DRBG counter (V) in big-endian order, rolling over as needed.
//
// The counter (V) is treated as a 128-bit unsigned integer in big-endian representation.
// Each call increments the counter by one, wrapping as appropriate. This function
// is used for advancing the DRBG keystream per SP 800-90A section on counter-mode.
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
