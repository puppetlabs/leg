package xtime

import "errors"

var (
	ErrNotImplemented         = errors.New("xtime: not implemented")
	ErrInvalidFormatString    = errors.New("xtime: invalid format string")
	ErrFormatStringTooComplex = errors.New("xtime: format string is too complex to be parsed")
	ErrInvalidTime            = errors.New("xtime: invalid time")
)
