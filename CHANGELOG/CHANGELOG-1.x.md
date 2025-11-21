# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Date format: `YYYY-MM-DD`

---
## [Unreleased]

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

---

## [1.57.0] - 2025-11-21

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security
- **risk:** Go module verification accepts a specific version of Go for checksum verification.
  - Check a specific tag: `TAG=v1.56.0 make module-verify`
  - Check the latest version: `make module-verify`

---

## [1.56.0] - 2025-11-20

### Added
- **risk**: Added signature verification make target to match the README instructions.
- **risk:** Added go module verification make target to verify module checksums.

### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
- **defect:** Fixed `README.md` instructions for verifying module checksums.

### Security
- **risk:** Upgraded `golang.org/x/crypto` to `v0.45.0` to address vulnerabilities.

---

## [1.55.0] - 2025-11-07

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **debt:** Updated documentation and Go-doc comments.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.54.0] - 2025-10-16

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **debt:** Updated documentation and Go-doc comments.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.53.0] - 2025-10-08

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **debt:** Updated documentation and Go-doc comments.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.52.0] - 2025-09-30

### Added
### Changed
- **debt:** Upgraded [PRNG-CHACHA](https://github.com/sixafter/prng-chacha) to latest stable version.
- **debt:** Upgraded [AES-CTR-DRBG](https://github.com/sixafter/aes-ctr-drbg) to latest stable version.
- **debt:** Upgraded dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.51.0] - 2025-09-15

### Added
### Changed
- **debt:** Upgraded [PRNG-CHACHA](https://github.com/sixafter/prng-chacha) to latest stable version.
- **debt:** Upgraded [AES-CTR-DRBG](https://github.com/sixafter/aes-ctr-drbg) to latest stable version.
- **debt:** Upgraded dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.50.0] - 2025-09-10

### Added
- **feature:** Added `WithAutoRandReader` option to `Config` to automatically select between the `AES-CTR-DRBG` and `PRNG-CHACHA` implementations based on FIPS-140 mode.
  - When FIPS-140 mode is [enabled](https://pkg.go.dev/crypto/fips140#Enabled), the [AES-CTR-DRBG](https://github.com/sixafter/aes-ctr-drbg) implementation is used.
  - When FIPS-140 mode is disabled, the [PRNG-CHACHA](https://github.com/sixafter/prng-chacha) implementation is used.
  - This option simplifies configuration by automatically selecting the appropriate random number generator based on the security requirements of the environment.

### Changed
### Deprecated
### Removed
### Fixed

### Security

---
## [1.49.0] - 2025-09-09

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Fixed release `v1.48.0` tag to point to the correct commit.

### Security

---
## [1.48.0] - 2025-09-09

### Added
### Changed
- **debt:** Upgraded [prng-chacha](https://github.com/sixafter/prng-chacha) to [v1.4.0](https://github.com/sixafter/prng-chacha/releases/tag/v1.4.0).
- **debt:** Upgraded dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.47.0] - 2025-09-01

### Added
### Changed
- **debt:** Upgraded [prng-chacha](https://github.com/sixafter/prng-chacha) to [v1.3.0](https://github.com/sixafter/prng-chacha/releases/tag/v1.3.0).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.46.0] - 2025-09-01

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **risk:** Updated copyright to reflect date range through present year.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.45.0] - 2025-08-14

### Added
### Changed
- **debt:** Updated to Go `1.25` to leverage the latest language features and performance improvements.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.44.0] - 2025-07-19

### Added
### Changed
- **debt:** Moved the `x/crypto/prng` package to the [prng-chacha](https://github.com/sixafter/prng-chacha) repository.
  - This change allows for more focused development and maintenance of the `PRNG-CHACHA` implementation, separating it from the Nano ID project.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.43.0] - 2025-07-18

### Added
### Changed
- **debt:** Moved the `x/crypto/ctrdrbg` package to the [aes-ctr-drbg](https://github.com/sixafter/aes-ctr-drbg) repository.
  - This change allows for more focused development and maintenance of the `AES-CTR-DRBG` implementation, separating it from the Nano ID project.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.42.0] - 2025-07-18

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Update `goreleaser` configuration to use the author's username in the changelog.

### Security

---
## [1.41.0] - 2025-07-18

### Added
- **feature:** Added `Config()` method to both `prng` and `ctrdrbg` implementations to retrieve the current configuration of the generator.

### Changed
- **debt:** Updated documentation to reflect consistency across `prng` and `ctrdrbg` implementations.
- **debt:** Updated documentation to clarify FIPS-140 usage and compliance.

### Deprecated
### Removed
### Fixed
- **defect**: Guarantee unique cryptographic stream per [`ctrdrbg.Reader`](../x/crypto/ctrdrbg) instance by persisting and synchronizing internal counter state across all `Read` operations.

### Security

---
## [1.40.0] - 2025-07-17

### Added
### Changed
- **feature:** - Added support configurable sharded `sync.Pool` instances in [`prng`](../x/crypto/prng) and [`ctrdrbg`](../x/crypto/ctrdrbg) for improved concurrency and throughput.
  - The number of shards defaults to `runtime.GOMAXPROCS(0)`, which is useful in containerized or CPU-constrained environments pending https://github.com/golang/go/issues/73193.

### Deprecated
### Removed
### Fixed
- **defect:** Fixed issue [64](https://github.com/sixafter/nanoid/issues/64) in shard `sync.Pool` logic where the index was off by 1, causing unexpected behavior when shards are greater than 1.

### Security

---
## [1.39.0] - 2025-07-16

### Added
- **feature:** FIPS‑140 mode compatibility.
  - Adds support for cryptographic operations in environments requiring FIPS‑140 validation.
  - All cryptographic primitives are sourced exclusively from the Go standard library.
  - No third-party or non-standard cryptography is included when FIPS mode is enabled.
  - See [FIPS‑140.md](FIPS-140.md) for configuration, deployment recommendations, and compliance details.
  - The [`x/crypto/ctrdrbg`](../x/crypto/ctrdrbg) subpackage provides a deterministic random bit generator compatible with FIPS‑140 mode.
- **feature:** Added dedicated Make targets for benchmarking `ctrdrbg`.
  - Added `bench-ctrdrbg` to run raw AES-CTR-DRBG benchmark tests with memory and CPU profiling.
  - Added `bench-ctrdrbg-uuid` for benchmarking UUID generation using the AES-CTR-DRBG with the Google UUID package.

### Changed
- **debt:** Enhanced documentation for [`x/crypto/prng`](../x/crypto/prng).
  - Updated and expanded GoDoc comments across all public and internal types and functions.
  - Improved usage examples and configuration guidance for the `prng` package.
  - Increased consistency and symmetry with the `ctrdrbg` package documentation.
  - Clarified package guarantees and recommended usage patterns.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.38.0] - 2025-07-14

### Added
- **feature:** Added `MaxRekeyBackoff` [config](../x/crypto/prng/config.go) option to clamp maximum exponential backoff duration between key rotation attempts.
- **debt:** Added `BenchmarkPRNG_Concurrent_SyncPool_Baseline` benchmark test to measure performance of the `sync.Pool` as a baseline for comparison against the `x/crypto/prng` implementation.

### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Corrected `bench-uuid` and `bench-csprng` parameters to ensure the correct `cpu.out` and `mem.out` files are generated.

### Security

---
## [1.37.0] - 2025-07-13

### Added
- **feature:** Added [`EnableKeyRotation`](../x/crypto/prng/config.go)] option to automatically rotate PRNG keys after a configurable number of bytes, improving key hygiene for long-lived instances.
- **feature:** Added [`UseZeroBuffer`](../x/crypto/prng/config.go) option to support legacy XORKeyStream behavior with a zero-filled buffer on each read.
- **feature:** Added [`DefaultBufferSize`](../x/crypto/prng/config.go) option to control initial allocation size for the internal zero buffer when UseZeroBuffer is enabled

### Changed
- **debt:** Upgraded dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.36.0] - 2025-07-12

### Added
- **risk:** Eagerly initialize PRNG pool in NewReader to surface initialization errors at construction time. This addresses issue [#57](https://github.com/sixafter/nanoid/issues/57).
- **risk:** Added refined steps to verify the signature of the release artifacts using [`cosign`](https://github.com/sigstore/cosign).

### Changed
- **debt:** Modified benchmark tests to favor use of [`b.Loop()`](https://go.dev/blog/testing-b-loop) over `b.N` to ensure consistent performance across measurements.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.35.0] - 2025-07-10

### Added
- **risk:** Added benchmark tests for [x/crypto/prng](../x/crypto/prng) to measure performance of Cryptographically Secure Pseudo-Random Number Generator (CSPRNG) Reader.
- **debt:** Added `bench-csprng` `Makefile` target for running [x/crypto/prng](../x/crypto/prng) benchmark tests.

### Changed
- **debt:** Updated [README](README.md) to include information about the new `bench-csprng` target and benchmark tests for [x/crypto/prng](../x/crypto/prng).

### Deprecated
### Removed
### Fixed
### Security

---
## [1.34.0] - 2025-07-09

### Added
### Changed
- **debt:** Refactored `Read` to `0 allocs/op` for `io.Reader` resulting in a materially significant performance improvement.
- **debt:** Refactored `New` and `NewWithLength` to use `Read` while maintaining `1 allocs/op`.
- **debt:** Refactored benchmark tests with more appropriate naming.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.33.0] - 2025-06-29

### Added
### Changed
- **debt:** GoReleaser configuration now includes only necessary files.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.32.1] - 2025-06-28

### Added
- **debt:** Added OpenSSF Best Practices badge to [README](README.md).

### Changed
### Deprecated
### Removed
### Fixed

### Security
- **risk:** Add digital signatures for release source and checksums files.

---
## [1.32.0] - 2025-06-28

### Added
- **debt:** Added `release-verify` make target to optimize break-fix cycle with GoReleaser testing.

### Changed
### Deprecated
### Removed
### Fixed

### Security
- **risk:** Add digital signatures for release source and checksums files.

---
## [1.31.2] - 2025-06-28

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Added archive stanza to Go Releaser config to ensure signature files are released.

### Security

---
## [1.31.1] - 2025-06-28

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Added Cosign step in release job.

### Security

---
## [1.31.0] - 2025-06-28

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
- **risk:** Add digital signatures for release source and checksum files.

---
## [1.30.0] - 2025-06-13

### Added
### Changed
- **debt**: Modified tests to reflect the new `Config()` method in the `Interface` interface.

### Deprecated
### Removed
### Fixed
### Security
- **risk:** Seed zeroing: ChaCha20 key and nonce material are immediately cleared from memory after cipher creation, preventing any residual secret data.
- **risk:** Fail-fast initialization: The pool uses a capped-retry loop and panics on repeated failures, guaranteeing that no insecure or uninitialized PRNG state is ever used.
- **risk:** Automatic key rotation: Keys rotate after a configurable byte threshold (default 1 GiB), bounding keystream lifetime and eliminating reuse vulnerabilities.
- **risk:** Jittered back-off: Retry delays for key rotation employ exponential back-off combined with cryptographically sourced jitter, avoiding synchronized retry storms under low-entropy conditions.
- **risk:** In-memory cipher wipe: After a successful rotation, the previous ChaCha20 instance is overwritten with its zero value, eradicating any leftover key schedule or counter state.
- **risk:** Concurrency safety: Each goroutine checks out its own prng from a sync.Pool, ensuring exclusive access to the ChaCha20 cipher and preventing internal state corruption.

---
## [1.29.0] - 2025-06-12

### Added
- **feature**: Added `Must()` and `MustWithLength()` functions to `Interface` interface for safe ID generation without error handling.
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.28.0] - 2025-05-31

### Added
- **feature**: Added SonarQube security rating badge to the README.

### Changed
- **debt**: Changed `GetConfig()` to `Config()` in the `Interface` interface to align with Go naming conventions.
- **risk**: Added nil receiver test for `ID`.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.27.0] - 2025-05-16

