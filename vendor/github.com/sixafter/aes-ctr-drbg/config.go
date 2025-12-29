// Copyright (c) 2024-2025 Six After, Inc
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

import (
	"runtime"
	"time"
)

// KeySize represents the valid AES key lengths supported by AES-CTR-DRBG.
//
// It enforces compile-time type safety for AES key selection and helps
// prevent accidental misuse of invalid key lengths. Only KeySize128,
// KeySize192, and KeySize256 are permitted values for cryptographic configuration.
//
// Usage:
//
//	// Configure for AES-256 (32 bytes)
//	cfg := ctrdrbg.DefaultConfig()
//	cfg.KeySize = ctrdrbg.KeySize256
type KeySize int

const (
	// KeySize128 specifies AES-128 (16-byte key).
	KeySize128 KeySize = 16

	// KeySize192 specifies AES-192 (24-byte key).
	KeySize192 KeySize = 24

	// KeySize256 specifies AES-256 (32-byte key).
	KeySize256 KeySize = 32
)

const (
	// maxReseedInterval is the NIST SP 800-90A maximum reseed interval for CTR_DRBG (2^48).
	maxReseedInterval uint64 = 1 << 48
)

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

	// ReseedInterval is the minimum time duration between automatic reseeds from system entropy.
	//
	// When set (non-zero), the DRBG will automatically reseed after this interval elapses,
	// regardless of key usage or bytes generated. Zero disables interval-based reseeding.
	ReseedInterval time.Duration

	// MaxBytesPerKey is the maximum number of bytes generated per key before triggering automatic rekeying.
	//
	// Rekeying after a fixed output window enforces forward secrecy and mitigates key exposure risk.
	// If set to zero, a default value of 1 GiB (1 << 30) is used.
	MaxBytesPerKey uint64

	// ReseedRequests is the maximum number of output requests (calls to Read) allowed before forcing a reseed.
	//
	// When set (non-zero), the DRBG will automatically reseed after this many Read or ReadWithAdditionalInput calls.
	// Zero disables reseed-on-request-count.
	ReseedRequests uint64

	// ForkDetectionInterval controls how often fork detection is performed.
	//
	// If 0 (default), fork detection runs on every output request (max safety, fully compliant).
	// If >0, fork detection is performed once every N output requests (advanced tuning; reduces overhead
	// at the cost of a negligible window of risk).
	//
	// WARNING: Setting this above zero is NOT recommended for compliance-sensitive environments.
	ForkDetectionInterval uint64

	// KeySize specifies the AES key length to use for this DRBG instance.
	//
	// Acceptable values:
	//   - KeySize128 (16 bytes, AES-128)
	//   - KeySize192 (24 bytes, AES-192)
	//   - KeySize256 (32 bytes, AES-256)
	//
	// Default: KeySize256 (AES-256, 32 bytes).
	KeySize KeySize

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

	// Shards control the number of pools (shards) to use for parallelism.
	//
	// If zero, defaults to runtime.GOMAXPROCS(0).
	// Increase this to improve throughput under high concurrency.
	Shards int

	// EnableKeyRotation controls whether instances automatically rotate their key after MaxBytesPerKey output.
	//
	// Automatic key rotation provides forward secrecy and aligns with cryptographic best practices.
	// Defaults to true.
	EnableKeyRotation bool

	// PredictionResistance enables NIST SP 800-90A prediction resistance mode for this DRBG instance.
	//
	// When set to true, the DRBG will automatically reseed from fresh system entropy before every
	// output generation (each call to Read or ReadWithAdditionalInput), as required for compliance
	// with "prediction resistance" in §9.3 of SP 800-90A. This defeats state compromise
	// extension and protects against backtracking even if the DRBG's internal state is exposed.
	//
	// When enabled:
	//   - All additionalInput passed to ReadWithAdditionalInput or Reseed is ignored.
	//   - Reseeding is performed before each output operation, guaranteeing that every call mixes in
	//     new system entropy and cannot be predicted from previous outputs, even if internal state is known.
	//
	// When false (default), the DRBG operates in the normal mode, and reseeding occurs only at
	// initialization, on demand, or after MaxBytesPerKey is exceeded (if key rotation is enabled).
	//
	// Set via Config directly or by using the WithPredictionResistance functional option.
	//
	// Security rationale:
	//   - Enables maximal resilience against state compromise and forward/backward prediction attacks.
	//   - May increase system entropy usage and affect throughput in high-rate scenarios.
	PredictionResistance bool

	// UseZeroBuffer determines whether each Read operation uses a zero-filled buffer for AES-CTR output.
	//
	// If true, Read uses an internal buffer of zeroes for XOR operations (if the underlying implementation requires).
	// If false, output may be generated in place, which is typically faster and allocation-free.
	// Defaults to false.
	UseZeroBuffer bool

	// EnableSelfTests controls whether FIPS 140-2 Known Answer Tests (KAT) are run
	// on first use of the DRBG to verify AES-CTR is functioning correctly.
	//
	// When true, RunSelfTests() is called automatically on the first NewReader() invocation.
	// When false (default), self-tests are skipped for performance.
	EnableSelfTests bool

	// EnableZeroization controls whether cryptographic key material is securely
	// erased from memory during key rotation operations.
	//
	// When true, old keys and counters are zeroized using crypto/subtle before
	// being replaced, providing forward secrecy hardening per FIPS 140-2 section 4.7.6.
	// When false (default), old key material is simply overwritten.
	EnableZeroization bool
}

