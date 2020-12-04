# encoding

## `transfer`

This package provides an interface to encode and decode values that will get
stored in a data store. This is required to ensure that values are consistently
stored safely, but also doesn't enforce an encoding that users must use when
storing secrets or outputs. The default algorithm is base64 and all encoded
strings generated will be prefixed with "base64:". If there is no encoding
prefix, there is `NoEncodingType` that will just return the original value
unencoded/decoded.
