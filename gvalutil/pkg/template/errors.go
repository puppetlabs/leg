package template

import "fmt"

type EvaluationError struct {
	Start string
	Cause error
}

func (e *EvaluationError) Unwrap() error {
	return e.Cause
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("evaluating template %q: %v", e.Start, e.Cause)
}
