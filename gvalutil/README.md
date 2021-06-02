# gvalutil

This module provides some extra functionality for working with
[Gval](https://github.com/PaesslerAG/gval), an expression evaluation library.

## Packages

### eval

An alternate implementation of a variable selector that is more strict than
Gval's default.

### langctx

This package provides Go context helpers that can be passed through Gval's
parsing and/or evaluation pipelines.

### langext

This package provides helpers to make writing complex languages less repetitive.

### template

This package provides support for template languages, where one more expression
languages are embedded into an outer string that should not itself be parsed.
