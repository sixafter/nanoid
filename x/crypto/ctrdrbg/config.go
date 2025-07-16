// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.
//
// Package ctrdrbg provides configuration types and functional options for the
// AES-CTR-DRBG (Deterministic Random Bit Generator) cryptographically secure pseudo-random number generator.
//
// The Config type exposes tunable parameters for the DRBG pool, instance management, and
// cryptographic behavior. These options support both security and operational flexibility.

package ctrdrbg

import "time"

// Config defines the tunable parameters for AES-CTR-DRBG instances and the DRBG pool.
//
// It supports fine-grained control over key size, key rotation, rekeying policies,
// backoff behavior, and instance personalization, enabling security-focused customization for a variety of use cases.
//
// Fields:
//   - KeySize: AES key length (16, 24, or 32 bytes for AES-128, -192, or -256).
//   - MaxBytesPerKey: Max output per key before automatic rekeying (forward secrecy).
//   - MaxInitRetries: Number of retries for DRBG pool initialization before panic.
//   - MaxRekeyAttempts: Max number of rekey attempts before giving up.
//   - MaxRekeyBackoff: Maximum backoff duration for exponential rekey retries.
//   - RekeyBackoff: Initial backoff for rekey attempts.
//   - EnableKeyRotation: Whether to enable automatic key rotation (default: true).
//   - Personalization: Optional per-instance byte string for domain separation.
type Config struct {
	// Personalization provides a per-instance personalization string, which is XOR-ed into the
	// DRBG’s initial seed to support domain separation or unique generator state.
	//
	// Purpose:
	// - Ensures cryptographic independence of DRBG streams even if seeds or environments overlap.
	// - Enables strong domain separation by context (service, user, tenant, device, etc.).
	//
	// Example:
	//   To ensure that two DRBGs used for "auth" and "billing" services are cryptographically isolated,
	//   pass unique byte strings (e.g., []byte("auth-service-v1") and []byte("billing-service-v1"))
	//   via WithPersonalization to their respective NewReader calls.
	//
	//   r1, _ := ctrdrbg.NewReader(ctrdrbg.WithPersonalization([]byte("auth-service-v1")))
	//   r2, _ := ctrdrbg.NewReader(ctrdrbg.WithPersonalization([]byte("billing-service-v1")))
	//
	// When unset (nil), no personalization is applied.
	Personalization []byte

	// RekeyBackoff is the initial delay before retrying a failed rekey operation.
	//
	// Exponential backoff doubles the delay for each failure up to MaxRekeyBackoff.
	// If set to zero, the default is 100 milliseconds.
	RekeyBackoff time.Duration

	// MaxRekeyBackoff specifies the maximum duration (clamped) for exponential backoff during rekey attempts.
	//
	// If set to zero, a default value of 2 seconds is used.
	MaxRekeyBackoff time.Duration

	// MaxBytesPerKey is the maximum number of bytes generated per key before triggering automatic rekeying.
	//
	// Rekeying after a fixed output window enforces forward secrecy and mitigates key exposure risk.
	// If set to zero, a default value of 1 GiB (1 << 30) is used.
	MaxBytesPerKey uint64

	// KeySize is the AES key length in bytes (16, 24, or 32).
	//
	// Valid values:
	//   - 16 (AES-128)
	//   - 24 (AES-192)
	//   - 32 (AES-256)
	//
	// Default: 32 (AES-256).
	KeySize int

	// MaxRekeyAttempts specifies the number of attempts to perform asynchronous rekeying.
	//
	// On failure, exponential backoff is used between attempts. If zero, a default of 5 is used.
	MaxRekeyAttempts int

	// MaxInitRetries is the maximum number of attempts to initialize a DRBG pool entry before giving up and panicking.
	//
	// Initialization can fail if system entropy is exhausted or if the cryptographic backend is unavailable.
	// If set to zero, a default of 3 is used.
	MaxInitRetries int

	// DefaultBufferSize specifies the initial capacity of the internal buffer used for zero-filled output operations.
	//
	// Only relevant if UseZeroBuffer is true. If zero, no preallocation is performed.
	DefaultBufferSize int

	// EnableKeyRotation controls whether DRBG instances automatically rotate their key after MaxBytesPerKey output.
	//
	// Automatic key rotation provides forward secrecy and aligns with cryptographic best practices.
	// Defaults to true.
	EnableKeyRotation bool

	// UseZeroBuffer determines whether each Read operation uses a zero-filled buffer for AES-CTR output.
	//
	// If true, Read uses an internal buffer of zeroes for XOR operations (if the underlying implementation requires).
	// If false, output may be generated in place, which is typically faster and allocation-free.
	// Defaults to false.
	UseZeroBuffer bool
}