// Default configuration constants for AES-CTR-DRBG.
const (
	// Default max bytes per key (1 GiB)
	defaultMaxBytes = 1 << 30

	// Default max initialization retries
	defaultInitRetries = 3

	// Default max rekey attempts
	defaultRekeyRetries = 5

	// Default max backoff for rekey (2 seconds)
	defaultMaxBackoff = 2 * time.Second

	// Default initial rekey backoff (100 ms)
	defaultRekeyBackoff = 100 * time.Millisecond
)

// DefaultConfig returns a Config struct populated with production-safe, NIST SP 800-90A §10.2.1-aligned defaults.
//
// This function provides a robust baseline configuration for AES-CTR-DRBG instances, suitable for general-purpose
// cryptographic use and high-concurrency workloads. All parameters are selected to ensure strong security, compliance
// with FIPS 140-2 and NIST SP 800-90A requirements, and operational reliability under diverse system conditions.
//
// The returned configuration enables fork-safety (via PID tracking), supports domain separation via personalization,
// and is compatible with Go's FIPS-140 mode (GODEBUG=fips140=on).
//
// Defaults:
//   - KeySize:            32 bytes (AES-256, recommended by NIST for most use cases)
//   - MaxBytesPerKey:     1 GiB (1 << 30); triggers key rotation for forward secrecy
//   - MaxInitRetries:     3 attempts to initialize each DRBG pool entry
//   - MaxRekeyAttempts:   5 attempts per automatic key rotation
//   - MaxRekeyBackoff:    2 seconds (maximum exponential backoff between failed rekey attempts)
//   - RekeyBackoff:       100 milliseconds (initial backoff for rekey attempts)
//   - EnableKeyRotation:  false (key rotation is disabled by default—set to true for forward secrecy)
//   - EnableSelfTests:    false (FIPS 140-2 KAT self-tests disabled; enable via WithSelfTests for compliance)
//   - EnableZeroization:  false (key zeroization disabled; enable via WithZeroization for FIPS 140-2 compliance)
//   - Personalization:    nil (no domain separation unless set by the caller)
//   - UseZeroBuffer:      false (random output generated directly into caller's buffer)
//   - DefaultBufferSize:  0 (no preallocation of zero-filled buffers)
//   - Shards:             runtime.GOMAXPROCS(0) (number of internal DRBG pools matches available CPUs)
//   - PredictionResistance: false (prediction resistance is disabled; enable only if required by policy)
//   - ForkDetectionInterval: 0 (fork detection performed on every output request for maximum safety)
//
// NIST Reference:
//   - See NIST SP 800-90A, §10.2.1 (CTR DRBG) for cryptographic construction details.
//
// Example usage:
//
//	cfg := ctrdrbg.DefaultConfig()
func DefaultConfig() Config {
	return Config{
		KeySize:               KeySize256,
		MaxBytesPerKey:        defaultMaxBytes,
		MaxInitRetries:        defaultInitRetries,
		MaxRekeyAttempts:      defaultRekeyRetries,
		MaxRekeyBackoff:       defaultMaxBackoff,
		RekeyBackoff:          defaultRekeyBackoff,
		EnableKeyRotation:     false,
		Personalization:       nil,
		UseZeroBuffer:         false,
		DefaultBufferSize:     0,
		Shards:                runtime.GOMAXPROCS(0),
		PredictionResistance:  false,
		ForkDetectionInterval: 0,
	}
}

