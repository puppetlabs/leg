package xflag

import (
	"flag"
)

// An XFlagSet is an extension of a flag.FlagSet that supports the additional
// functionality defined in this package. Changes to an XFlagSet modify the
// underlying FlagSet, and it is safe to use pointers to both simultaneously.
type XFlagSet struct {
	*flag.FlagSet
}

// Wraps an existing FlagSet to provide the additional functionality defined in
// this package.
func WrapFlagSet(f *flag.FlagSet) *XFlagSet {
	return &XFlagSet{f}
}

// Creates a new XFlagSet with the given name and error handling scheme. See
// the FlagSet documentation for more information.
func NewXFlagSet(name string, errorHandling flag.ErrorHandling) *XFlagSet {
	return WrapFlagSet(flag.NewFlagSet(name, errorHandling))
}

// An XFlagSet instance bound to the same FlagSet as flag.CommandLine.
var CommandLine = WrapFlagSet(flag.CommandLine)
