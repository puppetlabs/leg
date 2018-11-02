package xflag

import (
	"fmt"
	"net/url"
	"strings"
)

// Validation errors.

// An error representing a mismatch from the URL scheme passed to those
// expected by the flag. Returned by a ValidateURLSchemes validator.
type BadURLSchemeError struct {
	URL    *url.URL
	Wanted []string
}

func (e *BadURLSchemeError) Error() string {
	return fmt.Sprintf("url: invalid scheme: wanted one of '%s', but got '%s'", strings.Join(e.Wanted, "', '"), e.URL.Scheme)
}

// An error denoting that the given URL has a path component when none was
// expected. Returned by the ValidateURLHasNoPath validator.
type URLHasPathError struct {
	URL *url.URL
}

func (e *URLHasPathError) Error() string {
	return fmt.Sprintf("url: invalid path: wanted none, but got '%s'", e.URL.Path)
}

// Validation functions.

// The type of any URL validation function. User-specified validators must
// follow this function format.
type URLValidatorFunc func(u *url.URL) error

// A no-op validator for URLs. Always succeeds.
func ValidateURLNoOp(u *url.URL) error {
	return nil
}

// Creates a validator that checks whether the flag has one of the URL schemes
// passed.
func ValidateURLSchemes(schemes ...string) func(u *url.URL) error {
	if len(schemes) == 0 {
		return ValidateURLNoOp
	}

	return func(u *url.URL) error {
		for _, scheme := range schemes {
			if u.Scheme == scheme {
				return nil
			}
		}

		return &BadURLSchemeError{URL: u, Wanted: schemes}
	}
}

// Asserts that the flag has no path component.
func ValidateURLHasNoPath(u *url.URL) error {
	if strings.TrimLeft(u.Path, "/") != "" {
		return &URLHasPathError{URL: u}
	}

	return nil
}

// URL flag.

type urlValue struct {
	url        **url.URL
	validators []URLValidatorFunc
}

func (uv *urlValue) String() string {
	if uv.url == nil || *uv.url == nil {
		return ""
	}

	return (*uv.url).String()
}

func (uv *urlValue) Set(val string) (err error) {
	u, err := url.Parse(val)
	if err != nil {
		return
	}

	for _, fn := range uv.validators {
		if err = fn(u); err != nil {
			return
		}
	}

	*uv.url = u
	return
}

func (uv *urlValue) Get() interface{} {
	if uv.url == nil {
		return nil
	}

	return *uv.url
}

func newURLValue(value string, p **url.URL, validators []URLValidatorFunc) *urlValue {
	*p = mustParseURL(value)

	return &urlValue{
		url:        p,
		validators: validators,
	}
}

// Defines a URL-typed flag with the specified name, default value, usage
// string, and validators. The argument p points to a *url.URL variable in
// which to store the value of the flag.
//
// The default value must be able to be parsed as a URL (using url.Parse) or
// the call will panic.
func URLVar(p **url.URL, name string, value string, usage string, validators ...URLValidatorFunc) {
	CommandLine.Var(newURLValue(value, p, validators), name, usage)
}

// Defines a URL-typed flag with the specified name, default value, usage
// string, and validators. It returns a pointer to a *url.URL that stores the
// value of the flag.
//
// The default value must be able to be parsed as a URL (using url.Parse) or
// the call will panic.
func URL(name string, value string, usage string, validators ...URLValidatorFunc) **url.URL {
	return CommandLine.URL(name, value, usage, validators...)
}

// Defines a URL-typed flag with the specified name, default value, usage
// string, and validators. The argument p points to a *url.URL variable in
// which to store the value of the flag.
//
// The default value must be able to be parsed as a URL (using url.Parse) or
// the call will panic.
func (f *XFlagSet) URLVar(p **url.URL, name string, value string, usage string, validators ...URLValidatorFunc) {
	f.Var(newURLValue(value, p, validators), name, usage)
}

// Defines a URL-typed flag with the specified name, default value, usage
// string, and validators. It returns a pointer to a *url.URL that stores the
// value of the flag.
//
// The default value must be able to be parsed as a URL (using url.Parse) or
// the call will panic.
func (f *XFlagSet) URL(name string, value string, usage string, validators ...URLValidatorFunc) **url.URL {
	u := &url.URL{}
	p := &u

	f.URLVar(p, name, value, usage, validators...)
	return p
}

