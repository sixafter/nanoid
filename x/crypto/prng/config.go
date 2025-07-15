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
	"time"
)

// Config holds tunable parameters for the PRNG pool and per-instance behavior.
//
// All fields are optional; zero values will use the library defaults.
type Config struct {
	// MaxBytesPerKey specifies the maximum number of bytes generated with a single
	// ChaCha20 key/nonce pair before triggering automatic key rotation. Once the
	// threshold is reached, the PRNG asynchronously rekeys and uses a fresh cipher
	// instance. If zero, the default threshold is 1 GiB.
	MaxBytesPerKey uint64

	// MaxInitRetries defines the number of times pool.New will retry initialization
	// before panicking. If zero, the default is 3.
	MaxInitRetries int

	// MaxRekeyAttempts controls how many times the key-rotation process will retry
	// generating a new cipher instance before giving up on the current rekey cycle.
	// If zero, the default is 5.
	MaxRekeyAttempts int

	// MaxRekeyBackoff clamps the maximum backoff duration between rekey attempts.
	// During rekeying, the retry interval is doubled on each failed attempt, but will
	// never exceed this maximum duration. If zero, the default is 2 seconds.
	MaxRekeyBackoff time.Duration

	// RekeyBackoff specifies the initial wait duration between key-rotation retry
	// attempts. On each retry, this interval is doubled (exponential backoff).
	// If zero, the default is 100ms.
	RekeyBackoff time.Duration

	// EnableKeyRotation controls whether the PRNG automatically rotates keys
	// after generating a configured amount of key material (see MaxBytesPerKey).
	// If false, key usage is not tracked and keys are never automatically rotated.
	// Defaults to false for maximum performance.
	EnableKeyRotation bool

	// UseZeroBuffer, when set to true, causes each Read operation to use a
	// zero-filled buffer for ChaCha20's XORKeyStream, rather than performing
	// in-place XOR. This provides legacy or compatibility behavior, but may
	// reduce performance. Defaults to false for maximum performance.
	UseZeroBuffer bool

	// DefaultBufferSize specifies the initial capacity of the internal buffer
	// used for zero-filled XOR operations. This is only relevant if UseZeroBuffer
	// is true. If zero, no preallocation is performed.
	DefaultBufferSize int
}

const (
	// maxRekeyAttempts defines the maximum number of times the PRNG will attempt
	// to generate and install a fresh ChaCha20 cipher during an automatic key rotation
	// before abandoning the current rekey cycle. If all attempts fail, the existing
	// cipher remains in use until the next rekey is triggered.
	maxRekeyAttempts = 5

	// rekeyBackoff specifies the initial duration to wait between consecutive
	// rekey attempts. The back-off duration is doubled on each subsequent retry
	// (exponential backoff) to reduce contention and load on the random source.
	rekeyBackoff = 100 * time.Millisecond

	// maxRekeyBackoff defines the maximum backoff duration between rekey attempts.
	// If the exponential backoff exceeds this value, it is clamped.
	maxRekeyBackoff = 2 * time.Second

	// maxBytesPerKey sets the default maximum number of bytes that may be generated
	// using a single ChaCha20 key/nonce pair before automatic key rotation is triggered.
	// This value is set to 1 << 30, which is approximately 1 GiB.
	maxBytesPerKey = 1 << 30

	// defaultBufferSize specifies the initial capacity (in bytes) of the internal
	// zero-filled buffer used for XOR operations when UseZeroBuffer is enabled.
	// If no value is provided in the configuration, this default is used.
	defaultBufferSize = 64
)

// DefaultConfig returns a Config with sensible defaults.
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
	}
}

// Option is a functional option for tweaking Config.
type Option func(*Config)

// WithMaxBytesPerKey overrides the bytes-per-key threshold.
func WithMaxBytesPerKey(n uint64) Option {
	return func(cfg *Config) {
		cfg.MaxBytesPerKey = n
	}
}

// WithMaxInitRetries overrides the pool init retry count.
func WithMaxInitRetries(r int) Option {
	return func(cfg *Config) {
		cfg.MaxInitRetries = r
	}
}

// WithMaxRekeyAttempts overrides the retry count for key rotation.
func WithMaxRekeyAttempts(r int) Option {
	return func(cfg *Config) {
		cfg.MaxRekeyAttempts = r
	}
}

// WithMaxRekeyBackoff sets the maximum allowed backoff duration for key rotation retries.
func WithMaxRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.MaxRekeyBackoff = d
	}
}

// WithRekeyBackoff overrides the initial back-off duration for key rotation retries.
func WithRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.RekeyBackoff = d
	}
}

// WithZeroBuffer enables the use of a zero-filled buffer for XORKeyStream.
// If set, each Read() will use an internal zero buffer (legacy/compatibility mode).
func WithZeroBuffer(enable bool) Option {
	return func(cfg *Config) {
		cfg.UseZeroBuffer = enable
	}
}

// WithEnableKeyRotation enables or disables per-key byte tracking and async rekeying.
func WithEnableKeyRotation(enable bool) Option {
	return func(cfg *Config) {
		cfg.EnableKeyRotation = enable
	}
}

// WithDefaultBufferSize preallocates the internal zero buffer to at least n bytes.
func WithDefaultBufferSize(n int) Option {
	return func(cfg *Config) {
		cfg.DefaultBufferSize = n
	}
}