### Added
- **feature**: The `Interface` interface now includes a `NewWithlength` method to create a Nano ID with a specified length.

### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.26.0] - 2025-05-16

### Added
- **feature**: The `Interface` interface now includes a `GetConfig()` method to retrieve the current configuration of the generator.

### Changed
- **debt:** Upgraded dependencies to the latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.25.0] - 2025-05-13

### Added
### Changed
- **debt**: Refactored `ID.String` and `ID.Compare` methods to use value receivers for improved clarity and consistency with Go idioms. Pointer receivers remain for methods requiring nil checks or mutation.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.24.1] - 2025-05-04

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Fixed minor variable shadowing issues in several `x/crypto/prng` tests.

### Security

---
## [1.24.0] - 2025-04-14

### Added
- **risk:** Added copyright notice to GitHub Actions workflows.

### Changed
- **debt**: Upgraded all dependencies to the latest supported versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.23.0] - 2025-02-13

### Added
### Changed
- **debt**: Upgraded to Go 1.24.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.22.0] - 2024-12-26

### Added
### Changed
- **debt**: Upgraded all Go dependencies to the latest versions.
- **debt:** Upgraded the CI pipeline to use the new GitHub Action for SonarQube Cloud analysis.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.21.0] - 2024-12-07

### Added
### Changed
- **debt**: Minor refactoring to improve code readability and maintainability.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.20.1] - 2024-11-24