// Default configuration constants for AES-CTR-DRBG.
const (
	defaultKeySize      = 32                     // Default AES key size (32 bytes for AES-256)
	defaultMaxBytes     = 1 << 30                // Default max bytes per key (1 GiB)
	defaultInitRetries  = 3                      // Default max initialization retries
	defaultRekeyRetries = 5                      // Default max rekey attempts
	defaultMaxBackoff   = 2 * time.Second        // Default max backoff for rekey (2 seconds)
	defaultRekeyBackoff = 100 * time.Millisecond // Default initial rekey backoff (100 ms)
)

// DefaultConfig returns a Config struct populated with production-safe, recommended defaults.
//
// Defaults:
//   - KeySize: 32 bytes (AES-256)
//   - MaxBytesPerKey: 1 GiB (1 << 30)
//   - MaxInitRetries: 3
//   - MaxRekeyAttempts: 5
//   - MaxRekeyBackoff: 2 seconds
//   - RekeyBackoff: 100 milliseconds
//   - EnableKeyRotation: true
//   - Personalization: nil (no domain separation)
//
// Example usage:
//
//	cfg := ctrdrbg.DefaultConfig()
func DefaultConfig() Config {
	return Config{
		KeySize:           defaultKeySize,
		MaxBytesPerKey:    defaultMaxBytes,
		MaxInitRetries:    defaultInitRetries,
		MaxRekeyAttempts:  defaultRekeyRetries,
		MaxRekeyBackoff:   defaultMaxBackoff,
		RekeyBackoff:      defaultRekeyBackoff,
		EnableKeyRotation: true,
		Personalization:   nil,
		UseZeroBuffer:     false,
		DefaultBufferSize: 0,
	}
}

// Option defines a functional option for customizing a Config.
//
// Use Option values with NewReader or other constructors that accept variadic options.
//
// Example:
//
//	r, err := ctrdrbg.NewReader(
//	    ctrdrbg.WithKeySize(32),
//	    ctrdrbg.WithPersonalization([]byte("service-A")),
//	)
type Option func(*Config)

// WithKeySize returns an Option that sets the AES key length in bytes.
//
// Acceptable values: 16 (AES-128), 24 (AES-192), 32 (AES-256).
func WithKeySize(n int) Option { return func(cfg *Config) { cfg.KeySize = n } }

// WithMaxBytesPerKey returns an Option that sets the maximum output (in bytes) per key before rekeying.
//
// Recommended to lower for higher security or compliance regimes.
func WithMaxBytesPerKey(n uint64) Option { return func(cfg *Config) { cfg.MaxBytesPerKey = n } }

// WithMaxInitRetries returns an Option that sets the maximum number of DRBG pool initialization retries.
//
// Use for customizing startup reliability and error handling.
func WithMaxInitRetries(n int) Option { return func(cfg *Config) { cfg.MaxInitRetries = n } }

// WithMaxRekeyAttempts returns an Option that sets the maximum number of retries allowed for asynchronous rekeying.
//
// Applies exponential backoff (see WithMaxRekeyBackoff/WithRekeyBackoff).
func WithMaxRekeyAttempts(n int) Option { return func(cfg *Config) { cfg.MaxRekeyAttempts = n } }

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

// WithPersonalization returns an Option that sets a per-instance personalization string for DRBG state separation.
//
// Rationale:
//   - Ensures domain separation, i.e., two DRBG instances with the same system seed but different personalization
//     strings will output completely different random streams.
//   - Use for tenant, user, application, or service isolation.
//
// Example:
//
//	ctrdrbg.NewReader(
//	    ctrdrbg.WithPersonalization([]byte("tenant-42-prod")),
//	)
func WithPersonalization(p []byte) Option {
	return func(cfg *Config) { cfg.Personalization = p }
}

// WithUseZeroBuffer returns an Option to enable or disable use of a zero-filled buffer for output.
func WithUseZeroBuffer(enable bool) Option {
	return func(cfg *Config) { cfg.UseZeroBuffer = enable }
}

// WithDefaultBufferSize returns an Option to set the default buffer size for zero-filled output.
func WithDefaultBufferSize(n int) Option {
	return func(cfg *Config) { cfg.DefaultBufferSize = n }
}
