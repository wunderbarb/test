# changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
## [0.7.1] - 2024-11-24
### Changed 
- Removed dependency to `aws-sdk-go`
- Bumped to latest versions
## [0.7.0] - 2024-11-23
### Added
- `SwapCase` randomly swaps the case of the alphabetic characters in a string.  It is useful to test case-insensitive functions.
## [0.6.0] - 2024-07-11
### Changed
- Uses the new `math/rand/v2` that is concurrent-safe.
### Removed
- The type `Option` with its method `WithConcurrentSafe`.
## [0.5.5] - 2023-08-24
### Added
- `RandomEmail` generates a random email address.
### Fixed
- A raced condition with `RandomFileWithDir` when using With.RandomSafe() is fixed.
## [0.5.4] - 2023-08-03
### Removed
- `ErrReader` is now removed. Use `FaultyReader` instead.
## [0.5.3] - 2023-01-23
### Changed
- `FaultyReader` is now compliant with `io.ReadSeekCloser` interface.
## [0.5.2] - 2022-09-8
The initial released version