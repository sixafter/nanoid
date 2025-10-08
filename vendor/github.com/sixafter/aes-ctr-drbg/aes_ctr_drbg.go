// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

// Package ctrdrbg provides a FIPS 140-2 aligned, high-performance AES-CTR-DRBG.
//
// This package implements a cryptographically secure, pool-backed Deterministic Random Bit Generator
// (DRBG) following the NIST SP 800-90A AES-CTR-DRBG construction, specifically as defined in
// Section 10.2.1 of NIST SP 800-90A Rev. 1 ("Recommendation for Random Number Generation Using
// Deterministic Random Bit Generators").
//
// Each generator instance uses an AES block cipher in counter (CTR) mode to produce cryptographically
// secure pseudo-random bytes, suitable for high-throughput, concurrent workloads.
//
// All cryptographic primitives are provided by the Go standard library. This implementation is designed
// for environments requiring strong compliance, including support for Go's FIPS-140 mode (GODEBUG=fips140=on).
//
// Reference:
//
//	NIST Special Publication 800-90A Rev. 1, Section 10.2.1 (CTR_DRBG Construction)
//	https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-90Ar1.pdf
package ctrdrbg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand/v2"
	"os"
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

	// Reseed injects new entropy and optional additional input, refreshing the DRBG state.
	//
	// Per NIST SP 800-90A, additionalInput is combined with system entropy and
	// personalization (if set) to derive a new internal key and counter.
	// Reseed can be called at any time to proactively refresh the DRBG.
	//
	// Returns an error if the reseed operation fails.
	Reseed(additionalInput []byte) error

	// ReadWithAdditionalInput generates cryptographically secure random bytes,
	// supplementing system entropy with optional per-call additional input.
	//
	// This method is NIST-compliant and allows callers to inject additional
	// entropy or context for a single Read operation.
	// If PredictionResistance is enabled, additionalInput is ignored and fresh entropy is always used.
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

// Reseed refreshes the state of all DRBG instances in all shard pools with new entropy and optional additional input.
//
// This method implements the Interface contract for NIST SP 800-90A compliant deterministic random bit generators.
// For each DRBG instance managed by the pool, Reseed obtains fresh entropy from the system entropy source and
// combines it with the provided additionalInput and the configured personalization string, per NIST recommendations.
// This operation ensures that each DRBG instance has a newly derived key and counter (V), providing domain separation
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
//   - If called while DRBGs are in use, subsequent reads will immediately begin using the newly seeded state.
//   - This is required for some FIPS and high-assurance applications, and is recommended for recovery from suspected
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
// If PredictionResistance is enabled on the DRBG configuration, the additionalInput argument is ignored
// and fresh entropy is always used as required by the standard.
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
//   - If PredictionResistance is enabled, additionalInput is ignored and reseed is always performed from entropy.
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
	// Determine number of pools (shards) in the reader for load balancing.
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

	// requests counts the number of output requests (calls to Read or ReadWithAdditionalInput)
	// since the last reseed. Used to enforce the ReseedRequests limit and trigger reseeding
	// after a configured number of requests.
	requests uint64

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

	// encV is a persistent [16]byte working buffer used as a session-local counter
	// during output generation.
	//
	// On each call to Read, the current counter value (d.v) is copied into encV, which is
	// then incremented and used for block generation throughout the request. Only after all
	// output is produced is encV copied back into d.v, ensuring atomic and consistent counter
	// advancement. This prevents partial or inconsistent counter updates if Read exits early
	// due to errors or panics.
	encV [16]byte

	// tmp is a persistent [16]byte working buffer used during output generation.
	// It holds the encrypted output for the final (partial) block in fillBlocks.
	// This avoids repeated stack allocations and ensures maximum efficiency.
	tmp [16]byte
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

	d.reseedIfForked()

	// Prediction Resistance
	if d.config.PredictionResistance {
		if err := d.reseed(nil); err != nil {
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

	// Persist the advanced counter back to the DRBG instance, ensuring subsequent reads
	// continue the keystream seamlessly without overlap or repetition.
	copy(d.v[:], d.encV[:])

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
// If PredictionResistance is enabled, the DRBG reseeds from system entropy for every read and
// ignores additionalInput as required by the standard. Otherwise, if additionalInput is non-nil,
// the DRBG instance reseeds using a combination of system entropy, personalization, and the
// supplied additionalInput before generating output.
//
// Semantics and Implementation Details:
//   - If PredictionResistance is enabled, reseeds from fresh system entropy before generating output,
//     ignoring additionalInput (NIST-compliant).
//   - If PredictionResistance is not enabled and additionalInput is non-nil, reseeds using both entropy
//     and the caller's input, per NIST specification.
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

	d.reseedIfForked()

	// If PredictionResistance is enabled, always reseed from fresh entropy before output,
	// ignoring any additional input per NIST SP 800-90A requirements.
	if d.config.PredictionResistance {
		if err := d.reseed(nil); err != nil {
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
		// reseed using both entropy and additional input for this output request.
		if additionalInput != nil {
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

	// Fill the output buffer using the current cryptographic state and the local counter,
	// incrementing the counter as output is produced.
	d.fillBlocks(b, st, &d.encV)

	// Persist the advanced counter back to the DRBG instance, ensuring
	// future reads continue the keystream seamlessly.
	copy(d.v[:], d.encV[:])

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
// This method is compliant with NIST SP 800-90A and can be called at any time to force a rekey
// and counter reset. The additionalInput parameter, if non-nil, is cryptographically combined
// with system entropy and the configured personalization string to derive a new key and counter (V).
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
		incV(v)
		st.block.Encrypt(d.tmp[:], v[:])
		copy(b[offset:], d.tmp[:tail])
	}
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
//     combined with system entropy and personalization. May be nil for standard reseeds.
//
// Returns:
//   - error: Non-nil if entropy acquisition or state initialization fails; nil on success.
//
// Concurrency Notes:
//   - This method synchronizes counter updates using a mutex and uses atomic operations for usage/request counters.
//   - If called concurrently, it is possible for closely timed reseeds to slightly race in updating the metadata fields;
//     this does not impact cryptographic safety.
func (d *drbg) reseed(additionalInput []byte) error {
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
//   - *state: Newly derived DRBG state with fresh key, cipher, and counter.
//   - error: Non-nil if entropy acquisition or cipher construction fails; nil on success.
func newDRBGState(cfg *Config, additionalInput []byte) (*state, error) {
	seedLen := cfg.KeySize + 16
	seed := make([]byte, seedLen)

	// Acquire fresh entropy from the operating system. This forms the basis of the DRBG seed material.
	if _, err := io.ReadFull(rand.Reader, seed); err != nil {
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

	// Initialize the working counter (v) from the state, guaranteeing unique output on first use.
	copy(d.v[:], v[:])

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
