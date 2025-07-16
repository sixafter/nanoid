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

// TestConfig_DefaultConfig verifies that DefaultConfig returns a Config
// with the documented default field values, such as MaxBytesPerKey (1GiB) and MaxInitRetries (3).
func TestConfig_DefaultConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey, "DefaultConfig.MaxBytesPerKey should be 1GiB")
	is.Equal(3, cfg.MaxInitRetries, "DefaultConfig.MaxInitRetries should be 3")
}

// TestConfig_WithMaxBytesPerKey ensures that the WithMaxBytesPerKey option
// correctly overrides the MaxBytesPerKey field in the Config, while leaving other fields unchanged.
func TestConfig_WithMaxBytesPerKey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	base := DefaultConfig()
	opt := WithMaxBytesPerKey(42)
	opt(&base)

	is.Equal(uint64(42), base.MaxBytesPerKey, "WithMaxBytesPerKey should override MaxBytesPerKey")
	is.Equal(3, base.MaxInitRetries, "WithMaxBytesPerKey should not affect MaxInitRetries")
}

// TestConfig_WithMaxInitRetries ensures that the WithMaxInitRetries option
// correctly sets the MaxInitRetries field in the Config, without modifying other fields.
func TestConfig_WithMaxInitRetries(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	base := DefaultConfig()
	opt := WithMaxInitRetries(7)
	opt(&base)

	is.Equal(7, base.MaxInitRetries, "WithMaxInitRetries should override MaxInitRetries")
	is.Equal(uint64(1<<30), base.MaxBytesPerKey, "WithMaxInitRetries should not affect MaxBytesPerKey")
}

// TestConfig_WithMaxRekeyAttempts checks that WithMaxRekeyAttempts updates
// the MaxRekeyAttempts field in the Config, leaving all other default values unchanged.
func TestConfig_WithMaxRekeyAttempts(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxRekeyAttempts(10)(&cfg)
	is.Equal(10, cfg.MaxRekeyAttempts, "WithMaxRekeyAttempts should override MaxRekeyAttempts")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
	is.Equal(3, cfg.MaxInitRetries)
	is.Equal(100*time.Millisecond, cfg.RekeyBackoff)
}

// TestConfig_WithRekeyBackoff verifies that WithRekeyBackoff updates the
// RekeyBackoff field while all other configuration values remain at their defaults.
func TestConfig_WithRekeyBackoff(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithRekeyBackoff(500 * time.Millisecond)(&cfg)
	is.Equal(500*time.Millisecond, cfg.RekeyBackoff, "WithRekeyBackoff should override RekeyBackoff")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
	is.Equal(3, cfg.MaxInitRetries)
	is.Equal(5, cfg.MaxRekeyAttempts)
}

// TestConfig_CombinedOptions ensures that multiple option functions can be
// combined and applied in sequence, with each option updating its respective field.
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

// TestConfig_WithZeroBuffer checks that the WithZeroBuffer option sets the
// UseZeroBuffer field as intended, without affecting unrelated configuration values.
func TestConfig_WithZeroBuffer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithZeroBuffer(true)(&cfg)
	is.True(cfg.UseZeroBuffer, "WithZeroBuffer(true) should set UseZeroBuffer to true")
	WithZeroBuffer(false)(&cfg)
	is.False(cfg.UseZeroBuffer, "WithZeroBuffer(false) should set UseZeroBuffer to false")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
}

// TestConfig_WithEnableKeyRotation validates that WithEnableKeyRotation
// correctly updates the EnableKeyRotation field in the Config structure.
func TestConfig_WithEnableKeyRotation(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithEnableKeyRotation(true)(&cfg)
	is.True(cfg.EnableKeyRotation, "WithEnableKeyRotation(true) should set EnableKeyRotation to true")
	WithEnableKeyRotation(false)(&cfg)
	is.False(cfg.EnableKeyRotation, "WithEnableKeyRotation(false) should set EnableKeyRotation to false")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
}

// TestConfig_WithDefaultBufferSize ensures that the WithDefaultBufferSize
// option modifies only the DefaultBufferSize field, and does not affect UseZeroBuffer.
func TestConfig_WithDefaultBufferSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithDefaultBufferSize(128)(&cfg)
	is.Equal(128, cfg.DefaultBufferSize, "WithDefaultBufferSize should override DefaultBufferSize")
	is.False(cfg.UseZeroBuffer)
}

// TestConfig_WithMaxRekeyBackoff checks that WithMaxRekeyBackoff updates
// the MaxRekeyBackoff field as expected, with other defaults left unchanged.
func TestConfig_WithMaxRekeyBackoff(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxRekeyBackoff(777 * time.Millisecond)(&cfg)
	is.Equal(777*time.Millisecond, cfg.MaxRekeyBackoff, "WithMaxRekeyBackoff should override MaxRekeyBackoff")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey)
	is.Equal(3, cfg.MaxInitRetries)
	is.Equal(5, cfg.MaxRekeyAttempts)
	is.Equal(100*time.Millisecond, cfg.RekeyBackoff)
}

// TestConfig_WithShards ensures that WithShards updates only the Shards field.
func TestConfig_WithShards(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithShards(8)(&cfg)
	is.Equal(8, cfg.Shards, "WithShards should override Shards")
}

// TestConfig_AllOptions verifies that all option functions can be composed
// and applied together, each updating their corresponding field in the Config struct.
func TestConfig_AllOptions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	opts := []Option{
		WithMaxBytesPerKey(777),
		WithMaxInitRetries(2),
		WithMaxRekeyAttempts(1),
		WithRekeyBackoff(77 * time.Millisecond),
		WithZeroBuffer(true),
		WithEnableKeyRotation(true),
		WithDefaultBufferSize(321),
		WithMaxRekeyBackoff(1234 * time.Millisecond),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	is.Equal(uint64(777), cfg.MaxBytesPerKey)
	is.Equal(2, cfg.MaxInitRetries)
	is.Equal(1, cfg.MaxRekeyAttempts)
	is.Equal(77*time.Millisecond, cfg.RekeyBackoff)
	is.True(cfg.UseZeroBuffer)
	is.True(cfg.EnableKeyRotation)
	is.Equal(321, cfg.DefaultBufferSize)
	is.Equal(1234*time.Millisecond, cfg.MaxRekeyBackoff)
}