### Added
### Changed
- **debt**: Upgraded all Go dependencies to the latest versions.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.20.0] - 2024-11-16

### Added
### Changed
- **FEATURE:** Updated the [PRNG](../x/crypto/prng) reader to allow for each reader instance to have its own `sync.Pool` for buffer reuse.
- **DEBT:** Refactored to not use `strings.Builder` for ID generation in favor of `sync.Pool` for buffer reuse.
- **DEBT:** Modified [README.md](README.md) to include detailed information about `sync.Pool` usage for buffer reuse for both ASCII and Unicode ID generation.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.19.0] - 2024-11-16

### Added
- **FEATURE:** Added [README.md](../x/crypto/prng/README.md) to provide detailed information about the PRNG implementation.

### Changed
- **DEBT:** Updated [PRNG](../x/crypto/prng) benchmark tests to test the standard size of 21 characters for Nano ID generation.
- **DEBT:** The runtime configuration (`Config`) now uses pointer receivers for all methods to ensure consistent behavior and interface compliance.
- **DEBT:** Refactored Codebase: Split the `nanoid.go` file into multiple modular files within the `nanoid` package to enhance code organization, readability, and maintainability

### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed receiver types by updating all `ID` methods to use pointer receivers consistently, ensuring proper functionality and interface compliance.

