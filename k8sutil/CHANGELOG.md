# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2021-02-02

### Added

* Add [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/)-compatible reconcilers for generating self-signed TLS secrets and automatically setting the CA bundle for webhook configurations.
* Add support for passing event recorders through a context.
* Add object support for webhook configurations.
* Add support for impersonation to end-to-end testing framework.

### Fixed

* `RetryLoader`s now correctly report whether the retry operation succeeded as their boolean outcome.

### Changed

* Upgrade klog to v2.
* Improve and refactor object ownership helpers.

## [0.1.0] - 2021-01-06

### Added

* Initial release.

[Unreleased]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.2.0...HEAD
[0.2.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.1.0...k8sutil/v0.2.0
[0.1.0]: https://github.com/puppetlabs/leg/compare/c09b3cdca7104d5ea79152368de260f5d40316b6...k8sutil/v0.1.0
