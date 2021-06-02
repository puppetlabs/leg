package eval

import (
	"fmt"
)

type UnexpectedStringIndexError struct {
	RawIndex string
	Cause    error
}

func (e *UnexpectedStringIndexError) Unwrap() error {
	return e.Cause
}

func (e *UnexpectedStringIndexError) Error() string {
	return fmt.Sprintf("unexpected string index %q for slice, must be convertible to int: %v", e.RawIndex, e.Cause)
}

type UnexpectedIndexTypeError struct {
	RawIndex interface{}
}

func (e *UnexpectedIndexTypeError) Error() string {
	return fmt.Sprintf("unexpected index type %T for slice", e.RawIndex)
}

type IndexOutOfBoundsError struct {
	Index int
}

func (e *IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("index %d out of bounds", e.Index)
}

type UnexpectedKeyTypeError struct {
	RawKey interface{}
}

func (e *UnexpectedKeyTypeError) Error() string {
	return fmt.Sprintf("unexpected key type %T for map", e.RawKey)
}

type UnknownKeyError struct {
	Key string
}

func (e *UnknownKeyError) Error() string {
	return fmt.Sprintf("unknown key %s", e.Key)
}

type UnknownFieldError struct {
	Field string
}

func (e *UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %s", e.Field)
}

type UnsupportedValueTypeError struct {
	Value interface{}
	Field interface{}
}

func (e *UnsupportedValueTypeError) Error() string {
	return fmt.Sprintf("unsupported value type %T for select of field %s", e.Value, e.Field)
}