### Security

---
## [1.18.1] - 2024-11-15

### Added
### Changed
- **DEBT:** Added missing license header to the CodeQL analysis configuration file.
- **DEBT:** Refactored CHANGELOG date format to `YYYY-MM-DD`.

### Deprecated
### Removed
### Fixed

---
## [1.18.0] - 2024-11-15

### Added
- **FEATURE**: Added support for `fmt.Stringer`: Provides a string representation of the ID type.
- **FEATURE**: Added support for `encoding.TextMarshaler`: Supports marshaling ID into a text-based representation.
- **FEATURE**: Added support for `encoding.TextUnmarshaler`: Supports unmarshaling ID from a text-based representation.
- **FEATURE**: Added support for `encoding.BinaryMarshaler`: Supports marshaling ID into a binary representation.
- **FEATURE**: Added support for `encoding.BinaryUnmarshaler`: Supports unmarshaling ID from a binary representation.
- **FEATURE**: Added `DefaultRandReader` to provide a default random reader for generating IDs.
- **FEATURE**: Added `EmptyID` to provide an empty ID constant.
- **RISK**: Added support for CodeQL analysis when pushing to the `main` branch.

### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Addressed various documentation issues.

### Security

---
## [1.17.3] - 2024-11-14

### Added
### Changed
- **DEBT:** Refactored documentation to improve readability and clarity.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.17.2] - 2024-11-12

### Added
### Changed
- **DEFECT:** Corrected a few typos in the documentation.
- **DEBT:** Refactored documentation to improve readability and clarity.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.17.1] - 2024-11-12

### Added
### Changed
- **DEBT:** Refactored "godoc" comments to improve readability and clarity.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.17.0] - 2024-11-11

### Added
- **FEATURE:** Introduced new cryptographically secure Pseudo-Random Number Generator (PRNG) Reader based on the `ChaCha20` stream cypher to enhance random data generation capabilities.
  - The `crypto/rand.Reader` generates the key and nonce for each instance.
  - Reduced `ns/op` by approximately **93%** (e.g., ~323 `ns/op` to ~23 `ns/op` for a default ID of length 21).
  - Allocations per operation remain consistent at **1 allocs/op**. The new `io.Reader` maintains zero `allocs/op`.

### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.16.1] - 2024-11-09

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Removed unnecessary `//go:inline` directives that were left in inadvertently.
### Security

---
## [1.16.0] - 2024-11-09

### Added
### Changed
- **FEATURE:** Unicode Alphabet Optimization:
  - Reduced bytes per operation (`B/op`) by up to 33% for Unicode-based ID generation by introducing dynamic buffer sizing.
  - Reduced Execution Time (`ns/op`): Achieved up to 12.2% speed improvement in specific Unicode benchmark scenarios.
- **DEBT:** ASCII Alphabet Optimization:
  - Ensured consistent bytes per operation (`B/op`) aligned with ID lengths for ASCII alphabets, preserving memory efficiency without increasing allocations.
  - Reduced Execution Time (`ns/op`): Achieved up to 8.5% speed improvement in specific ASCII benchmark scenarios.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.15.0] - 2024-11-07

### Added
- **FEATURE:** Added documentation and comments.
  - Enhanced documentation for interfaces and functions by providing comprehensive comments on each method's purpose and usage.
  - Added detailed comments to the `buildRuntimeConfig` function variables for improved code clarity and understanding.
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.14.0] - 2024-11-06

### Added
- **FEATURE:** Added support for `io.Reader` Interface: 
  - The Nano ID generator now satisfies the `io.Reader` interface, allowing it to be used interchangeably with any `io.Reader` implementations.
  - Developers can now utilize the Nano ID generator in contexts such as streaming data processing, pipelines, and other I/O-driven operations.
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.13.4] - 2024-11-06

### Added
### Changed
- **DEBT:** Added `//go:inline` directive to optimize specific functions.
- **DEBT:** Updated benchmark tests to reflect optimizations.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.13.3] - 2024-11-0

### Added
- **DEBT:** Added test for `panic` in `MustWithLength` function.
### Changed
### Deprecated
### Removed
- **DEFECT:** ASCII-only detection now correctly uses `unicode.MaxASCII` instead of `0x7F` to ensure compatibility with all ASCII characters.
### Fixed
### Security

---
## [1.13.2] - 2024-11-06

### Added
### Changed
### Deprecated
### Removed
- **DEBT:** Removed unused `Buffer` interface:
    ```go
    // Buffer is a type constraint that allows either []byte or []rune.
    type Buffer interface {
       ~[]byte | ~[]rune
    }
    ```
### Fixed
### Security

---
## [1.13.1] - 2024-11-05

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Refactored ID generation functions `newASCII` and `newUnicode` by removing calls to `panic` in favor of returning an `error`.
### Security

---
## [1.13.0] - 2024-11-05

### Added
- **FEATURE:** Added documentation for the configuration options: `WithAlphabet`, `WithLengthHint`, and `WithRandReader`.
### Changed
- **DEBT:** ⚡ Performance Enhancements:
  - Unified `allocs/op`: Reduced allocations per operation to 1 for both ASCII and Unicode alphabets regardless of ID length, enhancing memory efficiency across all ID generations.
  - Decreased `ns/op`: Further optimized the Nano ID generation process to lower nanoseconds per operation, resulting in faster ID creation.
### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed various documentation errors.
### Security

---
## [1.12.0] - 2024-11-04

### Added
- **FEATURE**: Added dynamic buffer scaling:
  - Adaptive Buffer Size Growth based on alphabet size while ensuring sufficient randomness for smaller alphabets. 
  - Lower Memory Usage when the context allows, avoiding unnecessarily large buffers when the alphabet or ID length does not demand it.
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.11.0] - 2024-11-02

### Added
### Changed
- **RISK:** Adopted the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/).
### Deprecated
### Removed
### Fixed
### Security

---
## [1.10.2] - 2024-11-02

### Added
- **DEBT:** Added test for invalid UTF8 alphabet checks.
### Changed
- **DEBT:** Updated [README](README.md) to include additional details related to custom alphabet constraints.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.10.1] - 2024-11-02

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed minimum alphabet length logic and related tests and test functions for generating ASCII and Unicode alphabets.
### Security

