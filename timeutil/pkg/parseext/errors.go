package parseext

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrInvalidFormatString    = errors.New("invalid format string")
	ErrFormatStringTooComplex = errors.New("format string is too complex to be parsed")
	ErrInvalidTime            = errors.New("invalid time")
)
