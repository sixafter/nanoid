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

// Config holds tunable parameters for your PRNG pool and per-instance behavior.
type Config struct {
	// MaxBytesPerKey is the threshold (in bytes) at which a ChaCha20 instance
	// will rotate to a fresh key/nonce pair. Defaults to 1 GiB.
	MaxBytesPerKey uint64

	// MaxInitRetries is how many times pool.New will retry before
	// panicking. Defaults to 3.
	MaxInitRetries int

	// MaxRekeyAttempts defines how many times the key‐rotation process
	// will retry generating a fresh cipher before giving up on the current cycle.
	// Defaults to 5.
	MaxRekeyAttempts int

	// RekeyBackoff sets the initial wait duration between key‐rotation retry attempts.
	// This interval is doubled on each retry. Defaults to 100ms.
	RekeyBackoff time.Duration
}

const (
	// How many times to retry rekey before giving up on this cycle.
	maxRekeyAttempts = 5

	// Initial back-off duration between attempts.
	rekeyBackoff = 100 * time.Millisecond

	// maxBytesPerKey specifies the maximum number of bytes that may be generated
	// under a single ChaCha20 key/nonce pair before triggering a rotation.
	// It is set to 1 << 30 (approximately 1 GiB) by default.
	maxBytesPerKey = 1 << 30
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxBytesPerKey:   maxBytesPerKey,
		MaxInitRetries:   3,
		MaxRekeyAttempts: maxRekeyAttempts,
		RekeyBackoff:     rekeyBackoff,
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

// WithRekeyBackoff overrides the initial back-off duration for key rotation retries.
func WithRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.RekeyBackoff = d
	}
}