---
## [1.10.0] - 2024-11-01

### Added
- **FEATURE:** Added new `Must` function to simplify safe initialization of global variables.
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Addressed various documentation issues in the [README](../README.md).
### Security

---
## [1.9.0] - 2024-11-01

### Added
- **FEATURE:** Added new `MustGenerate` and `MustGenerateSize` functions to simplifies safe initialization of global variables.
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Addressed various documentation issues in the [README](../README.md).
### Security

---
## [1.8.2] - 2024-10-31

### Added
### Changed
- **RISK:** Excluded docs and scripts from SonarCloud analysis.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.1] - 2024-10-31

### Added
### Changed
- **DEBT:** Refactored to satisfy static code analysis failure.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.0] - 2024-10-31

### Added
- **FEATURE:** Added support for Unicode alphabets. 
- **FEATURE:** Added support for dynamic `bufferSize` calculation. 
  - The `bufferSize` is calculated by multiplying `bytesNeeded` (the number of bytes required to generate each character) by `bufferMultiplier`. This ensures that the buffer is appropriately sized to handle multiple characters per read, reducing the number of reads from the random source.
### Changed
- **DEBT:** Refactored benchmark tests to support ASCII and Unicode alphabets of varying length and the generation of Nano IDs of varying lengths.
  - Benchmarks the generation of Nano IDs across different alphabet types (`ASCII` and `Unicode`), alphabet lengths (2, 16, 32, 64, 95), and Nano ID lengths (8, 16, 21, 32, 64, 128).
    - `makeASCIIBasedAlphabet(length int) string`: Generates a printable ASCII alphabet of the specified length, starting from `'!'` (33) to `'~'` (126).
    - `makeUnicodeAlphabet(length int) string`: Generates a Unicode alphabet of the specified length, using a range of Unicode characters (e.g., from `'अ'` to `'ह'`).
  - Assesses the performance of Nano ID generation under concurrent (parallel) conditions.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.7.1] - 2024-10-30

### Added
### Changed
- **DEBT:** Refactored Comprehensive Benchmark Tests 
  - Added a suite of benchmark tests (`nanoid_bench_test.go`) within the `nanoid` package. 
  - Benchmarks cover both serial and concurrent generation of Nano IDs. 
  - Included tests for:
    - Default ID generation using `Generate()` and `GenerateSize()`. 
    - Various ID sizes: 8, 16, 21, 32, 64, 128 characters. 
    - Custom alphabets with varying lengths: 2, 16, 32, 64, 95 characters. 
    - Generation of IDs using custom alphabets both serially and in parallel. 
    - Performance measurement of creating new generator instances.
    - The `makeAlphabet` helper function to generate alphabets using only single-byte, printable ASCII characters.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.7.0] - 2024-10-30

### Added
### Changed
- **DEBT:** Optimized Struct Alignment for Reduced Memory Usage: 
  - Reordered fields in the `Config` struct from largest to smallest to minimize padding.
  - Changed data types of `AlphabetLen` and `Step` from `int` to `uint16`, reducing the `Config` struct size from 56 bytes to 32 bytes. 
  - Ensured that the `generator` struct fields are aligned efficiently, reducing its size from 72 bytes to 56 bytes.
- **DEBT:** Integrated IsPowerOfTwo Optimization in Generate Function:
  - Added an `IsPowerOfTwo` field to the `Config` struct to indicate if the alphabet length is a power of two. 
  - During initialization, calculated whether the alphabet length is a power of two. 
  - In the `Generate` function, included logic to skip the boundary check if `int(rnd) < g.config.AlphabetLen` when the alphabet length is a power of two, improving execution speed without increasing allocations.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.6.0] - 2024-10-29

### Added
- **FEATURE:** Added [Nano ID collision calculator](../docs/nanoid-collision-calculator.html).
### Changed
- **DEBT:** Check for duplicate characters using a bitmask with multiple `uint32`s. A `uint32` array can represent `256` bits (`32` bits per `uint32 × 8 = 256`). This allows us to track each possible byte value without the limitations of a single uint64
### Deprecated
### Removed
### Fixed
### Security

---
## [1.5.0] - 2024-10-28

