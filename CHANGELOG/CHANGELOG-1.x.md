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

[Unreleased]: https://github.com/scriptures-social/platform/compare/v1.5.0...HEAD
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
