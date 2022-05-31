# relspec

This module implements a data substition environment for JSON payloads. It includes a complete opinionated language, Relspec, as well as a number of interchangeable building blocks for customizing an environment to users' needs.

Relspec's goal is to make it easy for a user to write JSON that can be interpolated with external data that they don't know (or want to compute dynamically) at the time the JSON is authored. Its objective is similar to [CUE](https://cuelang.org), but instead of defining a new language, we use special hints in valid JSON to expand and manipulate data.

## Packages

### convert

The convert package contains processors to support conversions between different media formats. It is used by the `convert*` functions in the fnlib package.

### evaluate

The evaluate package contains the primitives of the evaluation engine that recursively processes JSON data. It includes a default evaluator that delegates to a customizable visitor and support for just-in-time node expansion (when implementing `Expandable`).

### fn

The fn package contains primitives for function definition and evaluation. Both the pathlang and relspec packages support calling any function that uses this package's framework.

### fnlib

The fnlib package is an opinionated standard library of functions that can be used by any package that supports `fn.Map`. It is automatically included in the relspec package, but can be overridden if desired.

### pathlang

The pathlang package implements a string-based expression language that can be used in place of the custom JSON tagging scheme to interpolate data into the result document.

### query

The query package provides APIs and several implementations for looking up data under a particular path within a document. It computes lazily so that data missing under a subtree that isn't requested won't be evaluated.

### ref

The ref package provides a mechanism for tracking the data that is used in a document or query. As a fully computed document is constructed, reference data from every subtree is merged together. The result object contains the full set of references needed to build the respective value.

### relspec

The relspec package combines each of the building blocks from the other packages in this module and exposes them through a convenient API.
