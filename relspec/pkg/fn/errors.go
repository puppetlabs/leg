package fn

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
)

var (
	ErrFunctionNotFound          = errors.New("function not found")
	ErrArgNotFound               = errors.New("arg not found")
	ErrPositionalArgsNotAccepted = errors.New("positional arguments cannot be used")
	ErrKeywordArgsNotAccepted    = errors.New("keyword arguments cannot be used")
)

type InvocationError struct {
	Name  string
	Cause error
}

var _ jsonpath.PropagatableError = &InvocationError{}

func (e *InvocationError) Unwrap() error {
	return e.Cause
}

func (e *InvocationError) Error() string {
	return fmt.Sprintf("invocation of function %q failed: %+v", e.Name, e.Cause)
}

func (e *InvocationError) Propagate() bool { return true }

type ArityError struct {
	Wanted   []int
	Variadic bool
	Got      int
}

func (e *ArityError) Error() string {
	wanted := make([]string, len(e.Wanted))
	for i, w := range e.Wanted {
		wanted[i] = strconv.FormatInt(int64(w), 10)
	}

	var variadic string
	if e.Variadic {
		variadic = " or more"
	}

	return fmt.Sprintf("unexpected number of arguments: %d (wanted %s%s)", e.Got, strings.Join(wanted, ", "), variadic)
}

type UnexpectedTypeError struct {
	Wanted []reflect.Type
	Got    reflect.Type
}

func (e *UnexpectedTypeError) Error() string {
	wanted := make([]string, len(e.Wanted))
	for i, w := range e.Wanted {
		wanted[i] = w.String()
	}

	return fmt.Sprintf("unexpected type %s (wanted %s)", e.Got, strings.Join(wanted, ", "))
}

type PositionalArgError struct {
	Arg   int
	Cause error
}

func (e *PositionalArgError) Unwrap() error {
	return e.Cause
}

func (e *PositionalArgError) Error() string {
	return fmt.Sprintf("arg %d: %+v", e.Arg, e.Cause)
}

type KeywordArgError struct {
	Arg   string
	Cause error
}

func (e *KeywordArgError) Unwrap() error {
	return e.Cause
}

func (e *KeywordArgError) Error() string {
	return fmt.Sprintf("arg %q: %+v", e.Arg, e.Cause)
}
