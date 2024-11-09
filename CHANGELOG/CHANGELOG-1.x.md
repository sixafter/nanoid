# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---
## [Unreleased]

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

---
## [1.16.0] - 2024-NOV-09

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
## [1.15.0] - 2024-NOV-07

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
## [1.14.0] - 2024-NOV-06

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
## [1.13.4] - 2024-NOV-06

### Added
### Changed
- **DEBT:** Added `//go:inline` directive to optimize specific functions.
- **DEBT:** Updated benchmark tests to reflect optimizations.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.13.3] - 2024-NOV-0

### Added
- **DEBT:** Added test for `panic` in `MustWithLength` function.
### Changed
### Deprecated
### Removed
- **DEFECT:** ASCII-only detection now correctly uses `unicode.MaxASCII` instead of `0x7F` to ensure compatibility with all ASCII characters.
### Fixed
### Security

---
## [1.13.2] - 2024-NOV-06

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
## [1.13.1] - 2024-NOV-05

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Refactored ID generation functions `newASCII` and `newUnicode` by removing calls to `panic` in favor of returning an `error`.
### Security

---
## [1.13.0] - 2024-NOV-05

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
## [1.12.0] - 2024-NOV-04

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
## [1.11.0] - 2024-NOV-02

### Added
### Changed
- **RISK:** Adopted the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/).
### Deprecated
### Removed
### Fixed
### Security

---
## [1.10.2] - 2024-NOV-02

### Added
- **DEBT:** Added test for invalid UTF8 alphabet checks.
### Changed
- **DEBT:** Updated [README](README.md) to include additional details related to custom alphabet constraints.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.10.1] - 2024-NOV-02

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed minimum alphabet length logic and related tests and test functions for generating ASCII and Unicode alphabets.
### Security

---
## [1.10.0] - 2024-NOV-01

### Added
- **FEATURE:** Added new `Must` function to simplify safe initialization of global variables.
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Addressed various documentation issues in the [README](../README.md).
### Security

---
## [1.9.0] - 2024-NOV-01

### Added
- **FEATURE:** Added new `MustGenerate` and `MustGenerateSize` functions to simplifies safe initialization of global variables.
### Changed
### Deprecated
### Removed
### Fixed
- **DEFECT:** Addressed various documentation issues in the [README](../README.md).
### Security

---
## [1.8.2] - 2024-OCT-31

### Added
### Changed
- **RISK:** Excluded docs and scripts from SonarCloud analysis.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.1] - 2024-OCT-31

### Added
### Changed
- **DEBT:** Refactored to satisfy static code analysis failure.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.8.0] - 2024-OCT-31

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
## [1.7.1] - 2024-OCT-30

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
## [1.7.0] - 2024-OCT-30

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
## [1.6.0] - 2024-OCT-29

### Added
- **FEATURE:** Added [Nano ID collision calculator](../docs/nanoid-collision-calculator.html).
### Changed
- **DEBT:** Check for duplicate characters using a bitmask with multiple `uint32`s. A `uint32` array can represent `256` bits (`32` bits per `uint32 × 8 = 256`). This allows us to track each possible byte value without the limitations of a single uint64
### Deprecated
### Removed
### Fixed
### Security

---
## [1.5.0] - 2024-OCT-28

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
## [1.4.0] - 2024-OCT-26

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
## [1.3.0] - 2024-OCT-26

### Added
- **FEATURE:** Added Unicode support for custom dictionaries.
### Changed
- **DEBT:** Modified implementation to be approximately 30% more efficient in terms of CPU complexity. See the `bench` make target.
### Deprecated
### Removed
### Fixed
### Security

---
## [1.2.0] - 2024-OCT-25

### Added
### Changed
- **DEBT:** Updated Go Report Card links in README. 
### Deprecated
### Removed
### Fixed
- **DEFECT:** Fixed version compare links in CHANGELOG.
### Security

---
## [1.0.0] - 2024-OCT-24

### Added
- **FEATURE:** Initial commit.
### Changed
### Deprecated
### Removed
### Fixed
### Security

[Unreleased]: https://github.com/scriptures-social/platform/compare/v1.16.0..HEAD
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
