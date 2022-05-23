# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2022-05-23

### Changed

* For added flexibility, the JSONPath template language now passes through the string formatter type from gvalutil.

### Build

* The minimum Go version compatible with this release is 1.18.

## [0.2.2] - 2021-06-02

### Changed

* Update internal implementation to use the gvalutil variable selector.

## [0.2.1] - 2021-05-25

### Changed

* Update internal implementation to use the new gvalutil module.

## [0.2.0] - 2021-02-10

### Added

* Add support for JSONPath and JSONPath templates from our customization of
  [Paessler AG's implementation](https://github.com/PaesslerAG/jsonpath).

### Changed

* Restructure the module to use our standard `pkg` package structure.

## [0.1.0] - 2020-12-04

### Changed

* Renamed project to Leg.

[Unreleased]: https://github.com/puppetlabs/leg/compare/jsonutil/v0.3.0...HEAD
[0.3.0]: https://github.com/puppetlabs/leg/compare/jsonutil/v0.2.2...jsonutil/v0.3.0
[0.2.2]: https://github.com/puppetlabs/leg/compare/jsonutil/v0.2.1...jsonutil/v0.2.2
[0.2.1]: https://github.com/puppetlabs/leg/compare/jsonutil/v0.2.0...jsonutil/v0.2.1
[0.2.0]: https://github.com/puppetlabs/leg/compare/jsonutil/v0.1.0...jsonutil/v0.2.0
[0.1.0]: https://github.com/puppetlabs/leg/compare/d290e8e835c3fa3ea4e93073bfe19e1958493d47...jsonutil/v0.1.0
