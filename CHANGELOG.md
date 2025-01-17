# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.2] - 2022-03-01
### Fixed
- Fix a bug where a single line would cause a panic

## [0.1.1] - 2022-02-09
### Fixed
- Don't consider empty files to have any lines

## [0.1.0] - 2022-02-07
### Changed
- Return error from `Lhdiff`
- CLI prints errors

## [0.0.5] - 2022-02-06
### Changed
- The `Lhdiff` function returns `[][]int`

## [0.0.4] - 2022-02-03
### Fixed
- Fix license owner

## [0.0.3] - 2022-02-03
### Changed
- Change module name from `github.com/aslakhellesoy/lhdiff` to `github.com/oselvar/lhdiff`

## [0.0.2] - 2022-02-03
### Added
- Add `-compact` option to exclude identical lines from output

### Changed
- Print new lines in right file

## [0.0.1]
### Added
- First functional version

[Unreleased]: https://github.com/oselvar/lhdiff/compare/v0.1.2...HEAD
[0.1.2]: https://github.com/oselvar/lhdiff/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/oselvar/lhdiff/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/oselvar/lhdiff/compare/v0.0.5...v0.1.0
[0.0.5]: https://github.com/oselvar/lhdiff/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/oselvar/lhdiff/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/oselvar/lhdiff/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/oselvar/lhdiff/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/oselvar/lhdiff/compare/6084d5de2ec3dbb25767433e79ab840d5941c2de...v0.0.1
