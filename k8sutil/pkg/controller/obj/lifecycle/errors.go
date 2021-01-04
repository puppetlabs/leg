package lifecycle

import (
	"fmt"
)

// RequiredError is an error reported by RequiredLoader if the target cannot be
// found.
type RequiredError struct {
	// Loader is the underlying loader that could not find the object.
	Loader Loader
}

// Error conforms a pointer to RequiredError to the error interface.
func (e *RequiredError) Error() string {
	return fmt.Sprintf("%T is required", e.Loader)
}
