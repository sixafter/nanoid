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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey, "DefaultConfig.MaxBytesPerKey should be 1GiB")
	is.Equal(3, cfg.MaxInitRetries, "DefaultConfig.MaxInitRetries should be 3")
}

func TestConfig_WithMaxBytesPerKey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	base := DefaultConfig()
	opt := WithMaxBytesPerKey(42)
	opt(&base)

	is.Equal(uint64(42), base.MaxBytesPerKey, "WithMaxBytesPerKey should override MaxBytesPerKey")
	// other field remains unchanged
	is.Equal(3, base.MaxInitRetries, "WithMaxBytesPerKey should not affect MaxInitRetries")
}

func TestConfig_WithMaxInitRetries(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	base := DefaultConfig()
	opt := WithMaxInitRetries(7)
	opt(&base)

	is.Equal(7, base.MaxInitRetries, "WithMaxInitRetries should override MaxInitRetries")
	// other field remains unchanged
	is.Equal(uint64(1<<30), base.MaxBytesPerKey, "WithMaxInitRetries should not affect MaxBytesPerKey")
}

func TestConfig_WithMaxRekeyAttempts(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxRekeyAttempts(10)(&cfg)
	is.Equal(10, cfg.MaxRekeyAttempts, "WithMaxRekeyAttempts should override MaxRekeyAttempts")
	// ensure other fields remain unchanged
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
	is.Equal(3, cfg.MaxInitRetries)
	is.Equal(100*time.Millisecond, cfg.RekeyBackoff)
}

func TestConfig_WithRekeyBackoff(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithRekeyBackoff(500 * time.Millisecond)(&cfg)
	is.Equal(500*time.Millisecond, cfg.RekeyBackoff, "WithRekeyBackoff should override RekeyBackoff")
	// ensure other fields remain unchanged
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
	is.Equal(3, cfg.MaxInitRetries)
	is.Equal(5, cfg.MaxRekeyAttempts)
}

func TestConfig_CombinedOptions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	opts := []Option{
		WithMaxBytesPerKey(99),
		WithMaxInitRetries(4),
		WithMaxRekeyAttempts(6),
		WithRekeyBackoff(250 * time.Millisecond),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	is.Equal(uint64(99), cfg.MaxBytesPerKey)
	is.Equal(4, cfg.MaxInitRetries)
	is.Equal(6, cfg.MaxRekeyAttempts)
	is.Equal(250*time.Millisecond, cfg.RekeyBackoff)
}
