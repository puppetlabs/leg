package relspec

import (
	"fmt"

	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
)

type InvalidTypeError struct {
	Type  string
	Cause error
}

var _ jsonpath.PropagatableError = &InvalidTypeError{}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("could not evaluate a %s type: %+v", e.Type, e.Cause)
}

func (e *InvalidTypeError) Propagate() bool { return true }

type InvalidEncodingError struct {
	Type  string
	Cause error
}

var _ jsonpath.PropagatableError = &InvalidEncodingError{}

func (e *InvalidEncodingError) Error() string {
	return fmt.Sprintf("could not evaluate encoding %q: %+v", e.Type, e.Cause)
}

func (e *InvalidEncodingError) Propagate() bool { return true }