// URLs flag.

// An alias for a []*url.URL with some utility methods.
type URLSlice []*url.URL

// Returns true if no flags were passed; false otherwise.
func (us URLSlice) Empty() bool {
	return len(us) == 0
}

// Returns a slice of just the host components of the flags. It is useful in
// conjunction with the ValidateURLSchemes and ValidateURLHasNoPath validators.
//
// Note that (as with the URL.Host field), the host component may contain a
// port, if specified in the flag.
func (us URLSlice) Hosts() []string {
	strs := make([]string, len(us))

	for i, u := range us {
		strs[i] = u.Host
	}

	return strs
}

type urlsValue struct {
	set        bool
	urls       *URLSlice
	validators []URLValidatorFunc
}

func (uv *urlsValue) String() string {
	if uv.urls == nil {
		return ""
	}

	strs := make([]string, len(*uv.urls))

	for i, u := range *uv.urls {
		strs[i] = u.String()
	}

	return strings.Join(strs, ", ")
}

func (uv *urlsValue) Set(val string) error {
	p := &url.URL{}
	u := newURLValue("", &p, uv.validators)
	if err := u.Set(val); err != nil {
		return err
	}

	if !uv.set {
		*uv.urls = URLSlice{}
		uv.set = true
	}

	*uv.urls = append(*uv.urls, p)
	return nil
}

func (uv *urlsValue) Get() interface{} {
	if uv.urls == nil {
		return URLSlice(nil)
	}

	return *uv.urls
}

func newURLsValue(value []string, p *URLSlice, validators []URLValidatorFunc) *urlsValue {
	*p = mustParseURLs(value)

	return &urlsValue{
		urls:       (*URLSlice)(p),
		validators: validators,
	}
}

// Defines a flag representing zero or more URLs mapped to a slice with the
// specified name, default value, usage string, and validators. Validators are
// applied to each instance of the flag individually.
//
// The argument p points to a URLSlice variable in which to store the values of
// the flags.
//
// The default value must be a slice of strings, each of which must be able to
// be parsed as a URL (using url.Parse) or the call will panic. It is
// acceptable to pass nil as the default value.
func URLsVar(p *URLSlice, name string, value []string, usage string, validators ...URLValidatorFunc) {
	CommandLine.Var(newURLsValue(value, p, validators), name, usage)
}

// Defines a flag representing zero or more URLs mapped to a slice with the
// specified name, default value, usage string, and validators. Validators are
// applied to each instance of the flag individually.
//
// It returns a pointer to a URLSlice that stores the value of the flag.
//
// The default value must be a slice of strings, each of which must be able to
// be parsed as a URL (using url.Parse) or the call will panic. It is
// acceptable to pass nil as the default value.
func URLs(name string, value []string, usage string, validators ...URLValidatorFunc) *URLSlice {
	return CommandLine.URLs(name, value, usage, validators...)
}

// Defines a flag representing zero or more URLs mapped to a slice with the
// specified name, default value, usage string, and validators. Validators are
// applied to each instance of the flag individually.
//
// The argument p points to a URLSlice variable in which to store the values of
// the flags.
//
// The default value must be a slice of strings, each of which must be able to
// be parsed as a URL (using url.Parse) or the call will panic. It is
// acceptable to pass nil as the default value.
func (f *XFlagSet) URLsVar(p *URLSlice, name string, value []string, usage string, validators ...URLValidatorFunc) {
	f.Var(newURLsValue(value, p, validators), name, usage)
}

// Defines a flag representing zero or more URLs mapped to a slice with the
// specified name, default value, usage string, and validators. Validators are
// applied to each instance of the flag individually.
//
// It returns a pointer to a URLSlice that stores the value of the flag.
//
// The default value must be a slice of strings, each of which must be able to
// be parsed as a URL (using url.Parse) or the call will panic. It is
// acceptable to pass nil as the default value.
func (f *XFlagSet) URLs(name string, value []string, usage string, validators ...URLValidatorFunc) *URLSlice {
	p := &URLSlice{}
	f.URLsVar(p, name, value, usage, validators...)
	return p
}

func mustParseURL(value string) *url.URL {
	if value == "" {
		return nil
	}

	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}

	return u
}

func mustParseURLs(values []string) URLSlice {
	us := make(URLSlice, len(values))

	for i, value := range values {
		us[i] = mustParseURL(value)
	}

	return us
}
