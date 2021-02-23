# timeutil

This module augments Go's standard library time package.

## Packages

### clock

The clock package provides an abstraction over most of the types and methods in `time` that depend on a current wall time.

### clock/k8sext

This package adapts the [Kubernetes API machinery clock fakes](https://pkg.go.dev/k8s.io/apimachinery/pkg/util/clock) to work with our clock package.

### clockctx

The clockctx package allows a clock fake to be passed through a [context](https://golang.org/pkg/context/). It also provides its own `WithTimeout` and `WithDeadline` functions that use a fake clock.

### backoff

The backoff package contains algorithms for determining how long to wait between attempting work.

### retry

The retry package builds on the backoff package, providing a utility to automatically attempt to perform work multiple times.

### iso8601

The iso8601 package contains types to work with ISO 8601 date/time formats not covered by the standard library: durations, intervals, and recurring intervals.
