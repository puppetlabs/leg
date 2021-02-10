# jsonutil

This module augments Go's native JSON library with additional types and support
for the JSONPath query language.

## Commands

### jpq

This small utility command is useful for testing JSONPath queries. It reads a
JSON document from standard input and runs a query against it read from the
program argument (similar to `jq`).

## Packages

### types

This package provides additional type adapters for Go's JSON library.

### jsonpath

This package provides a query language that is generally compatible with
[JSONPath](https://goessner.net/articles/JsonPath/). It also supports
placeholders in objects, a unique feature derived from the [implementation by
Paessler](https://github.com/PaesslerAG/jsonpath), which can put one or more
components from the path used to reach a selected value into an object key.

### jsonpath/template

This package adapts the JSONPath package to the [format used by
kubectl](https://kubernetes.io/docs/reference/kubectl/jsonpath/), which provides
additional support for looping constructs.