### Added
- **FEATURE**: Added Code of Conduct
- **FEATURE**: Added Contribution Guidelines
### Changed
- **DEBT:** Optimized overall implementation to reduce the allocations per operation to 2.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.4.0] - 2024-10-26

### Added
- **FEATURE:**: Added concurrent benchmark tests.
### Changed
- **DEBT:** Maintained Safety with Linter Suppression: Added `// nolint:gosec` with justification for safe conversions.
- **DEBT:** Refactored Slice Initialization: Initialized `idRunes` with zero length and pre-allocated capacity, using append to build the slice.
- **DEBT:** Ensured Comprehensive Testing: Reviewed and updated tests to handle all edge cases and ensure no runtime errors.
### Deprecated
### Removed
- **FEATURE:** Removed Unicode support for custom dictionaries.
### Fixed
- **DEFECT:** Fixed Operator Precedence: Changed `bits.Len(uint(alphabetLen - 1))` to `bits.Len(uint(alphabetLen) - 1)` to ensure safe conversion.
### Security

---
## [1.3.0] - 2024-10-26

### Added
- **FEATURE:** Added Unicode support for custom dictionaries.
### Changed
- **DEBT:** Modified implementation to be approximately 30% more efficient in terms of CPU complexity. See the `bench` make target.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.2.0] - 2024-10-25

### Added
### Changed
- **DEBT:** Updated Go Report Card links in README. 
### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed version compare links in CHANGELOG.
### Security

---
## [1.0.0] - 2024-10-24

### Added
- **FEATURE:** Initial commit.
### Changed
### Deprecated
### Removed
### Fixed
### Security

