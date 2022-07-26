package evaluate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
)

type UnsupportedValueError struct {
	Type reflect.Type
}

var _ jsonpath.PropagatableError = &UnsupportedValueError{}

func (e *UnsupportedValueError) Error() string {
	return fmt.Sprintf("could not evaluate a value of type %s, must be a YAML-compatible type", e.Type)
}

func (e *UnsupportedValueError) Propagate() bool { return true }

type PathEvaluationError struct {
	Path  string
	Cause error
}

var _ jsonpath.PropagatableError = &PathEvaluationError{}

func (e *PathEvaluationError) trace() ([]string, error) {
	var path []string
	for {
		path = append(path, e.Path)

		en, ok := e.Cause.(*PathEvaluationError)
		if !ok {
			return path, e.Cause
		}

		e = en
	}
}

func (e *PathEvaluationError) UnderlyingCause() error {
	_, err := e.trace()
	return err
}

func (e *PathEvaluationError) Unwrap() error {
	return e.Cause
}

func (e *PathEvaluationError) Error() string {
	path, err := e.trace()
	return fmt.Sprintf("path %q: %+v", strings.Join(path, "."), err)
}

func (e *PathEvaluationError) Propagate() bool { return true }