// Option defines a functional option for customizing a Config.
//
// Use Option values with NewReader or other constructors that accept variadic options.
//
// Example:
//
//	r, err := ctrdrbg.NewReader(
//	    ctrdrbg.WithKeySize(ctrdrbg.KeySize256),
//	    ctrdrbg.WithPersonalization([]byte("service-A")),
//	)
type Option func(*Config)

// WithKeySize returns an Option that sets the AES key size for this DRBG instance.
//
// Acceptable values are KeySize128 (16 bytes), KeySize192 (24 bytes), or KeySize256 (32 bytes).
// Any other value will cause NewReader to fail with an error at construction time.
func WithKeySize(k KeySize) Option { return func(cfg *Config) { cfg.KeySize = k } }

// WithMaxBytesPerKey returns an Option that sets the maximum number of bytes output per key before rekeying.
//
// This enforces a forward secrecy window: after MaxBytesPerKey random bytes are generated,
// the DRBG automatically performs a rekey operation (if key rotation is enabled) to derive a new key
// from fresh entropy. Lower this value to increase key rotation frequency for higher assurance.
func WithMaxBytesPerKey(n uint64) Option { return func(cfg *Config) { cfg.MaxBytesPerKey = n } }

// WithMaxInitRetries returns an Option that sets the maximum number of attempts to initialize a DRBG instance
// in the pool before failing.
//
// Increase this value if your system occasionally fails to gather entropy or encounters transient cryptographic errors
// at startup.
func WithMaxInitRetries(n int) Option { return func(cfg *Config) { cfg.MaxInitRetries = n } }

// WithMaxRekeyAttempts returns an Option that sets the maximum number of attempts allowed for
// asynchronous key rotation (rekey) in the DRBG.
//
// If all rekey attempts fail, the DRBG continues using the previous state. Exponential backoff is applied
// between attempts (see WithMaxRekeyBackoff and WithRekeyBackoff).
func WithMaxRekeyAttempts(n int) Option { return func(cfg *Config) { cfg.MaxRekeyAttempts = n } }

// WithMaxRekeyBackoff returns an Option that sets the maximum duration for exponential backoff between
// failed rekey attempts.
//
// When a rekey attempt fails, the DRBG waits with exponentially increasing intervals, up to this maximum duration,
// before retrying. If set to zero, a default (2s) is used.
func WithMaxRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) { cfg.MaxRekeyBackoff = d }
}

// WithRekeyBackoff returns an Option that sets the initial backoff duration before retrying a failed rekey operation.
//
// The first failure sleeps this duration, doubling for each subsequent failure, up to MaxRekeyBackoff.
func WithRekeyBackoff(d time.Duration) Option {
	return func(cfg *Config) { cfg.RekeyBackoff = d }
}

// WithEnableKeyRotation returns an Option that enables or disables automatic key rotation in the DRBG.
//
// When enabled (true), the DRBG automatically performs a key rotation after MaxBytesPerKey bytes are output.
// When disabled (false), key rotation is suppressed and the DRBG uses the original key indefinitely.
//
// For most use cases, leave this enabled for forward secrecy and NIST alignment. Disable only for
// compliance testing or special-purpose scenarios.
func WithEnableKeyRotation(enable bool) Option {
	return func(cfg *Config) { cfg.EnableKeyRotation = enable }
}

