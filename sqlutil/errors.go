package sqlutil

import "fmt"

type PanicError struct {
	Cause error
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %+v", e.Cause)
}

type RollbackError struct {
	Trigger error
	Cause   error
}

func (e *RollbackError) Error() string {
	return fmt.Sprintf("failed to roll back: %+v (triggered by error: %+v)", e.Cause, e.Trigger)
}
