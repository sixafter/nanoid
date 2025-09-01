// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.
//
// Package prng provides configuration types and functional options for the ChaCha20-based
// cryptographically secure pseudo-random number generator.
//
// The Config type exposes tunable parameters for the PRNG pool, instance management, and
// cryptographic behavior. These options support both security and operational flexibility.

package prng

import (
	"runtime"
	"time"
)

// Config defines the tunable parameters for ChaCha20-PRNG instances and the PRNG pool.
//
// It supports fine-grained control over key rotation, rekeying policies, buffer management,
// and operational backoff, enabling security-focused customization for a variety of use cases.
//
// Fields:
//   - MaxBytesPerKey: Max output per key before automatic rekeying (forward secrecy).
//   - MaxInitRetries: Number of retries for PRNG pool initialization before panic.
//   - MaxRekeyAttempts: Max number of rekey attempts before giving up.
//   - MaxRekeyBackoff: Maximum backoff duration for exponential rekey retries.
//   - RekeyBackoff: Initial backoff for rekey attempts.
//   - EnableKeyRotation: Whether to enable automatic key rotation (default: false).
//   - UseZeroBuffer: Whether to use a zero-filled buffer for ChaCha20 XORKeyStream.
//   - DefaultBufferSize: Initial internal buffer size for zero buffer operations.
type Config struct {
	// MaxBytesPerKey is the maximum number of bytes generated per key/nonce before triggering automatic rekeying.
	//
	// Rekeying after a fixed output window enforces forward secrecy and mitigates key exposure risk.
	// If set to zero, a default value of 1 GiB (1 << 30) is used.
	MaxBytesPerKey uint64

	// MaxInitRetries is the maximum number of attempts to initialize a PRNG pool entry before giving up and panicking.
	//
	// Initialization can fail if system entropy is exhausted or if the cryptographic backend is unavailable.
	// If set to zero, a default of 3 is used.
	MaxInitRetries int

	// MaxRekeyAttempts specifies the number of attempts to perform asynchronous rekeying.
	//
	// On failure, exponential backoff is used between attempts. If zero, a default of 5 is used.
	MaxRekeyAttempts int

	// MaxRekeyBackoff specifies the maximum duration (clamped) for exponential backoff during rekey attempts.
	//
	// If set to zero, a default value of 2 seconds is used.
	MaxRekeyBackoff time.Duration

	// RekeyBackoff is the initial delay before retrying a failed rekey operation.
	//
	// Exponential backoff doubles the delay for each failure up to MaxRekeyBackoff.
	// If set to zero, the default is 100 milliseconds.
	RekeyBackoff time.Duration

	// EnableKeyRotation controls whether PRNG instances automatically rotate their key/nonce after MaxBytesPerKey output.
	//
	// Automatic key rotation provides forward secrecy and aligns with cryptographic best practices.
	// Defaults to false for performance.
	EnableKeyRotation bool

	// UseZeroBuffer determines whether each Read operation uses a zero-filled buffer for ChaCha20's XORKeyStream.
	//
	// If true, Read uses an internal buffer of zeroes for output; if false, in-place XOR is used (faster).
	// Defaults to false.
	UseZeroBuffer bool

	// DefaultBufferSize specifies the initial capacity of the internal buffer used for zero-filled XOR operations.
	//
	// Only relevant if UseZeroBuffer is true. If zero, no preallocation is performed.
	DefaultBufferSize int

	// Shards controls the number of pools (shards) to use for parallelism.
	//
	// If zero, defaults to runtime.GOMAXPROCS(0).
	// Increase this to improve throughput under high concurrency.
	Shards int
}

// Default configuration constants for ChaCha20-PRNG.
const (
	maxRekeyAttempts  = 5                      // Default max rekey attempts
	rekeyBackoff      = 100 * time.Millisecond // Default initial rekey backoff (100 ms)
	maxRekeyBackoff   = 2 * time.Second        // Default max backoff for rekey (2 seconds)
	maxBytesPerKey    = 1 << 30                // Default max bytes per key (1 GiB)
	defaultBufferSize = 64                     // Default internal buffer size for XOR operations
)