// WithPersonalization returns an Option that sets a per-instance personalization string
// to be XOR-ed into the DRBG's seed for domain separation.
//
// Personalization ensures that two DRBG instances constructed with the same system seed but different
// personalization values produce independent random streams, even if instantiated simultaneously.
//
// Use for tenant, application, service, or hardware isolation as required by your security model.
func WithPersonalization(p []byte) Option {
	return func(cfg *Config) { cfg.Personalization = p }
}

// WithUseZeroBuffer returns an Option that enables or disables use of a zero-filled buffer
// for output in the DRBG.
//
// If enabled (true), the DRBG allocates and maintains a reusable zero-filled buffer
// for CTR-mode output. If disabled (false), output is written directly to the destination buffer.
//
// This option primarily affects performance tuning; it does not impact cryptographic security.
func WithUseZeroBuffer(enable bool) Option {
	return func(cfg *Config) { cfg.UseZeroBuffer = enable }
}

// WithDefaultBufferSize returns an Option that sets the initial capacity of the internal zero buffer
// used for output if UseZeroBuffer is enabled.
//
// This can reduce allocations when large or repeated output requests are expected.
func WithDefaultBufferSize(n int) Option {
	return func(cfg *Config) { cfg.DefaultBufferSize = n }
}

// WithShards returns an Option that sets the number of internal pool shards for the DRBG.
//
// Sharding improves parallelism and reduces contention under high concurrency, at the cost of increased memory use.
// If n <= 0, the shard count defaults to runtime.GOMAXPROCS(0).
func WithShards(n int) Option {
	return func(cfg *Config) {
		if n <= 0 {
			n = runtime.GOMAXPROCS(0)
		}
		cfg.Shards = n
	}
}

// WithPredictionResistance returns an Option that enables or disables NIST SP 800-90A prediction resistance mode.
//
// When enabled (true), the DRBG reseeds from fresh system entropy before every output generation,
// making the output stream robust against state compromise extension and backtracking attacks.
//
// This mode increases system entropy usage and can impact performance in high-throughput scenarios.
// Use only when required by compliance or application policy.
func WithPredictionResistance(enable bool) Option {
	return func(cfg *Config) { cfg.PredictionResistance = enable }
}

// WithReseedInterval returns an Option that sets the minimum duration between automatic reseeds from system entropy.
//
// When set to a non-zero value, the DRBG will reseed after this interval elapses, even if no key rotation
// or manual reseed occurs. Set to zero to disable interval-based reseeding.
func WithReseedInterval(d time.Duration) Option {
	return func(cfg *Config) { cfg.ReseedInterval = d }
}

// WithReseedRequests returns an Option that sets the maximum number of output requests (calls to Read)
// allowed before forcing an automatic reseed from system entropy.
//
// Values exceeding the NIST SP 800-90A maximum (2^48) are clamped to the maximum.
// Set to zero to disable reseed-on-request-count behavior.
func WithReseedRequests(n uint64) Option {
	return func(cfg *Config) {
		if n > maxReseedInterval {
			n = maxReseedInterval
		}
		cfg.ReseedRequests = n
	}
}

// WithForkDetectionInterval sets the number of output requests between fork detection checks.
//
// WARNING: Setting this above zero introduces a window where a fork may not be detected immediately.
// Only set for performance-tuned applications that do NOT require strict compliance!
func WithForkDetectionInterval(n uint64) Option {
	return func(cfg *Config) { cfg.ForkDetectionInterval = n }
}

// WithSelfTests returns an Option that enables or disables FIPS 140-2 Known Answer Tests (KAT).
//
// When enabled, RunSelfTests() is called automatically on first NewReader() invocation
// to verify AES-CTR is functioning correctly before generating output.
// Defaults to false.
func WithSelfTests(enable bool) Option {
	return func(cfg *Config) { cfg.EnableSelfTests = enable }
}

// WithZeroization returns an Option that enables or disables secure key zeroization.
//
// When enabled, old cryptographic key material is securely erased from memory
// during key rotation operations per FIPS 140-2 section 4.7.6.
// Defaults to false.
func WithZeroization(enable bool) Option {
	return func(cfg *Config) { cfg.EnableZeroization = enable }
}
