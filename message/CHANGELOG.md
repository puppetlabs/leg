# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2021-07-08

### Added

* Add a data deduplication interface with ephemeral and redis backends based on
  the bloom filter data structure for determining that a key has a high likely
  hood that it has been seen before.

[0.1.0]: https://github.com/puppetlabs/leg/compare/b51bab8d7199798fee1629bb2692f2621fcf7fa8...message/v0.1.0
