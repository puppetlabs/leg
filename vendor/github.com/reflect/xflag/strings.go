package xflag

import (
	"strings"
)

// Strings flag.

type stringsValue struct {
	set    bool
	values *[]string
}

func (sv *stringsValue) String() string {
	if sv.values == nil {
		return ""
	}

	return strings.Join(*sv.values, ",")
}

func (sv *stringsValue) Set(val string) error {
	if !sv.set {
		*sv.values = []string{}
		sv.set = true
	}

	*sv.values = append(*sv.values, val)
	return nil
}

func (sv *stringsValue) Get() interface{} {
	if sv.values == nil {
		return []string(nil)
	}

	return *sv.values
}

func newStringsValue(value []string, p *[]string) *stringsValue {
	*p = value

	return &stringsValue{
		values: p,
	}
}

// Defines a flag representing zero or more strings mapped to a slice with the
// specified name, default value, and usage string.
//
// The argument p points to a string slice in which to store the values of the
// flags.
func StringsVar(p *[]string, name string, value []string, usage string) {
	CommandLine.Var(newStringsValue(value, p), name, usage)
}

// Defines a flag representing zero or more strings mapped to a slice with the
// specified name, default value, and usage string.
//
// It returns a pointer to a string slice that stores the value of the flag.
func Strings(name string, value []string, usage string) *[]string {
	return CommandLine.Strings(name, value, usage)
}

// Defines a flag representing zero or more strings mapped to a slice with the
// specified name, default value, and usage string.
//
// The argument p points to a string slice in which to store the values of the
// flags.
func (f *XFlagSet) StringsVar(p *[]string, name string, value []string, usage string) {
	f.Var(newStringsValue(value, p), name, usage)
}

// Defines a flag representing zero or more strings mapped to a slice with the
// specified name, default value, and usage string.
//
// It returns a pointer to a string slice that stores the value of the flag.
func (f *XFlagSet) Strings(name string, value []string, usage string) *[]string {
	p := &[]string{}
	f.StringsVar(p, name, value, usage)
	return p
}
