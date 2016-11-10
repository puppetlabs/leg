package godat

import (
	"errors"
)

var (
	// Causes a ForEach() or ForEachInto() loop to terminate early.
	ErrStopIteration = errors.New("godat: stop iteration")

	// Raised in a panic() if a function passed to a ForEachInto() loop does
	// not conform to the expected interface.
	ErrInvalidFuncSignature = errors.New("godat: invalid function signature")
)