// DefaultConfig returns a Config struct populated with production-safe, recommended defaults.
//
// Defaults:
//   - MaxBytesPerKey: 1 GiB (1 << 30)
//   - MaxInitRetries: 3
//   - MaxRekeyAttempts: 5
//   - MaxRekeyBackoff: 2 seconds
//   - RekeyBackoff: 100 milliseconds
//   - EnableKeyRotation: false
//   - UseZeroBuffer: false
//   - DefaultBufferSize: 64
//
// Example usage:
//
//	cfg := prng.DefaultConfig()
func DefaultConfig() Config {
	return Config{
		MaxBytesPerKey:    maxBytesPerKey,
		MaxInitRetries:    3,
		MaxRekeyAttempts:  maxRekeyAttempts,
		MaxRekeyBackoff:   maxRekeyBackoff,
		RekeyBackoff:      rekeyBackoff,
		UseZeroBuffer:     false,
		EnableKeyRotation: false,
		DefaultBufferSize: defaultBufferSize,
		// Ref: Use of GOMAXPROCS is fine for now: https://github.com/golang/go/issues/73193
		Shards: runtime.GOMAXPROCS(0),
	}
}

// Option defines a functional option for customizing a Config.
//
// Use Option values with NewReader or other constructors that accept variadic options.
//
// Example:
//
//	r, err := prng.NewReader(
//	    prng.WithEnableKeyRotation(true),
//	    prng.WithDefaultBufferSize(128),
//	)
type Option func(*Config)

// WithMaxBytesPerKey returns an Option that sets the maximum output (in bytes) per key before rekeying.
//
// Recommended to lower for higher security or compliance regimes.
func WithMaxBytesPerKey(n uint64) Option {
	return func(cfg *Config) { cfg.MaxBytesPerKey = n }
}

// WithMaxInitRetries returns an Option that sets the maximum number of PRNG pool initialization retries.
//
// Use for customizing startup reliability and error handling.
func WithMaxInitRetries(r int) Option {
	return func(cfg *Config) { cfg.MaxInitRetries = r }
}

// WithMaxRekeyAttempts returns an Option that sets the maximum number of retries allowed for asynchronous rekeying.
//
// Applies exponential backoff (see WithMaxRekeyBackoff/WithRekeyBackoff).
func WithMaxRekeyAttempts(r int) Option {
	return func(cfg *Config) { cfg.MaxRekeyAttempts = r }
}

// WithMaxRekeyBackoff returns an Option that sets the maximum duration for rekey exponential backoff.
//
// Limits time spent in failed rekey attempts.
func WithMaxRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) { cfg.MaxRekeyBackoff = d }
}

// WithRekeyBackoff returns an Option that sets the initial backoff duration for rekey retries.
//
// Initial sleep interval before exponential growth on rekey failure.
func WithRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) { cfg.RekeyBackoff = d }
}

// WithEnableKeyRotation returns an Option that enables or disables automatic key rotation.
//
// Disable only if you understand and accept the security risk.
func WithEnableKeyRotation(enable bool) Option {
	return func(cfg *Config) { cfg.EnableKeyRotation = enable }
}

// WithZeroBuffer returns an Option that enables or disables use of a zero-filled buffer for XORKeyStream.
//
// Enable only if required for legacy compatibility.
func WithZeroBuffer(enable bool) Option {
	return func(cfg *Config) { cfg.UseZeroBuffer = enable }
}

// WithDefaultBufferSize returns an Option that sets the initial zero buffer size.
//
// Only relevant if UseZeroBuffer is true.
func WithDefaultBufferSize(n int) Option {
	return func(cfg *Config) { cfg.DefaultBufferSize = n }
}

// WithShards sets the number of independent sync.Pool shards to use.
// By default, a single shard is used. Sharding may reduce contention
// under high concurrency, but can increase overhead on most systems.
//
// Note: If n <= 0, the number of shards defaults to runtime.GOMAXPROCS(0),
// which is useful in containerized environments.
// See: https://github.com/golang/go/issues/73193
func WithShards(n int) Option {
	return func(cfg *Config) {
		cfg.Shards = n
	}
}
