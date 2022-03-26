# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.2] - 2022-03-26

* Adds `CheckNormalizeEngineMount` utility function.

## [0.1.1] - 2022-03-23

### Added

* Adds a new Vault initialization manager. Provides an opinionated initialization function that encompasses all common (configurable) Vault settings (i.e. plugins, secret engines, roles, policies, etc.)
* Adds a Vault kv-v2 client wrapper that provides fluid handling of Vault kv-v2 secrets.

## [0.1.0] - 2022-01-26

### Added

* Initial release.

[Unreleased]: https://github.com/puppetlabs/leg/compare/vaultutil/v0.1.2...HEAD
[0.1.2]: https://github.com/puppetlabs/leg/compare/vaultutil/v0.1.1...vaultutil/v0.1.2
[0.1.1]: https://github.com/puppetlabs/leg/compare/vaultutil/v0.1.0...vaultutil/v0.1.1
[0.1.0]: https://github.com/puppetlabs/leg/compare/7f07fbe0d917993e92fbebfe47490e044bbe16e0...vaultutil/v0.1.0
