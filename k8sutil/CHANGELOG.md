# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

* A new application, portforward, supports programmatic port forwarding to
  arbitrary Kubernetes pods and services.

## [0.6.7] - 2022-03-23

### Added

* Add fake controller client for testing

## [0.6.6] - 2022-03-07

### Fixed

* The `IgnoreNilLabelAnnotatableFrom` `LabelAnnotateFrom` method will redirect to the underlying `LabelAnnotatableFrom` method instead of recursively calling itself.

## [0.6.5] - 2022-02-02

### Added

* Add StatefulSet controller object.

## [0.6.4] - 2021-12-02

### Added

* Add SuffixObjectKey helper function to easily add a suffix to a NamespacedName
  object. Useful for controllers and operators that need to repeatedly prefix a
  set of managed objects.

## [0.6.3] - 2021-11-15

### Added

* Add ClusterRole and ClusterRoleBinding to lifecycle objects.

## [0.6.2] - 2021-11-01

### Changed

* `context.Canceled` will now be propagated automatically in the default error
  handler (the most common spurious error when upgrading a controller).

## [0.6.1] - 2021-09-21

### Changed

* Improve the manifest parsing logic based on standard handling.

## [0.6.0] - 2021-07-16

### Added

* Add support for the `Node` resource.
* Add support for handling Kubernetes qualified names (like `foo.example.com/bar`) in the norm package.

### Changed

* `norm.AnyDNSLabelName` can now make names conform to RFC 1035 or RFC 1123 depending on the application.

### Fixed

* The norm package now handles a variety of edge cases for DNS labels and subdomains correctly.

## [0.5.1] - 2021-07-12

### Build

* Use a distribution of the inlets project managed internally because the upstream version is no longer open-source software.

## [0.5.0] - 2021-07-09

### Added

* Add support for the `VolumeAttachment` resource.
* Add optional cache bypassing Kubernetes client that remains compatible with the controller-runtime client.

## [0.4.1] - 2021-05-25

### Fixed

* Updated inlets image to reflect external changes.

## [0.4.0] - 2021-03-19

### Added

* Add helpers to implement lifecycle methods for Kubernetes API objects.

### Changed

* `helper.Own` and `helper.OwnUncontrolled` now return an error when the owner has not been persisted.

## [0.3.2] - 2021-03-15

### Fixed

* An individual document in a parsed manifest can now be larger than 32kB.

## [0.3.1] - 2021-03-09

### Fixed

* The `AllowAll()` method of a `NetworkPolicy` object now correctly allows both ingress and egress traffic.

## [0.3.0] - 2021-03-09

### Added

* Add a configurable error handling reconciler.
* Port reconciler chaining and a single namespace filtering reconciler from Relay Core.
* Add support for conditionally loading Kubernetes objects based on boolean predicates or the existence of another object.

### Fixed

* Prevent the tunnel application from attempting to connect before the tunnel service endpoint is bound.

## [0.2.1] - 2021-02-24

### Fixed

* The tunnel application now waits for the proxy connection to be established before invoking its callback.

### Build

* Update timeutil package to v0.3.0.

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

[Unreleased]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.7...HEAD
[0.6.7]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.6...k8sutil/v0.6.7
[0.6.6]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.5...k8sutil/v0.6.6
[0.6.5]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.4...k8sutil/v0.6.5
[0.6.4]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.3...k8sutil/v0.6.4
[0.6.3]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.2...k8sutil/v0.6.3
[0.6.2]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.1...k8sutil/v0.6.2
[0.6.1]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.6.0...k8sutil/v0.6.1
[0.6.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.5.1...k8sutil/v0.6.0
[0.5.1]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.5.0...k8sutil/v0.5.1
[0.5.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.4.1...k8sutil/v0.5.0
[0.4.1]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.4.0...k8sutil/v0.4.1
[0.4.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.3.2...k8sutil/v0.4.0
[0.3.2]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.3.1...k8sutil/v0.3.2
[0.3.1]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.3.0...k8sutil/v0.3.1
[0.3.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.2.1...k8sutil/v0.3.0
[0.2.1]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.2.0...k8sutil/v0.2.1
[0.2.0]: https://github.com/puppetlabs/leg/compare/k8sutil/v0.1.0...k8sutil/v0.2.0
[0.1.0]: https://github.com/puppetlabs/leg/compare/c09b3cdca7104d5ea79152368de260f5d40316b6...k8sutil/v0.1.0
