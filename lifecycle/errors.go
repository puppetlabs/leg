package lifecycle

import "fmt"

type PanicError struct {
	Cause error
}

func (pe *PanicError) Error() string {
	return fmt.Sprintf("panic: %+v", pe.Cause)
}

type CloseError struct {
	Causes []error
}

func (ce *CloseError) Error() string {
	msg := "errors occurred when closing resources"
	if len(ce.Causes) > 0 {
		msg += ":"
		for _, cause := range ce.Causes {
			msg += "\n* " + cause.Error()
		}
	}

	return msg
}
