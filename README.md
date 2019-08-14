# Horsehead

Named after the [Horsehead Nebula](https://en.wikipedia.org/wiki/Horsehead_Nebula).
This repo provides Go packages that serve has helper functions and utility for
Go-based codebases at Puppet (mostly on the Nebula project).

## workdir package

This package provides utilties for creating and managing working directories.
It defaults to the XDG suite of directory standards from [freedesktop.org](https://www.freedesktop.org/wiki/Software/xdg-user-dirs/).

Help can be found by running `go doc -all github.com/puppetlabs/horsehead/workdir`.

The functionality in this package should work on Linux, MacOS and the BSDs.

### TODO

- add a mechanism for root interactions
- add Windows support

## encoding/transfer

This package provides an interface to encode and decode values that will get stored
in a data store. This is required to ensure that values are consistently stored safely,
but also doesn't enforce an encoding that users must use when storing secrets or outputs.
The default algorithm is base64 and all encoded strings generated will be prefixed with "base64:".
If there is no encoding prefix, there is `NoEncodingType` that will just return the original
value unencoded/decoded.

Help can be found by running `go doc -all github.com/puppetlabs/horsehead/encoding/transfer`.
