# errmap

This package focuses on extending [Go 1.13's new error
framework](https://blog.golang.org/go1.13-errors) to make it easier to inject
wrappers throughout the error handling process.

## Packages

### errmap

This package provides functions `MapFirst`, `MapBefore`, `MapAfter`, and
`MapLast` that allow a chain of error mappers to be constructed for a single
error. This makes it easier to control the ordering of error modifications.

For example, it is often desirable to first apply operations that significantly
alter an error (e.g., replacing an internal error with a different kind of
error), then perform non-mutating operations that gather information from the
error, and finally to apply convenient formatting for end-user consumption.

### errmark

This package provides an `errmap`-compatible wrapper that allows errors to be
classified into semantic groups. By marking an error, it can be inspected later
to see if it fits into one or more of these groups and handled as needed.
