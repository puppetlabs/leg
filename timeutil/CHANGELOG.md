# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2022-01-16

### Changed

* The clock/k8sext package now integrates with clock instances from the k8s.io/utils/clock package instead of the deprecated k8s.io/apimachinery/pkg/util/clock package.

### Build

* The minimum Go version compatible with this release is 1.18.

## [0.4.2] - 2021-07-12

### Fixed

* The clockctx timeout functions now correctly calculate the remaining duration of a timer against a clock's modeled time instead of the current wall time.

## [0.4.1] - 2021-06-23

### Fixed

* `retry.Wait` no longer masks real errors when a backoff error forces the function to exit early.
* Work around a bug in the Kubernetes fake clock implementation that causes a fake timer with zero duration not to fire until the clock is stepped.

## [0.4.0] - 2021-05-15

### Added

* Add helpers to indicate whether work should be retried more clearly.

## [0.3.0] - 2021-02-24

### Added

* Add clockctx package to enable a fake clock to be passed through a context.
* Import the parseext library from Reflect's xparse/xtime package.

## [0.2.0] - 2021-01-06

### Added

* Add backoff, clock, and retry packages.

### Changed

* Restructure the module to use our standard `pkg` package structure.

## [0.1.0] - 2020-12-04

### Changed

* Renamed project to Leg.

[Unreleased]: https://github.com/puppetlabs/leg/compare/timeutil/v0.5.0...HEAD
[0.5.0]: https://github.com/puppetlabs/leg/compare/timeutil/v0.4.2...timeutil/v0.5.0
[0.4.2]: https://github.com/puppetlabs/leg/compare/timeutil/v0.4.1...timeutil/v0.4.2
[0.4.1]: https://github.com/puppetlabs/leg/compare/timeutil/v0.4.0...timeutil/v0.4.1
[0.4.0]: https://github.com/puppetlabs/leg/compare/timeutil/v0.3.0...timeutil/v0.4.0
[0.3.0]: https://github.com/puppetlabs/leg/compare/timeutil/v0.2.0...timeutil/v0.3.0
[0.2.0]: https://github.com/puppetlabs/leg/compare/timeutil/v0.1.0...timeutil/v0.2.0
[0.1.0]: https://github.com/puppetlabs/leg/compare/d290e8e835c3fa3ea4e93073bfe19e1958493d47...timeutil/v0.1.0