[Unreleased]: https://github.com/sixafter/nanoid/compare/v1.57.0...HEAD
[1.57.0]: https://github.com/sixafter/nanoid/compare/v1.56.0...v1.57.0
[1.56.0]: https://github.com/sixafter/nanoid/compare/v1.55.0...v1.56.0
[1.54.0]: https://github.com/sixafter/nanoid/compare/v1.53.0...v1.54.0
[1.53.0]: https://github.com/sixafter/nanoid/compare/v1.52.0...v1.53.0
[1.52.0]: https://github.com/sixafter/nanoid/compare/v1.51.0...v1.52.0
[1.51.0]: https://github.com/sixafter/nanoid/compare/v1.50.0...v1.51.0
[1.50.0]: https://github.com/sixafter/nanoid/compare/v1.49.0...v1.50.0
[1.49.0]: https://github.com/sixafter/nanoid/compare/v1.48.0...v1.49.0
[1.48.0]: https://github.com/sixafter/nanoid/compare/v1.47.0...v1.48.0
[1.47.0]: https://github.com/sixafter/nanoid/compare/v1.46.0...v1.47.0
[1.46.0]: https://github.com/sixafter/nanoid/compare/v1.45.0...v1.46.0
[1.45.0]: https://github.com/sixafter/nanoid/compare/v1.44.0...v1.45.0
[1.44.0]: https://github.com/sixafter/nanoid/compare/v1.43.0...v1.44.0
[1.43.0]: https://github.com/sixafter/nanoid/compare/v1.42.0...v1.43.0
[1.42.0]: https://github.com/sixafter/nanoid/compare/v1.41.0...v1.42.0
[1.41.0]: https://github.com/sixafter/nanoid/compare/v1.40.0...v1.41.0
[1.40.0]: https://github.com/sixafter/nanoid/compare/v1.39.0...v1.40.0
[1.39.0]: https://github.com/sixafter/nanoid/compare/v1.38.0...v1.39.0
[1.38.0]: https://github.com/sixafter/nanoid/compare/v1.37.0...v1.38.0
[1.37.0]: https://github.com/sixafter/nanoid/compare/v1.36.0...v1.37.0
[1.36.0]: https://github.com/sixafter/nanoid/compare/v1.35.0...v1.36.0
[1.35.0]: https://github.com/sixafter/nanoid/compare/v1.34.0...v1.35.0
[1.34.0]: https://github.com/sixafter/nanoid/compare/v1.33.0...v1.34.0
[1.33.0]: https://github.com/sixafter/nanoid/compare/v1.32.1...v1.33.0
[1.32.1]: https://github.com/sixafter/nanoid/compare/v1.32.0...v1.32.1
[1.32.0]: https://github.com/sixafter/nanoid/compare/v1.31.2...v1.32.0
[1.31.2]: https://github.com/sixafter/nanoid/compare/v1.31.1...v1.31.2
[1.31.1]: https://github.com/sixafter/nanoid/compare/v1.31.0...v1.31.1
[1.31.0]: https://github.com/sixafter/nanoid/compare/v1.30.0...v1.31.0
[1.30.0]: https://github.com/sixafter/nanoid/compare/v1.29.0...v1.30.0
[1.29.0]: https://github.com/sixafter/nanoid/compare/v1.28.0...v1.29.0
[1.28.0]: https://github.com/sixafter/nanoid/compare/v1.27.0...v1.28.0
[1.27.0]: https://github.com/sixafter/nanoid/compare/v1.26.0...v1.27.0
[1.26.0]: https://github.com/sixafter/nanoid/compare/v1.25.0...v1.26.0
[1.25.0]: https://github.com/sixafter/nanoid/compare/v1.24.1...v1.25.0
[1.24.1]: https://github.com/sixafter/nanoid/compare/v1.24.0...v1.24.1
[1.24.0]: https://github.com/sixafter/nanoid/compare/v1.23.0...v1.24.0
[1.23.0]: https://github.com/sixafter/nanoid/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/sixafter/nanoid/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/sixafter/nanoid/compare/v1.20.1...v1.21.0
[1.20.1]: https://github.com/sixafter/nanoid/compare/v1.20.0...v1.20.1
[1.20.0]: https://github.com/sixafter/nanoid/compare/v1.19.0...v1.20.0
[1.19.0]: https://github.com/sixafter/nanoid/compare/v1.18.1...v1.19.0
[1.18.1]: https://github.com/sixafter/nanoid/compare/v1.18.0...v1.18.1
[1.18.0]: https://github.com/sixafter/nanoid/compare/v1.17.3...v1.18.0
[1.17.3]: https://github.com/sixafter/nanoid/compare/v1.17.2...v1.17.3
[1.17.2]: https://github.com/sixafter/nanoid/compare/v1.17.1...v1.17.2
[1.17.1]: https://github.com/sixafter/nanoid/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/sixafter/nanoid/compare/v1.16.1...v1.17.0
[1.16.1]: https://github.com/sixafter/nanoid/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/sixafter/nanoid/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/sixafter/nanoid/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/sixafter/nanoid/compare/v1.13.4...v1.14.0
[1.13.4]: https://github.com/sixafter/nanoid/compare/v1.13.3...v1.13.4
[1.13.3]: https://github.com/sixafter/nanoid/compare/v1.13.2...v1.13.3
[1.13.2]: https://github.com/sixafter/nanoid/compare/v1.13.1...v1.13.2
[1.13.1]: https://github.com/sixafter/nanoid/compare/v1.13.0...v1.13.1
[1.13.0]: https://github.com/sixafter/nanoid/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/sixafter/nanoid/compare/v1.11.1...v1.12.0
[1.11.1]: https://github.com/sixafter/nanoid/compare/v1.11.0...v1.11.1
[1.11.0]: https://github.com/sixafter/nanoid/compare/v1.10.2...v1.11.0
[1.10.2]: https://github.com/sixafter/nanoid/compare/v1.10.1...v1.10.2
[1.10.1]: https://github.com/sixafter/nanoid/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/sixafter/nanoid/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/sixafter/nanoid/compare/v1.8.2...v1.9.0
[1.8.2]: https://github.com/sixafter/nanoid/compare/v1.8.1...v1.8.2
[1.8.1]: https://github.com/sixafter/nanoid/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/sixafter/nanoid/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/sixafter/nanoid/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/sixafter/nanoid/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/sixafter/nanoid/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/sixafter/nanoid/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/sixafter/nanoid/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/sixafter/nanoid/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/sixafter/nanoid/compare/v1.0.0...v1.2.0
[1.0.0]: https://github.com/sixafter/nanoid/compare/a6a1eb74b61e518fd0216a17dfe5c9b4c432e6e8...v1.0.0

[MUST]: https://datatracker.ietf.org/doc/html/rfc2119
[MUST NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[MAY]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[REQUIRED]: https://datatracker.ietf.org/doc/html/rfc2119
[RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
[NOT RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
