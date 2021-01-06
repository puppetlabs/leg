# Leg

This repository contains several Go modules that make up Relay's set of common
utility components. Each module is versioned and managed independently.

## License

Please see the `LICENSE` file in each module for licensing information for that
module. The `LICENSE` file in this directory covers common scripts,
configuration, and documentation (like this file) only.

## Components

* [`datastructure`](datastructure): Support for data structures not part of the
  Go standard library.
* [`encoding`](encoding): Seamless transports for non-Unicode data over strings
  or JSON.
* [`graph`](graph): Directed and undirected graph data structures and
  algorithms.
* [`hashutil`](hashutil): Standardized structures for working with the standard
  library's hashing algorithms.
* [`httputil`](httputil): Standardized structures for HTTP requests and
  responses. Supporting algorithms for working with HTTP servers.
* [`instrumentation`](instrumentation): Integration with error reporting and
  metrics services.
* [`jsonutil`](jsonutil): Extra data types for JSON data.
* [`lifecycle`](lifecycle): Support for running and gracefully stopping
  services within a context boundary.
* [`logging`](logging): Standardized logging interface for Relay projects.
* [`mainutil`](mainutil): Support for managing a set of concurrent processes
  with automatic signal processing and process termination.
* [`mathutil`](mathutil): Additions to the Go standard library's math packages.
* [`netutil`](netutil): Additions to the Go standard library's networking
  packages.
* [`request`](request): Standardized support for passing rudimentary tracing
  information through a context.
* [`scheduler`](scheduler): Advanced management of Goroutines in process pools.
* [`sqlutil`](sqlutil): Additions to the Go standard library's SQL package.
* [`storage`](storage): Standardized interfaces and implementations for working
  with third-party storage services (like S3).
* [`stringutil`](stringutil): Additions to the Go standard library's strings
  package.
* [`timeutil`](timeutil): Additions to the Go standard library's time package,
  including support for ISO 8601 interval types.
* [`workdir`](workdir): Utilities for managing ephemeral or permanent
  application state directories.
