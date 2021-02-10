package jsonpath

import (
	"errors"
	"fmt"
	"strings"
)

// PropagatableError allows an error to be propagated even when a selector would
// otherwise drop it, indicating, e.g., problems with the underlying data.
type PropagatableError interface {
	error
	Propagate() bool
}

var (
	ErrMapperNotCombinable   = errors.New("jsonpath: mapper cannot be combined with range query")
	ErrRangeQueryOutOfBounds = errors.New("jsonpath: range query must be of the format [min:max:step]")
	ErrFilterOutOfBounds     = errors.New("jsonpath: filter needs exactly one key")
)

type UnexpectedSeparatorError struct {
	Separator rune
}

func (e *UnexpectedSeparatorError) Error() string {
	return fmt.Sprintf("jsonpath: unexpected separator %q", e.Separator)
}

type MixedSeparatorError struct {
	A, B rune
}

func (e *MixedSeparatorError) Error() string {
	return fmt.Sprintf("jsonpath: mixed %q and %q in JSON bracket", e.A, e.B)
}

type IndexParseError struct {
	Index int
	Cause error
}

func (e *IndexParseError) Error() string {
	return fmt.Sprintf("jsonpath: could not parse index %d: %+v", e.Index, e.Cause)
}

type KeyParseError struct {
	Key   string
	Cause error
}

func (e *KeyParseError) Error() string {
	return fmt.Sprintf("jsonpath: could not parse key %q: %+v", e.Key, e.Cause)
}

type PathResolutionError struct {
	Path  []string
	Cause error
}

func (e *PathResolutionError) Error() string {
	return fmt.Sprintf("jsonpath: error resolving path %q: %+v", strings.Join(e.Path, "."), e.Cause)
}

type UnexpectedStringIndexError struct {
	RawIndex string
	Cause    error
}

func (e *UnexpectedStringIndexError) Error() string {
	return fmt.Sprintf("jsonpath: unexpected string index %q for slice, must be convertible to int: %v", e.RawIndex, e.Cause)
}

type UnexpectedIndexTypeError struct {
	RawIndex interface{}
}

func (e *UnexpectedIndexTypeError) Error() string {
	return fmt.Sprintf("jsonpath: unexpected index type %T for slice", e.RawIndex)
}

type IndexOutOfBoundsError struct {
	Index int
}

func (e *IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("jsonpath: index %d out of bounds", e.Index)
}

type UnexpectedKeyTypeError struct {
	RawKey interface{}
}

func (e *UnexpectedKeyTypeError) Error() string {
	return fmt.Sprintf("jsonpath: unexpected key type %T for map", e.RawKey)
}

type UnknownKeyError struct {
	Key string
}

func (e *UnknownKeyError) Error() string {
	return fmt.Sprintf("jsonpath: unknown key %s", e.Key)
}

type UnsupportedValueTypeError struct {
	Value interface{}
}

func (e *UnsupportedValueTypeError) Error() string {
	return fmt.Sprintf("jsonpath: unsupported value type %T for select, expected map[string]interface{} or []interface{}", e.Value)
}

type UnknownVariableTypeError struct {
	Variable interface{}
}

var _ PropagatableError = &UnknownVariableTypeError{}

func (e *UnknownVariableTypeError) Error() string {
	return fmt.Sprintf("jsonpath: unknown variable type %T", e.Variable)
}

func (e *UnknownVariableTypeError) Propagate() bool {
	return true
}

type VarSelectorTypeError struct {
	Variable interface{}
}

var _ PropagatableError = &VarSelectorTypeError{}

func (e *VarSelectorTypeError) Error() string {
	return fmt.Sprintf("jsonpath: expected path variable to return values, but got %T", e.Variable)
}

func (e *VarSelectorTypeError) Propagate() bool {
	return true
}
